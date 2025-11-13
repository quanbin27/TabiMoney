from typing import List, Dict, Any
from datetime import datetime

import numpy as np
from sklearn.ensemble import IsolationForest

from app.core.database import get_db

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
        """Detect anomalies using IsolationForest over basic features.
        threshold is contamination cutoff for anomaly score (0-1, lower => fewer anomalies).
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
            return {"anomalies": [], "total_anomalies": 0, "detection_score": 0.0}

        # Build feature matrix: amount (log), day_of_week, month, category_id (as-is)
        X = []
        meta = []
        for r in rows:
            dt = datetime.strptime(str(r["transaction_date"]), "%Y-%m-%d")
            amt = float(r["amount"] or 0.0)
            amt_log = np.log1p(max(amt, 0.0))
            dow = dt.weekday()
            month = dt.month
            cat_id = int(r["category_id"] or 0)
            X.append([amt_log, dow, month, cat_id])
            meta.append(r)

        X = np.array(X, dtype=float)
        if len(X) < 10:
            # Too little data for IF, return empty to avoid noise
            return {"anomalies": [], "total_anomalies": 0, "detection_score": 0.0}

        contamination = min(max(threshold, 0.01), 0.4)  # clamp
        model = IsolationForest(n_estimators=200, contamination=contamination, random_state=42)
        model.fit(X)
        scores = model.decision_function(X)  # higher => normal; lower => anomalous
        preds = model.predict(X)  # -1 anomaly, 1 normal

        anomalies: List[Dict[str, Any]] = []
        for i, p in enumerate(preds):
            if p == -1:
                m = meta[i]
                anomalies.append({
                    "transaction_id": int(m["id"]),
                    "amount": float(m["amount"] or 0.0),
                    "category_name": m.get("category_name") or "",
                    "anomaly_score": float(-scores[i]),  # invert to make higher = more anomalous
                    "anomaly_type": "amount_pattern",
                    "description": "Giao dịch có mẫu khác thường theo mô hình IsolationForest",
                    "transaction_date": str(m["transaction_date"]),
                })

        # Aggregate score: average of anomaly scores (normalized)
        detection_score = float(np.clip(np.mean([-s for i, s in enumerate(scores) if preds[i] == -1]) if any(preds == -1) else 0.0, 0.0, 1.0))

        return {
            "anomalies": anomalies,
            "total_anomalies": len(anomalies),
            "detection_score": detection_score,
        }


