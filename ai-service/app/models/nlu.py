"""
NLU models
"""

from datetime import datetime
from typing import List, Optional
from pydantic import BaseModel


class Entity(BaseModel):
    type: str
    value: str
    confidence: float
    start_pos: int
    end_pos: int


class NLURequest(BaseModel):
    text: str
    user_id: int
    context: Optional[str] = ""


class NLUResponse(BaseModel):
    user_id: int
    intent: str
    entities: List[Entity]
    confidence: float
    suggested_action: str
    response: str
    needs_confirmation: bool = False
    generated_at: datetime


class ChatRequest(BaseModel):
    message: str
    user_id: int


class ChatResponse(BaseModel):
    user_id: int
    message: str
    response: str
    intent: str
    entities: List[Entity]
    suggestions: List[str]
    needs_confirmation: bool = False
    generated_at: datetime
