from fastapi import APIRouter, Depends
from app.core.dependencies import get_prediction_service
from app.models.prediction import ExpensePredictionRequest, ExpensePredictionResponse
from app.services.prediction_service import PredictionService

router = APIRouter()

@router.post("/expenses", response_model=ExpensePredictionResponse)
async def predict_expenses(
    request: ExpensePredictionRequest,
    prediction_service: PredictionService = Depends(get_prediction_service)
):
    """Predict future expenses based on historical data"""
    return await prediction_service.predict_expenses(request)


