#!/usr/bin/env python3
"""
Script để test các tính năng AI sau khi đã tạo mock data
Gọi qua Go backend API (port 8080 hoặc 3000)
"""

import requests
import json
from datetime import datetime, timedelta
import sys

# Configuration
# Frontend chạy trên port 3000 và proxy /api đến backend port 8080
# Có thể test qua:
# - Frontend proxy: http://localhost:3000 (giống như frontend thực tế)
# - Backend trực tiếp: http://localhost:8080 (nhanh hơn, không cần frontend)
BACKEND_URL = "http://localhost:3000"  # Mặc định qua frontend proxy (port 3000)
# Hoặc dùng backend trực tiếp: "http://localhost:8080"
USER_ID = 15
AUTH_TOKEN = None  # Sẽ được set sau khi login

def login_and_get_token(email, password):
    """Login và lấy auth token"""
    try:
        response = requests.post(
            f"{BACKEND_URL}/api/v1/auth/login",
            json={"email": email, "password": password},
            headers={"Content-Type": "application/json"},
            timeout=10
        )
        if response.status_code == 200:
            data = response.json()
            return data.get("access_token")
        else:
            print(f"⚠️  Login failed: {response.status_code}")
            print(f"   Response: {response.text}")
            return None
    except Exception as e:
        print(f"⚠️  Login error: {e}")
        return None

def test_prediction(user_id, token, months_back=6):
    """Test prediction service qua Go backend"""
    print("\n" + "="*60)
    print("🧮 Testing Prediction Service")
    print("="*60)
    
    end_date = datetime.now()
    start_date = end_date - timedelta(days=months_back*30)
    
    params = {
        "start_date": start_date.strftime("%Y-%m-%d"),
        "end_date": end_date.strftime("%Y-%m-%d")
    }
    
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.get(
            f"{BACKEND_URL}/api/v1/analytics/predictions",
            params=params,
            headers=headers,
            timeout=30
        )
        
        if response.status_code == 200:
            data = response.json()
            print(f"✅ Prediction successful!")
            print(f"\n📊 Results:")
            print(f"   Predicted Amount: {data.get('predicted_amount', 0):,.0f} VND")
            print(f"   Confidence Score: {data.get('confidence_score', 0):.2%}")
            print(f"   Generated At: {data.get('generated_at', 'N/A')}")
            
            # Category breakdown
            breakdown = data.get('category_breakdown', [])
            if breakdown:
                print(f"\n📈 Category Breakdown:")
                for item in breakdown[:5]:
                    print(f"   - {item.get('category_name', 'Unknown')}: "
                          f"{item.get('predicted_amount', 0):,.0f} VND "
                          f"(confidence: {item.get('confidence_score', 0):.2%})")
            
            # Trends
            trends = data.get('trends', [])
            if trends:
                print(f"\n📉 Trends ({len(trends)} periods):")
                for trend in trends[-3:]:  # Last 3 trends
                    print(f"   - {trend.get('period', 'N/A')}: "
                          f"{trend.get('amount', 0):,.0f} VND "
                          f"({trend.get('trend', 'stable')})")
            
            # Recommendations
            recommendations = data.get('recommendations', [])
            if recommendations:
                print(f"\n💡 Recommendations:")
                for rec in recommendations[:3]:
                    print(f"   - {rec}")
            
            return True
        else:
            print(f"❌ Prediction failed: {response.status_code}")
            print(f"   Response: {response.text}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Error calling prediction service: {e}")
        return False

def test_anomaly_detection(user_id, token, months_back=6, threshold=0.6):
    """Test anomaly detection service qua Go backend"""
    print("\n" + "="*60)
    print("🔍 Testing Anomaly Detection Service")
    print("="*60)
    
    end_date = datetime.now()
    start_date = end_date - timedelta(days=months_back*30)
    
    params = {
        "start_date": start_date.strftime("%Y-%m-%d"),
        "end_date": end_date.strftime("%Y-%m-%d"),
        "threshold": threshold
    }
    
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.get(
            f"{BACKEND_URL}/api/v1/analytics/anomalies",
            params=params,
            headers=headers,
            timeout=30
        )
        
        if response.status_code == 200:
            data = response.json()
            print(f"✅ Anomaly detection successful!")
            print(f"\n📊 Results:")
            print(f"   Total Anomalies: {data.get('total_anomalies', len(data.get('anomalies', [])))}")
            print(f"   Detection Score: {data.get('detection_score', 0):.2%}")
            
            # Show anomalies
            anomalies = data.get('anomalies', [])
            if anomalies:
                print(f"\n⚠️  Detected Anomalies:")
                for i, anomaly in enumerate(anomalies[:10], 1):  # Show first 10
                    print(f"\n   {i}. Transaction ID: {anomaly.get('transaction_id', 'N/A')}")
                    print(f"      Amount: {anomaly.get('amount', 0):,.0f} VND")
                    print(f"      Category: {anomaly.get('category_name', 'Unknown')}")
                    # Format date properly
                    date_str = anomaly.get('transaction_date', 'N/A')
                    if isinstance(date_str, str) and len(date_str) > 10:
                        date_str = date_str[:10]  # Take only YYYY-MM-DD part
                    print(f"      Date: {date_str}")
                    print(f"      Anomaly Score: {anomaly.get('anomaly_score', 0):.3f}")
                    print(f"      Type: {anomaly.get('anomaly_type', 'unknown')}")
                    print(f"      Description: {anomaly.get('description', 'N/A')}")
            else:
                print(f"\n✅ No anomalies detected (or threshold too high)")
            
            return True
        else:
            print(f"❌ Anomaly detection failed: {response.status_code}")
            print(f"   Response: {response.text}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Error calling anomaly detection service: {e}")
        return False

def main():
    """Main function"""
    global AUTH_TOKEN, BACKEND_URL, USER_ID
    
    print("\n" + "="*60)
    print("🧪 AI Features Test Script")
    print("="*60)
    print(f"\nConfiguration:")
    print(f"   Backend URL: {BACKEND_URL}")
    print(f"   User ID: {USER_ID}")
    print(f"\n💡 Usage:")
    print(f"   python3 test_ai_features.py [user_id] [backend_url]")
    print(f"   Examples:")
    print(f"     python3 test_ai_features.py 15 http://localhost:3000  # Via frontend proxy")
    print(f"     python3 test_ai_features.py 1 http://localhost:8080   # Direct backend")
    
    # Allow override via command line
    if len(sys.argv) > 1:
        USER_ID = int(sys.argv[1])
    if len(sys.argv) > 2:
        BACKEND_URL = sys.argv[2]
    
    # Check if backend is running
    try:
        response = requests.get(f"{BACKEND_URL}/health", timeout=5)
        print(f"\n✅ Backend is running")
    except requests.exceptions.RequestException:
        print(f"\n⚠️  Warning: Cannot connect to backend at {BACKEND_URL}")
        print(f"   Make sure the backend is running before testing")
        response = input("\nContinue anyway? (y/n): ")
        if response.lower() != 'y':
            sys.exit(1)
    
    # Login để lấy token
    print(f"\n🔐 Attempting to login...")
    print(f"   Testing with user ID: {USER_ID}")
    
    # Thử login với test user 15
    if USER_ID == 15:
        print(f"   Using test credentials: test15@tabimoney.com / test123456")
        AUTH_TOKEN = login_and_get_token("test15@tabimoney.com", "test123456")
    elif USER_ID == 1:
        print(f"   ⚠️  User 1 - Please provide credentials:")
        email = input("   Email: ").strip()
        password = input("   Password: ").strip()
        if email and password:
            AUTH_TOKEN = login_and_get_token(email, password)
    
    if not AUTH_TOKEN:
        print("\n⚠️  Could not get auth token automatically")
        print("   Options:")
        print("   1. Enter auth token manually (copy from browser DevTools > Application > Local Storage)")
        print("   2. Or login via frontend and copy the access_token")
        token = input("\n   Enter auth token (or press Enter to skip): ").strip()
        if token:
            AUTH_TOKEN = token
        else:
            print("   ⚠️  Skipping authentication - tests will fail without token")
            AUTH_TOKEN = ""
    
    if AUTH_TOKEN:
        print("✅ Got auth token")
    
    # Test prediction
    if AUTH_TOKEN:
        prediction_ok = test_prediction(USER_ID, AUTH_TOKEN)
    else:
        print("\n⚠️  Skipping prediction test (no auth token)")
        prediction_ok = False
    
    # Test anomaly detection
    if AUTH_TOKEN:
        anomaly_ok = test_anomaly_detection(USER_ID, AUTH_TOKEN)
    else:
        print("\n⚠️  Skipping anomaly detection test (no auth token)")
        anomaly_ok = False
    
    # Summary
    print("\n" + "="*60)
    print("📋 Test Summary")
    print("="*60)
    print(f"   Prediction: {'✅ PASS' if prediction_ok else '❌ FAIL'}")
    print(f"   Anomaly Detection: {'✅ PASS' if anomaly_ok else '❌ FAIL'}")
    print("="*60 + "\n")

if __name__ == "__main__":
    main()

