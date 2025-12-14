"""
Natural Language Understanding Service
NLU functionality for TabiMoney AI Service
"""

import logging
from typing import Dict, Any, List, Optional, Tuple
import re
from datetime import datetime, timedelta
import json

"""NLU service using Google Gemini."""

import aiohttp

from app.core.config import settings
from app.core.database import get_db
from app.models.nlu import ChatRequest, ChatResponse, Entity, NLURequest, NLUResponse
from app.services.transaction_service import TransactionService
from app.utils.json_utils import extract_json_block
from app.utils.llm import call_gemini

logger = logging.getLogger(__name__)


class NLUService:
    """Natural Language Understanding Service"""
    
    def __init__(self):
        self.client: Optional[object] = None
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
            if not settings.USE_GEMINI or not settings.GEMINI_API_KEY:
                raise ValueError("Gemini API key is required. Please set USE_GEMINI=true and GEMINI_API_KEY in environment variables.")
            
            await self._test_gemini_connection()
            logger.info("NLU Service initialized with Gemini")
            self.is_initialized = True
            
        except Exception as e:
            logger.error(f"Failed to initialize NLU Service: {e}")
            # Fallback to rule-based NLU
            self.is_initialized = True
            logger.warning("NLU Service initialized with rule-based fallback (Gemini not available)")
    
    async def cleanup(self):
        """Cleanup NLU service"""
        logger.info("Cleaning up NLU Service...")
        self.client = None
        self.is_initialized = False
    
    def is_ready(self) -> bool:
        """Check if NLU service is ready"""
        return self.is_initialized

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
    
    
    async def process_nlu(self, request: NLURequest) -> NLUResponse:
        """Process Natural Language Understanding request"""
        if not self.is_ready():
            raise RuntimeError("NLU Service not ready")
        
        try:
            # Use Gemini if available
            if settings.USE_GEMINI and settings.GEMINI_API_KEY:
                return await self._process_with_gemini(request)
            else:
                # Fallback to rule-based processing
                return await self._process_with_rules(request)
                
        except Exception as e:
            logger.error(f"Failed to process NLU: {e}")
            # Fallback to rule-based processing
            return await self._process_with_rules(request)

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
                        "maxOutputTokens": settings.GEMINI_MAX_TOKENS,
                        "response_mime_type": "application/json"
                    }
                }
                async with session.post(url, json=payload) as resp:
                    logger.info(f"Gemini response status: {resp.status}")
                    if resp.status != 200:
                        logger.error(f"Gemini returned status {resp.status}")
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
                        parsed = self._parse_gemini_response(json.dumps(content_dict), request.user_id, request.text)
                        parsed = await self._normalize_entities_add_category_id(parsed, request.user_id)
                        logger.info("Gemini parsing successful: intent=%s", parsed.intent)
                        return parsed
                    except Exception as e:
                        logger.warning("Failed to parse Gemini response: %s", e)
                        raise
        except Exception as e:
            logger.error(f"Gemini NLU processing failed: {e}")
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
    
    def _parse_gemini_response(self, content: str, user_id: int, original_text: str = "") -> NLUResponse:
        """Parse Gemini response"""
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
            logger.error(f"Failed to parse Gemini response: {e}")
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

    # deprecated helper removed
    
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
            elif nlu_response.intent == "general":
                await self._handle_general(request.user_id, nlu_response)
            
            # Generate suggestions based on intent
            suggestions = self._generate_chat_suggestions(nlu_response.intent)
            
            return ChatResponse(
                user_id=request.user_id,
                message=request.message,
                response=nlu_response.response,
                intent=nlu_response.intent,
                entities=nlu_response.entities,
                suggestions=suggestions,
                needs_confirmation=nlu_response.needs_confirmation,
                generated_at=datetime.now()
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
                needs_confirmation=False,
                generated_at=datetime.now()
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
    
    async def _generate_natural_response(self, data_summary: str, context: str, user_message: str = "") -> str:
        """Generate natural language response using Gemini AI"""
        try:
            if not settings.USE_GEMINI or not settings.GEMINI_API_KEY:
                # Fallback to simple format if Gemini not available
                return data_summary
            
            prompt = f"""Bạn là AI Assistant thân thiện cho ứng dụng quản lý tài chính TabiMoney. 
Người dùng đã hỏi: "{user_message}"

Dữ liệu thực tế từ database:
{data_summary}

Context: {context}

Nhiệm vụ: Tạo một phản hồi tự nhiên, thân thiện, dễ hiểu bằng tiếng Việt dựa trên dữ liệu trên.
- Sử dụng ngôn ngữ tự nhiên, không quá kỹ thuật
- Thêm emoji phù hợp để làm cho phản hồi sinh động
- Đưa ra insights và gợi ý hữu ích nếu có thể
- Giữ nguyên các con số chính xác từ dữ liệu
- Phản hồi nên ngắn gọn nhưng đầy đủ thông tin (khoảng 100-200 từ)

Phản hồi:"""

            result = await call_gemini(
                prompt,
                temperature=0.7,  # Higher temperature for more natural responses
                format_json=False,  # We want plain text, not JSON
                timeout=30.0
            )
            
            response = result.get("raw", "").strip()
            if response:
                return response
            else:
                # Fallback to data summary if AI fails
                return data_summary
                
        except Exception as e:
            logger.warning(f"Failed to generate natural response with AI: {e}, using fallback")
            return data_summary
    
    async def _handle_query_balance(self, user_id: int, nlu_response: NLUResponse):
        """Handle balance query directly"""
        try:
            result = await self.transaction_service.get_user_balance(user_id)
            if result["success"]:
                # Use AI to make response more natural
                data_summary = f"""
Tổng thu nhập tháng này: {result.get('total_income', 0):,.0f} VND
Tổng chi tiêu tháng này: {result.get('total_expense', 0):,.0f} VND
Số dư (chênh lệch): {result.get('net_amount', 0):,.0f} VND
"""
                natural_response = await self._generate_natural_response(
                    data_summary,
                    "Người dùng đang hỏi về số dư tài chính tháng hiện tại",
                    "Số dư của tôi thế nào?"
                )
                nlu_response.response = natural_response
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
                    # Tính toán dữ liệu
                    total_expense = sum(t['amount'] for t in transactions if t['transaction_type'] == 'expense')
                    total_income = sum(t['amount'] for t in transactions if t['transaction_type'] == 'income')
                    
                    # Phân tích theo category
                    category_spending = {}
                    for t in transactions:
                        if t['transaction_type'] == 'expense':
                            cat = t['category_name']
                            category_spending[cat] = category_spending.get(cat, 0) + t['amount']
                    
                    # Tạo data summary cho AI
                    category_list = "\n".join([f"- {cat}: {amt:,.0f} VND" for cat, amt in sorted(category_spending.items(), key=lambda x: x[1], reverse=True)[:5]])
                    
                    data_summary = f"""
Phân tích 30 ngày gần nhất:
- Tổng chi tiêu: {total_expense:,.0f} VND
- Tổng thu nhập: {total_income:,.0f} VND
- Số dư: {total_income - total_expense:,.0f} VND
- Số giao dịch: {len(transactions)}

Top 5 danh mục chi nhiều nhất:
{category_list if category_spending else "Chưa có dữ liệu"}
"""
                    
                    natural_response = await self._generate_natural_response(
                        data_summary,
                        "Người dùng muốn phân tích chi tiêu và thu nhập trong 30 ngày gần nhất",
                        "Phân tích chi tiêu của tôi"
                    )
                    nlu_response.response = natural_response
                else:
                    nlu_response.response = "Chưa có dữ liệu giao dịch để phân tích. Hãy thêm giao dịch để tôi có thể giúp bạn phân tích!"
                    
        except Exception as e:
            logger.error(f"Error handling data analysis: {e}")
            nlu_response.response = "Có lỗi xảy ra khi phân tích dữ liệu. Vui lòng thử lại sau."
    
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
                    expense_dict = {e['category_name']: e['total_spent'] for e in expenses}
                    budget_details = []
                    
                    for budget in budgets:
                        spent = expense_dict.get(budget['category_name'], 0)
                        remaining = budget['amount'] - spent
                        percentage = (spent / budget['amount']) * 100 if budget['amount'] > 0 else 0
                        
                        status = "vượt quá" if percentage > 100 else "sắp hết" if percentage > 90 else "đang ổn" if percentage > 70 else "còn nhiều"
                        
                        budget_details.append(
                            f"- {budget['category_name']}: "
                            f"Đã chi {spent:,.0f}/{budget['amount']:,.0f} VND ({percentage:.1f}%), "
                            f"Còn lại {remaining:,.0f} VND - {status}"
                        )
                    
                    data_summary = f"""
Tình hình ngân sách tháng này:
{chr(10).join(budget_details)}
"""
                    
                    natural_response = await self._generate_natural_response(
                        data_summary,
                        "Người dùng muốn kiểm tra tình hình ngân sách tháng hiện tại",
                        "Tình hình ngân sách của tôi thế nào?"
                    )
                    nlu_response.response = natural_response
                else:
                    nlu_response.response = "Bạn chưa có ngân sách nào. Hãy tạo ngân sách để quản lý chi tiêu tốt hơn! 💡"
                    
        except Exception as e:
            logger.error(f"Error handling budget management: {e}")
            nlu_response.response = "Có lỗi xảy ra khi kiểm tra ngân sách. Vui lòng thử lại sau."
    
    async def _handle_goal_tracking(self, user_id: int, nlu_response: NLUResponse):
        """AI theo dõi mục tiêu tài chính"""
        try:
            async with get_db() as db:
                # Lấy goals hiện tại với current_amount từ database
                goals_query = """
                SELECT id, title, description, target_amount, current_amount, target_date, 
                       goal_type, priority, is_achieved, created_at
                FROM financial_goals 
                WHERE user_id = %s AND is_achieved = false
                ORDER BY target_date ASC, created_at DESC
                """
                goals = await db.execute(goals_query, (user_id,))
                
                if goals:
                    goal_details = []
                    for goal in goals:
                        current_amount = float(goal.get('current_amount', 0) or 0)
                        target_amount = float(goal.get('target_amount', 0) or 0)
                        
                        if target_amount > 0:
                            progress_percentage = (current_amount / target_amount) * 100
                            remaining = target_amount - current_amount
                        else:
                            progress_percentage = 0
                            remaining = 0
                        
                        goal_title = goal.get('title', 'Mục tiêu không tên')
                        goal_type = goal.get('goal_type', 'savings')
                        
                        # Format date info
                        date_info = ""
                        if goal.get('target_date'):
                            from datetime import datetime
                            try:
                                target_date = goal['target_date']
                                if isinstance(target_date, str):
                                    target_date = datetime.strptime(target_date, '%Y-%m-%d').date()
                                elif hasattr(target_date, 'date'):
                                    target_date = target_date.date()
                                
                                today = datetime.now().date()
                                days_remaining = (target_date - today).days
                                
                                if days_remaining > 0:
                                    date_info = f", Còn {days_remaining} ngày"
                                elif days_remaining == 0:
                                    date_info = ", Hôm nay là hạn chót"
                                else:
                                    date_info = f", Đã quá hạn {abs(days_remaining)} ngày"
                            except Exception:
                                pass
                        
                        status = "đã đạt" if remaining <= 0 else f"còn thiếu {remaining:,.0f} VND"
                        
                        goal_details.append(
                            f"- {goal_title} ({goal_type}): "
                            f"{progress_percentage:.1f}% hoàn thành "
                            f"({current_amount:,.0f}/{target_amount:,.0f} VND), "
                            f"{status}{date_info}"
                        )
                    
                    data_summary = f"""
Tiến độ mục tiêu tài chính:
{chr(10).join(goal_details)}
"""
                    
                    natural_response = await self._generate_natural_response(
                        data_summary,
                        "Người dùng muốn kiểm tra tiến độ các mục tiêu tài chính",
                        "Tiến độ mục tiêu của tôi thế nào?"
                    )
                    nlu_response.response = natural_response
                else:
                    nlu_response.response = "Bạn chưa có mục tiêu tài chính nào. Hãy tạo mục tiêu để có động lực tiết kiệm! 🎯"
                    
        except Exception as e:
            logger.error(f"Error handling goal tracking: {e}", exc_info=True)
            nlu_response.response = "Có lỗi xảy ra khi kiểm tra mục tiêu. Vui lòng thử lại sau."
    
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
                LIMIT 5
                """
                analysis = await db.execute(analysis_query, (user_id,))
                
                if analysis:
                    # Tính tổng chi tiêu
                    total_spent_3m = sum(a['total_spent'] for a in analysis)
                    
                    # Format data
                    category_list = "\n".join([
                        f"- {a['category_name']}: {a['total_spent']:,.0f} VND "
                        f"({(a['total_spent']/total_spent_3m*100):.1f}%, {a['transaction_count']} giao dịch, "
                        f"trung bình {a['avg_amount']:,.0f} VND/giao dịch)"
                        for a in analysis
                    ])
                    
                    data_summary = f"""
Phân tích chi tiêu 3 tháng gần nhất:
Tổng chi tiêu: {total_spent_3m:,.0f} VND

Top 5 danh mục chi nhiều nhất:
{category_list}

Người dùng cần gợi ý cụ thể để tiết kiệm tiền dựa trên các danh mục này.
"""
                    
                    natural_response = await self._generate_natural_response(
                        data_summary,
                        "Người dùng muốn nhận gợi ý thông minh để tiết kiệm tiền dựa trên phân tích chi tiêu",
                        "Gợi ý tiết kiệm cho tôi"
                    )
                    nlu_response.response = natural_response
                else:
                    nlu_response.response = "Chưa có đủ dữ liệu để đưa ra gợi ý. Hãy thêm nhiều giao dịch hơn để tôi có thể phân tích và đưa ra gợi ý hữu ích! 📊"
                    
        except Exception as e:
            logger.error(f"Error handling smart recommendations: {e}")
            nlu_response.response = "Có lỗi xảy ra khi tạo gợi ý. Vui lòng thử lại sau."
    
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
                    # Tính toán
                    avg_expense = sum(e['monthly_expense'] for e in expenses) / len(expenses)
                    next_month_forecast = avg_expense * 1.1  # Tăng 10% để an toàn
                    
                    # Phân tích xu hướng
                    recent_avg = sum(e['monthly_expense'] for e in expenses[:3]) / 3
                    older_avg = sum(e['monthly_expense'] for e in expenses[3:]) / len(expenses[3:]) if len(expenses) > 3 else recent_avg
                    
                    trend_direction = "tăng" if recent_avg > older_avg else "giảm" if recent_avg < older_avg else "ổn định"
                    trend_percentage = abs((recent_avg - older_avg) / older_avg * 100) if older_avg > 0 else 0
                    
                    # Format monthly data
                    monthly_data = "\n".join([
                        f"- Tháng {e['month']}/{e['year']}: {e['monthly_expense']:,.0f} VND"
                        for e in expenses[:6]
                    ])
                    
                    data_summary = f"""
Dữ liệu chi tiêu 6 tháng gần nhất:
{monthly_data}

Thống kê:
- Trung bình 6 tháng: {avg_expense:,.0f} VND
- Trung bình 3 tháng gần nhất: {recent_avg:,.0f} VND
- Trung bình 3 tháng trước đó: {older_avg:,.0f} VND
- Xu hướng: {trend_direction} {trend_percentage:.1f}%

Dự đoán:
- Chi tiêu tháng tới (dự kiến): {next_month_forecast:,.0f} VND (tăng 10% so với trung bình để an toàn)
"""
                    
                    natural_response = await self._generate_natural_response(
                        data_summary,
                        "Người dùng muốn dự đoán chi tiêu tháng tới dựa trên dữ liệu lịch sử",
                        "Dự đoán chi tiêu tháng tới"
                    )
                    nlu_response.response = natural_response
                else:
                    nlu_response.response = "Cần ít nhất 3 tháng dữ liệu để dự đoán chính xác. Hãy sử dụng app thường xuyên hơn để tôi có thể đưa ra dự đoán tốt hơn! 📈"
                    
        except Exception as e:
            logger.error(f"Error handling expense forecasting: {e}")
            nlu_response.response = "Có lỗi xảy ra khi dự đoán chi tiêu. Vui lòng thử lại sau."
    
    async def _handle_general(self, user_id: int, nlu_response: NLUResponse):
        """AI xử lý các câu hỏi chung chung về tài chính"""
        try:
            # Intent "general" được dùng cho các câu hỏi không thuộc các intent cụ thể
            # AI đã tự tạo response phù hợp, chỉ cần đảm bảo response có ý nghĩa
            
            # Nếu response từ AI quá ngắn hoặc không rõ ràng, thêm thông tin hữu ích
            if not nlu_response.response or len(nlu_response.response.strip()) < 20:
                # Lấy thông tin tổng quan để cung cấp context
                async with get_db() as db:
                    # Lấy số giao dịch gần đây
                    count_query = """
                    SELECT COUNT(*) as total 
                    FROM transactions 
                    WHERE user_id = %s
                    """
                    result = await db.execute(count_query, (user_id,))
                    total_transactions = result[0]['total'] if result else 0
                    
                    if total_transactions > 0:
                        nlu_response.response = (
                            f"{nlu_response.response}\n\n"
                            f"💡 Bạn có {total_transactions} giao dịch trong hệ thống. "
                            f"Tôi có thể giúp bạn phân tích chi tiêu, quản lý ngân sách, theo dõi mục tiêu và đưa ra gợi ý tài chính."
                        )
                    else:
                        nlu_response.response = (
                            f"{nlu_response.response}\n\n"
                            f"💡 Bắt đầu bằng cách thêm giao dịch đầu tiên của bạn! "
                            f"Tôi có thể giúp bạn quản lý tài chính, phân tích chi tiêu và đưa ra gợi ý."
                        )
            
            logger.info(f"Handled general intent for user {user_id}")
            
        except Exception as e:
            logger.error(f"Error handling general intent: {e}")
            # Không cần thay đổi response nếu có lỗi, giữ nguyên response từ AI
    
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
