from fastapi import APIRouter

router = APIRouter()

@router.post("/detect")
async def detect_anomalies():
    return {"anomalies": [], "total_anomalies": 0, "detection_score": 0.0}


