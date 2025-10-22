#!/bin/bash

# TabiMoney Telegram Bot Setup Script

echo "🤖 Setting up TabiMoney Telegram Bot..."

# Check if Python 3 is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is not installed. Please install Python 3.8 or higher."
    exit 1
fi

# Check if pip is installed
if ! command -v pip3 &> /dev/null; then
    echo "❌ pip3 is not installed. Please install pip3."
    exit 1
fi

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "📦 Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source venv/bin/activate

# Install dependencies
echo "📥 Installing dependencies..."
pip install -r requirements.txt

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "⚙️ Creating .env file..."
    cp .env.example .env
    echo "📝 Please edit .env file and add your Telegram bot token and other configuration."
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p logs
mkdir -p database/migrations

# Run database migrations
echo "🗄️ Running database migrations..."
echo "Please make sure your database is running and run the migration manually:"
echo "mysql -u root -p tabimoney < database/migrations/001_telegram_integration.sql"

echo ""
echo "✅ Setup completed!"
echo ""
echo "📋 Next steps:"
echo "1. Edit .env file and add your Telegram bot token"
echo "2. Make sure your backend and AI service are running"
echo "3. Run database migration: mysql -u root -p tabimoney < database/migrations/001_telegram_integration.sql"
echo "4. Start the bot: python main.py"
echo ""
echo "🚀 To start the bot, run: python main.py"
