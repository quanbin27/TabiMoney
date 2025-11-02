# Hệ thống Notification TabiMoney

## Tổng quan

Hệ thống notification của TabiMoney cung cấp thông báo đa kênh (Email, Telegram, In-App) cho các tính năng từ urgent đến medium priority, giúp người dùng theo dõi và quản lý tài chính hiệu quả.

## Các loại Notification

### 🔴 Urgent Priority
- **Budget Exceeded**: Ngân sách vượt quá 100%
- **Goal Deadline Warning**: Cảnh báo hạn chót mục tiêu (30 ngày trước)

### 🟠 High Priority  
- **Budget Threshold Alert**: Ngân sách đạt ngưỡng cảnh báo (80%)
- **Anomaly Detection**: Phát hiện giao dịch bất thường
- **Goal Achievement**: Hoàn thành mục tiêu tài chính

### 🟡 Medium Priority
- **Goal Progress Updates**: Cập nhật tiến độ mục tiêu (25%, 50%, 75%, 90%)
- **Large Transaction Alert**: Giao dịch lớn (>1M VND)
- **Spending Prediction**: Dự đoán chi tiêu tháng tới
- **Financial Health Alert**: Cảnh báo sức khỏe tài chính

### 🟢 Low Priority
- **Monthly Reports**: Báo cáo tài chính hàng tháng
- **Budget Reminders**: Nhắc nhở ngân sách cuối tháng

## Kênh Notification

### 1. Email Notifications
- **Template**: HTML responsive với màu sắc theo loại notification
- **Cấu hình**: SMTP settings trong `.env`
- **Features**: 
  - Template động theo notification type
  - Action buttons với deep links
  - Branded design với TabiMoney logo

### 2. Telegram Notifications
- **Format**: Markdown với emoji và formatting
- **Features**:
  - Inline keyboard với action buttons
  - Rich formatting cho số tiền và phần trăm
  - Deep links đến web app

### 3. In-App Notifications
- **Storage**: Database với read/unread status
- **API**: RESTful endpoints cho CRUD operations
- **Features**: Real-time updates, pagination, filtering

## Cấu hình

### Environment Variables

```env
# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=TabiMoney

# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
```

### User Preferences

Người dùng có thể cấu hình:
- **Channels**: Bật/tắt email, telegram, in-app
- **Features**: Bật/tắt alerts cho budget, goals, AI, analytics
- **Priority**: Chọn priority levels muốn nhận
- **Frequency**: Daily/weekly/monthly digest
- **Quiet Hours**: Thời gian không nhận notification

## API Endpoints

### Notifications
```
GET    /api/notifications              # Lấy danh sách notifications
POST   /api/notifications/:id/read     # Đánh dấu đã đọc
```

### Notification Preferences
```
GET    /api/notification-preferences           # Lấy preferences
PUT    /api/notification-preferences           # Cập nhật preferences
GET    /api/notification-preferences/summary   # Lấy summary
POST   /api/notification-preferences/reset     # Reset về default
GET    /api/notification-preferences/channels  # Lấy enabled channels
POST   /api/notification-preferences/test      # Gửi test notification
```

## Trigger System

### Budget Management
- **Threshold Alert**: Khi ngân sách đạt ngưỡng cảnh báo
- **Exceeded Alert**: Khi ngân sách vượt quá 100%
- **Achievement Alert**: Khi hoàn thành tiết kiệm ngân sách

### Financial Goals
- **Progress Milestones**: 25%, 50%, 75%, 90%
- **Deadline Warning**: 30 ngày trước hạn chót
- **Achievement**: Khi hoàn thành mục tiêu

### AI Features
- **Anomaly Detection**: Giao dịch bất thường (amount, frequency, pattern)
- **Spending Prediction**: Dự đoán chi tiêu tháng tới
- **Category Suggestion**: AI đề xuất danh mục

### Analytics
- **Monthly Reports**: Báo cáo tài chính hàng tháng
- **Financial Health**: Cảnh báo sức khỏe tài chính
- **Spending Trends**: Xu hướng chi tiêu

### Transaction Management
- **Large Transactions**: Giao dịch >1M VND
- **Recurring Payments**: Thanh toán định kỳ đến hạn

## Scheduled Service

### Chạy mỗi giờ:
- Kiểm tra budget alerts
- Kiểm tra goal alerts
- Kiểm tra monthly reports
- Kiểm tra financial health alerts

### Chạy hàng ngày:
- Anomaly detection cho tất cả users
- Spending prediction cho tháng tới

## Database Schema

### notifications table
```sql
CREATE TABLE notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    notification_type ENUM('info', 'warning', 'success', 'error', 'reminder') NOT NULL,
    priority ENUM('low', 'medium', 'high', 'urgent') DEFAULT 'medium',
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    action_url VARCHAR(500),
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### user_profiles table (notification_settings)
```json
{
  "email_enabled": true,
  "telegram_enabled": true,
  "in_app_enabled": true,
  "budget_alerts": true,
  "goal_alerts": true,
  "ai_alerts": true,
  "transaction_alerts": true,
  "analytics_alerts": true,
  "urgent_notifications": true,
  "high_notifications": true,
  "medium_notifications": true,
  "low_notifications": false,
  "daily_digest": false,
  "weekly_digest": true,
  "monthly_digest": true,
  "real_time_alerts": true,
  "quiet_hours_start": "22:00",
  "quiet_hours_end": "08:00",
  "timezone": "Asia/Ho_Chi_Minh"
}
```

## Cách sử dụng

### 1. Cấu hình Email
```bash
# Gmail App Password
# 1. Bật 2FA cho Gmail
# 2. Tạo App Password tại: https://myaccount.google.com/apppasswords
# 3. Sử dụng App Password thay vì mật khẩu thường
```

### 2. Cấu hình Telegram Bot
```bash
# 1. Tạo bot với @BotFather
# 2. Lấy bot token
# 3. Thêm bot token vào .env
# 4. User cần link Telegram account với web app
```

### 3. Test Notification
```bash
curl -X POST "http://localhost:8080/api/notification-preferences/test?channel=email" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Monitoring & Debugging

### Logs
- Email sending: `log.Printf("Email sent successfully to %s", to)`
- Telegram sending: `log.Printf("Telegram message sent successfully to chat %d", chatID)`
- Notification creation: `log.Printf("notification created: id=%d user=%d title=%s", n.ID, userID, title)`

### Health Checks
- Email: Kiểm tra SMTP connection
- Telegram: Kiểm tra bot token validity
- Database: Kiểm tra notification table

## Best Practices

### 1. Rate Limiting
- Không gửi quá 10 notifications/giờ cho 1 user
- Sử dụng quiet hours để tránh spam
- Implement exponential backoff cho failed sends

### 2. Error Handling
- Graceful degradation khi service down
- Retry mechanism cho failed notifications
- Fallback to in-app khi external services fail

### 3. Performance
- Async processing cho email/telegram
- Batch processing cho scheduled tasks
- Database indexing cho queries

### 4. Security
- Validate user permissions trước khi gửi
- Sanitize notification content
- Rate limiting cho API endpoints

## Troubleshooting

### Email không gửi được
1. Kiểm tra SMTP credentials
2. Kiểm tra firewall/network
3. Kiểm tra Gmail App Password
4. Kiểm tra logs: `Failed to send email notification`

### Telegram không gửi được
1. Kiểm tra bot token
2. Kiểm tra user đã link Telegram account
3. Kiểm tra logs: `Failed to send telegram notification`

### Notification không tạo được
1. Kiểm tra database connection
2. Kiểm tra user permissions
3. Kiểm tra logs: `notification create failed`

## Future Enhancements

### Planned Features
- **Push Notifications**: Web push notifications
- **SMS Notifications**: SMS alerts cho urgent notifications
- **WhatsApp Integration**: WhatsApp Business API
- **Slack Integration**: Slack notifications cho teams
- **Advanced Scheduling**: Custom notification schedules
- **A/B Testing**: Test different notification formats
- **Analytics**: Notification open rates, click rates
- **Templates**: Custom notification templates
- **Multi-language**: Support multiple languages
- **Rich Media**: Images, charts trong notifications

