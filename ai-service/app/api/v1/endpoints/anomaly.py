from fastapi import APIRouter, Depends
from pydantic import BaseModel

from app.core.dependencies import get_anomaly_service

router = APIRouter()

class DetectRequest(BaseModel):
    user_id: int
    start_date: str
    end_date: str
    threshold: float | None = 0.6


@router.post("/detect")
async def detect_anomalies(req: DetectRequest, svc = Depends(get_anomaly_service)):
    return await svc.detect(req.user_id, req.start_date, req.end_date, req.threshold or 0.6)


