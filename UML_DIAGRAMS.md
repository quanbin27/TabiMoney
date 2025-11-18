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
@startuml SEQ_004_NLUTransaction
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
@startuml SEQ_009_ContributeGoal
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "Goal Service" as GoalService
participant "Database" as DB
participant "Notification Service" as Notification

User -> Frontend: Mở chi tiết mục tiêu
Frontend -> Backend: GET /api/v1/goals/:id
Backend -> GoalService: GetGoal(goalID, userID)
GoalService -> DB: SELECT financial_goals
alt Goal thuộc user
    DB --> GoalService: Goal data
    GoalService --> Backend: Goal detail
    Backend --> Frontend: Hiển thị tiến độ hiện tại
else Không tìm thấy
    DB --> GoalService: Goal not found
    GoalService --> Backend: 404 Not Found
    Backend --> Frontend: Error response
    Frontend --> User: Hiển thị lỗi
    stop
end

User -> Frontend: Nhập số tiền muốn đóng góp
Frontend -> Backend: POST /api/v1/goals/:id/contribute
Backend -> GoalService: AddContribution(goalID, amount, userID)
GoalService -> GoalService: Validate amount > 0
GoalService -> DB: UPDATE financial_goals SET current_amount += amount
DB --> GoalService: Updated
GoalService -> GoalService: Tính toán progress %
GoalService -> GoalService: Kiểm tra trạng thái đạt mục tiêu?
alt Đạt hoặc vượt 100%
    GoalService -> DB: UPDATE is_achieved = true
    GoalService -> Notification: SendGoalAchieved(userID, goalID)
    Notification --> User: Gửi thông báo chúc mừng
else Chưa đạt
    GoalService -> Notification: SendProgressUpdate (nếu cần)
    Notification --> User: Nhắc tiến độ
end
GoalService -> GoalService: Ghi log đóng góp mới
GoalService --> Backend: GoalProgressResponse
Backend --> Frontend: 200 OK + tiến độ mới
Frontend --> User: Hiển thị tiến độ sau khi đóng góp

@enduml
```

## 8. Biểu đồ Tuần tự - Phát hiện Bất thường

```
@startuml SEQ_013_AnomalyDetection
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
@startuml SEQ_012_AIChat
actor User
participant "Frontend (Web)" as Frontend
participant "AI Service (FastAPI)" as AIService
participant "Transaction Service" as TxService
participant "Database" as DB
participant "Gemini API" as Gemini

User -> Frontend: Gửi tin nhắn
Frontend -> AIService: POST {AI_SERVICE_URL}/api/v1/chat/process
AIService -> TxService: Lấy dữ liệu context (transactions/budgets/goals)
TxService -> DB: SELECT theo user_id
DB --> TxService: Data
TxService --> AIService: Context
AIService -> Gemini: Prompt + context
Gemini --> AIService: Natural response + entities
alt Intent = add_transaction
    AIService -> DB: INSERT INTO transactions
end
AIService --> Frontend: ChatResponse (response, intent, suggestions)
Frontend --> User: Hiển thị câu trả lời

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
@startuml SEQ_002_Login
actor User
participant "Frontend" as Frontend
participant "Backend API" as Backend
participant "Auth Service" as Auth
participant "Database" as DB
participant "JWT Service" as JWT

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
@startuml SEQ_011_DashboardAnalytics
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

