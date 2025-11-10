"""
Natural Language Understanding Service
NLU functionality for TabiMoney AI Service
"""

import asyncio
import logging
from typing import Dict, Any, List, Optional, Tuple
import re
from datetime import datetime, timedelta
import json

"""NLU service using Gemini or local Ollama; OpenAI removed."""

import aiohttp
import httpx

from app.core.config import settings
from app.core.database import get_db
from app.models.nlu import ChatRequest, ChatResponse, Entity, NLURequest, NLUResponse
from app.services.transaction_service import TransactionService
from app.utils.json_utils import extract_json_block
from app.utils.llm import call_ollama

logger = logging.getLogger(__name__)


class NLUService:
    """Natural Language Understanding Service"""
    
    def __init__(self):
        self.client: Optional[object] = None
        self.ollama_client: Optional[httpx.AsyncClient] = None
        self.is_initialized = False
        self.transaction_service = TransactionService()
        
        # Common patterns for entity extraction
        self.amount_patterns = [
            r'(\d+(?:\.\d+)?)\s*(?:k|nghìn|thousand)',
            r'(\d+(?:\.\d+)?)\s*(?:tr|triệu|million)',
            r'(\d+(?:\.\d+)?)\s*(?:tỷ|billion)',
            r'(\d+(?:\.\d+)?)\s*(?:đ|dong|vnd)',
            r'(\d+(?:\.\d+)?)'
        ]
        
        self.category_keywords = {
            'ăn uống': ['ăn', 'uống', 'cơm', 'phở', 'bún', 'cháo', 'cà phê', 'trà', 'nước', 'restaurant', 'food'],
            'giao thông': ['xe', 'taxi', 'bus', 'máy bay', 'tàu', 'xăng', 'dầu', 'transport', 'travel'],
            'mua sắm': ['mua', 'sắm', 'quần áo', 'giày', 'túi', 'shop', 'shopping', 'clothes'],
            'giải trí': ['phim', 'game', 'du lịch', 'karaoke', 'bar', 'club', 'entertainment', 'movie'],
            'y tế': ['bệnh viện', 'khám', 'thuốc', 'bác sĩ', 'health', 'medical', 'hospital'],
            'học tập': ['học', 'sách', 'khóa học', 'education', 'study', 'book', 'course'],
            'tiết kiệm': ['tiết kiệm', 'đầu tư', 'savings', 'investment', 'save'],
            'khác': ['khác', 'other', 'misc']
        }
    
    async def initialize(self):
        """Initialize NLU service"""
        logger.info("Initializing NLU Service...")
        
        try:
            if settings.USE_GEMINI and settings.GEMINI_API_KEY:
                # Prefer Gemini if configured
                self.ollama_client = None
                self.client = None
                await self._test_gemini_connection()
                logger.info("NLU Service initialized with Gemini")
            else:
                # Force use Ollama when USE_OPENAI is False
                self.ollama_client = httpx.AsyncClient(timeout=30.0)
                await self._test_ollama_connection()
                logger.info("NLU Service initialized with Ollama (forced)")
            
            self.is_initialized = True
            
        except Exception as e:
            logger.error(f"Failed to initialize NLU Service: {e}")
            # Fallback to rule-based NLU
            self.is_initialized = True
            logger.info("NLU Service initialized with rule-based fallback")
    
    async def cleanup(self):
        """Cleanup NLU service"""
        logger.info("Cleaning up NLU Service...")
        self.client = None
        if self.ollama_client:
            await self.ollama_client.aclose()
        self.is_initialized = False
    
    def is_ready(self) -> bool:
        """Check if NLU service is ready"""
        return self.is_initialized
    
    # OpenAI support removed

    async def _test_gemini_connection(self):
        """Test Gemini API connection"""
        try:
            async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=10)) as session:
                url = f"https://generativelanguage.googleapis.com/v1beta/models/{settings.GEMINI_MODEL}:generateContent?key={settings.GEMINI_API_KEY}"
                payload = {"contents": [{"parts": [{"text": "ping"}]}]}
                async with session.post(url, json=payload) as resp:
                    if resp.status != 200:
                        raise Exception(f"Gemini returned status {resp.status}")
        except Exception as e:
            logger.error(f"Gemini API connection failed: {e}")
            raise
    
    async def _test_ollama_connection(self):
        """Test Ollama connection"""
        try:
            response = await self.ollama_client.post(
                f"{settings.OLLAMA_BASE_URL}/api/generate",
                json={
                    "model": settings.LLM_MODEL,
                    "prompt": "Hello",
                    "stream": False
                }
            )
            if response.status_code == 200:
                logger.info("Ollama connection test successful")
            else:
                raise Exception(f"Ollama returned status {response.status_code}")
        except Exception as e:
            logger.error(f"Ollama connection test failed: {e}")
            raise
    
    async def process_nlu(self, request: NLURequest) -> NLUResponse:
        """Process Natural Language Understanding request"""
        if not self.is_ready():
            raise RuntimeError("NLU Service not ready")
        
        try:
            # Try AI service first if available
            if settings.USE_GEMINI and settings.GEMINI_API_KEY:
                return await self._process_with_gemini(request)
            elif self.ollama_client:
                return await self._process_with_ollama(request)
            else:
                return await self._process_with_rules(request)
                
        except Exception as e:
            logger.error(f"Failed to process NLU: {e}")
            # Fallback to rule-based processing
            return await self._process_with_rules(request)
    
    # OpenAI processing removed
    
    async def _process_with_ollama(self, request: NLURequest) -> NLUResponse:
        """Process NLU using Ollama"""
        try:
            prompt = await self._build_prompt_with_categories(request.text, request.user_id, request.context)
            
            result = await call_ollama(prompt, temperature=0.3, max_tokens=300, format_json=True)
            content_dict = result.get("json") or extract_json_block(result.get("raw", ""))
            if not content_dict:
                raise ValueError("Empty JSON payload from Ollama")

            parsed = self._parse_openai_response(json.dumps(content_dict), request.user_id, request.text)
            parsed = await self._normalize_entities_add_category_id(parsed, request.user_id)
            return parsed
            
        except Exception as e:
            logger.error(f"Ollama NLU processing failed: {e}")
            raise

    async def _process_with_gemini(self, request: NLURequest) -> NLUResponse:
        """Process NLU using Google Gemini"""
        try:
            logger.info(f"Processing with Gemini: {request.text}")
            prompt = await self._build_prompt_with_categories(request.text, request.user_id, request.context)
            async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=60)) as session:
                url = f"https://generativelanguage.googleapis.com/v1beta/models/{settings.GEMINI_MODEL}:generateContent?key={settings.GEMINI_API_KEY}"
                payload = {
                    "contents": [
                        {"parts": [
                            {"text": "Bạn là NLU, luôn trả về CHỈ JSON đúng schema, không văn bản khác."},
                            {"text": prompt}
                        ]}
                    ],
                    "generationConfig": {
                        "temperature": settings.GEMINI_TEMPERATURE,
                        "maxOutputTokens": 4000,  # Increase token limit
                        "response_mime_type": "application/json"
                    }
                }
                async with session.post(url, json=payload) as resp:
                    logger.info(f"Gemini response status: {resp.status}")
                    if resp.status != 200:
                        logger.error(f"Gemini returned status {resp.status}")
                        # fallback to Ollama
                        if self.ollama_client:
                            return await self._process_with_ollama(request)
                        raise Exception(f"Gemini returned status {resp.status}")
                    data = await resp.json()
                    logger.info(f"Gemini raw response: {json.dumps(data, indent=2)[:500]}...")
                    # Extract text
                    content = ""
                    try:
                        if "candidates" in data and data["candidates"]:
                            candidate = data["candidates"][0]
                            if "content" in candidate and "parts" in candidate["content"]:
                                content = candidate["content"]["parts"][0].get("text", "")
                            elif "text" in candidate:
                                content = candidate["text"]
                        elif "text" in data:
                            content = data["text"]
                        else:
                            content = json.dumps(data)
                    except Exception as e:
                        logger.error(f"Failed to extract Gemini content: {e}")
                        content = json.dumps(data)

                    try:
                        content_dict = extract_json_block(content)
                        parsed = self._parse_openai_response(json.dumps(content_dict), request.user_id, request.text)
                        parsed = await self._normalize_entities_add_category_id(parsed, request.user_id)
                        logger.info("Gemini parsing successful: intent=%s", parsed.intent)
                        return parsed
                    except Exception as e:
                        logger.warning("Falling back to Ollama due to parse failure (Gemini): %s", e)
                        return await self._process_with_ollama(request)
        except Exception as e:
            logger.error(f"Gemini NLU processing failed: {e}")
            # fallback to Ollama
            if self.ollama_client:
                return await self._process_with_ollama(request)
            raise
    
    async def _build_prompt_with_categories(self, text: str, user_id: int, context: str) -> str:
        """Build prompt including allowed categories (id|name) and strict schema returning category_id."""
        # Query top used categories for this user (and include system categories)
        allowed: List[Dict[str, Any]] = []
        try:
            async with get_db() as db:
                rows = await db.execute(
                    (
                        "SELECT c.id, c.name, c.name_en, COALESCE(COUNT(t.id),0) as usage_count "
                        "FROM categories c "
                        "LEFT JOIN transactions t ON t.category_id = c.id AND t.user_id = %s "
                        "WHERE c.is_active = TRUE AND (c.user_id IS NULL OR c.user_id = %s) "
                        "GROUP BY c.id, c.name, c.name_en "
                        "ORDER BY usage_count DESC, c.sort_order ASC, c.id ASC "
                        "LIMIT 30"
                    ),
                    (user_id, user_id)
                )
                allowed = rows or []
        except Exception as e:
            logger.warning(f"Failed to fetch categories for prompt: {e}")
            allowed = []

        def _fmt_cat(cat: Dict[str, Any]) -> str:
            name = cat.get("name") or ""
            name_en = cat.get("name_en") or ""
            if name_en and name_en.lower() != name.lower():
                return f"{cat['id']}|{name} ({name_en})"
            return f"{cat['id']}|{name}"

        allowed_lines = "\n".join(f"- {_fmt_cat(c)}" for c in allowed)

        # AI Agent tự chủ hoàn toàn
        prompt = f"""Bạn là AI Agent quản lý tài chính cá nhân. Phân tích câu và quyết định hành động:

Database có sẵn:
- Danh mục: {allowed_lines if allowed_lines else "9|Khác"}
- Bạn có thể: tạo transaction, query balance, phân tích dữ liệu

Schema:
{{"intent":"add_transaction|query_balance|analyze_data|budget_management|goal_tracking|smart_recommendations|expense_forecasting|general","entities":[{{"type":"amount|category_id|date|description|budget_amount|goal_amount","value":"...","confidence":0.9,"start_pos":0,"end_pos":5}}],"confidence":0.9,"needs_confirmation":false,"response":"Phản hồi tự nhiên","action":"mô tả hành động sẽ thực hiện"}}

QUAN TRỌNG - Parse amount:
- Khi extract amount entity, PHẢI chuyển đổi về số VND đầy đủ:
  * "16 triệu" hoặc "16 tr" → value: "16000000"
  * "400k" hoặc "400 nghìn" → value: "400000"
  * "2 tỷ" → value: "2000000000"
  * "100000" → value: "100000" (giữ nguyên)
- KHÔNG bao giờ trả về value là "16" nếu trong câu có "16 triệu"
- value phải là số thuần túy (không có text, không có đơn vị)

Quyết định:
- Tự phân tích intent từ context
- Tự chọn category phù hợp nhất
- Tự quyết định có cần xác nhận không
- Tự mô tả action sẽ làm
- Tự parse amount về VND đầy đủ

Input: "{text}"
Output:"""

        return prompt
    
    async def _process_with_rules(self, request: NLURequest) -> NLUResponse:
        """Process NLU using rule-based approach as fallback"""
        try:
            entities = self._extract_entities_rule_based(request.text)
            intent = self._determine_intent_rule_based(request.text, entities)
            response = self._generate_response_rule_based(intent, entities)
            suggested_action = self._get_suggested_action(intent)
            
            return NLUResponse(
                user_id=request.user_id,
                intent=intent,
                entities=entities,
                confidence=0.6,
                suggested_action=suggested_action,
                response=response,
                generated_at=datetime.now().isoformat() + "Z"
            )
        except Exception as e:
            logger.error(f"Rule-based NLU processing failed: {e}")
            # Return minimal fallback
            return NLUResponse(
                user_id=request.user_id,
                intent="general",
                entities=[],
                confidence=0.0,
                suggested_action="general_response",
                response="Xin lỗi, tôi không hiểu yêu cầu của bạn.",
                generated_at=datetime.now().isoformat() + "Z"
            )
    
    def _parse_openai_response(self, content: str, user_id: int, original_text: str = "") -> NLUResponse:
        """Parse OpenAI response"""
        try:
            import json
            data = json.loads(content)
            
            entities = []
            for entity_data in data.get('entities', []):
                value = entity_data.get('value')
                ent_type = entity_data.get('type')
                # AI model should already return normalized amount value (in VND)
                # Just ensure it's a valid number string
                if ent_type == 'amount' and isinstance(value, str):
                    # Try to parse as float to validate
                    try:
                        float(value)
                        # If it's a valid number, keep it as is (AI already normalized it)
                    except ValueError:
                        # If not a pure number, try to parse it (fallback for older AI responses)
                        num_match = re.search(r"\d+(?:[\.,]\d+)?", value)
                        if num_match:
                            num_text = num_match.group(0).replace(',', '.')
                            parsed_amount = self._parse_amount(num_text, value)
                            if parsed_amount > 0:
                                value = str(parsed_amount)
                entity = Entity(
                    type=ent_type,
                    value=value,
                    confidence=entity_data.get('confidence', 0.0),
                    start_pos=entity_data.get('start_pos', 0),
                    end_pos=entity_data.get('end_pos', 0)
                )
                entities.append(entity)
            
            return NLUResponse(
                user_id=user_id,
                intent=data.get('intent', 'general'),
                entities=entities,
                confidence=data.get('confidence', 0.0),
                suggested_action=data.get('suggested_action', 'general_response'),
                response=data.get('response', 'Xin lỗi, tôi không hiểu yêu cầu của bạn.'),
                needs_confirmation=data.get('needs_confirmation', False),
                generated_at=datetime.now().isoformat() + "Z"
            )
            
        except Exception as e:
            logger.error(f"Failed to parse OpenAI response: {e}")
            raise

    async def _normalize_entities_add_category_id(self, nlu: NLUResponse, user_id: int) -> NLUResponse:
        """If only a 'category' name is present, resolve and append 'category_id', then remove 'category'."""
        try:
            has_category_id = any(e.type == 'category_id' for e in nlu.entities)
            category_name_entity = next((e for e in nlu.entities if e.type == 'category'), None)
            
            # If already has category_id, remove any category entities
            if has_category_id:
                nlu.entities = [e for e in nlu.entities if e.type != 'category']
                logger.info(f"Removed category entities, kept category_id")
                return nlu
                
            if not category_name_entity or not isinstance(category_name_entity.value, str):
                return nlu
            # resolve via DB
            name_l = category_name_entity.value.strip().lower()
            async with get_db() as db:
                rows = await db.execute(
                    (
                        "SELECT c.id, c.name, COALESCE(c.name_en, '') as name_en FROM categories c "
                        "WHERE c.is_active = TRUE AND (c.user_id IS NULL OR c.user_id = %s)"
                    ),
                    (user_id,)
                )
            resolved_id: Optional[int] = None
            if rows:
                exact_vi = next((r for r in rows if r['name'].strip().lower() == name_l), None)
                if exact_vi:
                    resolved_id = int(exact_vi['id'])
                if resolved_id is None:
                    partial_vi = next((r for r in rows if name_l in r['name'].strip().lower()), None)
                    if partial_vi:
                        resolved_id = int(partial_vi['id'])
                if resolved_id is None:
                    exact_en = next((r for r in rows if r['name_en'] and r['name_en'].strip().lower() == name_l), None)
                    if exact_en:
                        resolved_id = int(exact_en['id'])
                if resolved_id is None:
                    partial_en = next((r for r in rows if r['name_en'] and name_l in r['name_en'].strip().lower()), None)
                    if partial_en:
                        resolved_id = int(partial_en['id'])
                if resolved_id is None:
                    other = next((r for r in rows if r['name'].strip().lower() == 'khác'), None)
                    if other:
                        resolved_id = int(other['id'])
            if resolved_id is not None:
                # Remove the old category entity and add category_id
                nlu.entities = [e for e in nlu.entities if e.type != 'category']
                nlu.entities.append(
                    Entity(
                        type='category_id',
                        value=str(resolved_id),
                        confidence=category_name_entity.confidence or 0.7,
                        start_pos=category_name_entity.start_pos,
                        end_pos=category_name_entity.end_pos,
                    )
                )
        except Exception as e:
            logger.warning(f"normalize_entities_add_category_id failed: {e}")
        return nlu

    def _extract_json_block(self, content: str) -> str:
        """Deprecated helper kept for backward compatibility."""
        parsed = extract_json_block(content)
        if not parsed:
            raise ValueError("empty content")
        return json.dumps(parsed)
    
    def _extract_entities_rule_based(self, text: str) -> List[Entity]:
        """Extract entities using rule-based approach"""
        entities = []
        
        # Extract amounts
        for pattern in self.amount_patterns:
            matches = re.finditer(pattern, text, re.IGNORECASE)
            for match in matches:
                amount_text = match.group(1)
                amount_value = self._parse_amount(amount_text, text[match.start():match.end()])
                
                entity = Entity(
                    type="amount",
                    value=str(amount_value),
                    confidence=0.8,
                    start_pos=match.start(),
                    end_pos=match.end()
                )
                entities.append(entity)
        
        # Extract categories
        for category, keywords in self.category_keywords.items():
            for keyword in keywords:
                if keyword in text:
                    start_pos = text.find(keyword)
                    # Map category name to ID
                    category_id = self._get_category_id(category)
                    entity = Entity(
                        type="category_id",
                        value=str(category_id),
                        confidence=0.7,
                        start_pos=start_pos,
                        end_pos=start_pos + len(keyword)
                    )
                    entities.append(entity)
                    break
        
        # Extract dates
        date_patterns = [
            r'hôm nay|today',
            r'hôm qua|yesterday',
            r'ngày mai|tomorrow',
            r'(\d{1,2})[\/\-](\d{1,2})[\/\-](\d{4})',  # DD/MM/YYYY
            r'(\d{1,2})[\/\-](\d{1,2})'  # DD/MM
        ]
        
        for pattern in date_patterns:
            matches = re.finditer(pattern, text, re.IGNORECASE)
            for match in matches:
                date_value = self._parse_date(match.group(0))
                entity = Entity(
                    type="date",
                    value=date_value,
                    confidence=0.6,
                    start_pos=match.start(),
                    end_pos=match.end()
                )
                entities.append(entity)
        
        return entities
    
    def _get_category_id(self, category_name: str) -> int:
        """Map category name to ID"""
        category_mapping = {
            "Ăn uống": 1,
            "Giao thông": 2,
            "Mua sắm": 3,
            "Giải trí": 4,
            "Sức khỏe": 5,
            "Học tập": 6,
            "Du lịch": 7,
            "Thu nhập": 8,
            "Khác": 9
        }
        return category_mapping.get(category_name, 9)  # Default to "Khác"
    
    def _determine_intent_rule_based(self, text: str, entities: List[Entity]) -> str:
        """Determine intent using rule-based approach"""
        # Check for query keywords first (more specific)
        query_keywords = ['bao nhiêu', 'tổng', 'tổng cộng', 'chi tiêu', 'thu nhập', 'số dư', 'muốn biết', 'kiểm tra']
        if any(keyword in text for keyword in query_keywords):
            return "query_balance"
        
        # Check for goals keywords
        goal_keywords = ['mục tiêu', 'goal', 'tiết kiệm', 'đầu tư', 'mua sắm', 'trả nợ']
        if any(keyword in text for keyword in goal_keywords):
            if 'tạo' in text or ('thêm' in text and 'mới' in text):
                return "create_goal"
            elif 'xem' in text or 'danh sách' in text or 'list' in text:
                return "list_goals"
            elif 'cập nhật' in text or 'đóng góp' in text or ('thêm' in text and 'tiền' in text) or ('tiết kiệm' in text and 'thêm' in text):
                return "update_goal"
            else:
                return "list_goals"  # Default to list if unclear
        
        # Check for budget keywords
        budget_keywords = ['ngân sách', 'budget', 'hạn mức', 'chi tiêu']
        if any(keyword in text for keyword in budget_keywords):
            if 'tạo' in text or ('thêm' in text and 'mới' in text):
                return "create_budget"
            elif 'xem' in text or 'danh sách' in text or 'list' in text or 'kiểm tra' in text:
                return "list_budgets"
            elif 'cập nhật' in text or 'trạng thái' in text or 'tình hình' in text:
                return "update_budget"
            else:
                return "list_budgets"  # Default to list if unclear
        
        # Check for transaction-related keywords
        transaction_keywords = ['mua', 'mua', 'chi', 'tiêu', 'ăn', 'uống', 'đi', 'mua sắm']
        if any(keyword in text for keyword in transaction_keywords):
            return "add_transaction"
        
        # Check for question keywords
        question_keywords = ['tại sao', 'như thế nào', 'làm sao', 'có thể', 'có nên']
        if any(keyword in text for keyword in question_keywords):
            return "ask_question"
        
        # Default intent
        return "general"
    
    def _generate_response_rule_based(self, intent: str, entities: List[Entity]) -> str:
        """Generate response using rule-based approach"""
        if intent == "add_transaction":
            amount_entity = next((e for e in entities if e.type == "amount"), None)
            category_entity = next((e for e in entities if e.type == "category"), None)
            
            if amount_entity and category_entity:
                return f"Tôi sẽ giúp bạn thêm giao dịch {amount_entity.value} VND cho danh mục {category_entity.value}."
            elif amount_entity:
                return f"Tôi sẽ giúp bạn thêm giao dịch {amount_entity.value} VND."
            else:
                return "Tôi sẽ giúp bạn thêm giao dịch mới."
        
        elif intent == "query_balance":
            return "Tôi sẽ kiểm tra số dư và chi tiêu của bạn. Bạn có thể xem chi tiết trong trang Analytics hoặc Dashboard."
        
        elif intent == "ask_question":
            return "Tôi sẽ cố gắng trả lời câu hỏi của bạn về tài chính."
        
        else:
            return "Tôi có thể giúp bạn quản lý tài chính cá nhân."
    
    def _get_suggested_action(self, intent: str) -> str:
        """Get suggested action based on intent"""
        action_map = {
            "add_transaction": "create_transaction",
            "query_balance": "get_balance",
            "ask_question": "answer_question",
            "general": "general_response"
        }
        return action_map.get(intent, "general_response")
    
    def _parse_amount(self, amount_text: str, full_match: str) -> float:
        """Parse amount from text"""
        try:
            amount = float(amount_text)
            
            # Handle Vietnamese number suffixes
            if 'k' in full_match.lower() or 'nghìn' in full_match.lower():
                amount *= 1000
            elif 'triệu' in full_match.lower() or 'tr' in full_match.lower():
                amount *= 1000000
            elif 'tỷ' in full_match.lower():
                amount *= 1000000000
            
            return amount
        except ValueError:
            return 0.0
    
    def _parse_date(self, date_text: str) -> str:
        """Parse date from text"""
        today = datetime.now()
        
        if 'hôm nay' in date_text.lower() or 'today' in date_text.lower():
            return today.strftime('%Y-%m-%d')
        elif 'hôm qua' in date_text.lower() or 'yesterday' in date_text.lower():
            yesterday = today - timedelta(days=1)
            return yesterday.strftime('%Y-%m-%d')
        elif 'ngày mai' in date_text.lower() or 'tomorrow' in date_text.lower():
            tomorrow = today + timedelta(days=1)
            return tomorrow.strftime('%Y-%m-%d')
        else:
            # Try to parse DD/MM/YYYY or DD/MM
            import re
            match = re.search(r'(\d{1,2})[\/\-](\d{1,2})(?:[\/\-](\d{4}))?', date_text)
            if match:
                day, month, year = match.groups()
                year = year or today.year
                return f"{year}-{month.zfill(2)}-{day.zfill(2)}"
        
        return today.strftime('%Y-%m-%d')
    
    async def process_chat(self, request: ChatRequest) -> ChatResponse:
        """Process chat message and return AI response"""
        if not self.is_ready():
            raise RuntimeError("NLU Service not ready")
        
        try:
            # Convert chat request to NLU request
            nlu_request = NLURequest(
                text=request.message,
                user_id=request.user_id,
                context="chat"
            )
            
            # Process with NLU
            nlu_response = await self.process_nlu(nlu_request)
            
            # AI tự quyết định action dựa trên intent
            if nlu_response.intent == "add_transaction" and not nlu_response.needs_confirmation:
                await self._handle_add_transaction(request.user_id, nlu_response)
            elif nlu_response.intent == "query_balance":
                await self._handle_query_balance(request.user_id, nlu_response)
            elif nlu_response.intent == "analyze_data":
                await self._handle_analyze_data(request.user_id, nlu_response)
            elif nlu_response.intent == "budget_management":
                await self._handle_budget_management(request.user_id, nlu_response)
            elif nlu_response.intent == "goal_tracking":
                await self._handle_goal_tracking(request.user_id, nlu_response)
            elif nlu_response.intent == "smart_recommendations":
                await self._handle_smart_recommendations(request.user_id, nlu_response)
            elif nlu_response.intent == "expense_forecasting":
                await self._handle_expense_forecasting(request.user_id, nlu_response)
            
            # Generate suggestions based on intent
            suggestions = self._generate_chat_suggestions(nlu_response.intent)
            
            return ChatResponse(
                user_id=request.user_id,
                message=request.message,
                response=nlu_response.response,
                intent=nlu_response.intent,
                entities=nlu_response.entities,
                suggestions=suggestions,
                generated_at=datetime.now().isoformat() + "Z"
            )
            
        except Exception as e:
            logger.error(f"Failed to process chat: {e}")
            # Return fallback response
            return ChatResponse(
                user_id=request.user_id,
                message=request.message,
                response="Xin lỗi, tôi không thể xử lý tin nhắn của bạn lúc này. Vui lòng thử lại sau.",
                intent="error",
                entities=[],
                suggestions=["Thử hỏi về chi tiêu", "Kiểm tra số dư", "Thêm giao dịch mới"],
                generated_at=datetime.now().isoformat() + "Z"
            )
    
    async def _handle_add_transaction(self, user_id: int, nlu_response: NLUResponse):
        """Handle adding transaction directly to database"""
        try:
            # Extract entities
            amount = None
            category_id = None
            
            for entity in nlu_response.entities:
                if entity.type == "amount":
                    amount = float(entity.value)
                elif entity.type == "category_id":
                    category_id = int(entity.value)
            
            if amount and category_id:
                # Determine transaction type based on category
                transaction_type = "expense"
                if category_id == 8:  # Thu nhập category
                    transaction_type = "income"
                
                # Create transaction
                result = await self.transaction_service.create_transaction(
                    user_id=user_id,
                    category_id=category_id,
                    amount=amount,
                    description=nlu_response.response,  # Use AI response as description
                    transaction_type=transaction_type
                )
                
                if result["success"]:
                    # Update the response with success message
                    nlu_response.response = result["message"]
                    logger.info(f"Successfully created transaction for user {user_id}")
                else:
                    # Update with error message
                    nlu_response.response = result["message"]
                    logger.error(f"Failed to create transaction: {result.get('error')}")
            else:
                logger.warning(f"Missing required entities for transaction: amount={amount}, category_id={category_id}")
                
        except Exception as e:
            logger.error(f"Error handling add transaction: {e}")
            nlu_response.response = "Có lỗi xảy ra khi thêm giao dịch."
    
    async def _handle_query_balance(self, user_id: int, nlu_response: NLUResponse):
        """Handle balance query directly"""
        try:
            result = await self.transaction_service.get_user_balance(user_id)
            if result["success"]:
                nlu_response.response = result["message"]
                logger.info(f"Successfully retrieved balance for user {user_id}")
            else:
                nlu_response.response = result["message"]
                logger.error(f"Failed to get balance: {result.get('error')}")
        except Exception as e:
            logger.error(f"Error handling balance query: {e}")
            nlu_response.response = "Có lỗi xảy ra khi lấy thông tin số dư."
    
    async def _handle_analyze_data(self, user_id: int, nlu_response: NLUResponse):
        """AI tự phân tích dữ liệu theo yêu cầu"""
        try:
            # AI có thể tự quyết định phân tích gì dựa trên context
            # Ví dụ: spending patterns, category analysis, trends, etc.
            
            # Lấy dữ liệu giao dịch gần đây
            async with get_db() as db:
                # Lấy transactions 30 ngày gần nhất
                transactions_query = """
                SELECT t.*, c.name as category_name 
                FROM transactions t 
                JOIN categories c ON t.category_id = c.id 
                WHERE t.user_id = %s 
                AND t.transaction_date >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
                ORDER BY t.transaction_date DESC
                LIMIT 50
                """
                
                transactions = await db.execute(transactions_query, (user_id,))
                
                if transactions:
                    # AI tự phân tích và đưa ra insights
                    total_expense = sum(t['amount'] for t in transactions if t['transaction_type'] == 'expense')
                    total_income = sum(t['amount'] for t in transactions if t['transaction_type'] == 'income')
                    
                    # Phân tích theo category
                    category_spending = {}
                    for t in transactions:
                        if t['transaction_type'] == 'expense':
                            cat = t['category_name']
                            category_spending[cat] = category_spending.get(cat, 0) + t['amount']
                    
                    # Tạo response tự nhiên
                    insights = []
                    if total_expense > 0:
                        insights.append(f"Tổng chi tiêu 30 ngày: {total_expense:,.0f} VND")
                    if total_income > 0:
                        insights.append(f"Tổng thu nhập 30 ngày: {total_income:,.0f} VND")
                    
                    if category_spending:
                        top_category = max(category_spending.items(), key=lambda x: x[1])
                        insights.append(f"Chi nhiều nhất: {top_category[0]} ({top_category[1]:,.0f} VND)")
                    
                    nlu_response.response = "Phân tích 30 ngày gần nhất:\n" + "\n".join(insights)
                else:
                    nlu_response.response = "Chưa có dữ liệu giao dịch để phân tích."
                    
        except Exception as e:
            logger.error(f"Error handling data analysis: {e}")
            nlu_response.response = "Có lỗi xảy ra khi phân tích dữ liệu."
    
    async def _handle_budget_management(self, user_id: int, nlu_response: NLUResponse):
        """AI quản lý ngân sách thông minh"""
        try:
            async with get_db() as db:
                # Lấy ngân sách hiện tại
                budget_query = """
                SELECT b.*, c.name as category_name 
                FROM budgets b 
                JOIN categories c ON b.category_id = c.id 
                WHERE b.user_id = %s AND b.is_active = true
                """
                budgets = await db.execute(budget_query, (user_id,))
                
                # Lấy chi tiêu tháng này
                expense_query = """
                SELECT c.name as category_name, SUM(t.amount) as total_spent
                FROM transactions t 
                JOIN categories c ON t.category_id = c.id 
                WHERE t.user_id = %s AND t.transaction_type = 'expense' 
                AND MONTH(t.transaction_date) = MONTH(CURDATE()) 
                AND YEAR(t.transaction_date) = YEAR(CURDATE())
                GROUP BY c.id, c.name
                """
                expenses = await db.execute(expense_query, (user_id,))
                
                if budgets:
                    insights = []
                    expense_dict = {e['category_name']: e['total_spent'] for e in expenses}
                    
                    for budget in budgets:
                        spent = expense_dict.get(budget['category_name'], 0)
                        remaining = budget['amount'] - spent
                        percentage = (spent / budget['amount']) * 100 if budget['amount'] > 0 else 0
                        
                        if percentage > 90:
                            insights.append(f"⚠️ {budget['category_name']}: {percentage:.0f}% ngân sách ({spent:,.0f}/{budget['amount']:,.0f} VND)")
                        elif percentage > 70:
                            insights.append(f"🟡 {budget['category_name']}: {percentage:.0f}% ngân sách ({spent:,.0f}/{budget['amount']:,.0f} VND)")
                        else:
                            insights.append(f"✅ {budget['category_name']}: {percentage:.0f}% ngân sách ({spent:,.0f}/{budget['amount']:,.0f} VND)")
                    
                    nlu_response.response = "📊 Tình hình ngân sách tháng này:\n" + "\n".join(insights)
                else:
                    nlu_response.response = "Bạn chưa có ngân sách nào. Hãy tạo ngân sách để quản lý chi tiêu tốt hơn!"
                    
        except Exception as e:
            logger.error(f"Error handling budget management: {e}")
            nlu_response.response = "Có lỗi xảy ra khi kiểm tra ngân sách."
    
    async def _handle_goal_tracking(self, user_id: int, nlu_response: NLUResponse):
        """AI theo dõi mục tiêu tài chính"""
        try:
            async with get_db() as db:
                # Lấy goals hiện tại
                goals_query = """
                SELECT * FROM financial_goals 
                WHERE user_id = %s AND is_achieved = false
                ORDER BY target_date ASC
                """
                goals = await db.execute(goals_query, (user_id,))
                
                if goals:
                    insights = []
                    for goal in goals:
                        # Tính progress
                        progress_query = """
                        SELECT SUM(amount) as saved_amount
                        FROM transactions 
                        WHERE user_id = %s AND category_id = 8 AND transaction_type = 'income'
                        AND transaction_date >= %s
                        """
                        progress = await db.execute(progress_query, (user_id, goal['created_at']))
                        saved = progress[0]['saved_amount'] if progress and progress[0]['saved_amount'] else 0
                        
                        progress_percentage = (saved / goal['target_amount']) * 100 if goal['target_amount'] > 0 else 0
                        remaining = goal['target_amount'] - saved
                        
                        insights.append(f"🎯 {goal['title']}: {progress_percentage:.0f}% ({saved:,.0f}/{goal['target_amount']:,.0f} VND)")
                        if remaining > 0:
                            insights.append(f"   Còn lại: {remaining:,.0f} VND")
                    
                    nlu_response.response = "🎯 Tiến độ mục tiêu:\n" + "\n".join(insights)
                else:
                    nlu_response.response = "Bạn chưa có mục tiêu tài chính nào. Hãy tạo mục tiêu để có động lực tiết kiệm!"
                    
        except Exception as e:
            logger.error(f"Error handling goal tracking: {e}")
            nlu_response.response = "Có lỗi xảy ra khi kiểm tra mục tiêu."
    
    async def _handle_smart_recommendations(self, user_id: int, nlu_response: NLUResponse):
        """AI đưa ra gợi ý thông minh"""
        try:
            async with get_db() as db:
                # Phân tích chi tiêu 3 tháng gần nhất
                analysis_query = """
                SELECT c.name as category_name, 
                       SUM(t.amount) as total_spent,
                       COUNT(*) as transaction_count,
                       AVG(t.amount) as avg_amount
                FROM transactions t 
                JOIN categories c ON t.category_id = c.id 
                WHERE t.user_id = %s AND t.transaction_type = 'expense' 
                AND t.transaction_date >= DATE_SUB(CURDATE(), INTERVAL 3 MONTH)
                GROUP BY c.id, c.name
                ORDER BY total_spent DESC
                """
                analysis = await db.execute(analysis_query, (user_id,))
                
                if analysis:
                    recommendations = []
                    
                    # Tìm category chi nhiều nhất
                    top_category = analysis[0]
                    recommendations.append(f"💡 Bạn chi nhiều nhất cho {top_category['category_name']} ({top_category['total_spent']:,.0f} VND)")
                    
                    # Gợi ý tiết kiệm
                    if top_category['category_name'] == 'Ăn uống':
                        recommendations.append("🍽️ Gợi ý: Nấu ăn ở nhà nhiều hơn, hạn chế giao đồ ăn")
                    elif top_category['category_name'] == 'Giao thông':
                        recommendations.append("🚗 Gợi ý: Sử dụng phương tiện công cộng, đi chung xe")
                    elif top_category['category_name'] == 'Mua sắm':
                        recommendations.append("🛍️ Gợi ý: Mua sắm có kế hoạch, tránh mua sắm bốc đồng")
                    
                    # Phân tích xu hướng
                    if len(analysis) > 1:
                        second_category = analysis[1]
                        recommendations.append(f"📈 Xu hướng: {second_category['category_name']} cũng chi khá nhiều ({second_category['total_spent']:,.0f} VND)")
                    
                    nlu_response.response = "🤖 Gợi ý thông minh:\n" + "\n".join(recommendations)
                else:
                    nlu_response.response = "Chưa có đủ dữ liệu để đưa ra gợi ý. Hãy thêm nhiều giao dịch hơn!"
                    
        except Exception as e:
            logger.error(f"Error handling smart recommendations: {e}")
            nlu_response.response = "Có lỗi xảy ra khi tạo gợi ý."
    
    async def _handle_expense_forecasting(self, user_id: int, nlu_response: NLUResponse):
        """AI dự đoán chi tiêu tương lai"""
        try:
            async with get_db() as db:
                # Phân tích chi tiêu 6 tháng gần nhất
                forecast_query = """
                SELECT MONTH(transaction_date) as month, 
                       YEAR(transaction_date) as year,
                       SUM(amount) as monthly_expense
                FROM transactions 
                WHERE user_id = %s AND transaction_type = 'expense' 
                AND transaction_date >= DATE_SUB(CURDATE(), INTERVAL 6 MONTH)
                GROUP BY YEAR(transaction_date), MONTH(transaction_date)
                ORDER BY YEAR(transaction_date) DESC, MONTH(transaction_date) DESC
                """
                expenses = await db.execute(forecast_query, (user_id,))
                
                if len(expenses) >= 3:
                    # Tính trung bình chi tiêu
                    avg_expense = sum(e['monthly_expense'] for e in expenses) / len(expenses)
                    
                    # Dự đoán tháng tới
                    next_month_forecast = avg_expense * 1.1  # Tăng 10% để an toàn
                    
                    # Phân tích xu hướng
                    recent_avg = sum(e['monthly_expense'] for e in expenses[:3]) / 3
                    older_avg = sum(e['monthly_expense'] for e in expenses[3:]) / len(expenses[3:]) if len(expenses) > 3 else recent_avg
                    
                    trend = "tăng" if recent_avg > older_avg else "giảm" if recent_avg < older_avg else "ổn định"
                    
                    insights = [
                        f"📊 Dự đoán chi tiêu tháng tới: {next_month_forecast:,.0f} VND",
                        f"📈 Xu hướng: {trend} ({recent_avg:,.0f} vs {older_avg:,.0f} VND)",
                        f"💰 Trung bình 6 tháng: {avg_expense:,.0f} VND"
                    ]
                    
                    nlu_response.response = "🔮 Dự đoán tài chính:\n" + "\n".join(insights)
                else:
                    nlu_response.response = "Cần ít nhất 3 tháng dữ liệu để dự đoán. Hãy sử dụng app thường xuyên hơn!"
                    
        except Exception as e:
            logger.error(f"Error handling expense forecasting: {e}")
            nlu_response.response = "Có lỗi xảy ra khi dự đoán chi tiêu."
    
    def _generate_chat_suggestions(self, intent: str) -> List[str]:
        """AI tự tạo suggestions dựa trên context"""
        # AI có thể tự quyết định suggestions phù hợp
        all_suggestions = [
            "Thêm giao dịch mới",
            "Xem số dư tháng này", 
            "Phân tích chi tiêu",
            "Xem xu hướng tài chính",
            "Gợi ý tiết kiệm",
            "Báo cáo chi tiết"
        ]
        
        # Tùy theo intent, AI có thể chọn suggestions phù hợp
        if intent == "add_transaction":
            return ["Thêm giao dịch khác", "Xem số dư", "Phân tích chi tiêu"]
        elif intent == "query_balance":
            return ["Thêm giao dịch", "Phân tích xu hướng", "Gợi ý tiết kiệm"]
        elif intent == "analyze_data":
            return ["Xem chi tiết", "Thêm giao dịch", "Báo cáo đầy đủ"]
        elif intent == "budget_management":
            return ["Tạo ngân sách mới", "Xem chi tiết ngân sách", "Điều chỉnh ngân sách"]
        elif intent == "goal_tracking":
            return ["Tạo mục tiêu mới", "Cập nhật tiến độ", "Xem lịch sử mục tiêu"]
        elif intent == "smart_recommendations":
            return ["Gợi ý tiết kiệm", "Phân tích chi tiêu", "Tối ưu ngân sách"]
        elif intent == "expense_forecasting":
            return ["Dự đoán dài hạn", "Phân tích xu hướng", "Lập kế hoạch tài chính"]
        else:
            return all_suggestions[:3]  # AI tự chọn 3 suggestions phù hợp
