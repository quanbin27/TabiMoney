# TÀI LIỆU USE CASE CHI TIẾT - HỆ THỐNG TABIMONEY

## 1. TỔNG QUAN HỆ THỐNG

**Tên hệ thống:** TabiMoney - AI-Powered Personal Finance Management System

**Mô tả:** Hệ thống quản lý chi tiêu cá nhân thông minh, tích hợp AI Agent để cung cấp phân tích tài chính, dự đoán chi tiêu và tư vấn cá nhân hóa qua Web App và Telegram Bot.

**Tác nhân chính:**
- **User (Người dùng):** Người sử dụng hệ thống để quản lý tài chính cá nhân
- **AI Agent:** Hệ thống AI tự động xử lý và phân tích dữ liệu
- **Telegram Bot:** Bot tự động xử lý yêu cầu từ Telegram
- **System:** Hệ thống tự động (thông báo, cảnh báo)

---

## 2. DANH SÁCH USE CASE CHI TIẾT

### UC-001: Đăng ký tài khoản

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User chưa có tài khoản trong hệ thống
- User có kết nối internet

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Hệ thống không tạo tài khoản trùng lặp
- Thông tin người dùng không bị lộ

**Điều kiện sau (Postcondition):**
- Tài khoản được tạo thành công
- User nhận được email xác thực (nếu có)
- User có thể đăng nhập vào hệ thống

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang đăng ký
2. User nhập thông tin: email, username, password, first_name, last_name (tùy chọn)
3. Hệ thống validate dữ liệu đầu vào
4. Hệ thống kiểm tra email và username chưa tồn tại
5. Hệ thống hash password bằng bcrypt
6. Hệ thống tạo user record trong database
7. Hệ thống tạo user_profile mặc định
8. Hệ thống trả về thông báo thành công
9. User được chuyển đến trang đăng nhập

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Email đã tồn tại**
- 3a. Hệ thống phát hiện email đã tồn tại
- 3b. Hệ thống trả về lỗi "Email already exists"
- 3c. Use case kết thúc

**A2: Username đã tồn tại**
- 4a. Hệ thống phát hiện username đã tồn tại
- 4b. Hệ thống trả về lỗi "Username already exists"
- 4c. Use case kết thúc

**A3: Dữ liệu không hợp lệ**
- 3a. Hệ thống phát hiện dữ liệu không hợp lệ (email format sai, password quá ngắn)
- 3b. Hệ thống trả về lỗi validation
- 3c. Use case kết thúc

---

### UC-002: Đăng nhập

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã có tài khoản trong hệ thống
- User chưa đăng nhập

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Mật khẩu không được lộ
- Session được quản lý an toàn

**Điều kiện sau (Postcondition):**
- User đã đăng nhập thành công
- JWT access token và refresh token được tạo
- User session được lưu trong database
- last_login_at được cập nhật

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang đăng nhập
2. User nhập email và password
3. Hệ thống validate dữ liệu đầu vào
4. Hệ thống tìm user theo email
5. Hệ thống so sánh password hash
6. Hệ thống tạo JWT access token (24h)
7. Hệ thống tạo refresh token (7 ngày)
8. Hệ thống lưu session vào database
9. Hệ thống cập nhật last_login_at
10. Hệ thống trả về access_token, refresh_token, user info
11. User được chuyển đến dashboard

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Email không tồn tại**
- 4a. Hệ thống không tìm thấy user
- 4b. Hệ thống trả về lỗi "Invalid credentials" (không tiết lộ email có tồn tại)
- 4c. Use case kết thúc

**A2: Password sai**
- 5a. Password hash không khớp
- 5b. Hệ thống trả về lỗi "Invalid credentials"
- 5c. Use case kết thúc

**A3: Tài khoản bị khóa**
- 4a. Hệ thống phát hiện tài khoản bị khóa (deleted_at != NULL)
- 4b. Hệ thống trả về lỗi "Account is disabled"
- 4c. Use case kết thúc

---

### UC-003: Nhập giao dịch thủ công

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- User có ít nhất một category

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Giao dịch được lưu chính xác
- Số tiền phải > 0

**Điều kiện sau (Postcondition):**
- Giao dịch được tạo trong database
- Dashboard được cập nhật
- Cache được làm mới

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang thêm giao dịch
2. User chọn loại giao dịch (income/expense/transfer)
3. User chọn category
4. User nhập số tiền
5. User nhập mô tả (tùy chọn)
6. User chọn ngày giao dịch
7. User nhập thông tin bổ sung: location, tags, metadata (tùy chọn)
8. User nhấn "Lưu"
9. Hệ thống validate dữ liệu
10. Hệ thống kiểm tra category tồn tại và thuộc về user
11. Hệ thống tạo transaction record
12. Hệ thống cập nhật cache dashboard
13. Hệ thống kiểm tra budget và tạo notification nếu vượt ngưỡng
14. Hệ thống trả về transaction đã tạo
15. User thấy giao dịch mới trong danh sách

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Category không tồn tại**
- 10a. Category không tồn tại hoặc không thuộc về user
- 10b. Hệ thống trả về lỗi "Category not found"
- 10c. Use case kết thúc

**A2: Số tiền không hợp lệ**
- 9a. Số tiền <= 0 hoặc không phải số
- 9b. Hệ thống trả về lỗi "Amount must be greater than 0"
- 9c. Use case kết thúc

**A3: Vượt budget**
- 13a. Hệ thống phát hiện chi tiêu vượt budget
- 13b. Hệ thống tạo notification cảnh báo
- 13c. Tiếp tục luồng chính

---

### UC-004: Nhập giao dịch bằng NLU (Natural Language Understanding)

**Tác nhân chính:** User, AI Agent

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- AI Service đang hoạt động
- Gemini API key hợp lệ

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Giao dịch được tạo chính xác
- AI confidence score được lưu

**Điều kiện sau (Postcondition):**
- Giao dịch được tạo từ câu lệnh tự nhiên
- AI confidence score được lưu
- User có thể chỉnh sửa nếu cần

**Luồng sự kiện chính (Main Flow):**
1. User nhập câu lệnh tự nhiên: "tôi vừa ăn bún bò 50k"
2. Hệ thống gửi request đến AI Service
3. AI Service gọi Gemini API để phân tích NLU
4. AI Service trích xuất: amount, category, description, date
5. AI Service tìm category phù hợp nhất
6. AI Service tính confidence score
7. AI Service trả về kết quả: {category_id, amount, description, confidence}
8. Hệ thống hiển thị preview giao dịch cho user xác nhận
9. User xác nhận hoặc chỉnh sửa
10. Hệ thống tạo transaction với ai_confidence và ai_suggested_category_id
11. Hệ thống cập nhật cache
12. Hệ thống trả về transaction đã tạo

**Luồng sự kiện phụ (Alternate Flow):**

**A1: AI không hiểu câu lệnh**
- 4a. AI không thể trích xuất đủ thông tin (thiếu amount hoặc category)
- 4b. AI Service trả về lỗi "Cannot parse transaction"
- 4c. Hệ thống yêu cầu user nhập lại hoặc nhập thủ công
- 4d. Use case kết thúc

**A2: Confidence score thấp**
- 6a. Confidence score < 0.5
- 6b. Hệ thống hiển thị cảnh báo và yêu cầu user xác nhận
- 6c. Tiếp tục luồng chính

**A3: Category không tìm thấy**
- 5a. AI không tìm thấy category phù hợp
- 5b. Hệ thống đề xuất các category gần nhất
- 5c. User chọn category thủ công
- 5d. Tiếp tục luồng chính

---

### UC-005: Xem danh sách giao dịch

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Chỉ hiển thị giao dịch của user hiện tại
- Dữ liệu được phân trang

**Điều kiện sau (Postcondition):**
- Danh sách giao dịch được hiển thị
- User có thể lọc và tìm kiếm

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang danh sách giao dịch
2. Hệ thống lấy user_id từ JWT token
3. Hệ thống đọc query parameters: page, limit, category_id, transaction_type, start_date, end_date, min_amount, max_amount, search
4. Hệ thống query database với filters
5. Hệ thống phân trang kết quả
6. Hệ thống trả về danh sách giao dịch với metadata (total, page, limit)
7. User thấy danh sách giao dịch

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Không có giao dịch**
- 5a. Không tìm thấy giao dịch nào
- 5b. Hệ thống trả về danh sách rỗng
- 5c. User thấy thông báo "Chưa có giao dịch"

**A2: Lọc không hợp lệ**
- 4a. Query parameters không hợp lệ (date format sai)
- 4b. Hệ thống bỏ qua filter không hợp lệ
- 4c. Tiếp tục với các filter hợp lệ

---

### UC-006: Cập nhật giao dịch

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- Giao dịch tồn tại và thuộc về user

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Chỉ user sở hữu mới được cập nhật
- Dữ liệu được validate

**Điều kiện sau (Postcondition):**
- Giao dịch được cập nhật
- Dashboard và cache được làm mới

**Luồng sự kiện chính (Main Flow):**
1. User chọn giao dịch cần cập nhật
2. User chỉnh sửa thông tin: amount, category, description, date, etc.
3. User nhấn "Lưu"
4. Hệ thống validate dữ liệu
5. Hệ thống kiểm tra giao dịch thuộc về user
6. Hệ thống cập nhật transaction record
7. Hệ thống cập nhật cache
8. Hệ thống kiểm tra budget và tạo notification nếu cần
9. Hệ thống trả về transaction đã cập nhật
10. User thấy giao dịch đã được cập nhật

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Giao dịch không tồn tại**
- 5a. Giao dịch không tồn tại
- 5b. Hệ thống trả về lỗi "Transaction not found"
- 5c. Use case kết thúc

**A2: Không có quyền**
- 5a. Giao dịch không thuộc về user
- 5b. Hệ thống trả về lỗi "Unauthorized"
- 5c. Use case kết thúc

---

### UC-007: Xóa giao dịch

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- Giao dịch tồn tại và thuộc về user

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Chỉ user sở hữu mới được xóa
- Dữ liệu được xóa an toàn

**Điều kiện sau (Postcondition):**
- Giao dịch được xóa khỏi database
- Dashboard và cache được cập nhật

**Luồng sự kiện chính (Main Flow):**
1. User chọn giao dịch cần xóa
2. User nhấn "Xóa"
3. Hệ thống hiển thị xác nhận
4. User xác nhận xóa
5. Hệ thống kiểm tra giao dịch thuộc về user
6. Hệ thống xóa transaction record
7. Hệ thống cập nhật cache
8. Hệ thống trả về thông báo thành công
9. Giao dịch biến mất khỏi danh sách

**Luồng sự kiện phụ (Alternate Flow):**

**A1: User hủy xóa**
- 4a. User nhấn "Hủy"
- 4b. Use case kết thúc, không có thay đổi

**A2: Giao dịch không tồn tại**
- 5a. Giao dịch không tồn tại
- 5b. Hệ thống trả về lỗi "Transaction not found"
- 5c. Use case kết thúc

---

### UC-008: Tạo mục tiêu tài chính

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Mục tiêu được lưu chính xác
- Target amount > 0

**Điều kiện sau (Postcondition):**
- Mục tiêu tài chính được tạo
- User có thể theo dõi tiến độ

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang quản lý mục tiêu
2. User nhấn "Tạo mục tiêu mới"
3. User nhập thông tin: title, description, target_amount, target_date, goal_type, priority
4. Hệ thống validate dữ liệu
5. Hệ thống tạo financial_goal record
6. Hệ thống tính toán progress (current_amount / target_amount)
7. Hệ thống trả về goal đã tạo
8. User thấy mục tiêu mới trong danh sách

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Target amount không hợp lệ**
- 4a. Target amount <= 0
- 4b. Hệ thống trả về lỗi "Target amount must be greater than 0"
- 4c. Use case kết thúc

**A2: Target date trong quá khứ**
- 4a. Target date < today
- 4b. Hệ thống cảnh báo nhưng vẫn cho phép tạo
- 4c. Tiếp tục luồng chính

---

### UC-009: Thêm tiền vào mục tiêu (Contribution)

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- Mục tiêu tồn tại và chưa đạt

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Số tiền được cập nhật chính xác
- Amount > 0

**Điều kiện sau (Postcondition):**
- current_amount được tăng lên
- Progress được cập nhật
- Nếu đạt mục tiêu, is_achieved = true

**Luồng sự kiện chính (Main Flow):**
1. User chọn mục tiêu cần thêm tiền
2. User nhấn "Thêm tiền"
3. User nhập số tiền và note (tùy chọn)
4. Hệ thống validate amount > 0
5. Hệ thống kiểm tra goal tồn tại và thuộc về user
6. Hệ thống cập nhật current_amount = current_amount + amount
7. Hệ thống tính progress = (current_amount / target_amount) * 100
8. Nếu current_amount >= target_amount:
   8a. Hệ thống set is_achieved = true
   8b. Hệ thống set achieved_at = now()
   8c. Hệ thống tạo notification chúc mừng
9. Hệ thống cập nhật goal record
10. Hệ thống trả về goal đã cập nhật
11. User thấy tiến độ được cập nhật

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Mục tiêu đã đạt**
- 5a. is_achieved = true
- 5b. Hệ thống cảnh báo "Goal already achieved"
- 5c. User vẫn có thể thêm tiền (over-achievement)
- 5d. Tiếp tục luồng chính

**A2: Số tiền không hợp lệ**
- 4a. Amount <= 0
- 4b. Hệ thống trả về lỗi "Amount must be greater than 0"
- 4c. Use case kết thúc

---

### UC-010: Tạo ngân sách (Budget)

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- User đã khai báo monthly_income (tùy chọn)

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Budget được lưu chính xác
- Amount > 0
- End date > Start date

**Điều kiện sau (Postcondition):**
- Budget được tạo
- Hệ thống sẽ theo dõi chi tiêu so với budget

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang quản lý ngân sách
2. User nhấn "Tạo ngân sách mới"
3. User nhập thông tin: name, amount, period (weekly/monthly/yearly), start_date, end_date, category_id (tùy chọn), alert_threshold
4. Hệ thống validate dữ liệu
5. Hệ thống kiểm tra end_date > start_date
6. Hệ thống tạo budget record
7. Hệ thống trả về budget đã tạo
8. User thấy budget mới trong danh sách

**Luồng sự kiện phụ (Alternate Flow):**

**A1: End date <= Start date**
- 5a. End date <= Start date
- 5b. Hệ thống trả về lỗi "End date must be after start date"
- 5c. Use case kết thúc

**A2: Category không tồn tại**
- 4a. Category_id được cung cấp nhưng không tồn tại
- 4b. Hệ thống trả về lỗi "Category not found"
- 4c. Use case kết thúc

---

### UC-011: Xem Dashboard Analytics

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Dữ liệu được cache để tối ưu performance
- Chỉ hiển thị dữ liệu của user

**Điều kiện sau (Postcondition):**
- Dashboard hiển thị các chỉ số tài chính
- User có cái nhìn tổng quan về tài chính

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang Dashboard
2. Hệ thống kiểm tra cache Redis
3. Nếu có cache và chưa hết hạn:
   3a. Hệ thống trả về dữ liệu từ cache
4. Nếu không có cache:
   4a. Hệ thống query database để tính toán:
       - Total income tháng này
       - Total expense tháng này
       - Net savings (income - expense)
       - Spending by category (pie chart data)
       - Spending trends (line chart data)
       - Top categories
       - Budget status
   4b. Hệ thống lưu kết quả vào cache (TTL: 5 phút)
5. Hệ thống trả về analytics data
6. Frontend hiển thị biểu đồ và số liệu

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Không có dữ liệu**
- 4a. User chưa có giao dịch nào
- 4b. Hệ thống trả về dữ liệu mặc định (0 cho tất cả)
- 4c. Frontend hiển thị "Chưa có dữ liệu"

---

### UC-012: AI Chat - Hỏi đáp tài chính

**Tác nhân chính:** User, AI Agent

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- AI Service đang hoạt động

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Câu trả lời dựa trên dữ liệu thực tế của user
- Privacy được đảm bảo

**Điều kiện sau (Postcondition):**
- User nhận được câu trả lời từ AI
- Lịch sử chat được lưu (tùy chọn)

**Luồng sự kiện chính (Main Flow):**
1. User mở chat interface
2. User nhập câu hỏi: "Tháng này tôi tiêu bao nhiêu cho ăn uống?"
3. Hệ thống gửi request đến AI Service với user_id và message
4. AI Service phân tích intent từ câu hỏi
5. AI Service query database để lấy dữ liệu liên quan:
   - Nếu hỏi về spending: query transactions
   - Nếu hỏi về goals: query financial_goals
   - Nếu hỏi về budget: query budgets
6. AI Service tạo context từ dữ liệu
7. AI Service gọi Gemini API với context và câu hỏi
8. Gemini trả về câu trả lời tự nhiên
9. AI Service format câu trả lời
10. Hệ thống trả về response cho user
11. User thấy câu trả lời trong chat

**Luồng sự kiện phụ (Alternate Flow):**

**A1: AI không hiểu câu hỏi**
- 4a. AI không thể xác định intent
- 4b. AI Service trả về "Tôi không hiểu câu hỏi. Bạn có thể hỏi về chi tiêu, mục tiêu, hoặc ngân sách."
- 4c. Use case kết thúc

**A2: Không có dữ liệu**
- 5a. Không tìm thấy dữ liệu liên quan
- 5b. AI Service trả về "Bạn chưa có dữ liệu về [topic]. Hãy thêm giao dịch để tôi có thể phân tích."

---

### UC-013: Phát hiện bất thường (Anomaly Detection)

**Tác nhân chính:** AI Agent, System

**Điều kiện trước (Precondition):**
- User đã có ít nhất 10 giao dịch
- AI Service đang hoạt động

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Chỉ phát hiện bất thường thực sự
- False positive rate thấp

**Điều kiện sau (Postcondition):**
- Giao dịch bất thường được đánh dấu
- Notification được gửi nếu cần

**Luồng sự kiện chính (Main Flow):**
1. User tạo giao dịch mới
2. Hệ thống trigger anomaly detection (async)
3. AI Service lấy lịch sử giao dịch của user (30 ngày gần nhất)
4. AI Service tính toán:
   - Mean và standard deviation của amount theo category
   - Pattern chi tiêu theo thời gian
5. AI Service so sánh giao dịch mới với pattern:
   - Nếu amount > mean + 2*std: FLAG
   - Nếu category khác pattern thường ngày: FLAG
   - Nếu thời gian khác pattern: FLAG
6. AI Service tính anomaly score (0-1)
7. Nếu score > 0.7:
   7a. Hệ thống tạo notification cảnh báo
   7b. Hệ thống gửi notification cho user
8. Giao dịch được đánh dấu (metadata)

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Không đủ dữ liệu**
- 3a. User có < 10 giao dịch
- 3b. AI Service bỏ qua anomaly detection
- 3c. Use case kết thúc

**A2: False positive**
- 7a. User xác nhận đây là giao dịch hợp lệ
- 7b. Hệ thống ghi nhận và điều chỉnh threshold

---

### UC-014: Dự đoán chi tiêu (Expense Prediction)

**Tác nhân chính:** AI Agent, User

**Điều kiện trước (Precondition):**
- User đã có ít nhất 3 tháng dữ liệu
- AI Service đang hoạt động

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Dự đoán dựa trên dữ liệu thực tế
- Confidence score được cung cấp

**Điều kiện sau (Postcondition):**
- Dự đoán chi tiêu tháng tới được tạo

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang "Dự đoán chi tiêu" hoặc hệ thống tự động chạy (cuối tháng)
2. AI Service lấy lịch sử giao dịch 3-6 tháng gần nhất
3. AI Service phân tích:
   - Trend chi tiêu theo tháng
   - Pattern theo category
   - Seasonal patterns
4. AI Service áp dụng ML model (Linear Regression hoặc LSTM)
5. AI Service dự đoán chi tiêu tháng tới:
   - Tổng chi tiêu
   - Chi tiêu theo category
   - Confidence score
6. Hệ thống hiển thị dự đoán cho user
7. User có thể xem chi tiết và so sánh với thực tế sau này

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Không đủ dữ liệu**
- 2a. User có < 3 tháng dữ liệu
- 2b. AI Service trả về "Cần thêm dữ liệu để dự đoán chính xác"
- 2c. Use case kết thúc

**A2: Confidence score thấp**
- 5a. Confidence score < 0.5
- 5b. Hệ thống cảnh báo "Dự đoán có độ tin cậy thấp"
- 5c. Tiếp tục hiển thị nhưng với disclaimer

---

### UC-015: Quản lý thông báo

**Tác nhân chính:** User, System

**Điều kiện trước (Precondition):**
- User đã đăng nhập

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Thông báo được gửi đúng user
- Thông báo không bị mất

**Điều kiện sau (Postcondition):**
- User xem được danh sách thông báo
- Thông báo được đánh dấu đã đọc

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang thông báo
2. Hệ thống query notifications của user, sắp xếp theo created_at DESC
3. Hệ thống phân trang (20 items/page)
4. Hệ thống trả về danh sách thông báo
5. User thấy danh sách thông báo
6. User click vào thông báo
7. Hệ thống đánh dấu is_read = true, read_at = now()
8. Nếu có action_url, user được chuyển đến trang đó

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Không có thông báo**
- 3a. Không có thông báo nào
- 3b. Hệ thống trả về danh sách rỗng
- 3c. Frontend hiển thị "Chưa có thông báo"

---

### UC-016: Liên kết tài khoản Telegram

**Tác nhân chính:** User, Telegram Bot

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- User có Telegram account

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Link code an toàn và có thời hạn
- Chỉ user sở hữu mới liên kết được

**Điều kiện sau (Postcondition):**
- Telegram account được liên kết với web account
- User có thể sử dụng Telegram Bot

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang Settings > Telegram Integration
2. User nhấn "Liên kết Telegram"
3. Hệ thống tạo link code (16 ký tự, unique)
4. Hệ thống lưu vào telegram_link_codes với expires_at = now() + 10 phút
5. Hệ thống hiển thị link code cho user
6. User mở Telegram Bot và gửi lệnh /link
7. Bot yêu cầu user nhập link code
8. User nhập link code
9. Bot validate code: kiểm tra tồn tại, chưa hết hạn, chưa được dùng
10. Bot lưu telegram_user_id và web_user_id vào telegram_accounts
11. Bot đánh dấu code đã dùng (used_at = now())
12. Bot thông báo "Liên kết thành công"
13. Web app cập nhật trạng thái "Đã liên kết"

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Code hết hạn**
- 9a. Code đã hết hạn (expires_at < now())
- 9b. Bot thông báo "Code đã hết hạn. Vui lòng tạo code mới."
- 9c. User quay lại bước 2

**A2: Code đã được dùng**
- 9a. Code đã có used_at != NULL
- 9b. Bot thông báo "Code đã được sử dụng"
- 9c. Use case kết thúc

**A3: Telegram account đã liên kết**
- 9a. Telegram user_id đã có trong telegram_accounts
- 9b. Bot thông báo "Tài khoản Telegram này đã được liên kết"
- 9c. Use case kết thúc

---

### UC-017: Nhập giao dịch qua Telegram Bot

**Tác nhân chính:** User, Telegram Bot

**Điều kiện trước (Precondition):**
- User đã liên kết Telegram account
- Telegram Bot đang hoạt động

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Giao dịch được tạo chính xác
- Chỉ user đã liên kết mới sử dụng được

**Điều kiện sau (Postcondition):**
- Giao dịch được tạo từ Telegram
- User nhận được xác nhận

**Luồng sự kiện chính (Main Flow):**
1. User mở Telegram Bot
2. User gửi tin nhắn: "ăn bún bò 50k"
3. Bot nhận diện telegram_user_id
4. Bot tìm web_user_id từ telegram_accounts
5. Bot gửi request đến Backend API với web_user_id và message
6. Backend gọi AI Service để xử lý NLU
7. AI Service trích xuất: amount, category, description
8. Backend tạo transaction
9. Backend trả về transaction đã tạo
10. Bot format và gửi xác nhận cho user: "✅ Đã thêm: Ăn uống - 50,000 VND"
11. User thấy xác nhận trong Telegram

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Chưa liên kết**
- 4a. Không tìm thấy web_user_id
- 4b. Bot thông báo "Bạn chưa liên kết tài khoản. Dùng /link để liên kết."
- 4c. Use case kết thúc

**A2: AI không hiểu**
- 7a. AI không thể trích xuất đủ thông tin
- 7b. Bot yêu cầu user nhập lại hoặc dùng format: "/add [số tiền] [mô tả]"
- 7c. Use case kết thúc

---

### UC-018: Cảnh báo vượt ngân sách

**Tác nhân chính:** System

**Điều kiện trước (Precondition):**
- User có budget đang active
- User có notification preferences bật

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Cảnh báo được gửi đúng lúc
- Không spam notification

**Điều kiện sau (Postcondition):**
- Notification được tạo
- User nhận được cảnh báo

**Luồng sự kiện chính (Main Flow):**
1. User tạo giao dịch expense
2. Hệ thống tính spent_amount cho budget liên quan
3. Hệ thống tính usage_percentage = (spent_amount / amount) * 100
4. Nếu usage_percentage >= alert_threshold:
   4a. Hệ thống kiểm tra đã gửi notification cho budget này chưa (trong 24h)
   4b. Nếu chưa:
       - Tạo notification với type='warning', priority='high'
       - Gửi notification qua các channels được bật (web, email, telegram)
   4c. Nếu đã gửi: bỏ qua
5. Nếu usage_percentage >= 100:
   5a. Hệ thống tạo notification với type='error', priority='urgent'
   5b. Gửi notification ngay lập tức

**Luồng sự kiện phụ (Alternate Flow):**

**A1: User tắt notification**
- 4a. User đã tắt budget notifications trong preferences
- 4b. Hệ thống bỏ qua, không gửi notification
- 4c. Use case kết thúc

---

### UC-019: Đề xuất ngân sách tự động

**Tác nhân chính:** AI Agent, User

**Điều kiện trước (Precondition):**
- User đã có ít nhất 1 tháng dữ liệu chi tiêu
- User đã khai báo monthly_income

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Đề xuất dựa trên lịch sử chi tiêu thực tế
- Tổng đề xuất không vượt quá monthly_income

**Điều kiện sau (Postcondition):**
- Danh sách đề xuất ngân sách được hiển thị
- User có thể chấp nhận hoặc chỉnh sửa

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang "Đề xuất ngân sách" hoặc nhấn "Tạo ngân sách tự động"
2. Hệ thống lấy monthly_income từ user_profile
3. Hệ thống phân tích chi tiêu 3 tháng gần nhất:
   - Tính trung bình chi tiêu theo category
   - Xác định các category chiếm % lớn
4. Hệ thống đề xuất budget cho từng category:
   - Suggested amount = avg_spending * 1.1 (thêm 10% buffer)
   - Nếu category chiếm > 30% tổng chi tiêu: đề xuất giảm
5. Hệ thống đảm bảo tổng suggested <= monthly_income * 0.9 (để lại 10% tiết kiệm)
6. Hệ thống hiển thị danh sách đề xuất với:
   - Category name
   - Suggested amount
   - Current average spending
   - Difference
7. User xem và có thể:
   - Chấp nhận tất cả
   - Chỉnh sửa từng mục
   - Từ chối
8. Nếu user chấp nhận:
   8a. Hệ thống tạo budgets cho các category được chọn
   8b. Hệ thống set period = 'monthly', start_date = đầu tháng, end_date = cuối tháng
   8c. Hệ thống trả về danh sách budgets đã tạo

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Chưa có monthly_income**
- 2a. monthly_income = 0 hoặc NULL
- 2b. Hệ thống yêu cầu user khai báo monthly_income trước
- 2c. Use case kết thúc

**A2: Không đủ dữ liệu**
- 3a. User có < 1 tháng dữ liệu
- 3b. Hệ thống thông báo "Cần thêm dữ liệu để đề xuất chính xác"
- 3c. Use case kết thúc

---

### UC-020: Xem báo cáo phân tích chi tiết

**Tác nhân chính:** User

**Điều kiện trước (Precondition):**
- User đã đăng nhập
- User có ít nhất một giao dịch

**Đảm bảo tối thiểu (Minimal Guarantee):**
- Dữ liệu được tính toán chính xác
- Performance tốt với dữ liệu lớn

**Điều kiện sau (Postcondition):**
- Báo cáo chi tiết được hiển thị
- User có thể export (tùy chọn)

**Luồng sự kiện chính (Main Flow):**
1. User truy cập trang "Báo cáo" hoặc "Phân tích"
2. User chọn khoảng thời gian: tháng này, quý này, năm này, hoặc custom range
3. Hệ thống query database với filters:
   - start_date, end_date
   - user_id
4. Hệ thống tính toán các chỉ số:
   - Total income, total expense, net savings
   - Spending by category (với biểu đồ)
   - Spending trends (daily/weekly/monthly)
   - Top 5 categories chi tiêu nhiều nhất
   - So sánh với kỳ trước (nếu có)
   - Average transaction amount
   - Number of transactions
5. Hệ thống tạo biểu đồ data:
   - Pie chart: Spending by category
   - Line chart: Spending trends
   - Bar chart: Comparison với kỳ trước
6. Hệ thống trả về report data
7. Frontend hiển thị báo cáo với các biểu đồ

**Luồng sự kiện phụ (Alternate Flow):**

**A1: Không có dữ liệu trong khoảng thời gian**
- 4a. Không có giao dịch nào trong khoảng thời gian chọn
- 4b. Hệ thống trả về report với tất cả giá trị = 0
- 4c. Frontend hiển thị "Không có dữ liệu trong khoảng thời gian này"

---

## 3. TỔNG KẾT USE CASE

**Tổng số Use Case:** 20

**Phân loại theo module:**
- **Authentication & User Management:** UC-001, UC-002
- **Transaction Management:** UC-003, UC-004, UC-005, UC-006, UC-007
- **Financial Goals:** UC-008, UC-009
- **Budget Management:** UC-010, UC-018, UC-019
- **Analytics & Reporting:** UC-011, UC-020
- **AI Features:** UC-012, UC-013, UC-014
- **Notifications:** UC-015, UC-018
- **Telegram Integration:** UC-016, UC-017

**Tác nhân:**
- **User:** 18 use cases
- **AI Agent:** 4 use cases (UC-004, UC-012, UC-013, UC-014)
- **System:** 2 use cases (UC-015, UC-018)
- **Telegram Bot:** 2 use cases (UC-016, UC-017)


