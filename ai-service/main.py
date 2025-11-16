"""
TabiMoney AI Service
AI-Powered Personal Finance Management - AI Agent Service
"""

import logging
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.trustedhost import TrustedHostMiddleware
from prometheus_client import Counter, Histogram, generate_latest, CollectorRegistry
from fastapi.responses import Response

from app.core.config import settings
from app.core.database import init_db
from app.api.v1.api import api_router
from app.core.logging import setup_logging
from app.core.dependencies import set_services
from app.services.nlu_service import NLUService
from app.services.prediction_service import PredictionService
from app.services.anomaly_service import AnomalyService

# Setup logging
setup_logging()
logger = logging.getLogger(__name__)

# Prometheus metrics with isolated registry to avoid duplicates
METRICS_REGISTRY = globals().get("METRICS_REGISTRY") or CollectorRegistry()
REQUEST_COUNT = globals().get("REQUEST_COUNT") or Counter(
    'ai_service_requests_total', 'Total requests', ['method', 'endpoint'], registry=METRICS_REGISTRY
)
REQUEST_DURATION = globals().get("REQUEST_DURATION") or Histogram(
    'ai_service_request_duration_seconds', 'Request duration', registry=METRICS_REGISTRY
)

# Global services
nlu_service: NLUService = None
prediction_service: PredictionService = None
anomaly_service: AnomalyService = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager"""
    global nlu_service, prediction_service, anomaly_service
    
    logger.info("Starting AI Service...")
    
    # Initialize database
    await init_db()
    logger.info("Database initialized")
    
    # Initialize services
    nlu_service = NLUService()
    await nlu_service.initialize()
    logger.info("NLU Service initialized")
    
    prediction_service = PredictionService()
    await prediction_service.initialize()
    logger.info("Prediction Service initialized")
    
    anomaly_service = AnomalyService()
    await anomaly_service.initialize()
    logger.info("Anomaly Service initialized")
    
    # Set global service instances for dependency injection
    set_services(nlu_service, prediction_service, anomaly_service)
    logger.info("Services registered for dependency injection")
    
    logger.info("AI Service started successfully")
    
    yield
    
    # Cleanup
    logger.info("Shutting down AI Service...")
    if nlu_service:
        await nlu_service.cleanup()
    if prediction_service:
        await prediction_service.cleanup()
    if anomaly_service:
        await anomaly_service.cleanup()
    logger.info("AI Service shutdown complete")


# Create FastAPI app
app = FastAPI(
    title="TabiMoney AI Service",
    description="AI-Powered Personal Finance Management - AI Agent Service",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
    lifespan=lifespan
)

# Add middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.add_middleware(
    TrustedHostMiddleware,
    allowed_hosts=settings.ALLOWED_HOSTS
)

# Include API routes
app.include_router(api_router, prefix="/api/v1")


@app.get("/health")
async def health_check():
    """Health check endpoint - returns healthy if service is running, even if models are still initializing"""
    try:
        # Basic health check - service is running
        # Models can still be initializing, that's OK for health check
        return {
            "status": "healthy",
            "service": "ai-service",
            "version": "1.0.0",
            "models": {
                "nlu": nlu_service.is_ready() if nlu_service else False,
                "prediction": prediction_service.is_ready() if prediction_service else False,
                "anomaly": anomaly_service.is_ready() if anomaly_service else False,
            }
        }
    except Exception as e:
        # Only fail if there's a critical error
        logger.error(f"Health check failed: {e}")
        # Return 200 with unhealthy status instead of raising exception
        # This allows service to start even if models are still initializing
        return {
            "status": "unhealthy",
            "service": "ai-service",
            "error": str(e)
        }


@app.get("/metrics")
async def metrics():
    """Prometheus metrics endpoint"""
    return Response(generate_latest(), media_type="text/plain")


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "TabiMoney AI Service",
        "version": "1.0.0",
        "docs": "/docs",
        "health": "/health"
    }


if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level="info"
    )
