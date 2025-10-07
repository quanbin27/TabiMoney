"""
Natural Language Understanding Service
NLU functionality for TabiMoney AI Service
"""

import asyncio
import logging
from typing import Dict, Any, List, Optional, Tuple
import re
from datetime import datetime, timedelta
import httpx
import json

import openai
from openai import AsyncOpenAI

from app.core.config import settings
from app.models.nlu import NLURequest, NLUResponse, Entity, ChatRequest, ChatResponse

logger = logging.getLogger(__name__)


class NLUService:
    """Natural Language Understanding Service"""
    
    def __init__(self):
        self.client: Optional[AsyncOpenAI] = None
        self.ollama_client: Optional[httpx.AsyncClient] = None
        self.is_initialized = False
        
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
            if settings.USE_OPENAI and settings.OPENAI_API_KEY:
                # Initialize OpenAI client
                self.client = AsyncOpenAI(api_key=settings.OPENAI_API_KEY)
                await self._test_openai_connection()
                logger.info("NLU Service initialized with OpenAI")
            else:
                # Initialize Ollama client (shorter timeout to avoid FE 502s)
                self.ollama_client = httpx.AsyncClient(timeout=10.0)
                await self._test_ollama_connection()
                logger.info("NLU Service initialized with Ollama")
            
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
    
    async def _test_openai_connection(self):
        """Test OpenAI API connection"""
        try:
            response = await self.client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=[{"role": "user", "content": "Hello"}],
                max_tokens=10
            )
            logger.info("OpenAI API connection successful")
        except Exception as e:
            logger.error(f"OpenAI API connection failed: {e}")
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
            if self.client:
                return await self._process_with_openai(request)
            elif self.ollama_client:
                return await self._process_with_ollama(request)
            else:
                return await self._process_with_rules(request)
                
        except Exception as e:
            logger.error(f"Failed to process NLU: {e}")
            # Fallback to rule-based processing
            return await self._process_with_rules(request)
    
    async def _process_with_openai(self, request: NLURequest) -> NLUResponse:
        """Process NLU using OpenAI"""
        try:
            # Build prompt for OpenAI
            prompt = self._build_nlu_prompt(request.text, request.context)
            
            response = await self.client.chat.completions.create(
                model=settings.OPENAI_MODEL,
                messages=[
                    {"role": "system", "content": prompt},
                    {"role": "user", "content": request.text}
                ],
                max_tokens=settings.OPENAI_MAX_TOKENS,
                temperature=settings.OPENAI_TEMPERATURE
            )
            
            # Parse OpenAI response
            content = response.choices[0].message.content
            return self._parse_openai_response(content, request.user_id)
            
        except Exception as e:
            logger.error(f"OpenAI NLU processing failed: {e}")
            raise
    
    async def _process_with_ollama(self, request: NLURequest) -> NLUResponse:
        """Process NLU using Ollama"""
        try:
            prompt = self._build_nlu_prompt(request.text, request.context)
            
            response = await self.ollama_client.post(
                f"{settings.OLLAMA_BASE_URL}/api/generate",
                json={
                    "model": settings.LLM_MODEL,
                    "prompt": prompt,
                    "stream": False,
                    "options": {
                        "temperature": 0.3,
                        "max_tokens": 1000
                    }
                }
            )
            
            if response.status_code != 200:
                raise Exception(f"Ollama returned status {response.status_code}")
            
            # Ollama returns JSON with { response: "..." }
            # Read as text first for resilience
            text_body = await response.aread()
            try:
                result = json.loads(text_body.decode("utf-8", errors="ignore"))
                content = result.get("response", "")
            except Exception:
                content = text_body.decode("utf-8", errors="ignore")

            # Try to parse structured JSON from model output
            try:
                return self._parse_openai_response(self._extract_json_block(content), request.user_id)
            except Exception:
                logger.warning("Falling back to rule-based NLU due to parse failure")
                return await self._process_with_rules(request)
            
        except Exception as e:
            logger.error(f"Ollama NLU processing failed: {e}")
            raise
    
    async def _process_with_rules(self, request: NLURequest) -> NLUResponse:
        """Process NLU using rule-based approach"""
        try:
            text = request.text.lower()
            
            # Extract entities
            entities = self._extract_entities_rule_based(text)
            
            # Determine intent
            intent = self._determine_intent_rule_based(text, entities)
            
            # Generate response
            response_text = self._generate_response_rule_based(intent, entities)
            
            return NLUResponse(
                user_id=request.user_id,
                intent=intent,
                entities=entities,
                confidence=0.7,  # Rule-based confidence
                suggested_action=self._get_suggested_action(intent),
                response=response_text,
                generated_at=datetime.now().isoformat() + "Z"
            )
            
        except Exception as e:
            logger.error(f"Rule-based NLU processing failed: {e}")
            raise
    
    def _build_nlu_prompt(self, text: str, context: str) -> str:
        """Build prompt for AI NLU processing"""
        return f"""
Bạn là trợ lý tài chính cho TabiMoney. Phân tích input của user và trích xuất:

1. Intent (add_transaction, query_balance, ask_question, general)
2. Entities (amounts, categories, dates, descriptions)
3. Confidence score (0-1)
4. Suggested action
5. Response (bằng tiếng Việt)

User input: {text}
Context: {context}

Return JSON format:
{{
  "intent": "add_transaction",
  "entities": [
    {{"type": "amount", "value": "50000", "confidence": 0.95, "start_pos": 0, "end_pos": 5}},
    {{"type": "category", "value": "ăn uống", "confidence": 0.9, "start_pos": 10, "end_pos": 17}}
  ],
  "confidence": 0.92,
  "suggested_action": "create_transaction",
  "response": "Tôi sẽ giúp bạn thêm giao dịch này..."
}}

Possible intents:
- add_transaction: user muốn thêm giao dịch mới
- query_balance: user muốn biết về số dư, chi tiêu, thu nhập
- ask_question: user hỏi câu hỏi về tài chính
- general: chào hỏi hoặc không xác định được

Possible entities:
- amount: số tiền
- category: danh mục chi tiêu
- date: ngày tháng
- description: mô tả giao dịch
"""
    
    def _parse_openai_response(self, content: str, user_id: int) -> NLUResponse:
        """Parse OpenAI response"""
        try:
            import json
            data = json.loads(content)
            
            entities = []
            for entity_data in data.get('entities', []):
                entity = Entity(
                    type=entity_data['type'],
                    value=entity_data['value'],
                    confidence=entity_data['confidence'],
                    start_pos=entity_data.get('start_pos', 0),
                    end_pos=entity_data.get('end_pos', 0)
                )
                entities.append(entity)
            
            return NLUResponse(
                user_id=user_id,
                intent=data['intent'],
                entities=entities,
                confidence=data['confidence'],
                suggested_action=data['suggested_action'],
                response=data['response'],
                generated_at=datetime.now().isoformat() + "Z"
            )
            
        except Exception as e:
            logger.error(f"Failed to parse OpenAI response: {e}")
            raise

    def _extract_json_block(self, content: str) -> str:
        """Extract JSON block from LLM output (handles code fences and extra text)."""
        if not content:
            raise ValueError("empty content")
        # Remove code fences
        content = content.strip()
        if content.startswith("```"):
            content = content.strip('`')
        # Find first {...} JSON object
        import re
        match = re.search(r"\{[\s\S]*\}", content)
        if match:
            return match.group(0)
        # If not found, assume content itself is JSON
        return content
    
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
                    entity = Entity(
                        type="category",
                        value=category,
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
    
    def _determine_intent_rule_based(self, text: str, entities: List[Entity]) -> str:
        """Determine intent using rule-based approach"""
        # Check for query keywords first (more specific)
        query_keywords = ['bao nhiêu', 'tổng', 'tổng cộng', 'chi tiêu', 'thu nhập', 'số dư', 'muốn biết', 'kiểm tra']
        if any(keyword in text for keyword in query_keywords):
            return "query_balance"
        
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
            elif 'tr' in full_match.lower() or 'triệu' in full_match.lower():
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
    
    def _generate_chat_suggestions(self, intent: str) -> List[str]:
        """Generate chat suggestions based on intent"""
        suggestion_map = {
            "add_transaction": [
                "Thêm giao dịch khác",
                "Xem danh mục chi tiêu",
                "Kiểm tra số dư"
            ],
            "query_balance": [
                "Xem chi tiết chi tiêu",
                "Phân tích xu hướng",
                "Dự đoán chi tiêu"
            ],
            "ask_question": [
                "Cách tiết kiệm tiền",
                "Quản lý ngân sách",
                "Đầu tư cá nhân"
            ],
            "general": [
                "Thêm giao dịch",
                "Xem báo cáo",
                "Phân tích tài chính"
            ]
        }
        
        return suggestion_map.get(intent, [
            "Thêm giao dịch",
            "Xem báo cáo",
            "Phân tích tài chính"
        ])
