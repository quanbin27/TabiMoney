"""
NLU API endpoints
"""

import logging

from fastapi import APIRouter, Depends

from app.services.nlu_service import NLUService
from app.core.dependencies import get_nlu_service

logger = logging.getLogger(__name__)

router = APIRouter()


@router.get("/health")
async def nlu_health(nlu_service: NLUService = Depends(get_nlu_service)):
    """
    Check NLU service health
    """
    return {
        "status": "healthy" if nlu_service.is_ready() else "unhealthy",
        "service": "nlu",
        "ready": nlu_service.is_ready()
    }
