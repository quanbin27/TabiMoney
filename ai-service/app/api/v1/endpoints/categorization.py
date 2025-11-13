from fastapi import APIRouter
from pydantic import BaseModel

from app.core.config import settings
from app.utils.llm import call_ollama
from app.utils.json_utils import extract_json_block

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
        existing_categories = getattr(payload, "existing_categories", []) or []
        categories_text = ""
        if existing_categories:
            categories_text = "\nDanh sách danh mục hiện có của bạn:\n" + "\n".join(
                f"- {cat.get('name', '')} ({cat.get('name_en', '')})"
                for cat in existing_categories
                if cat.get("name")
            )
        
        prompt = (
            "Bạn là trợ lý phân loại chi tiêu cho ứng dụng tài chính cá nhân.\n"
            f"Mô tả: {payload.description}\nSố tiền: {payload.amount}\n"
            f"Địa điểm: {payload.location or ''}\nTags: {', '.join(payload.tags or [])}\n"
            f"{categories_text}\n\n"
            "Hãy gợi ý danh mục phù hợp bằng TIẾNG VIỆT. Ưu tiên chọn từ danh mục hiện có nếu phù hợp, "
            "nếu không có danh mục nào phù hợp thì gợi ý danh mục mới.\n"
            "Trả về JSON với: suggestions (array of {category_name, confidence_score, reason, is_existing}) "
            "và confidence_score (tổng thể).\n"
            "is_existing: true nếu danh mục đã tồn tại, false nếu là danh mục mới.\n"
            "Ví dụ: {\"suggestions\": [{\"category_name\": \"Ăn uống\", \"confidence_score\": 0.9, \"reason\": \"Chi phí ăn uống\", \"is_existing\": true}], \"confidence_score\": 0.9}"
        )

        result = await call_ollama(prompt, temperature=0.2, max_tokens=400, format_json=True, timeout=120.0)
        parsed = result.get("json") or extract_json_block(result.get("raw", ""))
        if not parsed:
            parsed = {"suggestions": [], "confidence_score": 0.0}

        # Ensure schema integrity
        suggestions = parsed.get("suggestions") or []
        normalized = []
        for suggestion in suggestions:
            if not isinstance(suggestion, dict):
                continue
            normalized.append(
                {
                    "category_name": suggestion.get("category_name", ""),
                    "confidence_score": float(suggestion.get("confidence_score", 0.0) or 0.0),
                    "reason": suggestion.get("reason", ""),
                    "is_existing": bool(suggestion.get("is_existing", False)),
                }
            )

        return {
            "suggestions": normalized,
            "confidence_score": float(parsed.get("confidence_score", 0.0) or 0.0),
        }

    # Fallback OpenAI path (not used when USE_OPENAI=false)
    return {"suggestions": [], "confidence_score": 0.0}


