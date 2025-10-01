"""
Dependencies for FastAPI
"""

from app.services.nlu_service import NLUService
from app.services.prediction_service import PredictionService
from app.services.anomaly_service import AnomalyService
from app.services.ml_service import MLService

# Global service instances (will be set in main.py)
_nlu_service: NLUService = None
_prediction_service: PredictionService = None
_anomaly_service: AnomalyService = None
_ml_service: MLService = None


def set_services(
    nlu_service: NLUService,
    prediction_service: PredictionService,
    anomaly_service: AnomalyService,
    ml_service: MLService
):
    """Set global service instances"""
    global _nlu_service, _prediction_service, _anomaly_service, _ml_service
    _nlu_service = nlu_service
    _prediction_service = prediction_service
    _anomaly_service = anomaly_service
    _ml_service = ml_service


def get_nlu_service() -> NLUService:
    """Get NLU service instance"""
    if _nlu_service is None:
        raise RuntimeError("NLU service not initialized")
    return _nlu_service


def get_prediction_service() -> PredictionService:
    """Get prediction service instance"""
    if _prediction_service is None:
        raise RuntimeError("Prediction service not initialized")
    return _prediction_service


def get_anomaly_service() -> AnomalyService:
    """Get anomaly service instance"""
    if _anomaly_service is None:
        raise RuntimeError("Anomaly service not initialized")
    return _anomaly_service


def get_ml_service() -> MLService:
    """Get ML service instance"""
    if _ml_service is None:
        raise RuntimeError("ML service not initialized")
    return _ml_service
