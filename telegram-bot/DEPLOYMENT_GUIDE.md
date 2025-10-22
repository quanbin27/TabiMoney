# Hướng dẫn triển khai TabiMoney Telegram Bot

## Tổng quan

TabiMoney Telegram Bot đã được xây dựng hoàn chỉnh với các tính năng:

- ✅ **Liên kết tài khoản**: Sử dụng link-code mechanism để liên kết Telegram với web account
- ✅ **JWT Authentication**: Token không hết hạn (1 năm) cho Telegram users
- ✅ **Dashboard**: Xem thông tin tài chính cơ bản
- ✅ **AI Chat**: Tích hợp với AI service để chat và phân tích tài chính
- ✅ **Frontend Integration**: Component trong Settings để generate link-code

## Cấu trúc hệ thống

```
TabiMoney Telegram Bot
├── main.py                    # Entry point
├── handlers/                  # Bot command handlers
│   ├── start.py              # /start command
│   ├── link.py               # /link command & link code handling
│   ├── dashboard.py          # /dashboard command
│   └── chat.py               # AI chat handling
├── services/                 # Business logic
│   ├── auth_service.py       # Authentication & linking
│   └── api_service.py        # Backend API communication
├── middleware/               # Middleware components
│   └── auth_middleware.py    # Authentication middleware
├── utils/                    # Utilities
│   ├── config.py            # Configuration
│   └── logger.py            # Logging
├── database/                # Database migrations
│   └── migrations/          # SQL files
├── requirements.txt         # Dependencies
├── .env.example            # Environment template
├── setup.sh               # Setup script
├── test.sh                # Test script
└── README.md              # Documentation
```

## Triển khai

### 1. Chuẩn bị môi trường

```bash
# Di chuyển vào thư mục telegram-bot
cd telegram-bot

# Chạy script setup
./setup.sh

# Cấu hình environment
cp .env.example .env
# Chỉnh sửa .env với thông tin của bạn
```

### 2. Cấu hình .env

```env
# Telegram Bot Token (từ @BotFather)
TELEGRAM_BOT_TOKEN=your-telegram-bot-token-here

# Backend URLs
BACKEND_URL=http://localhost:8080
AI_SERVICE_URL=http://localhost:8001

# Database (cùng với hệ thống chính)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=tabimoney

# JWT Secret (cùng với hệ thống chính)
JWT_SECRET=your-super-secret-jwt-key-here
```

### 3. Database Migration

```bash
# Chạy migration để tạo bảng Telegram
mysql -u root -p tabimoney < database/migrations/001_telegram_integration.sql
```

### 4. Test hệ thống

```bash
# Chạy test script
./test.sh
```

### 5. Khởi động bot

```bash
# Activate virtual environment
source venv/bin/activate

# Start bot
python main.py
```

## Cách sử dụng

### Cho User

1. **Tìm bot**: Tìm @TabiMoneyBot trên Telegram
2. **Bắt đầu**: Gửi `/start`
3. **Liên kết tài khoản**:
   - Gửi `/link` để xem hướng dẫn
   - Vào web app → Settings → Telegram Integration
   - Tạo mã liên kết
   - Gửi mã cho bot
4. **Sử dụng tính năng**:
   - `/dashboard` - Xem dashboard
   - Gửi tin nhắn bất kỳ - Chat với AI

### Bot Commands

- `/start` - Tin nhắn chào mừng và hướng dẫn
- `/link` - Hướng dẫn liên kết tài khoản
- `/dashboard` - Xem dashboard tài chính
- Tin nhắn bất kỳ - Chat với AI assistant

## Tích hợp với Backend

### API Endpoints đã thêm

Backend đã được mở rộng với các endpoints:

```go
// Telegram integration endpoints
POST /auth/telegram/generate-link-code  // Tạo mã liên kết
GET  /auth/telegram/status              // Kiểm tra trạng thái liên kết
POST /auth/telegram/disconnect          // Hủy liên kết
POST /auth/telegram/link                // Liên kết với mã
```

### Database Schema

```sql
-- Bảng liên kết tài khoản Telegram
CREATE TABLE telegram_accounts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    telegram_user_id BIGINT NOT NULL UNIQUE,
    web_user_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (web_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Bảng mã liên kết tạm thời
CREATE TABLE telegram_link_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(16) NOT NULL UNIQUE,
    web_user_id BIGINT UNSIGNED NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (web_user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## Frontend Integration

### Settings Component

Frontend đã được cập nhật với component Telegram Integration trong Settings:

- ✅ Hiển thị trạng thái liên kết
- ✅ Tạo mã liên kết
- ✅ Copy mã liên kết
- ✅ Hủy liên kết
- ✅ Hướng dẫn sử dụng

## Bảo mật

### Link Code Mechanism

1. **Tạo mã**: User tạo mã trên web app
2. **Lưu trữ**: Mã được lưu trong database với thời hạn 10 phút
3. **Sử dụng**: User gửi mã cho bot
4. **Xác thực**: Bot validate mã và liên kết tài khoản
5. **Hủy mã**: Mã được đánh dấu đã sử dụng

### JWT Token

- **Loại**: `telegram_access`
- **Thời hạn**: 1 năm
- **Claims**: `user_id`, `telegram_user_id`, `type`
- **Sử dụng**: Gọi API backend và AI service

## Monitoring & Logging

### Logs

Bot ghi log ra console với format JSON:

```json
{
  "timestamp": "2024-01-01T00:00:00",
  "level": "INFO",
  "logger": "handlers.start",
  "message": "Start command executed for user 123456"
}
```

### Error Handling

- ✅ Graceful error handling cho tất cả operations
- ✅ User-friendly error messages
- ✅ Detailed logging cho debugging
- ✅ Fallback responses

## Troubleshooting

### Lỗi thường gặp

1. **Bot không phản hồi**:
   - Kiểm tra bot token
   - Kiểm tra backend/AI service đang chạy
   - Xem logs

2. **Mã liên kết không hoạt động**:
   - Kiểm tra database migration
   - Kiểm tra mã đã hết hạn (10 phút)
   - Kiểm tra mã đã được sử dụng

3. **Lỗi authentication**:
   - Kiểm tra JWT secret
   - Kiểm tra database connection
   - Kiểm tra user đã liên kết

### Debug Commands

```bash
# Test database connection
python3 -c "import pymysql; print('DB OK')"

# Test bot token
curl "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getMe"

# Test backend API
curl "$BACKEND_URL/health"

# Test AI service
curl "$AI_SERVICE_URL/health"
```

## Production Deployment

### Recommendations

1. **Process Management**: Sử dụng systemd hoặc PM2
2. **Logging**: Chuyển sang file logging
3. **Monitoring**: Thêm health checks
4. **Scaling**: Có thể chạy multiple instances
5. **Security**: Sử dụng HTTPS cho webhooks

### Systemd Service

```ini
[Unit]
Description=TabiMoney Telegram Bot
After=network.target

[Service]
Type=simple
User=tabimoney
WorkingDirectory=/path/to/telegram-bot
ExecStart=/path/to/telegram-bot/venv/bin/python main.py
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## Kết luận

TabiMoney Telegram Bot đã được xây dựng hoàn chỉnh với:

- ✅ **Authentication**: Link-code mechanism an toàn
- ✅ **Features**: Dashboard và AI chat
- ✅ **Integration**: Tích hợp hoàn chỉnh với backend/AI
- ✅ **Frontend**: Component để quản lý liên kết
- ✅ **Documentation**: Hướng dẫn chi tiết
- ✅ **Testing**: Script test tự động

Bot sẵn sàng để triển khai và sử dụng trong môi trường production!
