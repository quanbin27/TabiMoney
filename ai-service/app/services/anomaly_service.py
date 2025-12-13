from typing import List, Dict, Any
from datetime import datetime
import logging

import numpy as np
from sklearn.ensemble import IsolationForest
from sklearn.preprocessing import StandardScaler

from app.core.database import get_db

logger = logging.getLogger(__name__)

class AnomalyService:
    def __init__(self):
        self._ready = False

    async def initialize(self):
        self._ready = True

    async def cleanup(self):
        self._ready = False

    def is_ready(self):
        return self._ready

    async def detect(self, user_id: int, start_date: str, end_date: str, threshold: float = 0.6) -> Dict[str, Any]:
        """Detect anomalies using IsolationForest with improved feature engineering.
        
        Args:
            user_id: User ID
            start_date: Start date (YYYY-MM-DD)
            end_date: End date (YYYY-MM-DD)
            threshold: Contamination rate (0.01-0.5), lower = fewer anomalies detected
        
        Returns:
            Dictionary with anomalies list, total count, and detection score
        """
        # Load transactions
        query = (
            "SELECT t.id as id, t.amount, t.transaction_type, t.transaction_date, t.category_id, c.name as category_name "
            "FROM transactions t LEFT JOIN categories c ON c.id = t.category_id "
            "WHERE t.user_id = %s AND t.transaction_type = 'expense' AND t.transaction_date BETWEEN %s AND %s "
            "ORDER BY t.transaction_date ASC, t.id ASC"
        )
        params = (user_id, start_date[:10], end_date[:10])
        async with get_db() as db:
            rows = await db.execute(query, params)

        if not rows:
            return {"anomalies": [], "total_anomalies": 0, "detection_score": 0.0, "generated_at": self._now_iso()}

        if len(rows) < 10:
            # Too little data for reliable anomaly detection
            logger.warning(f"Insufficient data for anomaly detection: {len(rows)} transactions")
            return {"anomalies": [], "total_anomalies": 0, "detection_score": 0.0, "generated_at": self._now_iso()}

        # Build enhanced feature matrix
        X, meta = self._build_features(rows)
        
        if X.shape[0] < 10:
            return {"anomalies": [], "total_anomalies": 0, "detection_score": 0.0, "generated_at": self._now_iso()}

        # Normalize features for better model performance
        scaler = StandardScaler()
        X_scaled = scaler.fit_transform(X)

        # Convert threshold to contamination rate
        # Threshold interpretation:
        # - threshold 0.1 = very lenient (high contamination ~0.3 = 30% anomalies)
        # - threshold 0.5 = moderate (medium contamination ~0.1 = 10% anomalies)  
        # - threshold 0.9 = very strict (low contamination ~0.01 = 1% anomalies)
        # Formula: contamination = 0.3 * (1 - threshold) + 0.01
        # This maps threshold [0.0, 1.0] to contamination [0.31, 0.01]
        contamination_rate = max(0.01, min(0.3, 0.3 * (1.0 - threshold) + 0.01))
        
        logger.info(f"Anomaly detection: {len(rows)} transactions, contamination={contamination_rate:.3f}")

        # Train IsolationForest model
        model = IsolationForest(
            n_estimators=200,
            contamination=contamination_rate,
            random_state=42,
            max_samples='auto',
            max_features=1.0
        )
        model.fit(X_scaled)
        
        # Get predictions and scores
        # decision_function: negative values = anomalies, positive = normal
        # Lower (more negative) = more anomalous
        scores = model.decision_function(X_scaled)
        preds = model.predict(X_scaled)  # -1 = anomaly, 1 = normal

        # Build anomalies list
        anomalies: List[Dict[str, Any]] = []
        for i, p in enumerate(preds):
            if p == -1:  # Only include predicted anomalies
                m = meta[i]
                # Normalize anomaly score to 0-1 range (higher = more anomalous)
                # decision_function returns negative for anomalies, so we normalize
                raw_score = scores[i]
                # Normalize: convert negative scores to 0-1 scale
                # Most anomalous will have score close to -0.5, normal close to 0.5
                normalized_score = max(0.0, min(1.0, (0.5 - raw_score) / 1.0))
                
                # Determine anomaly type based on features
                anomaly_type = self._determine_anomaly_type(m, X[i], scores[i])
                
                anomalies.append({
                    "transaction_id": int(m["id"]),
                    "amount": float(m["amount"] or 0.0),
                    "category_name": m.get("category_name") or "Unknown",
                    "anomaly_score": float(normalized_score),
                    "anomaly_type": anomaly_type,
                    "description": self._generate_description(m, anomaly_type, normalized_score),
                    "transaction_date": str(m["transaction_date"]),
                })

        # Sort by anomaly score (highest first)
        anomalies.sort(key=lambda x: x["anomaly_score"], reverse=True)

        # Calculate aggregate detection score (average normalized score of detected anomalies)
        if anomalies:
            detection_score = float(np.mean([a["anomaly_score"] for a in anomalies]))
        else:
            detection_score = 0.0

        logger.info(f"Detected {len(anomalies)} anomalies out of {len(rows)} transactions")

        return {
            "user_id": user_id,
            "anomalies": anomalies,
            "total_anomalies": len(anomalies),
            "detection_score": detection_score,
            "generated_at": self._now_iso()
        }

    def _build_features(self, rows: List[Dict[str, Any]]) -> tuple:
        """Build enhanced feature matrix from transaction data."""
        X = []
        meta = []
        
        # Pre-compute statistics for relative features
        amounts = [float(r["amount"] or 0.0) for r in rows]
        mean_amount = np.mean(amounts) if amounts else 1.0
        std_amount = np.std(amounts) if len(amounts) > 1 else 1.0
        median_amount = np.median(amounts) if amounts else 1.0
        
        # Category frequency
        category_counts = {}
        for r in rows:
            cat_id = int(r["category_id"] or 0)
            category_counts[cat_id] = category_counts.get(cat_id, 0) + 1
        total_transactions = len(rows)
        
        for r in rows:
            dt = datetime.strptime(str(r["transaction_date"]), "%Y-%m-%d")
            amt = float(r["amount"] or 0.0)
            cat_id = int(r["category_id"] or 0)
            
            # Feature 1: Log amount (handles wide range)
            amt_log = np.log1p(max(amt, 0.0))
            
            # Feature 2: Relative amount (z-score)
            z_score = (amt - mean_amount) / max(std_amount, 1.0) if std_amount > 0 else 0.0
            
            # Feature 3: Relative to median (more robust)
            median_ratio = amt / max(median_amount, 1.0)
            
            # Feature 4: Day of week (0=Monday, 6=Sunday)
            dow = dt.weekday()
            
            # Feature 5: Day of month (1-31)
            dom = dt.day
            
            # Feature 6: Month (1-12)
            month = dt.month
            
            # Feature 7: Category frequency (normalized)
            cat_freq = category_counts.get(cat_id, 0) / max(total_transactions, 1.0)
            
            # Feature 8: Category ID (normalized to 0-1)
            # Assuming max category ID is reasonable, normalize it
            cat_id_norm = cat_id / max(max(category_counts.keys()) if category_counts else 1, 1.0)
            
            # Feature 9: Is weekend (0 or 1)
            is_weekend = 1.0 if dow >= 5 else 0.0
            
            # Feature 10: Is end of month (last 3 days)
            is_month_end = 1.0 if dom >= 28 else 0.0
            
            X.append([
                amt_log,
                z_score,
                median_ratio,
                dow / 6.0,  # Normalize to 0-1
                dom / 31.0,  # Normalize to 0-1
                month / 12.0,  # Normalize to 0-1
                cat_freq,
                cat_id_norm,
                is_weekend,
                is_month_end
            ])
            meta.append(r)
        
        return np.array(X, dtype=float), meta

    def _determine_anomaly_type(self, meta: Dict[str, Any], features: np.ndarray, score: float) -> str:
        """Determine the type of anomaly based on features."""
        # features: [amt_log, z_score, median_ratio, dow_norm, dom_norm, month_norm, cat_freq, cat_id_norm, is_weekend, is_month_end]
        z_score = features[1]
        median_ratio = features[2]
        cat_freq = features[6]
        
        if abs(z_score) > 2.0 or median_ratio > 3.0:
            return "amount_pattern"
        elif cat_freq < 0.05:  # Rare category
            return "category_pattern"
        else:
            return "behavioral_pattern"

    def _generate_description(self, meta: Dict[str, Any], anomaly_type: str, score: float) -> str:
        """Generate human-readable description for anomaly."""
        amount = float(meta.get("amount", 0))
        category = meta.get("category_name", "Unknown")
        
        if anomaly_type == "amount_pattern":
            if score > 0.7:
                return f"Giao dịch có số tiền bất thường cao ({category})"
            else:
                return f"Giao dịch có số tiền khác thường ({category})"
        elif anomaly_type == "category_pattern":
            return f"Giao dịch ở danh mục hiếm gặp ({category})"
        else:
            return f"Giao dịch có hành vi bất thường ({category})"

    def _now_iso(self) -> str:
        """Return current UTC time in ISO format."""
        from datetime import timezone
        return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")

