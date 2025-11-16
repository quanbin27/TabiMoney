# TabiMoney - UML Diagrams Documentation

## 1. Biểu đồ Use Case Tổng quát

```
@startuml UseCase_Overall
!define ACTOR_COLOR #FFD700
!define USECASE_COLOR #87CEEB
!define SYSTEM_COLOR #90EE90

actor User as user ACTOR_COLOR
actor "Telegram Bot" as telegram ACTOR_COLOR
actor "AI Service" as ai_service ACTOR_COLOR

package "TabiMoney System" SYSTEM_COLOR {
    usecase "Đăng ký tài khoản" as UC_Register USECASE_COLOR
    usecase "Đăng nhập" as UC_Login USECASE_COLOR
    usecase "Quản lý hồ sơ" as UC_ManageProfile USECASE_COLOR
    usecase "Nhập giao dịch thủ công" as UC_ManualTransaction USECASE_COLOR
    usecase "Nhập giao dịch bằng NLU" as UC_NLUTransaction USECASE_COLOR
    usecase "Xem danh sách giao dịch" as UC_ViewTransactions USECASE_COLOR
    usecase "Cập nhật giao dịch" as UC_UpdateTransaction USECASE_COLOR
    usecase "Xóa giao dịch" as UC_DeleteTransaction USECASE_COLOR
    usecase "Quản lý danh mục" as UC_ManageCategories USECASE_COLOR
    usecase "Tạo mục tiêu tài chính" as UC_CreateGoal USECASE_COLOR
    usecase "Theo dõi mục tiêu" as UC_TrackGoal USECASE_COLOR
    usecase "Đóng góp vào mục tiêu" as UC_ContributeGoal USECASE_COLOR
    usecase "Tạo ngân sách" as UC_CreateBudget USECASE_COLOR
    usecase "Theo dõi ngân sách" as UC_TrackBudget USECASE_COLOR
    usecase "Xem Dashboard" as UC_ViewDashboard USECASE_COLOR
    usecase "Xem phân tích chi tiêu" as UC_ViewAnalytics USECASE_COLOR
    usecase "Chat với AI" as UC_AIChat USECASE_COLOR
    usecase "Dự đoán chi tiêu" as UC_PredictExpense USECASE_COLOR
    usecase "Phát hiện bất thường" as UC_DetectAnomaly USECASE_COLOR
    usecase "Phân loại thông minh" as UC_SmartCategorize USECASE_COLOR
    usecase "Xem thông báo" as UC_ViewNotifications USECASE_COLOR
    usecase "Quản lý cài đặt thông báo" as UC_ManageNotificationSettings USECASE_COLOR
    usecase "Liên kết Telegram" as UC_LinkTelegram USECASE_COLOR
    usecase "Nhập giao dịch qua Telegram" as UC_TelegramTransaction USECASE_COLOR
}

user --> UC_Register
user --> UC_Login
user --> UC_ManageProfile
user --> UC_ManualTransaction
user --> UC_NLUTransaction
user --> UC_ViewTransactions
user --> UC_UpdateTransaction
user --> UC_DeleteTransaction
user --> UC_ManageCategories
user --> UC_CreateGoal
user --> UC_TrackGoal
user --> UC_ContributeGoal
user --> UC_CreateBudget
user --> UC_TrackBudget
user --> UC_ViewDashboard
user --> UC_ViewAnalytics
user --> UC_AIChat
user --> UC_PredictExpense
user --> UC_DetectAnomaly
user --> UC_ViewNotifications
user --> UC_ManageNotificationSettings
user --> UC_LinkTelegram

telegram --> UC_LinkTelegram
telegram --> UC_TelegramTransaction
telegram --> UC_AIChat

ai_service --> UC_NLUTransaction
ai_service --> UC_PredictExpense
ai_service --> UC_DetectAnomaly
ai_service --> UC_SmartCategorize
ai_service --> UC_AIChat

UC_NLUTransaction ..> UC_ManualTransaction : <<extends>>
UC_SmartCategorize ..> UC_ManualTransaction : <<extends>>
UC_ViewDashboard ..> UC_ViewAnalytics : <<includes>>
UC_TrackBudget ..> UC_ViewAnalytics : <<includes>>
UC_TrackGoal ..> UC_ViewAnalytics : <<includes>>

@enduml
```

## 2. Biểu đồ Use Case - Quản lý Giao dịch

```
@startuml UseCase_TransactionManagement
actor User as user

package "Transaction Management" {
    usecase "Nhập giao dịch thủ công" as UC_CreateManual
    usecase "Nhập giao dịch bằng NLU" as UC_CreateNLU
    usecase "Xem danh sách giao dịch" as UC_List
    usecase "Lọc giao dịch" as UC_Filter
    usecase "Tìm kiếm giao dịch" as UC_Search
    usecase "Cập nhật giao dịch" as UC_Update
    usecase "Xóa giao dịch" as UC_Delete
    usecase "Xem chi tiết giao dịch" as UC_ViewDetail
    usecase "Đánh dấu giao dịch định kỳ" as UC_MarkRecurring
    usecase "Phân loại tự động" as UC_AutoCategorize
}

user --> UC_CreateManual
user --> UC_CreateNLU
user --> UC_List
user --> UC_Filter
user --> UC_Search
user --> UC_Update
user --> UC_Delete
user --> UC_ViewDetail
user --> UC_MarkRecurring

UC_CreateNLU ..> UC_AutoCategorize : <<uses>>
UC_List ..> UC_Filter : <<includes>>
UC_List ..> UC_Search : <<includes>>
UC_List ..> UC_ViewDetail : <<includes>>

@enduml
```

## 3. Biểu đồ Use Case - Quản lý Mục tiêu Tài chính

```
@startuml UseCase_GoalManagement
actor User as user

package "Financial Goal Management" {
    usecase "Tạo mục tiêu" as UC_CreateGoal
    usecase "Xem danh sách mục tiêu" as UC_ListGoals
    usecase "Cập nhật mục tiêu" as UC_UpdateGoal
    usecase "Xóa mục tiêu" as UC_DeleteGoal
    usecase "Đóng góp vào mục tiêu" as UC_Contribute
    usecase "Xem tiến độ mục tiêu" as UC_ViewProgress
    usecase "Đánh dấu hoàn thành" as UC_MarkAchieved
    usecase "Nhận cảnh báo mục tiêu" as UC_GoalAlert
}

user --> UC_CreateGoal
user --> UC_ListGoals
user --> UC_UpdateGoal
user --> UC_DeleteGoal
user --> UC_Contribute
user --> UC_ViewProgress

UC_ListGoals ..> UC_ViewProgress : <<includes>>
UC_Contribute ..> UC_ViewProgress : <<includes>>
UC_ViewProgress ..> UC_GoalAlert : <<extends>>

@enduml
```

## 4. Biểu đồ Use Case - Quản lý Ngân sách

```
@startuml UseCase_BudgetManagement
actor User as user

package "Budget Management" {
    usecase "Tạo ngân sách" as UC_CreateBudget
    usecase "Xem danh sách ngân sách" as UC_ListBudgets
    usecase "Cập nhật ngân sách" as UC_UpdateBudget
    usecase "Xóa ngân sách" as UC_DeleteBudget
    usecase "Xem thống kê ngân sách" as UC_ViewBudgetStats
    usecase "Nhận cảnh báo vượt ngân sách" as UC_BudgetAlert
    usecase "Gợi ý ngân sách tự động" as UC_AutoSuggestBudget
    usecase "Tạo ngân sách từ gợi ý" as UC_CreateFromSuggestion
}

user --> UC_CreateBudget
user --> UC_ListBudgets
user --> UC_UpdateBudget
user --> UC_DeleteBudget
user --> UC_ViewBudgetStats
user --> UC_AutoSuggestBudget
user --> UC_CreateFromSuggestion

UC_ListBudgets ..> UC_ViewBudgetStats : <<includes>>
UC_ViewBudgetStats ..> UC_BudgetAlert : <<extends>>
UC_AutoSuggestBudget ..> UC_CreateFromSuggestion : <<extends>>

@enduml
```

## 5. Biểu đồ Use Case - AI Features

```
@startuml UseCase_AIFeatures
actor User as user
actor "AI Service" as ai_service

package "AI Features" {
    usecase "Xử lý ngôn ngữ tự nhiên" as UC_NLU
    usecase "Phân loại giao dịch" as UC_Categorize
    usecase "Dự đoán chi tiêu" as UC_Predict
    usecase "Phát hiện bất thường" as UC_Anomaly
    usecase "Chat với AI" as UC_Chat
    usecase "Phân tích chi tiêu" as UC_Analyze
    usecase "Gợi ý tài chính" as UC_Suggest
}

user --> UC_NLU
user --> UC_Categorize
user --> UC_Predict
user --> UC_Anomaly
user --> UC_Chat
user --> UC_Analyze

ai_service --> UC_NLU
ai_service --> UC_Categorize
ai_service --> UC_Predict
ai_service --> UC_Anomaly
ai_service --> UC_Chat
ai_service --> UC_Analyze

UC_NLU ..> UC_Categorize : <<uses>>
UC_Analyze ..> UC_Suggest : <<includes>>

@enduml
```

## 6. Biểu đồ Tuần tự - Nhập giao dịch bằng NLU

```
@startuml Sequence_NLUTransaction
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "AI Service" as AI
participant "NLU Service" as NLU
participant "Database" as DB

User -> Frontend: Nhập câu: "tôi vừa ăn bún bò 50k"
Frontend -> Backend: POST /api/v1/ai/nlu
Backend -> AI: Gửi yêu cầu NLU
AI -> NLU: Xử lý ngôn ngữ tự nhiên
NLU -> NLU: Trích xuất: amount, category, date
NLU --> AI: Trả về structured data
AI --> Backend: {amount: 50000, category: "Ăn uống", date: "2024-01-15"}
Backend -> DB: Kiểm tra category
DB --> Backend: Category info
Backend -> DB: Tạo transaction
DB --> Backend: Transaction created
Backend -> AI: Gợi ý category (nếu cần)
AI --> Backend: Category suggestion
Backend --> Frontend: Transaction response
Frontend --> User: Hiển thị giao dịch đã tạo

@enduml
```

## 7. Biểu đồ Tuần tự - Tạo và theo dõi Mục tiêu

```
@startuml Sequence_GoalManagement
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "Goal Service" as GoalService
participant "Database" as DB
participant "Notification Service" as Notification

User -> Frontend: Tạo mục tiêu mới
Frontend -> Backend: POST /api/v1/goals
Backend -> GoalService: CreateGoal(userID, request)
GoalService -> DB: INSERT INTO financial_goals
DB --> GoalService: Goal created
GoalService --> Backend: Goal response
Backend --> Frontend: Goal created

User -> Frontend: Đóng góp vào mục tiêu
Frontend -> Backend: POST /api/v1/goals/:id/contribute
Backend -> GoalService: AddContribution(goalID, amount)
GoalService -> DB: UPDATE financial_goals SET current_amount
DB --> GoalService: Updated
GoalService -> GoalService: Tính toán progress
GoalService -> GoalService: Kiểm tra đạt mục tiêu?
alt Mục tiêu đạt được
    GoalService -> DB: UPDATE is_achieved = true
    GoalService -> Notification: Gửi thông báo chúc mừng
    Notification --> User: Thông báo
end
GoalService --> Backend: Updated goal
Backend --> Frontend: Goal updated

@enduml
```

## 8. Biểu đồ Tuần tự - Phát hiện Bất thường

```
@startuml Sequence_AnomalyDetection
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "Analytics Handler" as Analytics
participant "AI Service" as AI
participant "Anomaly Service" as Anomaly
participant "Database" as DB
participant "Notification Service" as Notification

User -> Frontend: Xem phân tích bất thường
Frontend -> Backend: GET /api/v1/analytics/anomalies
Backend -> Analytics: GetAnomalies(userID)
Analytics -> DB: Lấy lịch sử giao dịch
DB --> Analytics: Transaction history
Analytics -> AI: POST /api/v1/ai/anomaly/detect
AI -> Anomaly: Phân tích pattern
Anomaly -> Anomaly: Tính toán statistical model
Anomaly -> Anomaly: So sánh với baseline
Anomaly --> AI: Danh sách anomalies
AI --> Analytics: Anomaly results
Analytics -> Notification: Kiểm tra cần gửi cảnh báo?
alt Có bất thường nghiêm trọng
    Notification -> DB: Tạo notification
    Notification --> User: Gửi thông báo
end
Analytics --> Backend: Anomaly list
Backend --> Frontend: Anomalies response
Frontend --> User: Hiển thị bất thường

@enduml
```

## 9. Biểu đồ Tuần tự - Chat với AI

```
@startuml Sequence_AIChat
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "AI Handler" as AIHandler
participant "AI Service" as AIService
participant "Chat Service" as ChatService
participant "Database" as DB

User -> Frontend: Gửi tin nhắn: "Tháng này tôi tiêu bao nhiêu cho ăn uống?"
Frontend -> Backend: POST /api/v1/ai/chat
Backend -> AIHandler: ProcessChat(request)
AIHandler -> AIService: ProcessChatMessage(userID, message)
AIService -> ChatService: Phân tích intent
ChatService -> ChatService: Xác định loại câu hỏi
alt Câu hỏi về dữ liệu
    ChatService -> DB: Query transactions
    DB --> ChatService: Transaction data
    ChatService -> AIService: Format data for LLM
end
AIService -> AIService: Gọi LLM với context
AIService -> AIService: Xử lý response
AIService --> AIHandler: AI response
AIHandler --> Backend: Chat response
Backend --> Frontend: Response message
Frontend --> User: Hiển thị câu trả lời

@enduml
```

## 10. Biểu đồ Hoạt động - Quy trình Nhập Giao dịch

```
@startuml Activity_TransactionEntry
start
:User chọn phương thức nhập;
if (Phương thức?) then (Thủ công)
    :Điền form: category, amount, date, description;
else (NLU)
    :Nhập câu tự nhiên;
    :Gửi đến AI Service;
    :AI xử lý và trích xuất thông tin;
    :Hiển thị preview;
    if (User xác nhận?) then (Không)
        :Cho phép chỉnh sửa;
    endif
endif
:Validate dữ liệu;
if (Dữ liệu hợp lệ?) then (Không)
    :Hiển thị lỗi;
    stop
endif
:Kiểm tra category;
if (Category hợp lệ?) then (Không)
    :Gợi ý category từ AI;
    :User chọn category;
endif
:Lưu transaction vào DB;
:Kiểm tra budget threshold;
if (Vượt ngân sách?) then (Có)
    :Gửi thông báo cảnh báo;
endif
:Kiểm tra giao dịch lớn;
if (Amount > 1M VND?) then (Có)
    :Gửi thông báo giao dịch lớn;
endif
:Xóa dashboard cache;
:Hiển thị thông báo thành công;
stop

@enduml
```

## 11. Biểu đồ Hoạt động - Quy trình Phân tích Tài chính

```
@startuml Activity_FinancialAnalysis
start
:User yêu cầu xem phân tích;
:Kiểm tra cache;
if (Có cache?) then (Có)
    :Trả về dữ liệu từ cache;
    stop
endif
:Lấy dữ liệu giao dịch từ DB;
:Tính toán tổng thu nhập;
:Tính toán tổng chi tiêu;
:Tính toán số dư;
:Phân tích theo danh mục;
:Tính toán tỷ lệ phần trăm;
:Tính toán financial health score;
:Phân tích xu hướng;
:Kiểm tra có đủ dữ liệu cho AI?;
if (Có đủ dữ liệu?) then (Có)
    :Gọi AI Service phân tích;
    :Nhận insights và recommendations;
endif
:Tổng hợp kết quả;
:Lưu vào cache;
:Trả về kết quả cho user;
stop

@enduml
```

## 12. Biểu đồ Hoạt động - Quy trình Quản lý Ngân sách

```
@startuml Activity_BudgetManagement
start
:User tạo/cập nhật ngân sách;
:Validate dữ liệu ngân sách;
if (Dữ liệu hợp lệ?) then (Không)
    :Hiển thị lỗi;
    stop
endif
:Lưu ngân sách vào DB;
:Tính toán spent amount;
:Tính toán remaining amount;
:Tính toán usage percentage;
if (Usage > Alert Threshold?) then (Có)
    :Gửi thông báo cảnh báo;
endif
if (Usage >= 100%?) then (Có)
    :Gửi thông báo vượt ngân sách;
endif
:Kiểm tra pacing;
if (Chi tiêu vượt tốc độ?) then (Có)
    :Gửi cảnh báo pacing;
endif
:Cập nhật dashboard cache;
:Hiển thị kết quả;
stop

@enduml
```

## 13. Biểu đồ Trạng thái - Giao dịch (Transaction)

```
@startuml State_Transaction
[*] --> Draft : Tạo mới
Draft --> Pending : Lưu nháp
Pending --> Confirmed : Xác nhận
Pending --> Cancelled : Hủy
Confirmed --> Updated : Cập nhật
Updated --> Confirmed : Xác nhận lại
Confirmed --> Deleted : Xóa
Updated --> Deleted : Xóa
Cancelled --> [*]
Deleted --> [*]

state Confirmed {
    [*] --> Active
    Active
    Recurring : Nếu is_recurring = true
}

@enduml
```

## 14. Biểu đồ Trạng thái - Mục tiêu Tài chính (Financial Goal)

```
@startuml State_FinancialGoal
[*] --> Created : Tạo mục tiêu
Created --> Active : Kích hoạt
Active --> InProgress : Bắt đầu đóng góp
InProgress --> InProgress : Tiếp tục đóng góp
InProgress --> Achieved : current_amount >= target_amount
InProgress --> Paused : Tạm dừng
Paused --> InProgress : Tiếp tục
Achieved --> [*] : Hoàn thành
Active --> Cancelled : Hủy
InProgress --> Cancelled : Hủy
Paused --> Cancelled : Hủy
Cancelled --> [*]

state InProgress {
    [*] --> Tracking
    Tracking : Theo dõi tiến độ
    Tracking --> Alert : Gần deadline
    Alert --> Tracking : Cập nhật
}

@enduml
```

## 15. Biểu đồ Trạng thái - Ngân sách (Budget)

```
@startuml State_Budget
[*] --> Created : Tạo ngân sách
Created --> Active : Kích hoạt
Active --> Monitoring : Bắt đầu theo dõi
Monitoring --> Warning : usage >= alert_threshold
Monitoring --> Exceeded : usage >= 100%
Warning --> Monitoring : Giảm chi tiêu
Warning --> Exceeded : Tiếp tục chi tiêu
Exceeded --> Monitoring : Điều chỉnh ngân sách
Monitoring --> Completed : Hết thời gian
Active --> Paused : Tạm dừng
Paused --> Active : Kích hoạt lại
Active --> Cancelled : Hủy
Monitoring --> Cancelled : Hủy
Completed --> [*]
Cancelled --> [*]

state Monitoring {
    [*] --> Normal
    Normal --> Alert : usage > threshold
    Alert --> Critical : usage > 90%
    Critical --> Alert : usage giảm
    Alert --> Normal : usage < threshold
}

@enduml
```

## 16. Biểu đồ ERD - Entity Relationship Diagram

```
@startuml ERD_TabiMoney
!define TABLE_COLOR #E1F5FF
!define PK_COLOR #FFD700

entity "users" as users TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    --
    * email : VARCHAR(255) <<UNIQUE>>
    * username : VARCHAR(100) <<UNIQUE>>
    * password_hash : VARCHAR(255)
    first_name : VARCHAR(100)
    last_name : VARCHAR(100)
    phone : VARCHAR(20)
    avatar_url : VARCHAR(500)
    is_verified : BOOLEAN
    verification_token : VARCHAR(255)
    reset_token : VARCHAR(255)
    reset_token_expires_at : TIMESTAMP
    last_login_at : TIMESTAMP
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
    deleted_at : TIMESTAMP
}

entity "user_profiles" as profiles TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    --
    monthly_income : DECIMAL(15,2)
    currency : VARCHAR(3)
    timezone : VARCHAR(50)
    language : VARCHAR(5)
    notification_settings : JSON
    ai_settings : JSON
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
}

entity "categories" as categories TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    user_id : BIGINT <<FK>>
    parent_id : BIGINT <<FK>>
    --
    * name : VARCHAR(100)
    name_en : VARCHAR(100)
    description : TEXT
    icon : VARCHAR(50)
    color : VARCHAR(7)
    is_system : BOOLEAN
    is_active : BOOLEAN
    sort_order : INT
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
}

entity "transactions" as transactions TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    * category_id : BIGINT <<FK>>
    parent_transaction_id : BIGINT <<FK>>
    ai_suggested_category_id : BIGINT <<FK>>
    --
    * amount : DECIMAL(15,2)
    description : TEXT
    * transaction_type : ENUM
    * transaction_date : DATE
    transaction_time : TIME
    location : VARCHAR(200)
    tags : JSON
    metadata : JSON
    is_recurring : BOOLEAN
    recurring_pattern : VARCHAR(50)
    ai_confidence : DECIMAL(3,2)
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
}

entity "financial_goals" as goals TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    --
    * title : VARCHAR(200)
    description : TEXT
    * target_amount : DECIMAL(15,2)
    current_amount : DECIMAL(15,2)
    target_date : DATE
    goal_type : ENUM
    priority : ENUM
    is_achieved : BOOLEAN
    achieved_at : TIMESTAMP
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
}

entity "budgets" as budgets TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    category_id : BIGINT <<FK>>
    --
    * name : VARCHAR(200)
    * amount : DECIMAL(15,2)
    * period : ENUM
    * start_date : DATE
    * end_date : DATE
    is_active : BOOLEAN
    alert_threshold : DECIMAL(5,2)
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
}

entity "notifications" as notifications TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    --
    * title : VARCHAR(200)
    * message : TEXT
    * notification_type : ENUM
    priority : ENUM
    is_read : BOOLEAN
    read_at : TIMESTAMP
    action_url : VARCHAR(500)
    metadata : JSON
    * created_at : TIMESTAMP
}

entity "ai_analysis" as ai_analysis TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    --
    * analysis_type : ENUM
    * data : JSON
    confidence_score : DECIMAL(3,2)
    model_version : VARCHAR(50)
    * created_at : TIMESTAMP
}

entity "ai_feedback" as ai_feedback TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    transaction_id : BIGINT <<FK>>
    --
    * feedback_type : ENUM
    original_prediction : JSON
    user_correction : JSON
    feedback_text : TEXT
    * created_at : TIMESTAMP
}

entity "user_sessions" as sessions TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * user_id : BIGINT <<FK>>
    --
    * token_hash : VARCHAR(255)
    * refresh_token_hash : VARCHAR(255)
    * expires_at : TIMESTAMP
    * refresh_expires_at : TIMESTAMP
    user_agent : TEXT
    ip_address : VARCHAR(45)
    is_active : BOOLEAN
    * created_at : TIMESTAMP
}

entity "telegram_accounts" as telegram TABLE_COLOR {
    * id : BIGINT <<PK>> PK_COLOR
    * telegram_user_id : BIGINT <<UNIQUE>>
    * web_user_id : BIGINT <<FK>>
    --
    * created_at : TIMESTAMP
    * updated_at : TIMESTAMP
}

users ||--o{ profiles : "has one"
users ||--o{ sessions : "has many"
users ||--o{ transactions : "has many"
users ||--o{ categories : "has many"
users ||--o{ goals : "has many"
users ||--o{ budgets : "has many"
users ||--o{ notifications : "has many"
users ||--o{ ai_analysis : "has many"
users ||--o{ ai_feedback : "has many"
users ||--o| telegram : "has one"

categories ||--o{ categories : "parent-child"
categories ||--o{ transactions : "categorizes"
categories ||--o{ budgets : "has budget"

transactions ||--o| transactions : "parent-child (recurring)"
transactions ||--o| categories : "ai_suggested"
transactions ||--o{ ai_feedback : "has feedback"

@enduml
```

## 17. Biểu đồ Tuần tự - Đăng ký và Đăng nhập

```
@startuml Sequence_Authentication
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "Auth Service" as Auth
participant "Database" as DB
participant "JWT Service" as JWT

== Đăng ký ==
User -> Frontend: Điền form đăng ký
Frontend -> Backend: POST /api/v1/auth/register
Backend -> Auth: Register(request)
Auth -> DB: Kiểm tra email/username tồn tại
alt Email/Username đã tồn tại
    DB --> Auth: Conflict
    Auth --> Backend: Error
    Backend --> Frontend: 409 Conflict
    Frontend --> User: Hiển thị lỗi
else Email/Username hợp lệ
    Auth -> Auth: Hash password
    Auth -> DB: INSERT INTO users
    DB --> Auth: User created
    Auth -> DB: INSERT INTO user_profiles
    DB --> Auth: Profile created
    Auth -> JWT: Generate tokens
    JWT --> Auth: AccessToken + RefreshToken
    Auth --> Backend: AuthResponse
    Backend --> Frontend: 201 Created
    Frontend --> User: Đăng ký thành công
end

== Đăng nhập ==
User -> Frontend: Nhập email/password
Frontend -> Backend: POST /api/v1/auth/login
Backend -> Auth: Login(request)
Auth -> DB: SELECT user WHERE email
DB --> Auth: User data
Auth -> Auth: Verify password
alt Mật khẩu sai
    Auth --> Backend: Error
    Backend --> Frontend: 401 Unauthorized
    Frontend --> User: Sai mật khẩu
else Mật khẩu đúng
    Auth -> DB: UPDATE last_login_at
    Auth -> JWT: Generate tokens
    JWT --> Auth: Tokens
    Auth -> DB: INSERT INTO user_sessions
    Auth --> Backend: AuthResponse
    Backend --> Frontend: 200 OK
    Frontend --> User: Đăng nhập thành công
end

@enduml
```

## 18. Biểu đồ Tuần tự - Dashboard và Analytics

```
@startuml Sequence_DashboardAnalytics
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "Analytics Handler" as Analytics
participant "Transaction Service" as TxService
participant "Database" as DB
participant "Redis Cache" as Cache

User -> Frontend: Truy cập Dashboard
Frontend -> Backend: GET /api/v1/analytics/dashboard
Backend -> Analytics: GetDashboardAnalytics(userID)
Analytics -> Cache: Kiểm tra cache
alt Cache hit
    Cache --> Analytics: Cached data
    Analytics --> Backend: Analytics data
    Backend --> Frontend: Response
    Frontend --> User: Hiển thị dashboard
else Cache miss
    Analytics -> TxService: GetMonthlySummary(userID)
    TxService -> DB: Query transactions
    DB --> TxService: Transaction list
    TxService -> TxService: Tính toán analytics
    TxService --> Analytics: Monthly summary
    Analytics -> Analytics: Tính toán financial health
    Analytics -> Analytics: Phân tích category breakdown
    Analytics -> Cache: Lưu vào cache
    Analytics --> Backend: Analytics data
    Backend --> Frontend: Response
    Frontend --> User: Hiển thị dashboard
end

User -> Frontend: Xem phân tích chi tiêu theo danh mục
Frontend -> Backend: GET /api/v1/analytics/category-spending
Backend -> Analytics: GetCategorySpending(userID, dateRange)
Analytics -> TxService: GetCategorySpending(userID, startDate, endDate)
TxService -> DB: Query với GROUP BY category
DB --> TxService: Category spending data
TxService -> TxService: Tính toán percentages
TxService --> Analytics: Category analytics
Analytics --> Backend: Category spending
Backend --> Frontend: Response
Frontend --> User: Hiển thị biểu đồ

@enduml
```

## 19. Biểu đồ Hoạt động - Quy trình Dự đoán Chi tiêu

```
@startuml Activity_ExpensePrediction
start
:User yêu cầu dự đoán chi tiêu;
:Lấy lịch sử giao dịch 3-6 tháng;
if (Đủ dữ liệu?) then (Không)
    :Trả về thông báo cần thêm dữ liệu;
    stop
endif
:Chuẩn bị dữ liệu training;
:Phân tích pattern chi tiêu;
:Tính toán xu hướng;
:Gọi AI Service;
:AI Service xử lý với ML model;
:Phân tích theo danh mục;
:Tính toán dự đoán cho từng danh mục;
:Tổng hợp kết quả;
:Tính toán confidence score;
if (Confidence > threshold?) then (Có)
    :Lưu prediction vào DB;
    :Trả về kết quả dự đoán;
else (Không)
    :Trả về với cảnh báo độ tin cậy thấp;
endif
stop

@enduml
```

## 20. Biểu đồ Hoạt động - Quy trình Thông báo

```
@startuml Activity_NotificationFlow
start
:Sự kiện trigger (transaction, budget, goal);
:Kiểm tra notification preferences;
if (User bật thông báo?) then (Không)
    stop
endif
:Xác định loại thông báo;
if (Loại thông báo?) then (Budget Alert)
    :Kiểm tra budget threshold;
    if (Vượt threshold?) then (Có)
        :Tạo notification;
    else (Không)
        stop
    endif
elseif (Goal Alert) then
    :Kiểm tra goal progress;
    if (Cần cảnh báo?) then (Có)
        :Tạo notification;
    else (Không)
        stop
    endif
elseif (Large Transaction) then
    :Kiểm tra amount > 1M;
    if (Lớn hơn?) then (Có)
        :Tạo notification;
    else (Không)
        stop
    endif
else (Scheduled)
    :Kiểm tra lịch;
    :Tạo notification;
endif
:Kiểm tra channel preferences;
if (Email enabled?) then (Có)
    :Gửi email;
endif
if (In-app enabled?) then (Có)
    :Lưu vào DB;
    :Push real-time notification;
endif
if (Telegram enabled?) then (Có)
    :Gửi qua Telegram bot;
endif
stop

@enduml
```

---

## Hướng dẫn sử dụng

Các biểu đồ trên được viết bằng định dạng PlantUML, có thể import vào:
- StarUML (thông qua plugin PlantUML)
- Visual Studio Code (với extension PlantUML)
- Online: http://www.plantuml.com/plantuml/uml/
- Các công cụ UML khác hỗ trợ PlantUML

Để chuyển đổi sang định dạng StarUML native, bạn có thể:
1. Import vào PlantUML online editor
2. Export sang định dạng khác
3. Hoặc sử dụng plugin PlantUML trong StarUML

