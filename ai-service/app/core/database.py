"""
Minimal MySQL access layer used by ML service.
Uses PyMySQL directly to run read-only queries for training data.
"""

import asyncio
import logging
import pymysql
from urllib.parse import urlparse, unquote
from contextlib import asynccontextmanager
from typing import Any, List, Dict, Tuple

from app.core.config import settings

logger = logging.getLogger(__name__)

_db_conf: Dict[str, Any] = {}


def _parse_mysql_url(url: str) -> Dict[str, Any]:
    parsed = urlparse(url)
    return {
        "host": parsed.hostname or "localhost",
        "port": parsed.port or 3306,
        "user": unquote(parsed.username or "root"),
        "password": unquote(parsed.password or ""),
        "db": (parsed.path or "/tabimoney").lstrip("/"),
        "charset": "utf8mb4",
        "cursorclass": pymysql.cursors.DictCursor,
        "autocommit": True,
    }


async def init_db():
    """Parse config and verify connection once."""
    global _db_conf
    _db_conf = _parse_mysql_url(settings.DATABASE_URL)
    # Verify connectivity in a thread to avoid blocking loop
    def _ping():
        conn = pymysql.connect(**_db_conf)
        try:
            with conn.cursor() as cur:
                cur.execute("SELECT 1")
        finally:
            conn.close()
    await asyncio.to_thread(_ping)
    logger.info("DB connection verified (%s:%s/%s)", _db_conf["host"], _db_conf["port"], _db_conf["db"])


class _DBSession:
    def __init__(self, conf: Dict[str, Any]):
        self._conf = conf
        self._conn = None

    async def __aenter__(self):
        self._conn = await asyncio.to_thread(pymysql.connect, **self._conf)
        return self

    async def __aexit__(self, exc_type, exc, tb):
        if self._conn:
            await asyncio.to_thread(self._conn.close)

    async def execute(self, query: str, params: Tuple[Any, ...] = ()) -> List[Dict[str, Any]]:
        def _exec():
            with self._conn.cursor() as cur:
                cur.execute(query, params)
                rows = cur.fetchall()
                return list(rows)
        return await asyncio.to_thread(_exec)


@asynccontextmanager
async def get_db():
    """Async context manager yielding a DB session with async execute()."""
    session = _DBSession(_db_conf)
    await session.__aenter__()
    try:
        yield session
    finally:
        await session.__aexit__(None, None, None)


