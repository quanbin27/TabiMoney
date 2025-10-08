from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from typing import List, Optional
import httpx
import json
from app.core.config import settings

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

    async with httpx.AsyncClient(timeout=30.0) as client:
        try:
            resp = await client.post(
                f"{settings.OLLAMA_BASE_URL}/api/generate",
                json={
                    "model": settings.LLM_MODEL,
                    "prompt": prompt,
                    "stream": False,
                    "options": {"temperature": 0.3, "max_tokens": 400},
                },
            )
            if resp.status_code != 200:
                raise HTTPException(status_code=500, detail="LLM generation failed")

            body = resp.json()
            content = body.get("response", "{}")
            # Extract JSON block
            try:
                data = json.loads(content)
            except Exception:
                # Best-effort: find first {...}
                import re
                m = re.search(r"\{[\s\S]*\}", content)
                data = json.loads(m.group(0)) if m else {"insights": [], "recommendations": []}

            insights = [str(x) for x in data.get("insights", [])]
            recs = [str(x) for x in data.get("recommendations", [])]
            return SpendingAnalysisResponse(insights=insights, recommendations=recs)
        except Exception as e:
            # Fallback minimal
            return SpendingAnalysisResponse(
                insights=["Chi tiêu của bạn tập trung ở một vài danh mục chính."],
                recommendations=["Cân nhắc đặt hạn mức và theo dõi các danh mục chi tiêu cao."],
            )


