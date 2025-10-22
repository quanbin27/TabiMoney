#!/bin/bash

# TabiMoney Telegram Bot Test Script

echo "🧪 Testing TabiMoney Telegram Bot Integration..."

# Check if virtual environment exists
if [ ! -d "venv" ]; then
    echo "❌ Virtual environment not found. Please run setup.sh first."
    exit 1
fi

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source venv/bin/activate

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "❌ .env file not found. Please copy .env.example to .env and configure it."
    exit 1
fi

# Load environment variables
export $(cat .env | grep -v '^#' | xargs)

# Check required environment variables
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "❌ TELEGRAM_BOT_TOKEN not set in .env file"
    exit 1
fi

if [ -z "$BACKEND_URL" ]; then
    echo "❌ BACKEND_URL not set in .env file"
    exit 1
fi

if [ -z "$AI_SERVICE_URL" ]; then
    echo "❌ AI_SERVICE_URL not set in .env file"
    exit 1
fi

echo "✅ Environment variables loaded"

# Test database connection
echo "🗄️ Testing database connection..."
python3 -c "
import pymysql
import os
try:
    conn = pymysql.connect(
        host=os.getenv('DB_HOST', 'localhost'),
        port=int(os.getenv('DB_PORT', '3306')),
        user=os.getenv('DB_USER', 'root'),
        password=os.getenv('DB_PASSWORD', 'password'),
        database=os.getenv('DB_NAME', 'tabimoney')
    )
    print('✅ Database connection successful')
    conn.close()
except Exception as e:
    print(f'❌ Database connection failed: {e}')
    exit(1)
"

if [ $? -ne 0 ]; then
    echo "❌ Database connection test failed"
    exit 1
fi

# Test backend API connection
echo "🌐 Testing backend API connection..."
python3 -c "
import requests
import os
try:
    response = requests.get(f'{os.getenv(\"BACKEND_URL\")}/health', timeout=5)
    if response.status_code == 200:
        print('✅ Backend API connection successful')
    else:
        print(f'⚠️ Backend API responded with status {response.status_code}')
except Exception as e:
    print(f'⚠️ Backend API connection test failed: {e}')
    print('This is expected if backend is not running')
"

# Test AI service connection
echo "🤖 Testing AI service connection..."
python3 -c "
import requests
import os
try:
    response = requests.get(f'{os.getenv(\"AI_SERVICE_URL\")}/health', timeout=5)
    if response.status_code == 200:
        print('✅ AI service connection successful')
    else:
        print(f'⚠️ AI service responded with status {response.status_code}')
except Exception as e:
    print(f'⚠️ AI service connection test failed: {e}')
    print('This is expected if AI service is not running')
"

# Test bot token
echo "🤖 Testing Telegram bot token..."
python3 -c "
import requests
import os
try:
    response = requests.get(f'https://api.telegram.org/bot{os.getenv(\"TELEGRAM_BOT_TOKEN\")}/getMe')
    if response.status_code == 200:
        bot_info = response.json()
        print(f'✅ Bot token valid - Bot: @{bot_info[\"result\"][\"username\"]}')
    else:
        print(f'❌ Bot token invalid: {response.status_code}')
        exit(1)
except Exception as e:
    print(f'❌ Bot token test failed: {e}')
    exit(1)
"

if [ $? -ne 0 ]; then
    echo "❌ Bot token test failed"
    exit 1
fi

echo ""
echo "✅ All tests passed!"
echo ""
echo "🚀 Ready to start the bot:"
echo "   python main.py"
echo ""
echo "📋 Make sure the following services are running:"
echo "   - MySQL database"
echo "   - TabiMoney backend (http://localhost:8080)"
echo "   - TabiMoney AI service (http://localhost:8001)"
echo ""
echo "🔗 Bot commands to test:"
echo "   /start - Welcome message"
echo "   /link - Link account instructions"
echo "   /dashboard - View dashboard (after linking)"
echo "   Send any message - Chat with AI (after linking)"
