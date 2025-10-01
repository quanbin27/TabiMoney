"""
Async Redis client initialization and helpers.
"""

import asyncio
import logging
from typing import Optional

from redis import asyncio as aioredis

from app.core.config import settings

logger = logging.getLogger(__name__)

_redis: Optional[aioredis.Redis] = None


async def init_redis():
    """Initialize Redis connection and verify with PING."""
    global _redis
    try:
        _redis = aioredis.from_url(
            settings.REDIS_URL,
            encoding="utf-8",
            decode_responses=True,
            socket_timeout=5,
            socket_connect_timeout=5,
        )
        pong = await _redis.ping()
        if pong is True:
            logger.info("Redis connected: %s", settings.REDIS_URL)
        else:
            raise RuntimeError("Redis ping returned non-True")
    except Exception as exc:
        logger.error("Redis initialization failed: %s", exc)
        raise


def get_redis() -> aioredis.Redis:
    if _redis is None:
        raise RuntimeError("Redis not initialized")
    return _redis


