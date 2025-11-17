# DANH SÁCH API CHI TIẾT - HỆ THỐNG TABIMONEY

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Tất cả API (trừ `/auth/register`, `/auth/login`, `/auth/telegram/link`) đều yêu cầu JWT token trong header:
```
Authorization: Bearer <access_token>
```

---

## 1. AUTHENTICATION & USER MANAGEMENT

### 1.1. Đăng ký tài khoản
- **Method:** `POST`
- **Endpoint:** `/auth/register`
- **Mô tả:** Tạo tài khoản người dùng mới
- **Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "First",
  "last_name": "Last",
  "phone": "0123456789"
}
```
- **Response (201):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "first_name": "First",
    "last_name": "Last",
    "created_at": "2024-01-15T10:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-16T10:00:00Z"
}
```

### 1.2. Đăng nhập
- **Method:** `POST`
- **Endpoint:** `/auth/login`
- **Mô tả:** Đăng nhập và nhận JWT tokens
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
- **Response (200):** Tương tự như đăng ký

### 1.3. Refresh Token
- **Method:** `POST`
- **Endpoint:** `/auth/refresh`
- **Mô tả:** Làm mới access token bằng refresh token
- **Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```
- **Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-16T10:00:00Z"
}
```

### 1.4. Đăng xuất
- **Method:** `POST`
- **Endpoint:** `/auth/logout`
- **Mô tả:** Đăng xuất và vô hiệu hóa session
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

### 1.5. Đổi mật khẩu
- **Method:** `POST`
- **Endpoint:** `/auth/change-password`
- **Mô tả:** Đổi mật khẩu của user hiện tại
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword123"
}
```
- **Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

### 1.6. Lấy thông tin profile
- **Method:** `GET`
- **Endpoint:** `/auth/profile`
- **Mô tả:** Lấy thông tin profile của user hiện tại
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "username",
  "first_name": "First",
  "last_name": "Last",
  "phone": "0123456789",
  "avatar_url": "https://...",
  "profile": {
    "monthly_income": 10000000,
    "currency": "VND",
    "timezone": "Asia/Ho_Chi_Minh",
    "language": "vi"
  }
}
```

### 1.7. Cập nhật profile
- **Method:** `PUT`
- **Endpoint:** `/auth/profile`
- **Mô tả:** Cập nhật thông tin profile
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "first_name": "New First",
  "last_name": "New Last",
  "phone": "0987654321",
  "avatar_url": "https://..."
}
```
- **Response (200):** Tương tự như GET profile

### 1.8. Lấy thu nhập hàng tháng
- **Method:** `GET`
- **Endpoint:** `/auth/income`
- **Mô tả:** Lấy thu nhập hàng tháng của user
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "monthly_income": 10000000,
  "currency": "VND"
}
```

### 1.9. Thiết lập thu nhập hàng tháng
- **Method:** `PUT`
- **Endpoint:** `/auth/income`
- **Mô tả:** Thiết lập hoặc cập nhật thu nhập hàng tháng
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "monthly_income": 15000000,
  "currency": "VND"
}
```
- **Response (200):** Tương tự như GET income

---

## 2. TELEGRAM INTEGRATION

### 2.1. Tạo link code để liên kết Telegram
- **Method:** `POST`
- **Endpoint:** `/auth/telegram/generate-link-code`
- **Mô tả:** Tạo mã liên kết để liên kết tài khoản Telegram
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "code": "ABC123XYZ456",
  "expires_at": "2024-01-15T10:10:00Z"
}
```

### 2.2. Kiểm tra trạng thái liên kết Telegram
- **Method:** `GET`
- **Endpoint:** `/auth/telegram/status`
- **Mô tả:** Kiểm tra xem tài khoản đã liên kết Telegram chưa
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "linked": true,
  "telegram_user_id": 123456789,
  "linked_at": "2024-01-15T09:00:00Z"
}
```

### 2.3. Ngắt liên kết Telegram
- **Method:** `POST`
- **Endpoint:** `/auth/telegram/disconnect`
- **Mô tả:** Ngắt liên kết tài khoản Telegram
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "message": "Telegram account disconnected"
}
```

### 2.4. Liên kết tài khoản Telegram
- **Method:** `POST`
- **Endpoint:** `/auth/telegram/link`
- **Mô tả:** Liên kết tài khoản Telegram với web account (gọi từ Telegram Bot)
- **Request Body:**
```json
{
  "code": "ABC123XYZ456",
  "telegram_user_id": 123456789
}
```
- **Response (200):**
```json
{
  "message": "Telegram account linked successfully"
}
```

---

## 3. TRANSACTIONS (GIAO DỊCH)

### 3.1. Lấy danh sách giao dịch
- **Method:** `GET`
- **Endpoint:** `/transactions`
- **Mô tả:** Lấy danh sách giao dịch với phân trang và lọc
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `page` (int, default: 1): Số trang
  - `limit` (int, default: 20): Số items mỗi trang
  - `category_id` (uint64, optional): Lọc theo category
  - `transaction_type` (string: income/expense/transfer, optional): Lọc theo loại
  - `start_date` (date: YYYY-MM-DD, optional): Ngày bắt đầu
  - `end_date` (date: YYYY-MM-DD, optional): Ngày kết thúc
  - `min_amount` (float, optional): Số tiền tối thiểu
  - `max_amount` (float, optional): Số tiền tối đa
  - `search` (string, optional): Tìm kiếm trong description
  - `sort_by` (string: created_at/transaction_date/amount, optional)
  - `sort_order` (string: asc/desc, optional)
- **Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "category_id": 5,
      "amount": 50000,
      "description": "Ăn bún bò",
      "transaction_type": "expense",
      "transaction_date": "2024-01-15",
      "category": {
        "id": 5,
        "name": "Ăn uống",
        "icon": "restaurant",
        "color": "#FF6B6B"
      }
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

### 3.2. Tạo giao dịch mới
- **Method:** `POST`
- **Endpoint:** `/transactions`
- **Mô tả:** Tạo giao dịch mới
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "category_id": 5,
  "amount": 50000,
  "description": "Ăn bún bò",
  "transaction_type": "expense",
  "transaction_date": "2024-01-15",
  "transaction_time": "12:30:00",
  "location": "Quán bún bò Huế",
  "tags": ["lunch", "food"],
  "metadata": {
    "payment_method": "cash"
  },
  "is_recurring": false,
  "recurring_pattern": null
}
```
- **Response (201):** Transaction object

### 3.3. Cập nhật giao dịch
- **Method:** `PUT`
- **Endpoint:** `/transactions/:id`
- **Mô tả:** Cập nhật giao dịch đã tồn tại
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của giao dịch
- **Request Body:** Tương tự như tạo mới (không có is_recurring, recurring_pattern)
- **Response (200):** Transaction object đã cập nhật

### 3.4. Xóa giao dịch
- **Method:** `DELETE`
- **Endpoint:** `/transactions/:id`
- **Mô tả:** Xóa giao dịch
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của giao dịch
- **Response (200):**
```json
{
  "message": "Deleted"
}
```

---

## 4. CATEGORIES (DANH MỤC)

### 4.1. Lấy danh sách categories
- **Method:** `GET`
- **Endpoint:** `/categories`
- **Mô tả:** Lấy danh sách categories (system + user custom)
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Ăn uống",
      "name_en": "Food & Dining",
      "icon": "restaurant",
      "color": "#FF6B6B",
      "is_system": true,
      "is_active": true
    }
  ]
}
```

### 4.2. Tạo category mới
- **Method:** `POST`
- **Endpoint:** `/categories`
- **Mô tả:** Tạo category tùy chỉnh của user
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "name": "Học phí",
  "name_en": "Tuition",
  "description": "Chi phí học tập",
  "icon": "school",
  "color": "#4ECDC4",
  "parent_id": null
}
```
- **Response (201):** Category object

### 4.3. Cập nhật category
- **Method:** `PUT`
- **Endpoint:** `/categories/:id`
- **Mô tả:** Cập nhật category (chỉ category của user)
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của category
- **Request Body:** Tương tự như tạo mới
- **Response (200):** Category object đã cập nhật

### 4.4. Xóa category
- **Method:** `DELETE`
- **Endpoint:** `/categories/:id`
- **Mô tả:** Xóa category (chỉ category của user, không thể xóa system categories)
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của category
- **Response (200):**
```json
{
  "message": "Category deleted"
}
```

---

## 5. FINANCIAL GOALS (MỤC TIÊU TÀI CHÍNH)

### 5.1. Lấy danh sách mục tiêu
- **Method:** `GET`
- **Endpoint:** `/goals`
- **Mô tả:** Lấy tất cả mục tiêu tài chính của user
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Mua xe máy",
      "description": "Tiết kiệm để mua xe máy mới",
      "target_amount": 50000000,
      "current_amount": 10000000,
      "target_date": "2024-12-31",
      "goal_type": "purchase",
      "priority": "high",
      "is_achieved": false,
      "progress": 20.0
    }
  ]
}
```

### 5.2. Tạo mục tiêu mới
- **Method:** `POST`
- **Endpoint:** `/goals`
- **Mô tả:** Tạo mục tiêu tài chính mới
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "title": "Mua xe máy",
  "description": "Tiết kiệm để mua xe máy mới",
  "target_amount": 50000000,
  "target_date": "2024-12-31",
  "goal_type": "purchase",
  "priority": "high"
}
```
- **Response (201):** Goal object

### 5.3. Cập nhật mục tiêu
- **Method:** `PUT`
- **Endpoint:** `/goals/:id`
- **Mô tả:** Cập nhật mục tiêu
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của goal
- **Request Body:**
```json
{
  "title": "Mua xe máy mới",
  "target_amount": 60000000,
  "current_amount": 15000000
}
```
- **Response (200):** Goal object đã cập nhật

### 5.4. Xóa mục tiêu
- **Method:** `DELETE`
- **Endpoint:** `/goals/:id`
- **Mô tả:** Xóa mục tiêu
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của goal
- **Response (200):**
```json
{
  "message": "Goal deleted successfully"
}
```

### 5.5. Thêm tiền vào mục tiêu
- **Method:** `POST`
- **Endpoint:** `/goals/:id/contribute`
- **Mô tả:** Thêm tiền vào mục tiêu (contribution)
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của goal
- **Request Body:**
```json
{
  "amount": 5000000,
  "note": "Tiết kiệm tháng 1"
}
```
- **Response (200):** Goal object đã cập nhật (với current_amount mới)

---

## 6. BUDGETS (NGÂN SÁCH)

### 6.1. Lấy danh sách ngân sách
- **Method:** `GET`
- **Endpoint:** `/budgets`
- **Mô tả:** Lấy tất cả ngân sách của user
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "category_id": 5,
      "name": "Ngân sách ăn uống tháng 1",
      "amount": 2000000,
      "period": "monthly",
      "start_date": "2024-01-01",
      "end_date": "2024-01-31",
      "is_active": true,
      "alert_threshold": 80.0,
      "spent_amount": 1500000,
      "remaining_amount": 500000,
      "usage_percentage": 75.0,
      "category": {
        "id": 5,
        "name": "Ăn uống"
      }
    }
  ]
}
```

### 6.2. Tạo ngân sách mới
- **Method:** `POST`
- **Endpoint:** `/budgets`
- **Mô tả:** Tạo ngân sách mới
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "category_id": 5,
  "name": "Ngân sách ăn uống tháng 1",
  "amount": 2000000,
  "period": "monthly",
  "start_date": "2024-01-01",
  "end_date": "2024-01-31",
  "alert_threshold": 80.0
}
```
- **Response (201):** Budget object

### 6.3. Cập nhật ngân sách
- **Method:** `PUT`
- **Endpoint:** `/budgets/:id`
- **Mô tả:** Cập nhật ngân sách
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của budget
- **Request Body:** Tương tự như tạo mới, thêm `is_active`
- **Response (200):** Budget object đã cập nhật

### 6.4. Xóa ngân sách
- **Method:** `DELETE`
- **Endpoint:** `/budgets/:id`
- **Mô tả:** Xóa ngân sách
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của budget
- **Response (200):**
```json
{
  "message": "Budget deleted successfully"
}
```

### 6.5. Lấy insights ngân sách
- **Method:** `GET`
- **Endpoint:** `/budgets/insights`
- **Mô tả:** Lấy thông tin phân tích ngân sách (safe-to-spend, pacing)
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "data": {
    "user_id": 1,
    "period": "monthly",
    "as_of": "2024-01-15T10:00:00Z",
    "safe_to_spend_daily": 166666.67,
    "safe_to_spend_weekly": 1166666.67,
    "total_remaining": 5000000,
    "days_left": 16,
    "budgets": [
      {
        "budget_id": 1,
        "name": "Ngân sách ăn uống",
        "amount": 2000000,
        "spent_amount": 1500000,
        "remaining_amount": 500000,
        "usage_percentage": 75.0,
        "allowed_pace_pct": 50.0,
        "actual_pace_pct": 75.0,
        "is_over_pace": true
      }
    ],
    "projected_end_usage_pct": 150.0,
    "risk_budget_ids": [1]
  }
}
```

### 6.6. Lấy đề xuất ngân sách tự động
- **Method:** `GET`
- **Endpoint:** `/budgets/auto/suggestions`
- **Mô tả:** AI đề xuất ngân sách dựa trên lịch sử chi tiêu
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "data": {
    "user_id": 1,
    "monthly_income": 10000000,
    "period": "monthly",
    "start_date": "2024-02-01",
    "end_date": "2024-02-29",
    "suggestions": [
      {
        "category_id": 5,
        "name": "Ăn uống",
        "suggested_amount": 2200000
      }
    ],
    "total_suggested": 9000000,
    "notes": [
      "Dựa trên chi tiêu trung bình 3 tháng gần nhất",
      "Đã để lại 10% cho tiết kiệm"
    ]
  }
}
```

### 6.7. Tạo ngân sách từ đề xuất
- **Method:** `POST`
- **Endpoint:** `/budgets/auto/create`
- **Mô tả:** Tạo nhiều ngân sách từ danh sách đề xuất
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "period": "monthly",
  "start_date": "2024-02-01",
  "end_date": "2024-02-29",
  "budgets": [
    {
      "category_id": 5,
      "suggested_amount": 2200000
    }
  ],
  "alert_threshold": 80.0
}
```
- **Response (201):**
```json
{
  "data": [
    {
      "id": 2,
      "name": "Ngân sách ăn uống tháng 2",
      "amount": 2200000
    }
  ]
}
```

---

## 7. NOTIFICATIONS (THÔNG BÁO)

### 7.1. Lấy danh sách thông báo
- **Method:** `GET`
- **Endpoint:** `/notifications`
- **Mô tả:** Lấy danh sách thông báo của user
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `page` (int, default: 1)
  - `limit` (int, default: 20)
  - `is_read` (boolean, optional): Lọc theo trạng thái đã đọc
- **Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Cảnh báo vượt ngân sách",
      "message": "Bạn đã sử dụng 85% ngân sách ăn uống",
      "notification_type": "warning",
      "priority": "high",
      "is_read": false,
      "action_url": "/budgets/1",
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 20
}
```

### 7.2. Đánh dấu đã đọc
- **Method:** `POST`
- **Endpoint:** `/notifications/:id/read`
- **Mô tả:** Đánh dấu thông báo đã đọc
- **Headers:** `Authorization: Bearer <token>`
- **Path Parameters:**
  - `id` (uint64): ID của notification
- **Response (200):**
```json
{
  "message": "Notification marked as read"
}
```

---

## 8. NOTIFICATION PREFERENCES (TÙY CHỈNH THÔNG BÁO)

### 8.1. Lấy preferences
- **Method:** `GET`
- **Endpoint:** `/notification-preferences`
- **Mô tả:** Lấy cài đặt thông báo của user
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):** Notification preferences object

### 8.2. Cập nhật preferences
- **Method:** `PUT`
- **Endpoint:** `/notification-preferences`
- **Mô tả:** Cập nhật cài đặt thông báo
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:** Preferences object
- **Response (200):** Updated preferences

### 8.3. Lấy summary preferences
- **Method:** `GET`
- **Endpoint:** `/notification-preferences/summary`
- **Mô tả:** Lấy tóm tắt cài đặt thông báo
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):** Summary object

### 8.4. Reset về mặc định
- **Method:** `POST`
- **Endpoint:** `/notification-preferences/reset`
- **Mô tả:** Reset preferences về mặc định
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "message": "Preferences reset to defaults"
}
```

### 8.5. Lấy danh sách channels được bật
- **Method:** `GET`
- **Endpoint:** `/notification-preferences/channels`
- **Mô tả:** Lấy danh sách notification channels đang được bật
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "channels": ["web", "email", "telegram"]
}
```

### 8.6. Test notification
- **Method:** `POST`
- **Endpoint:** `/notification-preferences/test`
- **Mô tả:** Gửi thông báo test để kiểm tra cài đặt
- **Headers:** `Authorization: Bearer <token>`
- **Response (200):**
```json
{
  "message": "Test notification sent"
}
```

---

## 9. AI ENDPOINTS

### 9.1. Đề xuất category
- **Method:** `POST`
- **Endpoint:** `/ai/suggest-category`
- **Mô tả:** AI đề xuất category cho giao dịch từ câu lệnh tự nhiên
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "message": "tôi vừa ăn bún bò 50k",
  "amount": 50000,
  "description": "ăn bún bò"
}
```
- **Response (200):**
```json
{
  "category_id": 5,
  "category_name": "Ăn uống",
  "amount": 50000,
  "description": "ăn bún bò",
  "confidence": 0.95,
  "suggested_date": "2024-01-15"
}
```

### 9.2. AI Chat
- **Method:** `POST`
- **Endpoint:** `/ai/chat`
- **Mô tả:** Chat với AI để hỏi đáp về tài chính
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "message": "Tháng này tôi tiêu bao nhiêu cho ăn uống?",
  "conversation_id": "optional-conversation-id"
}
```
- **Response (200):**
```json
{
  "response": "Tháng này bạn đã chi 2,500,000 VND cho ăn uống, chiếm 25% tổng chi tiêu.",
  "conversation_id": "conv-123",
  "sources": [
    {
      "type": "transactions",
      "data": {
        "total": 2500000,
        "count": 45
      }
    }
  ]
}
```

---

## 10. ANALYTICS (PHÂN TÍCH)

### 10.1. Dashboard Analytics
- **Method:** `GET`
- **Endpoint:** `/analytics/dashboard`
- **Mô tả:** Lấy dữ liệu analytics cho dashboard
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `start_date` (date, optional): Ngày bắt đầu
  - `end_date` (date, optional): Ngày kết thúc
- **Response (200):**
```json
{
  "data": {
    "period": {
      "start_date": "2024-01-01",
      "end_date": "2024-01-31"
    },
    "summary": {
      "total_income": 10000000,
      "total_expense": 7500000,
      "net_savings": 2500000,
      "savings_rate": 25.0
    },
    "spending_by_category": [
      {
        "category_id": 5,
        "category_name": "Ăn uống",
        "amount": 2000000,
        "percentage": 26.67
      }
    ],
    "trends": {
      "daily": [
        {"date": "2024-01-01", "income": 0, "expense": 50000}
      ],
      "weekly": [],
      "monthly": []
    },
    "top_categories": [
      {
        "category_id": 5,
        "category_name": "Ăn uống",
        "total": 2000000,
        "count": 30
      }
    ],
    "budget_status": {
      "total_budgets": 5,
      "over_budget": 1,
      "at_risk": 2,
      "on_track": 2
    }
  }
}
```

### 10.2. Category Spending
- **Method:** `GET`
- **Endpoint:** `/analytics/category-spending`
- **Mô tả:** Lấy chi tiết chi tiêu theo category
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `start_date` (date, optional)
  - `end_date` (date, optional)
  - `category_id` (uint64, optional)
- **Response (200):**
```json
{
  "data": [
    {
      "category_id": 5,
      "category_name": "Ăn uống",
      "total_amount": 2000000,
      "transaction_count": 30,
      "average_amount": 66666.67,
      "percentage": 26.67
    }
  ]
}
```

### 10.3. Spending Patterns
- **Method:** `GET`
- **Endpoint:** `/analytics/spending-patterns`
- **Mô tả:** Phân tích pattern chi tiêu
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `period` (string: daily/weekly/monthly, default: monthly)
  - `months` (int, default: 6): Số tháng để phân tích
- **Response (200):**
```json
{
  "data": {
    "patterns": [
      {
        "category_id": 5,
        "category_name": "Ăn uống",
        "trend": "increasing",
        "average_monthly": 2000000,
        "variance": 500000
      }
    ],
    "insights": [
      "Chi tiêu ăn uống tăng 15% so với tháng trước",
      "Chi tiêu giao thông ổn định"
    ]
  }
}
```

### 10.4. Anomalies
- **Method:** `GET`
- **Endpoint:** `/analytics/anomalies`
- **Mô tả:** Lấy danh sách giao dịch bất thường
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `start_date` (date, optional)
  - `end_date` (date, optional)
- **Response (200):**
```json
{
  "data": [
    {
      "transaction_id": 123,
      "amount": 5000000,
      "description": "Mua laptop",
      "category_id": 3,
      "category_name": "Mua sắm",
      "anomaly_score": 0.85,
      "reason": "Amount significantly higher than average for this category",
      "detected_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

### 10.5. Predictions
- **Method:** `GET`
- **Endpoint:** `/analytics/predictions`
- **Mô tả:** Lấy dự đoán chi tiêu
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `period` (string: next_month/next_quarter/next_year, default: next_month)
- **Response (200):**
```json
{
  "data": {
    "period": {
      "start_date": "2024-02-01",
      "end_date": "2024-02-29"
    },
    "predictions": {
      "total_expense": 7800000,
      "confidence": 0.82,
      "by_category": [
        {
          "category_id": 5,
          "category_name": "Ăn uống",
          "predicted_amount": 2100000,
          "confidence": 0.85
        }
      ]
    },
    "insights": [
      "Dự đoán chi tiêu tháng tới tăng 4% so với tháng này",
      "Chi tiêu ăn uống có xu hướng tăng"
    ],
    "model_version": "v1.2",
    "generated_at": "2024-01-15T10:00:00Z"
  }
}
```

---

## ERROR RESPONSES

Tất cả API có thể trả về các lỗi sau:

### 400 Bad Request
```json
{
  "error": "Invalid request",
  "message": "Validation error details"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": "Forbidden",
  "message": "You don't have permission to access this resource"
}
```

### 404 Not Found
```json
{
  "error": "Not found",
  "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "message": "An unexpected error occurred"
}
```

---

## NOTES

1. Tất cả dates sử dụng format ISO 8601: `YYYY-MM-DD` hoặc `YYYY-MM-DDTHH:mm:ssZ`
2. Tất cả amounts là số thực (float), đơn vị theo currency của user
3. Pagination: mặc định page=1, limit=20, tối đa limit=100
4. Rate limiting: 100 requests/phút cho mỗi user
5. CORS: Chỉ cho phép từ origins được cấu hình


