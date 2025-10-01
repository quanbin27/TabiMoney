"""
NLU API endpoints
"""

import logging
from typing import Dict, Any

from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel

from app.models.nlu import NLURequest
from app.services.nlu_service import NLUService
from app.core.dependencies import get_nlu_service

logger = logging.getLogger(__name__)

router = APIRouter()


class NLURequestModel(BaseModel):
    text: str
    user_id: int
    context: str = ""


class NLUResponseModel(BaseModel):
    user_id: int
    intent: str
    entities: list
    confidence: float
    suggested_action: str
    response: str
    generated_at: str


@router.post("/process", response_model=NLUResponseModel)
async def process_nlu(
    request: NLURequestModel,
    nlu_service: NLUService = Depends(get_nlu_service)
):
    """
    Process Natural Language Understanding request
    """
    try:
        # Convert to internal model
        nlu_request = NLURequest(
            text=request.text,
            user_id=request.user_id,
            context=request.context
        )
        
        # Process with NLU service
        response = await nlu_service.process_nlu(nlu_request)
        
        # Convert to response model
        return NLUResponseModel(
            user_id=response.user_id,
            intent=response.intent,
            entities=[entity.dict() for entity in response.entities],
            confidence=response.confidence,
            suggested_action=response.suggested_action,
            response=response.response,
            generated_at=response.generated_at.isoformat()
        )
        
    except Exception as e:
        logger.error(f"NLU processing failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


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
