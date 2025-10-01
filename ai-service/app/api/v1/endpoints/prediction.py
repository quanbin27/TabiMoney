from fastapi import APIRouter

router = APIRouter()

@router.post("/expenses")
async def predict_expenses():
    return {"predicted_amount": 0, "confidence_score": 0.0, "category_breakdown": [], "trends": [], "recommendations": []}


