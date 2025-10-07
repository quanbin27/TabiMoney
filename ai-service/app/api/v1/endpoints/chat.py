from fastapi import APIRouter, Depends
from app.core.dependencies import get_nlu_service
from app.models.nlu import ChatRequest, ChatResponse
from app.services.nlu_service import NLUService

router = APIRouter()

@router.post("/process", response_model=ChatResponse)
async def process_chat(
    request: ChatRequest,
    nlu_service: NLUService = Depends(get_nlu_service)
):
    """Process chat message and return AI response"""
    return await nlu_service.process_chat(request)


