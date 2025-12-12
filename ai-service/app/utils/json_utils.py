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
        except (json.JSONDecodeError, TypeError) as e:
            logger.debug("Failed to parse JSON block extracted via regex: %s", str(e))
            # Try to fix common issues: incomplete JSON, trailing commas, etc.
            candidate_fixed = candidate.strip()
            
            # Remove trailing commas before closing braces/brackets
            candidate_fixed = re.sub(r',(\s*[}\]])', r'\1', candidate_fixed)
            
            # Try to close incomplete JSON objects
            if candidate_fixed.startswith("{") and not candidate_fixed.endswith("}"):
                # Count braces to see if we can auto-close
                open_braces = candidate_fixed.count("{")
                close_braces = candidate_fixed.count("}")
                missing_closes = open_braces - close_braces
                if missing_closes > 0:
                    # Try to intelligently close the JSON
                    # If the last character is a comma, remove it first
                    if candidate_fixed.rstrip().endswith(","):
                        candidate_fixed = candidate_fixed.rstrip()[:-1]
                    candidate_fixed += "}" * missing_closes
                    try:
                        return json.loads(candidate_fixed)
                    except (json.JSONDecodeError, TypeError) as e2:
                        logger.debug("Failed to parse auto-closed JSON: %s", str(e2))
            
            # Try parsing again after fixing trailing commas
            try:
                return json.loads(candidate_fixed)
            except (json.JSONDecodeError, TypeError):
                pass

    # Log more details for debugging
    content_preview = content[:500] + ("..." if len(content) > 500 else "")
    logger.warning(
        "Unable to extract JSON block from content (length: %d): %s",
        len(content),
        content_preview
    )
    return default or {}


def ensure_string_list(values: Iterable[Any]) -> list[str]:
    """Best-effort conversion of an iterable into a list of strings."""
    if values is None:
        return []
    return [str(item) for item in values]



