"""
Link command handler for Telegram bot
"""

from telegram import Update
from telegram.ext import ContextTypes
from services.auth_service import TelegramAuthService
from utils.logger import setup_logger

logger = setup_logger(__name__)

async def link_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /link command"""
    try:
        user = update.effective_user
        auth_service = TelegramAuthService()
        
        # Check if already linked
        if await auth_service.is_telegram_linked(user.id):
            await update.message.reply_text(
                "✅ Tài khoản của bạn đã được liên kết với TabiMoney!\n\n"
                "Bạn có thể sử dụng các tính năng:\n"
                "/dashboard - Xem dashboard\n"
                "Gửi tin nhắn để chat với AI"
            )
            return
        
        instructions = f"""
🔗 Liên kết tài khoản TabiMoney

Để liên kết tài khoản Telegram với TabiMoney, hãy làm theo các bước sau:

1️⃣ Truy cập website TabiMoney: http://localhost:3000
2️⃣ Đăng nhập vào tài khoản của bạn
3️⃣ Vào Settings → Telegram Integration
4️⃣ Tạo mã liên kết
5️⃣ Sao chép mã liên kết và gửi cho bot này

💡 Sau khi liên kết thành công, bạn sẽ có thể:
• Xem dashboard tài chính
• Chat với AI để phân tích chi tiêu
• Quản lý ngân sách và mục tiêu

📝 Lưu ý: Mã liên kết có hiệu lực trong 10 phút.
        """
        
        await update.message.reply_text(instructions)
        logger.info(f"Link instructions sent to user {user.id}")
        
    except Exception as e:
        logger.error(f"Error in link command: {e}")
        await update.message.reply_text("❌ Có lỗi xảy ra. Vui lòng thử lại sau.")

async def handle_link_code(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle link code message"""
    try:
        user = update.effective_user
        message_text = update.message.text.strip()
        
        # Extract code from message
        if message_text.startswith("LINK_"):
            code = message_text
        else:
            code = message_text
        
        auth_service = TelegramAuthService()
        
        # Validate link code and get web user ID
        web_user_id = await auth_service.validate_link_code(code)
        
        if not web_user_id:
            await update.message.reply_text(
                "❌ Mã liên kết không hợp lệ hoặc đã hết hạn.\n\n"
                "Vui lòng:\n"
                "1. Kiểm tra lại mã liên kết\n"
                "2. Đảm bảo bạn đã nhập mã trên website\n"
                "3. Sử dụng lệnh /link để tạo mã mới"
            )
            return
        
        # Link Telegram account with web user
        success = await auth_service.link_telegram_account(user.id, web_user_id)
        
        if success:
            await update.message.reply_text(
                "✅ Liên kết tài khoản thành công!\n\n"
                "Bây giờ bạn có thể sử dụng các tính năng:\n"
                "/dashboard - Xem dashboard tài chính\n"
                "Gửi tin nhắn để chat với AI\n\n"
                "Chúc bạn sử dụng TabiMoney Bot hiệu quả! 🎉"
            )
            
            logger.info(f"Account linked successfully for user {user.id}")
        else:
            await update.message.reply_text(
                "❌ Có lỗi xảy ra khi liên kết tài khoản. Vui lòng thử lại sau."
            )
        
    except Exception as e:
        logger.error(f"Error handling link code: {e}")
        await update.message.reply_text("❌ Có lỗi xảy ra khi xử lý mã liên kết. Vui lòng thử lại sau.")
