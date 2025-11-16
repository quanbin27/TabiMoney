# CHI TIẾT CÁC BẢNG DATABASE - HỆ THỐNG TABIMONEY

## 1. BẢNG users

**Mô tả:** Lưu trữ thông tin người dùng chính

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất của user |
| email | VARCHAR(255) | UNIQUE, NOT NULL | Email đăng nhập (duy nhất) |
| username | VARCHAR(100) | UNIQUE, NOT NULL | Tên đăng nhập (duy nhất) |
| password_hash | VARCHAR(255) | NOT NULL | Mật khẩu đã hash bằng bcrypt |
| first_name | VARCHAR(100) | NULL | Tên |
| last_name | VARCHAR(100) | NULL | Họ |
| phone | VARCHAR(20) | NULL | Số điện thoại |
| avatar_url | VARCHAR(500) | NULL | URL ảnh đại diện |
| is_verified | BOOLEAN | DEFAULT FALSE | Trạng thái xác thực email |
| verification_token | VARCHAR(255) | NULL | Token để xác thực email |
| reset_token | VARCHAR(255) | NULL | Token để reset password |
| reset_token_expires_at | TIMESTAMP | NULL | Thời gian hết hạn reset token |
| last_login_at | TIMESTAMP | NULL | Thời gian đăng nhập cuối cùng |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Thời gian cập nhật |
| deleted_at | TIMESTAMP | NULL | Soft delete timestamp |

**Indexes:**
- PRIMARY KEY: `id`
- UNIQUE: `email`, `username`
- INDEX: `idx_email`, `idx_username`, `idx_created_at`

**Quan hệ:**
- 1:1 với `user_profiles`
- 1:N với `transactions`, `categories`, `financial_goals`, `budgets`, `notifications`, `ai_analysis`
- 1:1 với `telegram_accounts`

---

## 2. BẢNG user_profiles

**Mô tả:** Lưu trữ thông tin cấu hình tài chính của user

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, UNIQUE, NOT NULL | ID của user (1:1) |
| monthly_income | DECIMAL(15,2) | DEFAULT 0.00 | Thu nhập hàng tháng |
| currency | VARCHAR(3) | DEFAULT 'VND' | Đơn vị tiền tệ (VND, USD, ...) |
| timezone | VARCHAR(50) | DEFAULT 'Asia/Ho_Chi_Minh' | Múi giờ |
| language | VARCHAR(5) | DEFAULT 'vi' | Ngôn ngữ (vi, en, ...) |
| notification_settings | JSON | NULL | Cài đặt thông báo (JSON) |
| ai_settings | JSON | NULL | Cài đặt AI (JSON) |
| telegram_enabled | BOOLEAN | DEFAULT FALSE | Bật tích hợp Telegram |
| telegram_notifications | BOOLEAN | DEFAULT TRUE | Bật thông báo Telegram |
| telegram_language | VARCHAR(5) | DEFAULT 'vi' | Ngôn ngữ Telegram |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Thời gian cập nhật |

**Indexes:**
- PRIMARY KEY: `id`
- UNIQUE: `unique_user_profile (user_id)`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)

**Quan hệ:**
- N:1 với `users` (mỗi user có 1 profile)

---

## 3. BẢNG categories

**Mô tả:** Lưu trữ danh mục chi tiêu (system + user custom)

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NULL | ID của user (NULL = system category) |
| name | VARCHAR(100) | NOT NULL | Tên danh mục (tiếng Việt) |
| name_en | VARCHAR(100) | NULL | Tên danh mục (tiếng Anh) |
| description | TEXT | NULL | Mô tả |
| icon | VARCHAR(50) | NULL | Icon name (Material Icons) |
| color | VARCHAR(7) | NULL | Màu sắc (HEX) |
| parent_id | BIGINT UNSIGNED | FOREIGN KEY → categories.id, NULL | ID danh mục cha (hierarchical) |
| is_system | BOOLEAN | DEFAULT FALSE | Là danh mục hệ thống |
| is_active | BOOLEAN | DEFAULT TRUE | Đang hoạt động |
| sort_order | INT | DEFAULT 0 | Thứ tự sắp xếp |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Thời gian cập nhật |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- FOREIGN KEY: `parent_id` → `categories.id` (ON DELETE SET NULL)
- INDEX: `idx_user_id`, `idx_parent_id`, `idx_is_system`

**Quan hệ:**
- N:1 với `users` (user_id có thể NULL cho system categories)
- N:1 với `categories` (self-referential, parent-child)
- 1:N với `transactions`, `budgets`

**Danh mục hệ thống mặc định:**
- Ăn uống (Food & Dining)
- Giao thông (Transportation)
- Mua sắm (Shopping)
- Giải trí (Entertainment)
- Y tế (Healthcare)
- Học tập (Education)
- Tiết kiệm (Savings)
- Thu nhập (Income)
- Khác (Other)

---

## 4. BẢNG transactions

**Mô tả:** Lưu trữ tất cả giao dịch thu/chi của user

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID của user |
| category_id | BIGINT UNSIGNED | FOREIGN KEY → categories.id, NOT NULL | ID danh mục |
| amount | DECIMAL(15,2) | NOT NULL | Số tiền (> 0) |
| description | TEXT | NULL | Mô tả giao dịch |
| transaction_type | ENUM | NOT NULL | Loại: 'income', 'expense', 'transfer' |
| transaction_date | DATE | NOT NULL | Ngày giao dịch |
| transaction_time | TIME | NULL | Giờ giao dịch |
| location | VARCHAR(200) | NULL | Địa điểm |
| tags | JSON | NULL | Tags (array of strings) |
| metadata | JSON | NULL | Dữ liệu bổ sung (payment method, etc.) |
| is_recurring | BOOLEAN | DEFAULT FALSE | Giao dịch định kỳ |
| recurring_pattern | VARCHAR(50) | NULL | Pattern: 'daily', 'weekly', 'monthly', 'yearly' |
| parent_transaction_id | BIGINT UNSIGNED | FOREIGN KEY → transactions.id, NULL | ID giao dịch cha (cho recurring) |
| ai_confidence | DECIMAL(3,2) | NULL | Độ tin cậy AI (0.00-1.00) |
| ai_suggested_category_id | BIGINT UNSIGNED | FOREIGN KEY → categories.id, NULL | Category AI đề xuất |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Thời gian cập nhật |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- FOREIGN KEY: `category_id` → `categories.id` (ON DELETE RESTRICT)
- FOREIGN KEY: `parent_transaction_id` → `transactions.id` (ON DELETE SET NULL)
- FOREIGN KEY: `ai_suggested_category_id` → `categories.id` (ON DELETE SET NULL)
- INDEX: `idx_user_id`, `idx_category_id`, `idx_transaction_date`, `idx_transaction_type`, `idx_amount`, `idx_created_at`
- COMPOSITE: `idx_transactions_user_date`, `idx_transactions_user_type`

**Quan hệ:**
- N:1 với `users`
- N:1 với `categories`
- N:1 với `transactions` (self-referential, parent-child cho recurring)

---

## 5. BẢNG financial_goals

**Mô tả:** Lưu trữ mục tiêu tài chính của user

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID của user |
| title | VARCHAR(200) | NOT NULL | Tiêu đề mục tiêu |
| description | TEXT | NULL | Mô tả chi tiết |
| target_amount | DECIMAL(15,2) | NOT NULL | Số tiền mục tiêu (> 0) |
| current_amount | DECIMAL(15,2) | DEFAULT 0.00 | Số tiền hiện tại |
| target_date | DATE | NULL | Ngày đạt mục tiêu |
| goal_type | ENUM | DEFAULT 'savings' | Loại: 'savings', 'debt_payment', 'investment', 'purchase', 'other' |
| priority | ENUM | DEFAULT 'medium' | Độ ưu tiên: 'low', 'medium', 'high', 'urgent' |
| is_achieved | BOOLEAN | DEFAULT FALSE | Đã đạt mục tiêu |
| achieved_at | TIMESTAMP | NULL | Thời gian đạt mục tiêu |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Thời gian cập nhật |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- INDEX: `idx_user_id`, `idx_goal_type`, `idx_target_date`

**Quan hệ:**
- N:1 với `users`

**Tính toán:**
- `progress` = (current_amount / target_amount) * 100 (calculated field, không lưu trong DB)

---

## 6. BẢNG budgets

**Mô tả:** Lưu trữ ngân sách của user

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID của user |
| category_id | BIGINT UNSIGNED | FOREIGN KEY → categories.id, NULL | ID danh mục (NULL = tổng ngân sách) |
| name | VARCHAR(200) | NOT NULL | Tên ngân sách |
| amount | DECIMAL(15,2) | NOT NULL | Số tiền ngân sách (> 0) |
| period | ENUM | DEFAULT 'monthly' | Chu kỳ: 'weekly', 'monthly', 'yearly' |
| start_date | DATE | NOT NULL | Ngày bắt đầu |
| end_date | DATE | NOT NULL | Ngày kết thúc (phải > start_date) |
| is_active | BOOLEAN | DEFAULT TRUE | Đang hoạt động |
| alert_threshold | DECIMAL(5,2) | DEFAULT 80.00 | Ngưỡng cảnh báo (%) |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Thời gian cập nhật |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- FOREIGN KEY: `category_id` → `categories.id` (ON DELETE CASCADE)
- INDEX: `idx_user_id`, `idx_category_id`, `idx_period`, `idx_start_date`

**Quan hệ:**
- N:1 với `users`
- N:1 với `categories` (có thể NULL)

**Tính toán (calculated fields, không lưu trong DB):**
- `spent_amount`: Tổng chi tiêu trong khoảng thời gian
- `remaining_amount`: amount - spent_amount
- `usage_percentage`: (spent_amount / amount) * 100

---

## 7. BẢNG notifications

**Mô tả:** Lưu trữ thông báo cho user

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID của user |
| title | VARCHAR(200) | NOT NULL | Tiêu đề thông báo |
| message | TEXT | NOT NULL | Nội dung thông báo |
| notification_type | ENUM | NOT NULL | Loại: 'info', 'warning', 'success', 'error', 'reminder' |
| priority | ENUM | DEFAULT 'medium' | Độ ưu tiên: 'low', 'medium', 'high', 'urgent' |
| is_read | BOOLEAN | DEFAULT FALSE | Đã đọc |
| read_at | TIMESTAMP | NULL | Thời gian đọc |
| action_url | VARCHAR(500) | NULL | URL để chuyển đến khi click |
| metadata | JSON | NULL | Dữ liệu bổ sung |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- INDEX: `idx_user_id`, `idx_is_read`, `idx_created_at`

**Quan hệ:**
- N:1 với `users`

---

## 8. BẢNG ai_analysis

**Mô tả:** Lưu trữ kết quả phân tích AI

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID của user |
| analysis_type | ENUM | NOT NULL | Loại: 'expense_prediction', 'anomaly_detection', 'category_suggestion', 'spending_pattern', 'goal_analysis' |
| data | JSON | NOT NULL | Kết quả phân tích (JSON) |
| confidence_score | DECIMAL(3,2) | NULL | Độ tin cậy (0.00-1.00) |
| model_version | VARCHAR(50) | NULL | Phiên bản model AI |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- INDEX: `idx_user_id`, `idx_analysis_type`, `idx_created_at`

**Quan hệ:**
- N:1 với `users`

**Ví dụ data JSON:**
- `expense_prediction`: `{"total_expense": 7800000, "by_category": [...], "confidence": 0.82}`
- `anomaly_detection`: `{"transaction_id": 123, "score": 0.85, "reason": "..."}`

---

## 9. BẢNG ai_feedback

**Mô tả:** Lưu trữ feedback của user cho AI để học hỏi

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID của user |
| transaction_id | BIGINT UNSIGNED | FOREIGN KEY → transactions.id, NULL | ID giao dịch liên quan |
| feedback_type | ENUM | NOT NULL | Loại: 'category_correct', 'category_incorrect', 'prediction_accurate', 'prediction_inaccurate', 'suggestion_helpful', 'suggestion_not_helpful' |
| original_prediction | JSON | NULL | Dự đoán ban đầu của AI |
| user_correction | JSON | NULL | Sửa đổi của user |
| feedback_text | TEXT | NULL | Feedback dạng text |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- FOREIGN KEY: `transaction_id` → `transactions.id` (ON DELETE SET NULL)
- INDEX: `idx_user_id`, `idx_feedback_type`

**Quan hệ:**
- N:1 với `users`
- N:1 với `transactions` (có thể NULL)

---

## 10. BẢNG user_sessions

**Mô tả:** Lưu trữ session JWT của user

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID của user |
| token_hash | VARCHAR(255) | NOT NULL | Hash của access token |
| refresh_token_hash | VARCHAR(255) | NOT NULL | Hash của refresh token |
| expires_at | TIMESTAMP | NOT NULL | Thời gian hết hạn access token |
| refresh_expires_at | TIMESTAMP | NOT NULL | Thời gian hết hạn refresh token |
| user_agent | TEXT | NULL | User agent của client |
| ip_address | VARCHAR(45) | NULL | IP address |
| is_active | BOOLEAN | DEFAULT TRUE | Session đang hoạt động |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |

**Indexes:**
- PRIMARY KEY: `id`
- FOREIGN KEY: `user_id` → `users.id` (ON DELETE CASCADE)
- INDEX: `idx_user_id`, `idx_token_hash`, `idx_expires_at`

**Quan hệ:**
- N:1 với `users`

---

## 11. BẢNG telegram_accounts

**Mô tả:** Liên kết tài khoản Telegram với web account

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| telegram_user_id | BIGINT | UNIQUE, NOT NULL | Telegram User ID |
| web_user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID user trên web |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Thời gian cập nhật |

**Indexes:**
- PRIMARY KEY: `id`
- UNIQUE: `telegram_user_id`
- FOREIGN KEY: `web_user_id` → `users.id` (ON DELETE CASCADE)
- INDEX: `idx_telegram_user_id`, `idx_web_user_id`

**Quan hệ:**
- N:1 với `users` (1:1 thực tế)

---

## 12. BẢNG telegram_link_codes

**Mô tả:** Mã liên kết tạm thời để liên kết Telegram

| Tên cột | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
|---------|--------------|-----------|---------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | ID duy nhất |
| code | VARCHAR(16) | UNIQUE, NOT NULL | Mã liên kết (16 ký tự) |
| telegram_user_id | BIGINT | NULL | Telegram User ID (khi user nhập code) |
| web_user_id | BIGINT UNSIGNED | FOREIGN KEY → users.id, NOT NULL | ID user trên web |
| expires_at | TIMESTAMP | NOT NULL | Thời gian hết hạn (10 phút) |
| used_at | TIMESTAMP | NULL | Thời gian sử dụng |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Thời gian tạo |

**Indexes:**
- PRIMARY KEY: `id`
- UNIQUE: `code`
- FOREIGN KEY: `web_user_id` → `users.id` (ON DELETE CASCADE)
- INDEX: `idx_code`, `idx_telegram_user_id`, `idx_expires_at`

**Quan hệ:**
- N:1 với `users`

---

## VIEWS (VIEWS)

### 13. VIEW user_monthly_summary

**Mô tả:** Tổng hợp chi tiêu theo tháng cho mỗi user

**Các cột:**
- `user_id`: ID user
- `username`: Tên đăng nhập
- `year`: Năm
- `month`: Tháng
- `total_income`: Tổng thu nhập
- `total_expense`: Tổng chi tiêu
- `income_count`: Số giao dịch thu
- `expense_count`: Số giao dịch chi

**Sử dụng:** Dashboard analytics, báo cáo

---

### 14. VIEW category_spending

**Mô tả:** Tổng hợp chi tiêu theo category

**Các cột:**
- `user_id`: ID user
- `category_name`: Tên category
- `icon`: Icon
- `color`: Màu sắc
- `total_amount`: Tổng số tiền
- `transaction_count`: Số giao dịch
- `avg_amount`: Số tiền trung bình

**Sử dụng:** Biểu đồ pie chart, phân tích category

---

## RÀNG BUỘC VÀ QUY TẮC NGHIỆP VỤ

1. **Soft Delete:** Bảng `users` sử dụng soft delete (`deleted_at`)
2. **Cascade Delete:** Khi xóa user, tất cả dữ liệu liên quan được xóa (CASCADE)
3. **Restrict Delete:** Không thể xóa category nếu đang có transactions sử dụng (RESTRICT)
4. **Unique Constraints:** 
   - Email và username phải unique
   - Mỗi user chỉ có 1 profile
   - Telegram user_id phải unique
5. **Default Values:** 
   - `is_verified = FALSE`
   - `is_active = TRUE` cho categories
   - `is_read = FALSE` cho notifications
6. **Timestamps:** Tất cả bảng có `created_at`, một số có `updated_at`

---

## INDEXES TỐI ƯU

1. **Composite Indexes:**
   - `(user_id, transaction_date)` cho transactions
   - `(user_id, transaction_type)` cho transactions
2. **Single Column Indexes:** Tất cả foreign keys và các cột thường query
3. **Unique Indexes:** Email, username, code, telegram_user_id

---

## MIGRATION NOTES

- Database: MySQL 8.0+
- Character Set: utf8mb4
- Collation: utf8mb4_unicode_ci
- Engine: InnoDB (hỗ trợ transactions và foreign keys)

