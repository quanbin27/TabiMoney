"""
Chat handler for Telegram bot
"""

from telegram import Update
from telegram.ext import ContextTypes
from services.auth_service import TelegramAuthService
from services.api_service import APIService
from utils.logger import setup_logger

logger = setup_logger(__name__)

async def handle_chat_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle chat messages for AI interaction"""
    try:
        user = update.effective_user
        message_text = update.message.text.strip()
        
        auth_service = TelegramAuthService()
        api_service = APIService()
        
        # Check if user is linked
        if not await auth_service.is_telegram_linked(user.id):
            await update.message.reply_text(
                "❌ Bạn chưa liên kết tài khoản với hệ thống.\n"
                "Vui lòng sử dụng lệnh /link để liên kết tài khoản của bạn."
            )
            return
        
        # Get JWT token and web user ID
        jwt_token = await auth_service.get_telegram_jwt_token(user.id)
        web_user_id = await auth_service.get_web_user_id(user.id)
        
        if not jwt_token or not web_user_id:
            await update.message.reply_text(
                "❌ Không thể xác thực tài khoản. Vui lòng liên kết lại bằng lệnh /link."
            )
            return
        
        # Show typing indicator
        try:
            await context.bot.send_chat_action(chat_id=update.effective_chat.id, action="typing")
        except Exception:
            pass
        
        # Send message to AI service
        ai_response = await api_service.send_chat_message(jwt_token, message_text, web_user_id)
        
        # Debug logging
        logger.info(f"AI response received: {ai_response}")
        
        if ai_response and isinstance(ai_response, dict):
            # AI service đã xử lý hoàn toàn, chỉ cần hiển thị response text
            response_text = ai_response.get('response', 'Xin lỗi, tôi không thể xử lý yêu cầu của bạn.')

            # Format response for Telegram (chỉ format hiển thị, không modify content)
            formatted_response = await format_ai_response(response_text, ai_response)

            await update.message.reply_text(formatted_response, parse_mode='HTML')
            logger.info(f"AI response sent to user {user.id}")
        else:
            await update.message.reply_text(
                "❌ Không thể kết nối với AI service. Vui lòng thử lại sau."
            )
        
    except Exception as e:
        logger.error(f"Error handling chat message: {e}")
        await update.message.reply_text("❌ Có lỗi xảy ra khi xử lý tin nhắn. Vui lòng thử lại sau.")

async def format_ai_response(response_text: str, ai_response: dict) -> str:
    """Format AI response for Telegram"""
    try:
        # Basic formatting - response_text đã được AI tạo tự nhiên, chỉ cần format cho Telegram
        formatted = response_text
        
        # Map intent names to Vietnamese labels (cập nhật theo 8 intent mới)
        intent_labels = {
            'add_transaction': 'Thêm giao dịch',
            'query_balance': 'Truy vấn số dư',
            'analyze_data': 'Phân tích dữ liệu',
            'budget_management': 'Quản lý ngân sách',
            'goal_tracking': 'Theo dõi mục tiêu',
            'smart_recommendations': 'Gợi ý thông minh',
            'expense_forecasting': 'Dự đoán chi tiêu',
            'general': 'Câu hỏi chung',
            'error': 'Lỗi'
        }
        
        # Add intent badge if available (optional, không bắt buộc)
        intent = ai_response.get('intent')
        if intent and intent != 'general' and intent != 'error':
            intent_label = intent_labels.get(intent, intent)
            formatted = f"🤖 <b>AI Assistant</b> • {intent_label}\n\n{formatted}"
        
        # Add suggestions if available (AI đã tạo suggestions tự nhiên)
        suggestions = ai_response.get('suggestions', [])
        if suggestions and len(suggestions) > 0:
            formatted += "\n\n💡 <b>Gợi ý tiếp theo:</b>\n"
            for suggestion in suggestions[:3]:  # Limit to 3 suggestions
                formatted += f"• {suggestion}\n"
        
        # Add quick actions
        formatted += "\n\n⚡ <b>Thao tác nhanh:</b>\n"
        formatted += "• /dashboard - Xem dashboard\n"
        formatted += "• Gửi tin nhắn khác để tiếp tục chat\n"
        
        return formatted
        
    except Exception as e:
        logger.error(f"Error formatting AI response: {e}")
        return response_text
