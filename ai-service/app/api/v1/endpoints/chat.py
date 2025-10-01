from fastapi import APIRouter

router = APIRouter()

@router.post("/process")
async def process_chat():
    return {"response": "Hello from AI service", "intent": "general", "entities": [], "suggestions": []}


