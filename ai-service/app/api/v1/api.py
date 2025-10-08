"""
API v1 router
"""

from fastapi import APIRouter

from app.api.v1.endpoints import nlu, prediction, anomaly, categorization, chat, analysis

api_router = APIRouter()

# Include all endpoint routers
api_router.include_router(nlu.router, prefix="/nlu", tags=["NLU"])
api_router.include_router(prediction.router, prefix="/prediction", tags=["Prediction"])
api_router.include_router(anomaly.router, prefix="/anomaly", tags=["Anomaly Detection"])
api_router.include_router(categorization.router, prefix="/categorization", tags=["Categorization"])
api_router.include_router(chat.router, prefix="/chat", tags=["Chat"])
api_router.include_router(analysis.router, prefix="/analysis", tags=["Analysis"])
