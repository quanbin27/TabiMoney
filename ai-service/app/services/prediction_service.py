"""
Prediction Service for TabiMoney AI Service
Handles expense prediction and financial forecasting
"""

import asyncio
import logging
from typing import Dict, Any, List, Optional
from datetime import datetime, timedelta
import pandas as pd
import numpy as np
from sklearn.ensemble import RandomForestRegressor
from sklearn.preprocessing import LabelEncoder
import json

from app.core.config import settings
from app.core.database import get_db
from app.models.prediction import ExpensePredictionRequest, ExpensePredictionResponse

logger = logging.getLogger(__name__)


class PredictionService:
    """Prediction Service for financial forecasting"""
    
    def __init__(self, ml_service):
        self.ml_service = ml_service
        self._ready = False
        self.model = None
        self.category_encoder = LabelEncoder()
        
    async def initialize(self):
        """Initialize Prediction Service"""
        logger.info("Initializing Prediction Service...")
        
        try:
            # Initialize ML model
            self.model = RandomForestRegressor(
                n_estimators=100,
                random_state=42,
                max_depth=10
            )
            
            self._ready = True
            logger.info("Prediction Service initialized successfully")
            
        except Exception as e:
            logger.error(f"Failed to initialize Prediction Service: {e}")
            self._ready = False
    
    async def cleanup(self):
        """Cleanup Prediction Service"""
        logger.info("Cleaning up Prediction Service...")
        self.model = None
        self._ready = False
    
    def is_ready(self) -> bool:
        """Check if Prediction Service is ready"""
        return self._ready
    
    async def predict_expenses(self, request: ExpensePredictionRequest) -> ExpensePredictionResponse:
        """Predict future expenses based on historical data"""
        if not self.is_ready():
            raise RuntimeError("Prediction Service not ready")
        
        try:
            logger.info(f"Starting prediction for user {request.user_id}")
            
            # Get historical data from database
            historical_data = await self._get_historical_data(request.user_id, request.start_date, request.end_date)
            logger.info(f"Retrieved {len(historical_data)} historical transactions")
            
            if len(historical_data) < 5:  # Need at least 5 transactions for prediction
                logger.warning(f"Not enough data for prediction: {len(historical_data)} transactions")
                return ExpensePredictionResponse(
                    user_id=request.user_id,
                    predicted_amount=0,
                    confidence_score=0.0,
                    category_breakdown=[],
                    trends=[],
                    recommendations=["Cần thêm dữ liệu giao dịch để đưa ra dự đoán chính xác"],
                    generated_at=datetime.now()
                )
            
            # Prepare data for ML model
            features, target = self._prepare_training_data(historical_data)
            
            if len(features) < 3:  # Need at least 3 data points
                return ExpensePredictionResponse(
                    user_id=request.user_id,
                    predicted_amount=0,
                    confidence_score=0.0,
                    category_breakdown=[],
                    trends=[],
                    recommendations=["Dữ liệu không đủ để đưa ra dự đoán"],
                    generated_at=datetime.now()
                )
            
            # Train model
            self.model.fit(features, target)
            
            # Make prediction for next period (ML)
            next_period_features = self._prepare_prediction_features(historical_data)
            ml_pred = float(self.model.predict([next_period_features])[0])

            # Time-series EMA per-user on daily expenses as a complementary signal
            ts_pred = self._predict_with_ema(historical_data)

            # Blend predictions when both available; otherwise fallback to whichever exists
            if ts_pred is not None:
                predicted_amount = 0.5 * ml_pred + 0.5 * ts_pred
            else:
                predicted_amount = ml_pred
            
            # Calculate confidence score (increase slightly when both signals agree)
            base_conf = min(0.95, len(historical_data) / 100.0)
            if ts_pred is not None:
                # agreement factor: closer predictions => higher confidence
                denom = max(1.0, abs(ml_pred) + abs(ts_pred))
                agree = 1.0 - min(1.0, abs(ml_pred - ts_pred) / denom)
                confidence_score = min(0.99, base_conf * (0.9 + 0.1 * agree))
            else:
                confidence_score = base_conf
            
            # Generate category breakdown
            category_breakdown = self._generate_category_breakdown(historical_data)
            
            # Generate trends
            trends = self._generate_trends(historical_data)
            
            # Generate recommendations
            recommendations = self._generate_recommendations(historical_data, predicted_amount)
            
            return ExpensePredictionResponse(
                user_id=request.user_id,
                predicted_amount=max(0, predicted_amount),
                confidence_score=confidence_score,
                category_breakdown=category_breakdown,
                trends=trends,
                recommendations=recommendations,
                generated_at=datetime.now().isoformat() + "Z"
            )
            
        except Exception as e:
            logger.error(f"Failed to predict expenses: {e}")
            return ExpensePredictionResponse(
                user_id=request.user_id,
                predicted_amount=0,
                confidence_score=0.0,
                category_breakdown=[],
                trends=[],
                recommendations=[f"Lỗi trong quá trình dự đoán: {str(e)}"],
                generated_at=datetime.now().isoformat() + "Z"
            )

    def _predict_with_ema(self, historical_data: List[Dict[str, Any]]) -> Optional[float]:
        """Compute an EMA-based projection for next-period total expenses.
        Uses simple daily aggregation and alpha based on window size.
        Returns None if insufficient data.
        """
        try:
            if not historical_data:
                return None
            df = pd.DataFrame(historical_data)
            df['date'] = pd.to_datetime(df['date'])
            df = df[df['type'] == 'expense']
            if df.empty:
                return None
            # Daily total expenses
            daily = df.groupby(df['date'].dt.date)['amount'].sum().astype(float)
            if len(daily) < 5:
                return None
            # Exponential moving average
            span = max(5, min(20, len(daily)//2))
            ema = daily.ewm(span=span, adjust=False).mean()
            # Project next period as last EMA value scaled by simple factor for month length
            last_ema = float(ema.iloc[-1])
            # Approximate next-month total from daily EMA
            projected = last_ema * 30.0
            return max(0.0, projected)
        except Exception:
            return None
    
    async def _get_historical_data(self, user_id: int, start_date: str, end_date: str) -> List[Dict[str, Any]]:
        """Get historical transaction data from database"""
        try:
            query = (
                "SELECT t.amount, t.transaction_type AS type, t.transaction_date AS date, "
                "COALESCE(c.name, 'Unknown') AS category "
                "FROM transactions t "
                "LEFT JOIN categories c ON c.id = t.category_id "
                "WHERE t.user_id = %s AND t.transaction_date BETWEEN %s AND %s "
                "ORDER BY t.transaction_date ASC, t.id ASC"
            )
            params = (user_id, start_date[:10], end_date[:10])
            async with get_db() as db:
                rows = await db.execute(query, params)
            data: List[Dict[str, Any]] = []
            for r in rows:
                # Normalize fields to expected schema
                rec = {
                    "amount": float(r.get("amount", 0) or 0),
                    "category": r.get("category") or "Unknown",
                    "date": str(r.get("date")),
                    "type": r.get("type") or "expense",
                }
                data.append(rec)
            return data
            
        except Exception as e:
            logger.error(f"Failed to get historical data: {e}")
            return []
    
    def _prepare_training_data(self, historical_data: List[Dict[str, Any]]) -> tuple:
        """Prepare training data for ML model"""
        try:
            # Convert to DataFrame
            df = pd.DataFrame(historical_data)
            
            # Filter only expenses
            df = df[df['type'] == 'expense']
            
            if len(df) < 3:
                return [], []
            
            # Create features
            features = []
            target = []
            
            for i in range(len(df) - 1):
                # Features: month, day of week, category encoded
                month = pd.to_datetime(df.iloc[i]['date']).month
                day_of_week = pd.to_datetime(df.iloc[i]['date']).dayofweek
                
                # Encode category
                category_encoded = hash(df.iloc[i]['category']) % 100  # Simple encoding
                
                features.append([month, day_of_week, category_encoded])
                target.append(df.iloc[i]['amount'])
            
            return np.array(features), np.array(target)
            
        except Exception as e:
            logger.error(f"Failed to prepare training data: {e}")
            return [], []
    
    def _prepare_prediction_features(self, historical_data: List[Dict[str, Any]]) -> List[float]:
        """Prepare features for prediction"""
        try:
            # Use the most recent transaction as base
            if not historical_data:
                return [0, 0, 0]
            
            latest = historical_data[-1]
            month = pd.to_datetime(latest['date']).month
            day_of_week = pd.to_datetime(latest['date']).dayofweek
            category_encoded = hash(latest['category']) % 100
            
            return [month, day_of_week, category_encoded]
            
        except Exception as e:
            logger.error(f"Failed to prepare prediction features: {e}")
            return [0, 0, 0]
    
    def _generate_category_breakdown(self, historical_data: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
        """Generate category breakdown"""
        try:
            df = pd.DataFrame(historical_data)
            df = df[df['type'] == 'expense']
            
            if len(df) == 0:
                return []
            
            category_summary = df.groupby('category')['amount'].agg(['sum', 'count', 'mean']).reset_index()
            category_summary.columns = ['category', 'total_amount', 'transaction_count', 'average_amount']
            
            # Sort by total amount
            category_summary = category_summary.sort_values('total_amount', ascending=False)
            
            breakdown = []
            for _, row in category_summary.head(5).iterrows():
                breakdown.append({
                    'category_id': 0,  # Default category ID
                    'category_name': row['category'],
                    'predicted_amount': float(row['total_amount']),
                    'confidence_score': 0.8,
                    'trend': 'stable'
                })
            
            return breakdown
            
        except Exception as e:
            logger.error(f"Failed to generate category breakdown: {e}")
            return []
    
    def _generate_trends(self, historical_data: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
        """Generate spending trends"""
        try:
            df = pd.DataFrame(historical_data)
            df = df[df['type'] == 'expense']
            
            if len(df) < 3:
                return []
            
            df['date'] = pd.to_datetime(df['date'])
            df['month'] = df['date'].dt.to_period('M')
            
            monthly_spending = df.groupby('month')['amount'].sum().reset_index()
            monthly_spending['month'] = monthly_spending['month'].astype(str)
            
            trends = []
            for _, row in monthly_spending.iterrows():
                trends.append({
                    'period': row['month'],
                    'amount': float(row['amount']),
                    'change_percentage': 0.0,  # Simplified
                    'trend': 'increasing' if len(trends) > 0 and row['amount'] > trends[-1]['amount'] else 'stable'
                })
            
            return trends
            
        except Exception as e:
            logger.error(f"Failed to generate trends: {e}")
            return []
    
    def _generate_recommendations(self, historical_data: List[Dict[str, Any]], predicted_amount: float) -> List[str]:
        """Generate financial recommendations"""
        try:
            recommendations = []
            
            df = pd.DataFrame(historical_data)
            df_expenses = df[df['type'] == 'expense']
            df_income = df[df['type'] == 'income']
            
            if len(df_expenses) == 0:
                return ["Chưa có dữ liệu chi tiêu để đưa ra khuyến nghị"]
            
            # Calculate average monthly spending
            avg_monthly_spending = df_expenses['amount'].mean()
            
            # Spending analysis
            if predicted_amount > avg_monthly_spending * 1.2:
                recommendations.append("Chi tiêu dự đoán cao hơn mức trung bình. Hãy cân nhắc cắt giảm chi tiêu không cần thiết.")
            
            # Category analysis
            category_summary = df_expenses.groupby('category')['amount'].sum().sort_values(ascending=False)
            top_category = category_summary.index[0]
            top_amount = category_summary.iloc[0]
            
            if top_amount > df_expenses['amount'].sum() * 0.4:
                recommendations.append(f"Chi tiêu cho '{top_category}' chiếm tỷ lệ cao. Hãy cân nhắc phân bổ lại ngân sách.")
            
            # Income vs Expense
            if len(df_income) > 0:
                total_income = df_income['amount'].sum()
                total_expenses = df_expenses['amount'].sum()
                
                if total_expenses > total_income * 0.8:
                    recommendations.append("Chi tiêu gần bằng thu nhập. Hãy tăng cường tiết kiệm.")
                elif total_expenses < total_income * 0.5:
                    recommendations.append("Bạn đang tiết kiệm tốt! Hãy tiếp tục duy trì thói quen này.")
            
            # General recommendations
            recommendations.append("Hãy theo dõi chi tiêu hàng ngày để kiểm soát tài chính tốt hơn.")
            recommendations.append("Đặt mục tiêu tiết kiệm cụ thể và theo dõi tiến độ.")
            
            return recommendations[:5]  # Limit to 5 recommendations
            
        except Exception as e:
            logger.error(f"Failed to generate recommendations: {e}")
            return ["Không thể tạo khuyến nghị do lỗi xử lý dữ liệu"]
