# TabiMoney Telegram Bot

Telegram bot integration for TabiMoney personal finance management system.

## Features

- 🔗 **Account Linking**: Link Telegram account with TabiMoney web account using secure link codes
- 📊 **Dashboard**: View financial dashboard directly in Telegram
- 🤖 **AI Chat**: Chat with AI assistant for financial analysis and advice
- 🔐 **Secure Authentication**: JWT-based authentication with permanent tokens for Telegram users
- 📱 **User-friendly Interface**: Intuitive commands and responses

## Prerequisites

- Python 3.8 or higher
- MySQL database (same as main TabiMoney system)
- Redis server
- TabiMoney backend running
- TabiMoney AI service running
- Telegram bot token from [@BotFather](https://t.me/botfather)

## Installation

1. **Clone and setup**:
   ```bash
   cd telegram-bot
   ./setup.sh
   ```

2. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

3. **Database migration**:
   ```bash
   mysql -u root -p tabimoney < database/migrations/001_telegram_integration.sql
   ```

4. **Start the bot**:
   ```bash
   source venv/bin/activate
   python main.py
   ```

## Configuration

Edit `.env` file with your settings:

```env
# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=your-telegram-bot-token-here

# Backend API Configuration
BACKEND_URL=http://localhost:8080
AI_SERVICE_URL=http://localhost:8001

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=tabimoney

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
```

## Usage

### For Users

1. **Start the bot**: Send `/start` to @TabiMoneyBot
2. **Link account**: 
   - Go to TabiMoney web app → Settings → Telegram Integration
   - Generate a link code
   - Send `/link` to the bot
   - Enter the link code
3. **Use features**:
   - `/dashboard` - View financial dashboard
   - Send any message to chat with AI assistant

### Bot Commands

- `/start` - Welcome message and instructions
- `/link` - Generate link code for account linking
- `/dashboard` - View financial dashboard
- Any text message - Chat with AI assistant

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Telegram Bot  │    │   Backend API   │    │   AI Service    │
│                 │    │                 │    │                 │
│ • Handlers      │◄──►│ • Auth Service  │◄──►│ • Chat API      │
│ • Auth Middleware│    │ • Telegram APIs │    │ • NLU Service   │
│ • API Service   │    │ • Database      │    │ • ML Service    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Database       │
                    │                 │
                    │ • Users         │
                    │ • Telegram      │
                    │ • Link Codes    │
                    │ • Sessions      │
                    └─────────────────┘
```

## Security

- **Link Codes**: Temporary codes (10 minutes expiry) for secure account linking
- **JWT Tokens**: Permanent tokens (1 year expiry) for Telegram users
- **Authentication**: Middleware checks authentication before processing requests
- **Rate Limiting**: Built-in rate limiting to prevent abuse

## Development

### Project Structure

```
telegram-bot/
├── main.py                 # Main bot entry point
├── handlers/               # Command and message handlers
│   ├── start.py           # /start command
│   ├── link.py            # /link command and link code handling
│   ├── dashboard.py       # /dashboard command
│   └── chat.py            # AI chat message handling
├── services/              # Business logic services
│   ├── auth_service.py    # Authentication and linking
│   └── api_service.py     # Backend API communication
├── middleware/            # Middleware components
│   └── auth_middleware.py # Authentication middleware
├── utils/                 # Utility modules
│   ├── config.py         # Configuration management
│   └── logger.py         # Logging setup
├── database/             # Database migrations
│   └── migrations/       # SQL migration files
├── requirements.txt      # Python dependencies
├── .env.example         # Environment configuration template
└── setup.sh            # Setup script
```

### Adding New Features

1. **Create handler**: Add new handler in `handlers/` directory
2. **Add service methods**: Extend services in `services/` directory
3. **Update main.py**: Register new handlers
4. **Test**: Test with local bot instance

## Troubleshooting

### Common Issues

1. **Bot not responding**:
   - Check if bot token is correct
   - Verify backend and AI service are running
   - Check logs for errors

2. **Link code not working**:
   - Ensure database migration is run
   - Check if link code is expired (10 minutes)
   - Verify Redis is running

3. **Authentication errors**:
   - Check JWT secret configuration
   - Verify database connection
   - Check if user is properly linked

### Logs

Bot logs are written to console. For production, consider using a proper logging service.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is part of TabiMoney and follows the same license terms.
