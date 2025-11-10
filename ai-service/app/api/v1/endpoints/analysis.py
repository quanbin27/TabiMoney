from typing import List, Optional

from fastapi import APIRouter
from pydantic import BaseModel

from app.utils.llm import call_ollama
from app.utils.json_utils import extract_json_block, ensure_string_list

router = APIRouter()


class CategorySummary(BaseModel):
    category_id: int
    category_name: str
    total_amount: float
    transaction_count: int


class SpendingAnalysisRequest(BaseModel):
    user_id: int
    start_date: Optional[str] = None
    end_date: Optional[str] = None
    patterns: List[CategorySummary]


class SpendingAnalysisResponse(BaseModel):
    insights: List[str]
    recommendations: List[str]


@router.post("/spending", response_model=SpendingAnalysisResponse)
async def analyze_spending(req: SpendingAnalysisRequest):
    # Build concise Vietnamese prompt
    categories_text = "\n".join([
        f"- {p.category_name}: {int(p.total_amount)} VND ({p.transaction_count} giao dịch)"
        for p in req.patterns
    ])

    prompt = (
        "Bạn là trợ lý tài chính. Dựa trên dữ liệu chi tiêu sau, hãy TRẢ VỀ JSON với 2 mảng: 'insights' và 'recommendations'.\n"
        "Mỗi phần tử là một câu ngắn gọn tiếng Việt. KHÔNG giải thích thêm.\n\n"
        f"Danh sách danh mục:\n{categories_text}\n\n"
        "JSON mẫu:\n{\n  \"insights\": [\"Nhận định 1\", \"Nhận định 2\"],\n  \"recommendations\": [\"Khuyến nghị 1\", \"Khuyến nghị 2\"]\n}"
    )

    try:
        result = await call_ollama(prompt, temperature=0.3, max_tokens=400, format_json=True)
        payload = result.get("json") or extract_json_block(result.get("raw", ""))
        insights = ensure_string_list(payload.get("insights"))
        recommendations = ensure_string_list(payload.get("recommendations"))

        return SpendingAnalysisResponse(insights=insights, recommendations=recommendations)
    except Exception:
        return SpendingAnalysisResponse(
            insights=["Chi tiêu của bạn tập trung ở một vài danh mục chính."],
            recommendations=["Cân nhắc đặt hạn mức và theo dõi các danh mục chi tiêu cao."],
        )


