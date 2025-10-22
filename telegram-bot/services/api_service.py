"""
API service for communicating with backend and AI service
"""

import aiohttp
import asyncio
from typing import Optional, Dict, Any, List
import json

from utils.config import Config
from utils.logger import setup_logger

logger = setup_logger(__name__)

class APIService:
    """Service for API communication"""
    
    def __init__(self):
        self.config = Config()
        self.backend_url = self.config.BACKEND_URL
        self.backend_api_base = f"{self.backend_url}/api/v1"
        self.ai_service_url = self.config.AI_SERVICE_URL
        
    async def _make_request(self, method: str, url: str, headers: Optional[Dict] = None, 
                          data: Optional[Dict] = None, params: Optional[Dict] = None) -> Optional[Dict]:
        """Make HTTP request"""
        try:
            async with aiohttp.ClientSession() as session:
                async with session.request(
                    method=method,
                    url=url,
                    headers=headers,
                    json=data,
                    params=params,
                    timeout=aiohttp.ClientTimeout(total=30)
                ) as response:
                    text = await response.text()
                    ct = response.headers.get('Content-Type', '')
                    if response.status == 200:
                        try:
                            if 'application/json' in ct:
                                return await response.json()
                            return json.loads(text)
                        except Exception:
                            logger.error(f"API JSON parse error for {url}: {text}")
                            return None
                    else:
                        logger.error(f"API request failed {method} {url}: {response.status} - {text}")
                        return None
                        
        except Exception as e:
            logger.error(f"Error making API request: {e}")
            return None
    
    async def get_user_profile(self, jwt_token: str) -> Optional[Dict]:
        """Get user profile from backend"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/auth/profile"
        
        return await self._make_request("GET", url, headers=headers)
    
    async def get_dashboard_data(self, jwt_token: str) -> Optional[Dict]:
        """Get dashboard data from backend"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/analytics/dashboard"
        
        return await self._make_request("GET", url, headers=headers)
    
    async def get_transactions(self, jwt_token: str, limit: int = 10) -> Optional[List[Dict]]:
        """Get recent transactions"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/transactions"
        params = {"limit": limit}
        
        return await self._make_request("GET", url, headers=headers, params=params)
    
    async def get_categories(self, jwt_token: str) -> Optional[List[Dict]]:
        """Get user categories"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/categories"
        
        return await self._make_request("GET", url, headers=headers)
    
    async def get_budgets(self, jwt_token: str) -> Optional[List[Dict]]:
        """Get user budgets"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/budgets"
        
        return await self._make_request("GET", url, headers=headers)
    
    async def get_goals(self, jwt_token: str) -> Optional[List[Dict]]:
        """Get user financial goals"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/goals"
        
        return await self._make_request("GET", url, headers=headers)
    
    async def send_chat_message(self, jwt_token: str, message: str, user_id: int) -> Optional[Dict]:
        """Send chat message to AI service"""
        headers = {
            "Authorization": f"Bearer {jwt_token}",
            "Content-Type": "application/json"
        }
        url = f"{self.ai_service_url}/api/v1/chat/process"
        
        data = {
            "message": message,
            "user_id": user_id
        }
        
        return await self._make_request("POST", url, headers=headers, data=data)
    
    async def add_transaction(self, jwt_token: str, transaction_data: Dict) -> Optional[Dict]:
        """Add new transaction"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/transactions"
        
        return await self._make_request("POST", url, headers=headers, data=transaction_data)
    
    async def get_monthly_income(self, jwt_token: str) -> Optional[float]:
        """Get user's monthly income"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/auth/income"
        
        response = await self._make_request("GET", url, headers=headers)
        if response:
            return response.get("monthly_income", 0)
        return None
    
    async def set_monthly_income(self, jwt_token: str, amount: float) -> bool:
        """Set user's monthly income"""
        headers = {"Authorization": f"Bearer {jwt_token}"}
        url = f"{self.backend_api_base}/auth/income"
        
        data = {"monthly_income": amount}
        response = await self._make_request("PUT", url, headers=headers, data=data)
        return response is not None
