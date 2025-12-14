"""
NLU models
"""

from datetime import datetime
from typing import List, Optional
from pydantic import BaseModel, field_validator


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
    
    @field_validator('message')
    @classmethod
    def validate_message(cls, v):
        if not v or not v.strip():
            raise ValueError('Message cannot be empty')
        if len(v) > 2000:
            raise ValueError('Message is too long (max 2000 characters)')
        return v.strip()
    
    @field_validator('user_id')
    @classmethod
    def validate_user_id(cls, v):
        if v <= 0:
            raise ValueError('user_id must be positive')
        return v


class ChatResponse(BaseModel):
    user_id: int
    message: str
    response: str
    intent: str
    entities: List[Entity]
    suggestions: List[str]
    needs_confirmation: bool = False
    generated_at: datetime
    
    class Config:
        json_encoders = {
            datetime: lambda v: v.isoformat() + "Z" if v else None
        }
