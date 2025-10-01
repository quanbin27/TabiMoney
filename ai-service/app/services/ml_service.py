"""
Machine Learning Service
Core ML functionality for TabiMoney AI Service
"""

import asyncio
import logging
import pickle
import os
from typing import Dict, Any, List, Optional, Tuple
from datetime import datetime, timedelta

import numpy as np
import pandas as pd
from sklearn.ensemble import IsolationForest, RandomForestClassifier
from sklearn.preprocessing import StandardScaler, LabelEncoder
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score, classification_report
import joblib

from app.core.config import settings
from app.core.database import get_db

logger = logging.getLogger(__name__)


class MLService:
    """Machine Learning Service for financial data analysis"""
    
    def __init__(self):
        self.models: Dict[str, Any] = {}
        self.scalers: Dict[str, Any] = {}
        self.encoders: Dict[str, Any] = {}
        self.is_initialized = False
        self.model_cache_dir = settings.MODEL_CACHE_DIR
        
        # Ensure model cache directory exists
        os.makedirs(self.model_cache_dir, exist_ok=True)
    
    async def initialize(self):
        """Initialize ML service"""
        logger.info("Initializing ML Service...")
        
        try:
            # Load or train models
            await self._load_or_train_models()
            self.is_initialized = True
            logger.info("ML Service initialized successfully")
        except Exception as e:
            logger.error(f"Failed to initialize ML Service: {e}")
            raise
    
    async def cleanup(self):
        """Cleanup ML service"""
        logger.info("Cleaning up ML Service...")
        self.models.clear()
        self.scalers.clear()
        self.encoders.clear()
        self.is_initialized = False
    
    def is_ready(self) -> bool:
        """Check if ML service is ready"""
        return self.is_initialized and len(self.models) > 0
    
    async def _load_or_train_models(self):
        """Load existing models or train new ones"""
        # Try to load existing models
        if await self._load_models():
            logger.info("Loaded existing ML models")
            return
        
        # Train new models
        logger.info("Training new ML models...")
        await self._train_models()
        await self._save_models()
        logger.info("ML models trained and saved")
    
    async def _load_models(self) -> bool:
        """Load models from cache"""
        try:
            model_files = [
                "category_classifier.pkl",
                "anomaly_detector.pkl",
                "expense_predictor.pkl",
                "category_scaler.pkl",
                "category_encoder.pkl"
            ]
            
            for model_file in model_files:
                model_path = os.path.join(self.model_cache_dir, model_file)
                if not os.path.exists(model_path):
                    return False
                
                if model_file.endswith("_scaler.pkl"):
                    self.scalers[model_file.replace("_scaler.pkl", "")] = joblib.load(model_path)
                elif model_file.endswith("_encoder.pkl"):
                    self.encoders[model_file.replace("_encoder.pkl", "")] = joblib.load(model_path)
                else:
                    model_name = model_file.replace(".pkl", "")
                    self.models[model_name] = joblib.load(model_path)
            
            return True
        except Exception as e:
            logger.error(f"Failed to load models: {e}")
            return False
    
    async def _save_models(self):
        """Save models to cache"""
        try:
            # Save models
            for model_name, model in self.models.items():
                model_path = os.path.join(self.model_cache_dir, f"{model_name}.pkl")
                joblib.dump(model, model_path)
            
            # Save scalers
            for scaler_name, scaler in self.scalers.items():
                scaler_path = os.path.join(self.model_cache_dir, f"{scaler_name}_scaler.pkl")
                joblib.dump(scaler, scaler_path)
            
            # Save encoders
            for encoder_name, encoder in self.encoders.items():
                encoder_path = os.path.join(self.model_cache_dir, f"{encoder_name}_encoder.pkl")
                joblib.dump(encoder, encoder_path)
            
            logger.info("Models saved successfully")
        except Exception as e:
            logger.error(f"Failed to save models: {e}")
    
    async def _train_models(self):
        """Train ML models"""
        # Get training data
        training_data = await self._get_training_data()
        
        if training_data.empty:
            logger.warning("No training data available, using default models")
            await self._create_default_models()
            return
        
        # Train category classifier
        await self._train_category_classifier(training_data)
        
        # Train anomaly detector
        await self._train_anomaly_detector(training_data)
        
        # Train expense predictor
        await self._train_expense_predictor(training_data)
    
    async def _get_training_data(self) -> pd.DataFrame:
        """Get training data from database"""
        try:
            async with get_db() as db:
                # Get transactions from last 6 months
                six_months_ago = datetime.now() - timedelta(days=180)
                
                transactions = await db.execute(
                    """
                    SELECT 
                        t.amount,
                        t.description,
                        t.transaction_type,
                        t.transaction_date,
                        t.location,
                        c.name as category_name,
                        c.id as category_id
                    FROM transactions t
                    JOIN categories c ON t.category_id = c.id
                    WHERE t.transaction_date >= %s
                    ORDER BY t.transaction_date DESC
                    """,
                    (six_months_ago,)
                )
                
                if not transactions:
                    return pd.DataFrame()
                
                # Convert to DataFrame
                df = pd.DataFrame(transactions)
                
                # Feature engineering
                df = self._engineer_features(df)
                
                return df
                
        except Exception as e:
            logger.error(f"Failed to get training data: {e}")
            return pd.DataFrame()
    
    def _engineer_features(self, df: pd.DataFrame) -> pd.DataFrame:
        """Engineer features for ML models"""
        # Extract features from description
        df['description_length'] = df['description'].str.len()
        df['has_numbers'] = df['description'].str.contains(r'\d+', regex=True)
        df['word_count'] = df['description'].str.split().str.len()
        
        # Extract time features
        df['transaction_date'] = pd.to_datetime(df['transaction_date'])
        df['day_of_week'] = df['transaction_date'].dt.dayofweek
        df['day_of_month'] = df['transaction_date'].dt.day
        df['month'] = df['transaction_date'].dt.month
        df['hour'] = df['transaction_date'].dt.hour
        
        # Amount features
        df['amount_log'] = np.log1p(df['amount'])
        df['amount_scaled'] = df['amount'] / df['amount'].std()
        
        # Location features
        df['has_location'] = df['location'].notna()
        
        return df
    
    async def _train_category_classifier(self, df: pd.DataFrame):
        """Train category classification model"""
        try:
            # Prepare features and target
            feature_columns = [
                'amount', 'description_length', 'has_numbers', 'word_count',
                'day_of_week', 'day_of_month', 'month', 'hour',
                'amount_log', 'amount_scaled', 'has_location'
            ]
            
            X = df[feature_columns].fillna(0)
            y = df['category_id']
            
            # Encode target
            label_encoder = LabelEncoder()
            y_encoded = label_encoder.fit_transform(y)
            
            # Scale features
            scaler = StandardScaler()
            X_scaled = scaler.fit_transform(X)
            
            # Train model
            classifier = RandomForestClassifier(
                n_estimators=100,
                max_depth=10,
                random_state=42,
                n_jobs=-1
            )
            
            classifier.fit(X_scaled, y_encoded)
            
            # Store model and preprocessors
            self.models['category_classifier'] = classifier
            self.scalers['category'] = scaler
            self.encoders['category'] = label_encoder
            
            logger.info("Category classifier trained successfully")
            
        except Exception as e:
            logger.error(f"Failed to train category classifier: {e}")
            raise
    
    async def _train_anomaly_detector(self, df: pd.DataFrame):
        """Train anomaly detection model"""
        try:
            # Prepare features for anomaly detection
            feature_columns = [
                'amount', 'description_length', 'word_count',
                'day_of_week', 'day_of_month', 'month', 'hour'
            ]
            
            X = df[feature_columns].fillna(0)
            
            # Scale features
            scaler = StandardScaler()
            X_scaled = scaler.fit_transform(X)
            
            # Train anomaly detector
            anomaly_detector = IsolationForest(
                contamination=0.1,  # 10% of data is considered anomalous
                random_state=42,
                n_jobs=-1
            )
            
            anomaly_detector.fit(X_scaled)
            
            # Store model and scaler
            self.models['anomaly_detector'] = anomaly_detector
            self.scalers['anomaly'] = scaler
            
            logger.info("Anomaly detector trained successfully")
            
        except Exception as e:
            logger.error(f"Failed to train anomaly detector: {e}")
            raise
    
    async def _train_expense_predictor(self, df: pd.DataFrame):
        """Train expense prediction model"""
        try:
            # Filter expense transactions
            expense_df = df[df['transaction_type'] == 'expense'].copy()
            
            if expense_df.empty:
                logger.warning("No expense data for prediction model")
                return
            
            # Prepare features
            feature_columns = [
                'day_of_week', 'day_of_month', 'month', 'hour',
                'description_length', 'word_count', 'has_location'
            ]
            
            X = expense_df[feature_columns].fillna(0)
            y = expense_df['amount']
            
            # Scale features
            scaler = StandardScaler()
            X_scaled = scaler.fit_transform(X)
            
            # Train model (using Random Forest for regression)
            predictor = RandomForestClassifier(
                n_estimators=100,
                max_depth=10,
                random_state=42,
                n_jobs=-1
            )
            
            predictor.fit(X_scaled, y)
            
            # Store model and scaler
            self.models['expense_predictor'] = predictor
            self.scalers['expense'] = scaler
            
            logger.info("Expense predictor trained successfully")
            
        except Exception as e:
            logger.error(f"Failed to train expense predictor: {e}")
            raise
    
    async def _create_default_models(self):
        """Create default models when no training data is available"""
        # Default category classifier
        self.models['category_classifier'] = RandomForestClassifier(
            n_estimators=10,
            random_state=42
        )
        
        # Default anomaly detector
        self.models['anomaly_detector'] = IsolationForest(
            contamination=0.1,
            random_state=42
        )
        
        # Default expense predictor
        self.models['expense_predictor'] = RandomForestClassifier(
            n_estimators=10,
            random_state=42
        )
        
        # Default scalers and encoders
        self.scalers['category'] = StandardScaler()
        self.scalers['anomaly'] = StandardScaler()
        self.scalers['expense'] = StandardScaler()
        self.encoders['category'] = LabelEncoder()
        
        logger.info("Default models created")
    
    async def predict_category(self, transaction_data: Dict[str, Any]) -> Tuple[int, float]:
        """Predict category for a transaction"""
        if not self.is_ready():
            raise RuntimeError("ML Service not ready")
        
        try:
            # Prepare features
            features = self._prepare_category_features(transaction_data)
            
            # Scale features
            scaler = self.scalers['category']
            features_scaled = scaler.transform([features])
            
            # Predict
            classifier = self.models['category_classifier']
            prediction = classifier.predict(features_scaled)[0]
            probability = classifier.predict_proba(features_scaled)[0].max()
            
            # Decode prediction
            encoder = self.encoders['category']
            category_id = encoder.inverse_transform([prediction])[0]
            
            return int(category_id), float(probability)
            
        except Exception as e:
            logger.error(f"Failed to predict category: {e}")
            raise
    
    async def detect_anomaly(self, transaction_data: Dict[str, Any]) -> Tuple[bool, float]:
        """Detect if a transaction is anomalous"""
        if not self.is_ready():
            raise RuntimeError("ML Service not ready")
        
        try:
            # Prepare features
            features = self._prepare_anomaly_features(transaction_data)
            
            # Scale features
            scaler = self.scalers['anomaly']
            features_scaled = scaler.transform([features])
            
            # Detect anomaly
            detector = self.models['anomaly_detector']
            anomaly_score = detector.decision_function(features_scaled)[0]
            is_anomaly = detector.predict(features_scaled)[0] == -1
            
            return bool(is_anomaly), float(anomaly_score)
            
        except Exception as e:
            logger.error(f"Failed to detect anomaly: {e}")
            raise
    
    async def predict_expense(self, user_data: Dict[str, Any]) -> float:
        """Predict future expenses"""
        if not self.is_ready():
            raise RuntimeError("ML Service not ready")
        
        try:
            # Prepare features
            features = self._prepare_expense_features(user_data)
            
            # Scale features
            scaler = self.scalers['expense']
            features_scaled = scaler.transform([features])
            
            # Predict
            predictor = self.models['expense_predictor']
            prediction = predictor.predict(features_scaled)[0]
            
            return float(prediction)
            
        except Exception as e:
            logger.error(f"Failed to predict expense: {e}")
            raise
    
    def _prepare_category_features(self, data: Dict[str, Any]) -> List[float]:
        """Prepare features for category classification"""
        features = [
            float(data.get('amount', 0)),
            float(len(data.get('description', ''))),
            float(1 if any(c.isdigit() for c in data.get('description', '')) else 0),
            float(len(data.get('description', '').split())),
            float(data.get('day_of_week', 0)),
            float(data.get('day_of_month', 0)),
            float(data.get('month', 0)),
            float(data.get('hour', 0)),
            float(np.log1p(data.get('amount', 0))),
            float(data.get('amount', 0) / 1000000),  # Simple scaling
            float(1 if data.get('location') else 0)
        ]
        return features
    
    def _prepare_anomaly_features(self, data: Dict[str, Any]) -> List[float]:
        """Prepare features for anomaly detection"""
        features = [
            float(data.get('amount', 0)),
            float(len(data.get('description', ''))),
            float(len(data.get('description', '').split())),
            float(data.get('day_of_week', 0)),
            float(data.get('day_of_month', 0)),
            float(data.get('month', 0)),
            float(data.get('hour', 0))
        ]
        return features
    
    def _prepare_expense_features(self, data: Dict[str, Any]) -> List[float]:
        """Prepare features for expense prediction"""
        features = [
            float(data.get('day_of_week', 0)),
            float(data.get('day_of_month', 0)),
            float(data.get('month', 0)),
            float(data.get('hour', 0)),
            float(len(data.get('description', ''))),
            float(len(data.get('description', '').split())),
            float(1 if data.get('location') else 0)
        ]
        return features
