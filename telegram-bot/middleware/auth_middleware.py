"""
Authentication middleware for Telegram bot
"""

from typing import Optional
from services.auth_service import TelegramAuthService
from utils.logger import setup_logger

logger = setup_logger(__name__)

class AuthMiddleware:
    """Middleware for handling authentication"""
    
    def __init__(self, auth_service: TelegramAuthService):
        self.auth_service = auth_service
    
    async def is_authenticated(self, telegram_user_id: int) -> bool:
        """Check if Telegram user is authenticated"""
        try:
            return await self.auth_service.is_telegram_linked(telegram_user_id)
        except Exception as e:
            logger.error(f"Error checking authentication: {e}")
            return False
    
    async def get_user_token(self, telegram_user_id: int) -> Optional[str]:
        """Get JWT token for authenticated user"""
        try:
            if await self.is_authenticated(telegram_user_id):
                return await self.auth_service.get_telegram_jwt_token(telegram_user_id)
            return None
        except Exception as e:
            logger.error(f"Error getting user token: {e}")
            return None
    
    async def get_web_user_id(self, telegram_user_id: int) -> Optional[int]:
        """Get web user ID for Telegram user"""
        try:
            if await self.is_authenticated(telegram_user_id):
                return await self.auth_service.get_web_user_id(telegram_user_id)
            return None
        except Exception as e:
            logger.error(f"Error getting web user ID: {e}")
            return None
