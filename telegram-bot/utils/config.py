"""
Configuration management for Telegram bot
"""

import os
from typing import Optional
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

class Config:
    """Configuration class for Telegram bot"""
    
    def __init__(self):
        # Telegram Bot Configuration
        self.TELEGRAM_BOT_TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")
        if not self.TELEGRAM_BOT_TOKEN:
            raise ValueError("TELEGRAM_BOT_TOKEN environment variable is required")
        
        # Backend API Configuration
        self.BACKEND_URL = os.getenv("BACKEND_URL", "http://localhost:8080")
        self.AI_SERVICE_URL = os.getenv("AI_SERVICE_URL", "http://localhost:8001")
        
        # Database Configuration (for link codes)
        self.DB_HOST = os.getenv("DB_HOST", "localhost")
        self.DB_PORT = int(os.getenv("DB_PORT", "3306"))
        self.DB_USER = os.getenv("DB_USER", "root")
        self.DB_PASSWORD = os.getenv("DB_PASSWORD", "password")
        self.DB_NAME = os.getenv("DB_NAME", "tabimoney")
        
        # Redis Configuration
        self.REDIS_HOST = os.getenv("REDIS_HOST", "localhost")
        self.REDIS_PORT = int(os.getenv("REDIS_PORT", "6379"))
        self.REDIS_PASSWORD = os.getenv("REDIS_PASSWORD", "")
        self.REDIS_DB = int(os.getenv("REDIS_DB", "0"))
        
        # JWT Configuration
        self.JWT_SECRET = os.getenv("JWT_SECRET", "your-super-secret-jwt-key-here")
        
        # Link Code Configuration
        self.LINK_CODE_EXPIRE_MINUTES = int(os.getenv("LINK_CODE_EXPIRE_MINUTES", "10"))
        self.LINK_CODE_LENGTH = int(os.getenv("LINK_CODE_LENGTH", "8"))
        
        # Bot Configuration
        self.BOT_NAME = "TabiMoney Bot"
        self.SUPPORTED_LANGUAGES = ["vi", "en"]
        self.DEFAULT_LANGUAGE = "vi"
        
        # Rate Limiting
        self.RATE_LIMIT_MESSAGES = int(os.getenv("RATE_LIMIT_MESSAGES", "30"))
        self.RATE_LIMIT_WINDOW_MINUTES = int(os.getenv("RATE_LIMIT_WINDOW_MINUTES", "1"))
        
        # Logging
        self.LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO")
        self.LOG_FORMAT = os.getenv("LOG_FORMAT", "json")
        
        # Environment
        self.ENV = os.getenv("ENV", "development")
        self.DEBUG = self.ENV == "development"
    
    def get_database_url(self) -> str:
        """Get database connection URL"""
        return f"mysql+pymysql://{self.DB_USER}:{self.DB_PASSWORD}@{self.DB_HOST}:{self.DB_PORT}/{self.DB_NAME}"
    
    def get_redis_url(self) -> str:
        """Get Redis connection URL"""
        if self.REDIS_PASSWORD:
            return f"redis://:{self.REDIS_PASSWORD}@{self.REDIS_HOST}:{self.REDIS_PORT}/{self.REDIS_DB}"
        return f"redis://{self.REDIS_HOST}:{self.REDIS_PORT}/{self.REDIS_DB}"
