"""
Prediction models for TabiMoney AI Service
"""

from datetime import datetime
from typing import List, Dict, Any, Optional
from pydantic import BaseModel


class ExpensePredictionRequest(BaseModel):
    """Request model for expense prediction"""
    user_id: int
    start_date: str
    end_date: str


class CategoryBreakdown(BaseModel):
    """Category breakdown model - matches Go CategoryPrediction"""
    category_id: int
    category_name: str
    predicted_amount: float
    confidence_score: float
    trend: str


class Trend(BaseModel):
    """Trend model - matches Go ExpenseTrend"""
    period: str
    amount: float
    change_percentage: float
    trend: str


class ExpensePredictionResponse(BaseModel):
    """Response model for expense prediction"""
    user_id: int
    predicted_amount: float
    confidence_score: float
    category_breakdown: List[CategoryBreakdown]
    trends: List[Trend]
    recommendations: List[str]
    generated_at: datetime
