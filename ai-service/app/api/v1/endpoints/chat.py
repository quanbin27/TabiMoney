from fastapi import APIRouter, Depends, HTTPException
from app.core.dependencies import get_nlu_service
from app.models.nlu import ChatRequest, ChatResponse
from app.services.nlu_service import NLUService
import logging

logger = logging.getLogger(__name__)

router = APIRouter()

@router.post("/process", response_model=ChatResponse)
async def process_chat(
    request: ChatRequest,
    nlu_service: NLUService = Depends(get_nlu_service)
):
    """Process chat message and return AI response"""
    try:
        # Validate request
        if not request.message or not request.message.strip():
            raise HTTPException(status_code=400, detail="Message cannot be empty")
        
        if not request.user_id or request.user_id <= 0:
            raise HTTPException(status_code=400, detail="Invalid user_id")
        
        # Process chat
        response = await nlu_service.process_chat(request)
        return response
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error processing chat: {e}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Failed to process chat: {str(e)}")


