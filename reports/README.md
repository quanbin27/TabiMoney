# TÀI LIỆU THIẾT KẾ HỆ THỐNG TABIMONEY

Thư mục này chứa toàn bộ tài liệu thiết kế hệ thống chi tiết cho **TabiMoney - AI-Powered Personal Finance Management System**.

## 📋 DANH SÁCH TÀI LIỆU

### 1. Use Case Chi Tiết
**File:** `01_USE_CASE_DETAILED.md`

Tài liệu mô tả chi tiết 20 Use Case chính của hệ thống, bao gồm:
- Tác nhân chính
- Điều kiện trước/sau
- Luồng sự kiện chính và phụ
- Đảm bảo tối thiểu

**Các Use Case chính:**
- Authentication & User Management (UC-001, UC-002)
- Transaction Management (UC-003 đến UC-007)
- Financial Goals (UC-008, UC-009)
- Budget Management (UC-010, UC-018, UC-019)
- Analytics & Reporting (UC-011, UC-020)
- AI Features (UC-012, UC-013, UC-014)
- Notifications (UC-015, UC-018)
- Telegram Integration (UC-016, UC-017)

---

### 2. Sơ Đồ Use Case
**File:** `02_USE_CASE_DIAGRAM.drawio`

Sơ đồ Use Case Diagram dạng XML Draw.io, có thể mở trực tiếp trong Draw.io để xem và chỉnh sửa.

**Các thành phần:**
- 4 Tác nhân: User, Telegram Bot, AI Agent, System
- 20 Use Cases được phân loại theo màu sắc
- Các mối quan hệ giữa tác nhân và use cases

---

### 3. Sơ Đồ Tuần Tự (Sequence Diagrams)
**File:** `03_SEQUENCE_DIAGRAMS.drawio`

Chứa các sơ đồ tuần tự cho các use case quan trọng:
- **Sequence_NLU_Transaction:** Quy trình nhập giao dịch bằng NLU
- **Sequence_Login:** Quy trình đăng nhập

**Các thành phần:**
- User, Frontend, Backend API, AI Service, Gemini API, Database
- Các message flows chi tiết
- Activation boxes

---

### 4. Sơ Đồ Hoạt Động (Activity Diagrams)
**File:** `04_ACTIVITY_DIAGRAMS.drawio`

Sơ đồ Activity Diagram mô tả quy trình nghiệp vụ:
- **Activity_TransactionEntry:** Quy trình nhập giao dịch (thủ công và NLU)

**Các thành phần:**
- Start/End nodes
- Decision nodes
- Action nodes
- Flow paths

---

### 5. Sơ Đồ Trạng Thái (State Diagrams)
**File:** `05_STATE_DIAGRAMS.drawio`

Sơ đồ State Diagram cho các đối tượng quan trọng:
- **State_FinancialGoal:** Trạng thái mục tiêu tài chính (Created → Active → Achieved/Cancelled)
- **State_Transaction:** Trạng thái giao dịch (Draft → Pending → Completed → Updated/Deleted)

---

### 6. Danh Sách API Chi Tiết
**File:** `06_API_LIST.md`

Tài liệu đầy đủ về tất cả API endpoints, bao gồm:

**10 nhóm API:**
1. Authentication & User Management (9 endpoints)
2. Telegram Integration (4 endpoints)
3. Transactions (4 endpoints)
4. Categories (4 endpoints)
5. Financial Goals (5 endpoints)
6. Budgets (7 endpoints)
7. Notifications (2 endpoints)
8. Notification Preferences (6 endpoints)
9. AI Endpoints (2 endpoints)
10. Analytics (5 endpoints)

**Mỗi API bao gồm:**
- Method (GET, POST, PUT, DELETE)
- Endpoint URL
- Mô tả chức năng
- Request/Response examples
- Error handling

---

### 7. Sơ Đồ ERD (Entity-Relationship Diagram)
**File:** `07_ERD_DIAGRAM.drawio`

Sơ đồ ERD đầy đủ của database, bao gồm:

**12 bảng chính:**
- users
- user_profiles
- categories
- transactions
- financial_goals
- budgets
- notifications
- ai_analyses
- chat_messages
- user_sessions
- telegram_accounts
- telegram_link_codes

**Các quan hệ:**
- 1:1, 1:N, N:1
- Foreign keys
- Self-referential (categories, transactions)

---

### 8. Chi Tiết Các Bảng Database
**File:** `08_DATABASE_TABLES.md`

Tài liệu chi tiết về cấu trúc database:

**Mỗi bảng bao gồm:**
- Mô tả chức năng
- Bảng chi tiết các cột (tên, kiểu dữ liệu, ràng buộc, ý nghĩa)
- Indexes
- Quan hệ với các bảng khác
- Ràng buộc và quy tắc nghiệp vụ

**Views:**
- user_monthly_summary
- category_spending

---

### 9. Cơ Chế Xử Lý & Thuật Toán
**File:** `09_ALGORITHMS_AND_PROCESSING.md`

Tài liệu chi tiết về các thuật toán và cơ chế xử lý:

**8 thuật toán chính:**
1. **NLU Processing:** Xử lý ngôn ngữ tự nhiên để trích xuất thông tin giao dịch
2. **Anomaly Detection:** Phát hiện giao dịch bất thường (Z-Score + Isolation Forest)
3. **Expense Prediction:** Dự đoán chi tiêu (Linear Regression + Time Series)
4. **Budget Suggestions:** Đề xuất ngân sách tự động
5. **Budget Alerts:** Kiểm tra và cảnh báo vượt ngân sách
6. **Dashboard Analytics:** Tính toán các chỉ số tài chính
7. **Cache Strategy:** Chiến lược cache với Redis
8. **Error Handling:** Xử lý lỗi và edge cases

**Mỗi thuật toán bao gồm:**
- Mô tả chi tiết
- Luồng xử lý
- Code mẫu (pseudo-code)
- Ví dụ minh họa

---

## 🎯 CÁCH SỬ DỤNG

### Xem sơ đồ Draw.io
1. Truy cập https://app.diagrams.net (hoặc Draw.io desktop)
2. File → Open from → Device
3. Chọn file `.drawio` trong thư mục `reports/`
4. Sơ đồ sẽ hiển thị và có thể chỉnh sửa

### Đọc tài liệu Markdown
- Mở file `.md` bằng bất kỳ Markdown viewer nào
- Hoặc xem trực tiếp trên GitHub/GitLab

---

## 📊 TỔNG QUAN HỆ THỐNG

**TabiMoney** là hệ thống quản lý chi tiêu cá nhân thông minh với các tính năng:

### Tính năng chính:
- ✅ Nhập giao dịch thủ công và bằng NLU
- ✅ Quản lý mục tiêu tài chính
- ✅ Quản lý ngân sách với cảnh báo tự động
- ✅ Dashboard analytics với biểu đồ
- ✅ AI Chat hỏi đáp tài chính
- ✅ Phát hiện bất thường trong chi tiêu
- ✅ Dự đoán chi tiêu tháng tới
- ✅ Tích hợp Telegram Bot
- ✅ Thông báo real-time

### Kiến trúc:
- **Backend:** Golang + Echo Framework
- **Frontend:** Vue.js 3 + Vuetify
- **AI Service:** Python + Google Gemini API
- **Database:** MySQL 8.0
- **Cache:** Redis 7.0
- **Telegram Bot:** Python

---

## 📝 GHI CHÚ

- Tất cả sơ đồ XML có thể mở trực tiếp trong Draw.io
- Tài liệu được viết bằng tiếng Việt để dễ hiểu
- Code examples sử dụng pseudo-code (Python-like)
- Tất cả API endpoints đều có examples request/response

---

## 🔄 CẬP NHẬT

Tài liệu này được tạo vào: **2024-01-15**

Khi hệ thống có thay đổi, vui lòng cập nhật các tài liệu tương ứng.

---

## 📞 LIÊN HỆ

Nếu có câu hỏi hoặc cần làm rõ về tài liệu, vui lòng liên hệ team phát triển.

---

**Chúc bạn sử dụng tài liệu hiệu quả! 🚀**



