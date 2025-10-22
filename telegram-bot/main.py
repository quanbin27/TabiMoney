#!/usr/bin/env python3
"""
TabiMoney Telegram Bot
Main entry point for the Telegram bot integration
"""

import asyncio
import logging
import os
from typing import Optional

from telegram import Update
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes

from handlers.start import start_command
from handlers.link import link_command, handle_link_code
from handlers.dashboard import dashboard_command
from handlers.chat import handle_chat_message
from middleware.auth_middleware import AuthMiddleware
from services.auth_service import TelegramAuthService
from services.api_service import APIService
from utils.config import Config
from utils.logger import setup_logger

# Setup logging
logger = setup_logger(__name__)

class TabiMoneyBot:
    def __init__(self):
        self.config = Config()
        self.auth_service = TelegramAuthService()
        self.api_service = APIService()
        self.auth_middleware = AuthMiddleware(self.auth_service)
        
    async def start(self):
        """Start the bot"""
        logger.info("Starting TabiMoney Telegram Bot...")
        
        # Create application
        application = Application.builder().token(self.config.TELEGRAM_BOT_TOKEN).build()
        
        # Add handlers
        application.add_handler(CommandHandler("start", start_command))
        application.add_handler(CommandHandler("link", link_command))
        application.add_handler(CommandHandler("dashboard", dashboard_command))
        application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, self.handle_message))
        
        # Add middleware
        application.add_error_handler(self.error_handler)
        
        # Start bot
        logger.info("Bot started successfully!")
        await application.initialize()
        await application.start()
        await application.updater.start_polling()
        
        # Keep running
        try:
            await asyncio.Event().wait()
        except KeyboardInterrupt:
            logger.info("Bot stopping...")
        finally:
            await application.updater.stop()
            await application.stop()
            await application.shutdown()
    
    async def handle_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle incoming messages"""
        try:
            message_text = (update.message.text or "").strip()

            # 1) Accept link codes BEFORE auth check
            # Accept both formats: "LINK_xxx" or 8-16 uppercase alphanumeric codes
            is_prefixed_code = message_text.upper().startswith("LINK_")
            is_plain_code = 8 <= len(message_text) <= 16 and message_text.isalnum() and message_text.upper() == message_text
            if is_prefixed_code or is_plain_code:
                await handle_link_code(update, context)
                return

            # 2) Require authentication for other actions
            user_id = update.effective_user.id
            if not await self.auth_middleware.is_authenticated(user_id):
                await update.message.reply_text(
                    "❌ Bạn chưa liên kết tài khoản với hệ thống.\n"
                    "Vui lòng dùng /link và gửi mã liên kết để liên kết tài khoản."
                )
                return

            # 3) Authenticated: treat as chat message
            await handle_chat_message(update, context)

        except Exception as e:
            logger.error(f"Error handling message: {e}")
            await update.message.reply_text("❌ Có lỗi xảy ra. Vui lòng thử lại sau.")
    
    async def error_handler(self, update: Optional[Update], context: ContextTypes.DEFAULT_TYPE):
        """Handle errors"""
        logger.error(f"Update {update} caused error {context.error}")
        if update and update.effective_message:
            await update.effective_message.reply_text(
                "❌ Có lỗi xảy ra. Vui lòng thử lại sau."
            )

def main():
    """Main function"""
    bot = TabiMoneyBot()
    try:
        asyncio.run(bot.start())
    except KeyboardInterrupt:
        logger.info("Bot stopped by user")
    except Exception as e:
        logger.error(f"Bot error: {e}")

if __name__ == "__main__":
    main()