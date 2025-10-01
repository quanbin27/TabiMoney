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
            # Get historical data from database
            historical_data = await self._get_historical_data(request.user_id, request.start_date, request.end_date)
            
            if len(historical_data) < 5:  # Need at least 5 transactions for prediction
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
            
            # Make prediction for next period
            next_period_features = self._prepare_prediction_features(historical_data)
            predicted_amount = self.model.predict([next_period_features])[0]
            
            # Calculate confidence score
            confidence_score = min(0.95, len(historical_data) / 100.0)
            
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
                generated_at=datetime.now()
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
                generated_at=datetime.now()
            )
    
    async def _get_historical_data(self, user_id: int, start_date: str, end_date: str) -> List[Dict[str, Any]]:
        """Get historical transaction data from database"""
        try:
            # This would normally query the database
            # For now, we'll return mock data based on the transactions we created
            
            # Mock historical data based on our created transactions
            mock_data = [
                {"amount": 50000, "category": "Ăn uống", "date": "2025-09-15", "type": "expense"},
                {"amount": 75000, "category": "Ăn uống", "date": "2025-09-14", "type": "expense"},
                {"amount": 25000, "category": "Ăn uống", "date": "2025-09-13", "type": "expense"},
                {"amount": 80000, "category": "Giao thông", "date": "2025-09-10", "type": "expense"},
                {"amount": 45000, "category": "Giao thông", "date": "2025-09-09", "type": "expense"},
                {"amount": 500000, "category": "Mua sắm", "date": "2025-09-05", "type": "expense"},
                {"amount": 15000000, "category": "Mua sắm", "date": "2025-09-04", "type": "expense"},
                {"amount": 150000, "category": "Giải trí", "date": "2025-08-31", "type": "expense"},
                {"amount": 200000, "category": "Giải trí", "date": "2025-08-30", "type": "expense"},
                {"amount": 500000, "category": "Y tế", "date": "2025-08-26", "type": "expense"},
                {"amount": 150000, "category": "Y tế", "date": "2025-08-25", "type": "expense"},
                {"amount": 2000000, "category": "Giáo dục", "date": "2025-08-21", "type": "expense"},
                {"amount": 500000, "category": "Giáo dục", "date": "2025-08-20", "type": "expense"},
                {"amount": 5000000, "category": "Du lịch", "date": "2025-08-16", "type": "expense"},
                {"amount": 2000000, "category": "Du lịch", "date": "2025-08-15", "type": "expense"},
                {"amount": 400000, "category": "Tiện ích", "date": "2025-08-11", "type": "expense"},
                {"amount": 200000, "category": "Tiện ích", "date": "2025-08-10", "type": "expense"},
                {"amount": 15000000, "category": "Ăn uống", "date": "2025-08-01", "type": "income"},
                {"amount": 5000000, "category": "Ăn uống", "date": "2025-07-15", "type": "income"},
                {"amount": 8000000, "category": "Ăn uống", "date": "2025-07-01", "type": "income"},
            ]
            
            return mock_data
            
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
                    'category': row['category'],
                    'total_amount': float(row['total_amount']),
                    'transaction_count': int(row['transaction_count']),
                    'average_amount': float(row['average_amount']),
                    'percentage': float(row['total_amount'] / df['amount'].sum() * 100)
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




