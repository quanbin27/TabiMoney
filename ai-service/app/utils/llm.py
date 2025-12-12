"""
Shared helpers for calling Google Gemini API.
"""

from __future__ import annotations

import logging
from typing import Any, Dict, Optional

import aiohttp

from app.core.config import settings
from app.utils.json_utils import extract_json_block

logger = logging.getLogger(__name__)


async def call_gemini(
    prompt: str,
    *,
    temperature: float = 0.3,
    max_tokens: Optional[int] = None,
    format_json: bool = True,
    timeout: Optional[float] = None,
) -> Dict[str, Any]:
    """
    Call the Google Gemini API and optionally parse the JSON response.

    Args:
        prompt: Prompt sent to the model.
        temperature: Sampling temperature.
        max_tokens: Override the global max token setting.
        format_json: Whether to request JSON output and parse it.
        timeout: Override request timeout (seconds).

    Returns:
        Dictionary with keys:
            - `raw` (str): raw text response from Gemini.
            - `json` (dict): parsed JSON payload when `format_json=True`.
    """
    if not settings.USE_GEMINI or not settings.GEMINI_API_KEY:
        raise ValueError("Gemini API is not configured. Please set USE_GEMINI=true and GEMINI_API_KEY.")
    
    timeout = timeout or 60.0
    max_tokens = max_tokens or settings.GEMINI_MAX_TOKENS

    url = f"https://generativelanguage.googleapis.com/v1beta/models/{settings.GEMINI_MODEL}:generateContent?key={settings.GEMINI_API_KEY}"
    
    payload: Dict[str, Any] = {
        "contents": [
            {"parts": [{"text": prompt}]}
        ],
        "generationConfig": {
            "temperature": temperature,
            "maxOutputTokens": max_tokens,
        }
    }
    
    if format_json:
        payload["generationConfig"]["response_mime_type"] = "application/json"

    logger.debug("Calling Gemini model=%s", settings.GEMINI_MODEL)
    
    timeout_obj = aiohttp.ClientTimeout(total=timeout)
    async with aiohttp.ClientSession(timeout=timeout_obj) as session:
        async with session.post(url, json=payload) as resp:
            resp.raise_for_status()
            data = await resp.json()

    # Extract text from Gemini response
    content = ""
    try:
        if "candidates" in data and data["candidates"]:
            candidate = data["candidates"][0]
            # Check for finish reason (might indicate incomplete response)
            finish_reason = candidate.get("finishReason", "")
            if finish_reason == "MAX_TOKENS":
                logger.warning("Gemini response was truncated (MAX_TOKENS) - consider increasing max_tokens")
            elif finish_reason not in ("STOP", ""):
                logger.warning("Gemini finish reason: %s", finish_reason)
            
            if "content" in candidate and "parts" in candidate["content"]:
                content = candidate["content"]["parts"][0].get("text", "")
            elif "text" in candidate:
                content = candidate["text"]
        elif "text" in data:
            content = data["text"]
        else:
            content = str(data)
            logger.warning("Unexpected Gemini response structure: %s", list(data.keys())[:10])
    except Exception as e:
        logger.error(f"Failed to extract Gemini content: {e}, response: {str(data)[:500]}")
        content = str(data)

    if not content:
        logger.error("Empty content from Gemini API")
        return {"raw": "", "json": {}}

    parsed = extract_json_block(content) if format_json else {}
    
    # Log if JSON parsing failed but format_json was requested
    if format_json and not parsed:
        logger.warning(
            "Failed to parse JSON from Gemini response. Raw content (first 500 chars): %s",
            content[:500]
        )
    
    return {"raw": content, "json": parsed}
