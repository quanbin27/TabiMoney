"""
Utility helpers for working with JSON payloads returned from LLMs.
"""

from __future__ import annotations

import json
import logging
import re
from typing import Any, Dict, Iterable, Optional

logger = logging.getLogger(__name__)


CODE_FENCE_PATTERN = re.compile(r"```(?:json)?\s*([\s\S]*?)```", re.IGNORECASE)
JSON_BLOCK_PATTERN = re.compile(r"\{[\s\S]*\}")


def extract_json_block(content: str, *, default: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
    """
    Try to extract a JSON object from the given string. Handles responses that
    contain Markdown code fences or other surrounding text.

    Args:
        content: Raw text returned from an LLM.
        default: Value to return when parsing fails. Defaults to empty dict.

    Returns:
        Parsed JSON dictionary (best-effort).
    """
    if not content:
        return default or {}

    # 1. Try plain JSON parse first.
    try:
        return json.loads(content)
    except (json.JSONDecodeError, TypeError):
        pass

    # 2. Look for Markdown code fences.
    fence_match = CODE_FENCE_PATTERN.search(content)
    if fence_match:
        fenced = fence_match.group(1).strip()
        try:
            return json.loads(fenced)
        except (json.JSONDecodeError, TypeError):
            logger.debug("Failed to parse fenced JSON block")

    # 3. Find the first {...} block as a last resort.
    block_match = JSON_BLOCK_PATTERN.search(content)
    if block_match:
        candidate = block_match.group(0)
        try:
            return json.loads(candidate)
        except (json.JSONDecodeError, TypeError):
            logger.debug("Failed to parse JSON block extracted via regex")

    logger.warning("Unable to extract JSON block from content: %s", content[:200] + ("..." if len(content) > 200 else ""))
    return default or {}


def ensure_string_list(values: Iterable[Any]) -> list[str]:
    """Best-effort conversion of an iterable into a list of strings."""
    if values is None:
        return []
    return [str(item) for item in values]



