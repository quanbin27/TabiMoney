"""
Start command handler for Telegram bot
"""

from telegram import Update
from telegram.ext import ContextTypes
from services.auth_service import TelegramAuthService
from utils.logger import setup_logger

logger = setup_logger(__name__)

async def start_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /start command"""
    try:
        user = update.effective_user
        auth_service = TelegramAuthService()

        if await auth_service.is_telegram_linked(user.id):
            # Already linked: show quick menu
            message = f"""
🎉 Chào mừng trở lại, {user.first_name}! 👋

Tài khoản của bạn đã được liên kết với TabiMoney.
Bạn có thể sử dụng các tính năng ngay:

• /dashboard - Xem dashboard tài chính
• Gửi tin nhắn để chat với AI
• /link - Quản lý liên kết (hủy liên kết nếu cần)
            """
        else:
            # Not linked: show instructions
            message = f"""
🎉 Chào mừng bạn đến với TabiMoney Bot!

Xin chào {user.first_name}! 👋

TabiMoney Bot sẽ giúp bạn:
📊 Xem dashboard tài chính
🤖 Chat với AI để phân tích chi tiêu
💰 Quản lý ngân sách và mục tiêu
📈 Theo dõi giao dịch

🔗 Để bắt đầu, bạn cần liên kết tài khoản với hệ thống TabiMoney:
1. Truy cập website TabiMoney: http://localhost:3000
2. Vào phần Settings → Telegram Integration
3. Tạo mã liên kết
4. Gửi mã đó cho bot bằng lệnh /link

📝 Các lệnh có sẵn:
/link - Liên kết tài khoản
/dashboard - Xem dashboard (sau khi liên kết)

Chúc bạn sử dụng TabiMoney Bot hiệu quả! 💪
            """

        await update.message.reply_text(message)
        logger.info(f"Start command executed for user {user.id}")

    except Exception as e:
        logger.error(f"Error in start command: {e}")
        await update.message.reply_text("❌ Có lỗi xảy ra. Vui lòng thử lại sau.")
