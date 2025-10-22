"""
Logging configuration for Telegram bot
"""

import logging
import sys
import os
from typing import Optional

def setup_logger(name: str, level: Optional[str] = None) -> logging.Logger:
    """Setup logger with proper configuration"""
    
    # Get log level from environment or parameter
    log_level = level or os.getenv("LOG_LEVEL", "INFO")
    
    # Create logger
    logger = logging.getLogger(name)
    logger.setLevel(getattr(logging, log_level.upper()))
    
    # Avoid duplicate handlers
    if logger.handlers:
        return logger
    
    # Create formatter
    if os.getenv("LOG_FORMAT", "json") == "json":
        formatter = logging.Formatter(
            '{"timestamp": "%(asctime)s", "level": "%(levelname)s", "logger": "%(name)s", "message": "%(message)s"}'
        )
    else:
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
    
    # Create console handler
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)
    
    return logger
