"""
Transaction service for AI agent to interact directly with database
"""

import logging
from datetime import datetime
from typing import Optional, Dict, Any
from app.core.database import get_db

logger = logging.getLogger(__name__)


class TransactionService:
    """Service for AI agent to handle transactions directly"""
    
    async def create_transaction(
        self, 
        user_id: int, 
        category_id: int, 
        amount: float, 
        description: str,
        transaction_type: str = "expense"
    ) -> Dict[str, Any]:
        """Create a new transaction directly in database"""
        try:
            async with get_db() as db:
                # Get category name
                category_result = await db.execute(
                    "SELECT name FROM categories WHERE id = %s AND (user_id = %s OR is_system = true)",
                    (category_id, user_id)
                )
                
                if not category_result:
                    raise ValueError(f"Category {category_id} not found")
                
                category_name = category_result[0]["name"]
                
                # Insert transaction
                insert_query = """
                INSERT INTO transactions 
                (user_id, category_id, amount, description, transaction_type, transaction_date, created_at, updated_at)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                """
                
                now = datetime.now()
                today = now.strftime("%Y-%m-%d")
                
                await db.execute(
                    insert_query,
                    (user_id, category_id, amount, description, transaction_type, today, now, now)
                )
                
                # Get the inserted transaction
                transaction_result = await db.execute(
                    "SELECT * FROM transactions WHERE user_id = %s ORDER BY id DESC LIMIT 1",
                    (user_id,)
                )
                
                if transaction_result:
                    transaction = transaction_result[0]
                    logger.info(f"Created transaction {transaction['id']} for user {user_id}")
                    
                    return {
                        "success": True,
                        "transaction_id": transaction["id"],
                        "amount": amount,
                        "category_name": category_name,
                        "transaction_type": transaction_type,
                        "message": f"Đã thêm giao dịch {amount:,.0f} VND cho danh mục {category_name}."
                    }
                else:
                    raise ValueError("Failed to retrieve created transaction")
                    
        except Exception as e:
            logger.error(f"Failed to create transaction: {e}")
            return {
                "success": False,
                "error": str(e),
                "message": "Có lỗi xảy ra khi thêm giao dịch."
            }
    
    async def get_user_balance(self, user_id: int) -> Dict[str, Any]:
        """Get user's current month balance"""
        try:
            async with get_db() as db:
                # Get current month totals
                now = datetime.now()
                start_date = now.replace(day=1, hour=0, minute=0, second=0, microsecond=0)
                end_date = start_date.replace(month=start_date.month + 1) if start_date.month < 12 else start_date.replace(year=start_date.year + 1, month=1)
                
                # Format dates as strings for MySQL
                start_date_str = start_date.strftime("%Y-%m-%d")
                end_date_str = end_date.strftime("%Y-%m-%d")
                
                balance_query = """
                SELECT 
                    SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) as total_income,
                    SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) as total_expense
                FROM transactions 
                WHERE user_id = %s AND transaction_date BETWEEN %s AND %s
                """
                
                result = await db.execute(balance_query, (user_id, start_date_str, end_date_str))
                
                if result:
                    row = result[0]
                    total_income = float(row["total_income"] or 0)
                    total_expense = float(row["total_expense"] or 0)
                    net_amount = total_income - total_expense
                    
                    return {
                        "success": True,
                        "total_income": total_income,
                        "total_expense": total_expense,
                        "net_amount": net_amount,
                        "message": f"Tháng {now.month:02d}/{now.year}: Tổng chi tiêu {total_expense:,.0f} VND, tổng thu {total_income:,.0f} VND, chênh lệch {net_amount:,.0f} VND."
                    }
                else:
                    return {
                        "success": True,
                        "total_income": 0,
                        "total_expense": 0,
                        "net_amount": 0,
                        "message": f"Tháng {now.month:02d}/{now.year}: Chưa có giao dịch nào."
                    }
                    
        except Exception as e:
            logger.error(f"Failed to get user balance: {e}")
            return {
                "success": False,
                "error": str(e),
                "message": "Có lỗi xảy ra khi lấy thông tin số dư."
            }
    
