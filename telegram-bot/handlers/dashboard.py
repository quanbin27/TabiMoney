"""
Dashboard command handler for Telegram bot
"""

from telegram import Update
from telegram.ext import ContextTypes
from services.auth_service import TelegramAuthService
from services.api_service import APIService
from utils.logger import setup_logger

logger = setup_logger(__name__)

async def dashboard_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /dashboard command"""
    try:
        user = update.effective_user
        auth_service = TelegramAuthService()
        api_service = APIService()
        
        # Check if user is linked
        if not await auth_service.is_telegram_linked(user.id):
            await update.message.reply_text(
                "❌ Bạn chưa liên kết tài khoản với hệ thống.\n"
                "Vui lòng sử dụng lệnh /link để liên kết tài khoản của bạn."
            )
            return
        
        # Get JWT token
        jwt_token = await auth_service.get_telegram_jwt_token(user.id)
        if not jwt_token:
            await update.message.reply_text(
                "❌ Không thể xác thực tài khoản. Vui lòng liên kết lại bằng lệnh /link."
            )
            return
        
        # Show loading message
        loading_message = await update.message.reply_text("🔄 Đang tải dữ liệu dashboard...")
        
        # Get dashboard data
        dashboard_data = await api_service.get_dashboard_data(jwt_token)
        user_profile = await api_service.get_user_profile(jwt_token)
        recent_transactions = await api_service.get_transactions(jwt_token, 5)
        monthly_income = await api_service.get_monthly_income(jwt_token)
        
        # Delete loading message
        await loading_message.delete()
        
        # Format dashboard message
        dashboard_message = await format_dashboard_message(
            user_profile or {}, dashboard_data or {}, recent_transactions or [], monthly_income
        )
        
        await update.message.reply_text(dashboard_message, parse_mode='HTML')
        logger.info(f"Dashboard displayed for user {user.id}")
        
    except Exception as e:
        logger.error(f"Error in dashboard command: {e}")
        await update.message.reply_text("❌ Có lỗi xảy ra khi tải dashboard. Vui lòng thử lại sau.")

async def format_dashboard_message(user_profile, dashboard_data, recent_transactions, monthly_income):
    """Format dashboard message"""
    try:
        # Debug logging
        logger.info(f"Dashboard data received:")
        logger.info(f"  user_profile: {user_profile}")
        logger.info(f"  dashboard_data: {dashboard_data}")
        logger.info(f"  recent_transactions: {recent_transactions}")
        logger.info(f"  monthly_income: {monthly_income}")
        
        message = "📊 <b>DASHBOARD TÀI CHÍNH</b>\n\n"
        
        # User info (support different casing from BE)
        if user_profile:
            first_name = user_profile.get('first_name') or user_profile.get('FirstName') or ''
            last_name = user_profile.get('last_name') or user_profile.get('LastName') or ''
            email = user_profile.get('email') or user_profile.get('Email') or ''
            full_name = (f"{first_name} {last_name}").strip()
            if full_name:
                message += f"👤 <b>Người dùng:</b> {full_name}\n"
            if email:
                message += f"📧 <b>Email:</b> {email}\n\n"
        
        # Monthly income
        if monthly_income is not None:
            message += f"💰 <b>Thu nhập hàng tháng:</b> {monthly_income:,.0f} VND\n\n"
        
        # Dashboard summary
        if dashboard_data:
            message += "📈 <b>TỔNG QUAN THÁNG NÀY</b>\n"
            
            # Total income
            total_income = dashboard_data.get('total_income') or dashboard_data.get('totalIncome') or 0
            message += f"💵 <b>Tổng thu:</b> {total_income:,.0f} VND\n"
            
            # Total expenses
            total_expenses = dashboard_data.get('total_expense') or dashboard_data.get('total_expenses') or dashboard_data.get('totalExpense') or 0
            message += f"💸 <b>Tổng chi:</b> {total_expenses:,.0f} VND\n"
            
            # Balance
            balance = dashboard_data.get('net_amount') or dashboard_data.get('netAmount')
            if balance is None:
                balance = total_income - total_expenses
            balance_emoji = "📈" if balance >= 0 else "📉"
            message += f"{balance_emoji} <b>Số dư:</b> {balance:,.0f} VND\n\n"
            
            # Top categories
            top_categories = dashboard_data.get('top_categories') or dashboard_data.get('topCategories') or []
            if top_categories and isinstance(top_categories, list):
                message += "🏆 <b>DANH MỤC CHI TIÊU NHIỀU NHẤT</b>\n"
                for i, category in enumerate(top_categories[:3], 1):
                    name = category.get('name') or category.get('Name') or 'Unknown'
                    amount = category.get('amount') or category.get('Amount') or 0
                    message += f"{i}. {name}: {amount:,.0f} VND\n"
                message += "\n"
        
        # Recent transactions
        if recent_transactions and isinstance(recent_transactions, list):
            message += "📝 <b>GIAO DỊCH GẦN ĐÂY</b>\n"
            for transaction in recent_transactions[:5]:
                amount = transaction.get('amount') or transaction.get('Amount') or 0
                description = transaction.get('description') or transaction.get('Description') or 'Không có mô tả'
                catObj = transaction.get('category') or transaction.get('Category') or {}
                category = (catObj.get('name') or catObj.get('Name') or 'Không phân loại') if isinstance(catObj, dict) else str(catObj)
                transaction_type = transaction.get('type') or transaction.get('Type') or 'expense'
                
                emoji = "💸" if transaction_type == 'expense' else "💵"
                amount_str = f"-{amount:,.0f}" if transaction_type == 'expense' else f"+{amount:,.0f}"
                
                message += f"{emoji} {description}\n"
                message += f"   💰 {amount_str} VND | 📂 {category}\n\n"
        
        # Quick actions
        message += "⚡ <b>THAO TÁC NHANH</b>\n"
        message += "• Gửi tin nhắn để chat với AI\n"
        message += "• /dashboard - Làm mới dashboard\n"
        message += "• /link - Quản lý liên kết tài khoản\n"
        
        return message
        
    except Exception as e:
        logger.error(f"Error formatting dashboard message: {e}", exc_info=True)
        return f"❌ Có lỗi khi định dạng dữ liệu dashboard: {str(e)}"
