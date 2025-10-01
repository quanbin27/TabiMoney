#!/usr/bin/env python3
"""
Script to generate mock transaction data for testing AI service
"""
import requests
import json
import random
from datetime import datetime, timedelta
import time

# Configuration
BASE_URL = "http://localhost:3000/api/v1"
EMAIL = "test+1@example.com"
PASSWORD = "secret123"

# Mock transaction data
CATEGORIES = [
    {"id": 1, "name": "Ăn uống"},
    {"id": 2, "name": "Giao thông"},
    {"id": 3, "name": "Mua sắm"},
    {"id": 4, "name": "Giải trí"},
    {"id": 5, "name": "Y tế"},
    {"id": 6, "name": "Giáo dục"},
    {"id": 7, "name": "Du lịch"},
    {"id": 8, "name": "Tiện ích"},
]

TRANSACTION_TEMPLATES = [
    # Food & Dining
    {"category_id": 1, "descriptions": ["Ăn trưa", "Ăn tối", "Cà phê", "Đồ ăn nhanh", "Nhà hàng"], "amount_range": (20000, 200000)},
    # Transportation
    {"category_id": 2, "descriptions": ["Xăng xe", "Taxi", "Grab", "Xe bus", "Đậu xe"], "amount_range": (15000, 100000)},
    # Shopping
    {"category_id": 3, "descriptions": ["Quần áo", "Điện thoại", "Laptop", "Sách", "Đồ gia dụng"], "amount_range": (50000, 5000000)},
    # Entertainment
    {"category_id": 4, "descriptions": ["Xem phim", "Game", "Karaoke", "Cafe", "Concert"], "amount_range": (50000, 500000)},
    # Healthcare
    {"category_id": 5, "descriptions": ["Khám bệnh", "Thuốc", "Bảo hiểm", "Nha khoa", "Khám sức khỏe"], "amount_range": (100000, 2000000)},
    # Education
    {"category_id": 6, "descriptions": ["Học phí", "Sách giáo khoa", "Khóa học", "Thi cử", "Dụng cụ học tập"], "amount_range": (200000, 10000000)},
    # Travel
    {"category_id": 7, "descriptions": ["Vé máy bay", "Khách sạn", "Du lịch", "Vé tàu", "Thuê xe"], "amount_range": (500000, 20000000)},
    # Utilities
    {"category_id": 8, "descriptions": ["Điện", "Nước", "Internet", "Điện thoại", "Gas"], "amount_range": (100000, 1000000)},
]

def get_auth_token():
    """Get authentication token"""
    response = requests.post(f"{BASE_URL}/auth/login", json={
        "email": EMAIL,
        "password": PASSWORD
    })
    
    if response.status_code == 200:
        return response.json()["access_token"]
    else:
        print(f"Login failed: {response.text}")
        return None

def create_transaction(token, transaction_data):
    """Create a transaction"""
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.post(f"{BASE_URL}/transactions", 
                           json=transaction_data, 
                           headers=headers)
    
    if response.status_code == 201:
        return True
    else:
        print(f"Failed to create transaction: {response.text}")
        return False

def generate_mock_transactions(token, num_transactions=50):
    """Generate mock transactions"""
    print(f"Generating {num_transactions} mock transactions...")
    
    # Generate transactions for the last 6 months
    end_date = datetime.now()
    start_date = end_date - timedelta(days=180)
    
    created_count = 0
    
    for i in range(num_transactions):
        # Random date within the last 6 months
        random_days = random.randint(0, 180)
        transaction_date = start_date + timedelta(days=random_days)
        
        # Random time
        random_hours = random.randint(8, 22)
        random_minutes = random.randint(0, 59)
        transaction_time = f"{random_hours:02d}:{random_minutes:02d}"
        
        # Random category and template
        template = random.choice(TRANSACTION_TEMPLATES)
        description = random.choice(template["descriptions"])
        
        # Random amount within range
        min_amount, max_amount = template["amount_range"]
        amount = random.randint(min_amount, max_amount)
        
        # Random transaction type (mostly expenses, some income)
        transaction_type = "expense" if random.random() < 0.9 else "income"
        
        # Create transaction data
        transaction_data = {
            "category_id": template["category_id"],
            "amount": amount,
            "description": description,
            "transaction_type": transaction_type,
            "transaction_date": transaction_date.strftime("%Y-%m-%d"),
            "transaction_time": transaction_time,
            "location": random.choice(["Hà Nội", "TP.HCM", "Đà Nẵng", "Hải Phòng", ""]),
            "tags": random.sample(["urgent", "work", "personal", "family", "business"], random.randint(0, 2))
        }
        
        # Create transaction
        if create_transaction(token, transaction_data):
            created_count += 1
            print(f"Created transaction {i+1}/{num_transactions}: {description} - {amount:,} VND")
        
        # Small delay to avoid overwhelming the API
        time.sleep(0.1)
    
    print(f"Successfully created {created_count}/{num_transactions} transactions")
    return created_count

def test_ai_service():
    """Test AI service with the new data"""
    print("\nTesting AI service...")
    
    # Test prediction endpoint
    response = requests.post("http://localhost:8001/api/v1/prediction/expenses", json={
        "user_id": 2,
        "start_date": "2025-04-01T00:00:00Z",
        "end_date": "2025-10-01T00:00:00Z"
    })
    
    if response.status_code == 200:
        data = response.json()
        print("AI Prediction Response:")
        print(f"  Predicted Amount: {data.get('predicted_amount', 0):,} VND")
        print(f"  Confidence Score: {data.get('confidence_score', 0):.2%}")
        print(f"  Trends: {len(data.get('trends', []))}")
        print(f"  Recommendations: {len(data.get('recommendations', []))}")
    else:
        print(f"AI service test failed: {response.text}")

def main():
    print("🚀 TabiMoney Mock Data Generator")
    print("=" * 50)
    
    # Get authentication token
    print("1. Authenticating...")
    token = get_auth_token()
    if not token:
        print("❌ Authentication failed!")
        return
    
    print("✅ Authentication successful!")
    
    # Generate mock transactions
    print("\n2. Generating mock transactions...")
    created_count = generate_mock_transactions(token, num_transactions=100)
    
    if created_count > 0:
        print(f"\n✅ Successfully created {created_count} transactions!")
        
        # Wait a bit for data to be processed
        print("\n3. Waiting for data processing...")
        time.sleep(5)
        
        # Test AI service
        test_ai_service()
        
        print("\n🎉 Mock data generation complete!")
        print("You can now test the Analytics Dashboard at: http://localhost:3000/analytics")
    else:
        print("❌ No transactions were created!")

if __name__ == "__main__":
    main()




