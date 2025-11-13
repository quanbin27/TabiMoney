"""
Shared helpers for calling local LLM backends (e.g. Ollama).
"""

from __future__ import annotations

import asyncio
import logging
from typing import Any, Dict, Optional

import httpx

from app.core.config import settings
from app.utils.json_utils import extract_json_block

logger = logging.getLogger(__name__)


async def call_ollama(
    prompt: str,
    *,
    temperature: float = 0.3,
    max_tokens: Optional[int] = None,
    format_json: bool = True,
    timeout: Optional[float] = None,
) -> Dict[str, Any]:
    """
    Call the configured Ollama endpoint and optionally parse the JSON response.

    Args:
        prompt: Prompt sent to the model.
        temperature: Sampling temperature.
        max_tokens: Override the global max token setting.
        format_json: Whether to request JSON output and parse it.
        timeout: Override request timeout (seconds).

    Returns:
        Dictionary with keys:
            - `raw` (str): raw text response from Ollama.
            - `json` (dict): parsed JSON payload when `format_json=True`.
    """
    timeout = timeout or settings.LLM_TIMEOUT
    max_tokens = max_tokens or settings.LLM_MAX_TOKENS

    payload: Dict[str, Any] = {
        "model": settings.LLM_MODEL,
        "prompt": prompt,
        "stream": False,
        "options": {
            "temperature": temperature,
            "max_tokens": max_tokens,
        },
    }
    if format_json:
        payload["format"] = "json"

    logger.debug("Calling Ollama model=%s", settings.LLM_MODEL)
    async with httpx.AsyncClient(timeout=timeout) as client:
        response = await client.post(f"{settings.OLLAMA_BASE_URL}/api/generate", json=payload)
        response.raise_for_status()
        body = response.json()

    raw = body.get("response", "")
    parsed = extract_json_block(raw) if format_json else {}
    return {"raw": raw, "json": parsed}



