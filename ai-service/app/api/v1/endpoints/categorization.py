from fastapi import APIRouter, Depends
from pydantic import BaseModel
import httpx
import json
from app.core.config import settings

router = APIRouter()

class SuggestRequest(BaseModel):
    user_id: int | None = None
    description: str
    amount: float
    location: str | None = None
    tags: list[str] | None = None
    existing_categories: list[dict] | None = None

@router.post("/suggest")
async def suggest_category(payload: SuggestRequest):
    if not settings.ENABLE_CATEGORIZATION:
        return {"suggestions": [], "confidence_score": 0.0}

    if not settings.USE_OPENAI:
        # Use Ollama local LLM with user's existing categories
        # Get user's existing categories from the request
        existing_categories = getattr(payload, 'existing_categories', [])
        categories_text = ""
        if existing_categories:
            categories_text = f"\nDanh sách danh mục hiện có của bạn:\n" + "\n".join([f"- {cat.get('name', '')} ({cat.get('name_en', '')})" for cat in existing_categories if cat.get('name')])
        
        prompt = (
            "Bạn là trợ lý phân loại chi tiêu cho ứng dụng tài chính cá nhân.\n"
            f"Mô tả: {payload.description}\nSố tiền: {payload.amount}\n"
            f"Địa điểm: {payload.location or ''}\nTags: {', '.join(payload.tags or [])}\n"
            f"{categories_text}\n\n"
            "Hãy gợi ý danh mục phù hợp bằng TIẾNG VIỆT. Ưu tiên chọn từ danh mục hiện có nếu phù hợp, nếu không có danh mục nào phù hợp thì gợi ý danh mục mới.\n"
            "Trả về JSON với: suggestions (array of {category_name, confidence_score, reason, is_existing}) và confidence_score (tổng thể).\n"
            "is_existing: true nếu danh mục đã tồn tại, false nếu là danh mục mới.\n"
            "Ví dụ: {\"suggestions\": [{\"category_name\": \"Ăn uống\", \"confidence_score\": 0.9, \"reason\": \"Chi phí ăn uống\", \"is_existing\": true}], \"confidence_score\": 0.9}"
        )
        async with httpx.AsyncClient(timeout=120) as client:
            resp = await client.post(
                f"{settings.OLLAMA_BASE_URL}/api/generate",
                json={
                    "model": settings.LLM_MODEL,
                    "prompt": prompt,
                    "options": {"temperature": 0.2},
                    "format": "json",
                    "stream": False
                },
            )
        resp.raise_for_status()
        data = resp.json()
        text = data.get("response", "{}")
        try:
            parsed = json.loads(text)
        except Exception:
            parsed = {"suggestions": [], "confidence_score": 0.0}
        return parsed

    # Fallback OpenAI path (not used when USE_OPENAI=false)
    return {"suggestions": [], "confidence_score": 0.0}


