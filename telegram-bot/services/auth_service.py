"""
Authentication service for Telegram bot
Handles link codes and JWT token management
"""

import asyncio
import secrets
import time
from datetime import datetime, timedelta
from typing import Optional, Dict, Any
import jwt
import pymysql
from pymysql.cursors import DictCursor

from utils.config import Config
from utils.logger import setup_logger

logger = setup_logger(__name__)

class TelegramAuthService:
    """Service for handling Telegram authentication"""
    
    def __init__(self):
        self.config = Config()
        self.db_connection = None
        
    def _get_db_connection(self):
        """Get database connection"""
        if not self.db_connection or not self.db_connection.open:
            self.db_connection = pymysql.connect(
                host=self.config.DB_HOST,
                port=self.config.DB_PORT,
                user=self.config.DB_USER,
                password=self.config.DB_PASSWORD,
                database=self.config.DB_NAME,
                cursorclass=DictCursor,
                charset='utf8mb4'
            )
        return self.db_connection
    
    async def generate_link_code(self, user_id: int) -> str:
        """Generate a link code for Telegram user"""
        try:
            # Generate random code
            code = secrets.token_hex(self.config.LINK_CODE_LENGTH // 2).upper()
            
            # Store in database with expiration
            db = self._get_db_connection()
            with db.cursor() as cursor:
                cursor.execute(
                    "INSERT INTO telegram_link_codes (code, web_user_id, expires_at, created_at) VALUES (%s, %s, %s, NOW())",
                    (code, user_id, datetime.now() + timedelta(minutes=self.config.LINK_CODE_EXPIRE_MINUTES))
                )
                db.commit()
            
            logger.info(f"Generated link code for Telegram user {user_id}: {code}")
            return code
            
        except Exception as e:
            logger.error(f"Error generating link code: {e}")
            raise
    
    async def validate_link_code(self, code: str) -> Optional[int]:
        """Validate link code and return Telegram user ID"""
        try:
            db = self._get_db_connection()
            
            with db.cursor() as cursor:
                # Get unused, non-expired code
                cursor.execute(
                    "SELECT web_user_id FROM telegram_link_codes WHERE code = %s AND expires_at > NOW() AND used_at IS NULL",
                    (code,)
                )
                result = cursor.fetchone()
                
                if result:
                    web_user_id = result['web_user_id']
                    
                    # Mark as used
                    cursor.execute(
                        "UPDATE telegram_link_codes SET used_at = NOW() WHERE code = %s",
                        (code,)
                    )
                    db.commit()
                    
                    return web_user_id
                
                return None
                
        except Exception as e:
            logger.error(f"Error validating link code: {e}")
            return None
    
    async def link_telegram_account(self, telegram_user_id: int, web_user_id: int) -> bool:
        """Link Telegram account with web user account"""
        try:
            db = self._get_db_connection()
            
            with db.cursor() as cursor:
                # Check if already linked
                cursor.execute(
                    "SELECT id FROM telegram_accounts WHERE telegram_user_id = %s",
                    (telegram_user_id,)
                )
                existing = cursor.fetchone()
                
                if existing:
                    # Update existing link
                    cursor.execute(
                        "UPDATE telegram_accounts SET web_user_id = %s, updated_at = NOW() WHERE telegram_user_id = %s",
                        (web_user_id, telegram_user_id)
                    )
                else:
                    # Create new link
                    cursor.execute(
                        "INSERT INTO telegram_accounts (telegram_user_id, web_user_id, created_at, updated_at) VALUES (%s, %s, NOW(), NOW())",
                        (telegram_user_id, web_user_id)
                    )
                
                db.commit()
                
                # Generate permanent JWT token for Telegram user
                await self.generate_telegram_jwt_token(telegram_user_id, web_user_id)
                
                logger.info(f"Linked Telegram user {telegram_user_id} with web user {web_user_id}")
                return True
                
        except Exception as e:
            logger.error(f"Error linking Telegram account: {e}")
            return False
    
    async def generate_telegram_jwt_token(self, telegram_user_id: int, web_user_id: int) -> str:
        """Generate permanent JWT token for Telegram user"""
        try:
            # Create claims with 1 year expiration for Telegram users
            claims = {
                "user_id": web_user_id,
                "telegram_user_id": telegram_user_id,
                "type": "telegram_access",
                "iat": int(time.time()),
                "exp": int(time.time()) + (365 * 24 * 60 * 60)  # 1 year expiration
            }
            
            # Generate token
            token = jwt.encode(claims, self.config.JWT_SECRET, algorithm="HS256")
            
            logger.info(f"Generated JWT token for Telegram user {telegram_user_id}")
            return token
            
        except Exception as e:
            logger.error(f"Error generating JWT token: {e}")
            raise
    
    async def get_telegram_jwt_token(self, telegram_user_id: int) -> Optional[str]:
        """Get JWT token for Telegram user"""
        try:
            db = self._get_db_connection()
            
            with db.cursor() as cursor:
                cursor.execute(
                    "SELECT web_user_id FROM telegram_accounts WHERE telegram_user_id = %s",
                    (telegram_user_id,)
                )
                result = cursor.fetchone()
                
                if result:
                    web_user_id = result['web_user_id']
                    return await self.generate_telegram_jwt_token(telegram_user_id, web_user_id)
            
            return None
            
        except Exception as e:
            logger.error(f"Error getting JWT token: {e}")
            return None
    
    async def is_telegram_linked(self, telegram_user_id: int) -> bool:
        """Check if Telegram user is linked to web account"""
        try:
            db = self._get_db_connection()
            
            with db.cursor() as cursor:
                cursor.execute(
                    "SELECT id FROM telegram_accounts WHERE telegram_user_id = %s",
                    (telegram_user_id,)
                )
                result = cursor.fetchone()
                
                return result is not None
                
        except Exception as e:
            logger.error(f"Error checking Telegram link status: {e}")
            return False
    
    async def get_web_user_id(self, telegram_user_id: int) -> Optional[int]:
        """Get web user ID for Telegram user"""
        try:
            db = self._get_db_connection()
            
            with db.cursor() as cursor:
                cursor.execute(
                    "SELECT web_user_id FROM telegram_accounts WHERE telegram_user_id = %s",
                    (telegram_user_id,)
                )
                result = cursor.fetchone()
                
                if result:
                    return result['web_user_id']
                
                return None
                
        except Exception as e:
            logger.error(f"Error getting web user ID: {e}")
            return None
    
    async def unlink_telegram_account(self, telegram_user_id: int) -> bool:
        """Unlink Telegram account"""
        try:
            db = self._get_db_connection()
            
            with db.cursor() as cursor:
                cursor.execute(
                    "DELETE FROM telegram_accounts WHERE telegram_user_id = %s",
                    (telegram_user_id,)
                )
                db.commit()
            
            logger.info(f"Unlinked Telegram user {telegram_user_id}")
            return True
            
        except Exception as e:
            logger.error(f"Error unlinking Telegram account: {e}")
            return False
