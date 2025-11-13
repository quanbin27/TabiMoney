"""
Prediction Service for TabiMoney AI Service
Handles expense prediction and financial forecasting
"""

import logging
from typing import Dict, Any, List, Optional, Tuple
from datetime import datetime, timezone
import pandas as pd
import numpy as np
from sklearn.ensemble import RandomForestRegressor

from app.core.config import settings
from app.core.database import get_db
from app.models.prediction import ExpensePredictionRequest, ExpensePredictionResponse

logger = logging.getLogger(__name__)


class PredictionService:
    """Prediction Service for financial forecasting"""
    
    def __init__(self):
        self._ready = False
        self.model = None
        # Per-user model cache and series fingerprint to avoid retraining
        self._user_models: Dict[int, RandomForestRegressor] = {}
        self._user_series_fp: Dict[int, str] = {}
        
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
            
            if len(historical_data) < 5:  # Need at least minimal transactions
                logger.warning(f"Not enough data for prediction: {len(historical_data)} transactions")
                return ExpensePredictionResponse(
                    user_id=request.user_id,
                    predicted_amount=0,
                    confidence_score=0.0,
                    category_breakdown=[],
                    trends=[],
                    recommendations=["Cần thêm dữ liệu giao dịch để đưa ra dự đoán chính xác"],
                    generated_at=self._now_iso()
                )
            
            # Aggregate to monthly totals and build training data (monthly-level)
            monthly_df = self._build_monthly_series(historical_data)
            if monthly_df.empty or len(monthly_df) < 3:
                return ExpensePredictionResponse(
                    user_id=request.user_id,
                    predicted_amount=0,
                    confidence_score=0.0,
                    category_breakdown=[],
                    trends=[],
                    recommendations=["Dữ liệu không đủ để dự đoán theo tháng"],
                    generated_at=self._now_iso()
                )
            X, y = self._prepare_monthly_training_data(monthly_df)
            if len(X) < 2:
                return ExpensePredictionResponse(
                    user_id=request.user_id,
                    predicted_amount=0,
                    confidence_score=0.0,
                    category_breakdown=[],
                    trends=[],
                    recommendations=["Dữ liệu không đủ để dự đoán theo tháng"],
                    generated_at=self._now_iso()
                )
            
            # Train or reuse per-user model
            model = self._get_or_train_user_model(request.user_id, X, y, monthly_df)
            # Prepare features to predict next month total
            next_features = self._prepare_next_month_features(monthly_df)
            ml_pred = float(model.predict([next_features])[0])

            # Time-series EMA per-user on daily expenses as a complementary signal
            ts_pred = self._predict_with_ema(historical_data)

            # Blend predictions when both available; otherwise fallback to whichever exists
            if ts_pred is not None:
                # ts_pred is monthly-level estimate; ml_pred is monthly-level too
                predicted_amount = max(0.0, 0.6 * ml_pred + 0.4 * ts_pred)
            else:
                predicted_amount = max(0.0, ml_pred)
            
            # Calculate confidence score (increase slightly when both signals agree)
            # Base on number of months rather than raw transactions
            base_conf = min(0.95, len(monthly_df) / 36.0)  # scale to ~3 years
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
                generated_at=self._now_iso()
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
                generated_at=self._now_iso()
            )

    def _get_or_train_user_model(
        self, user_id: int, X: np.ndarray, y: np.ndarray, monthly_df: pd.DataFrame
    ) -> RandomForestRegressor:
        """Return cached model for user or train if missing/outdated."""
        fp = self._fingerprint_monthly_series(monthly_df)
        mdl = self._user_models.get(user_id)
        last_fp = self._user_series_fp.get(user_id)
        if mdl is None or last_fp != fp:
            mdl = RandomForestRegressor(n_estimators=200, random_state=42, max_depth=12)
            mdl.fit(X, y)
            self._user_models[user_id] = mdl
            self._user_series_fp[user_id] = fp
        return mdl

    def _now_iso(self) -> str:
        """Return current UTC time in RFC3339 format with trailing Z."""
        return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")

    def _fingerprint_monthly_series(self, monthly_df: pd.DataFrame) -> str:
        """Compute a lightweight fingerprint of monthly series to detect changes."""
        try:
            # Use last 24 months totals to fingerprint
            s = monthly_df['total_expense'].tail(24).round(2).astype(str).tolist()
            return "|".join(s)
        except Exception:
            return str(len(monthly_df))

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

    def _build_monthly_series(self, historical_data: List[Dict[str, Any]]) -> pd.DataFrame:
        """Aggregate transactions to monthly totals and compute helper columns."""
        df = pd.DataFrame(historical_data)
        if df.empty:
            return pd.DataFrame()
        df['date'] = pd.to_datetime(df['date'])
        # Focus on expenses only
        df = df[df['type'] == 'expense'].copy()
        if df.empty:
            return pd.DataFrame()
        df['year_month'] = df['date'].dt.to_period('M').astype(str)
        monthly = df.groupby('year_month')['amount'].sum().reset_index()
        monthly = monthly.sort_values('year_month')
        monthly['total_expense'] = monthly['amount'].astype(float)
        # Add temporal indices
        monthly['month'] = pd.to_datetime(monthly['year_month']).dt.month
        monthly['year'] = pd.to_datetime(monthly['year_month']).dt.year
        # Rolling stats
        monthly['roll_mean_3'] = monthly['total_expense'].rolling(window=3, min_periods=1).mean()
        monthly['roll_mean_6'] = monthly['total_expense'].rolling(window=6, min_periods=1).mean()
        monthly['roll_std_6'] = monthly['total_expense'].rolling(window=6, min_periods=1).std().fillna(0.0)
        monthly['count_seen'] = np.arange(1, len(monthly) + 1)
        return monthly[['year_month', 'year', 'month', 'total_expense', 'roll_mean_3', 'roll_mean_6', 'roll_std_6', 'count_seen']]

    def _prepare_monthly_training_data(self, monthly_df: pd.DataFrame) -> Tuple[np.ndarray, np.ndarray]:
        """Build (X, y) to predict next month's total based on prior months."""
        # Use all but last month for training; predict last-known as validation
        # Features: month, roll_mean_3 (t-1), roll_mean_6 (t-1), roll_std_6 (t-1), count_seen
        # Target: total_expense of current month
        if len(monthly_df) < 2:
            return np.array([]), np.array([])
        df = monthly_df.copy().reset_index(drop=True)
        # Shift rolling stats to avoid leakage
        df['rm3_prev'] = df['roll_mean_3'].shift(1)
        df['rm6_prev'] = df['roll_mean_6'].shift(1)
        df['rs6_prev'] = df['roll_std_6'].shift(1)
        df['count_prev'] = (df['count_seen'] - 1).clip(lower=0)
        # Drop first row due to shift
        df = df.iloc[1:].reset_index(drop=True)
        X = df[['month', 'rm3_prev', 'rm6_prev', 'rs6_prev', 'count_prev']].fillna(0.0).astype(float).values
        y = df['total_expense'].astype(float).values
        return X, y

    def _prepare_next_month_features(self, monthly_df: pd.DataFrame) -> List[float]:
        """Prepare features for the next month based on latest known month."""
        df = monthly_df.copy().reset_index(drop=True)
        if df.empty:
            return [0.0, 0.0, 0.0, 0.0, 0.0]
        last = df.iloc[-1]
        # Estimate next month index
        next_month = int(last['month']) + 1
        if next_month == 13:
            next_month = 1
        rm3 = float(df['total_expense'].tail(3).mean())
        rm6 = float(df['total_expense'].tail(6).mean()) if len(df) >= 6 else rm3
        rs6 = float(df['total_expense'].tail(6).std()) if len(df) >= 6 else 0.0
        count_prev = float(last['count_seen'])
        return [float(next_month), rm3, rm6, rs6, count_prev]
    
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
