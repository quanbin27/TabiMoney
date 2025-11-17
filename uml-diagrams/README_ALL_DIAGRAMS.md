# Tổng hợp đầy đủ tất cả biểu đồ UML - TabiMoney System

## Tổng quan

Tài liệu này liệt kê **đầy đủ tất cả các biểu đồ UML** đã được tạo cho hệ thống TabiMoney, bao gồm:
- **20 Use Case Diagrams** - Mỗi use case có một biểu đồ riêng
- **18+ Activity Diagrams** - Quy trình chi tiết cho các use cases
- **15+ Sequence Diagrams** - Luồng tương tác giữa các components
- **3 State Diagrams** - Trạng thái của các entities chính
- **1 ERD Diagram** - Mô hình dữ liệu

---

## 1. USE CASE DIAGRAMS (20 files)

Mỗi use case có một biểu đồ Use Case riêng mô tả các sub-use cases và relationships:

| File | Use Case | Mô tả |
|------|----------|-------|
| `UC_001_Register.puml` | Đăng ký tài khoản | Validate, CheckExists, HashPassword, CreateUser, CreateProfile |
| `UC_002_Login.puml` | Đăng nhập | Validate, FindUser, VerifyPassword, GenerateTokens, SaveSession |
| `UC_003_ManualTransaction.puml` | Nhập giao dịch thủ công | SelectType, SelectCategory, EnterAmount, Validate, SaveTransaction |
| `UC_004_NLUTransaction.puml` | Nhập giao dịch bằng NLU | EnterNaturalLanguage, SendToAI, ProcessNLU, ExtractInfo, FindCategory |
| `UC_005_ViewTransactions.puml` | Xem danh sách giao dịch | FilterTransactions, SearchTransactions, SortTransactions, Paginate |
| `UC_006_UpdateTransaction.puml` | Cập nhật giao dịch | SelectTransaction, EditInfo, Validate, CheckOwnership, UpdateRecord |
| `UC_007_DeleteTransaction.puml` | Xóa giao dịch | SelectTransaction, ConfirmDelete, CheckOwnership, DeleteRecord |
| `UC_008_CreateGoal.puml` | Tạo mục tiêu tài chính | EnterGoalInfo, Validate, CreateRecord, CalculateProgress |
| `UC_009_ContributeGoal.puml` | Đóng góp vào mục tiêu | SelectGoal, EnterAmount, Validate, UpdateAmount, CheckAchieved |
| `UC_010_CreateBudget.puml` | Tạo ngân sách | EnterBudgetInfo, Validate, CheckDateRange, CreateRecord |
| `UC_011_DashboardAnalytics.puml` | Xem Dashboard Analytics | CheckCache, CalculateIncome, CalculateExpense, AnalyzeByCategory |
| `UC_012_AIChat.puml` | Chat với AI | EnterQuestion, AnalyzeIntent, QueryData, CreateContext, CallLLM |
| `UC_013_AnomalyDetection.puml` | Phát hiện bất thường | GetHistory, CalculatePattern, CompareBaseline, CalculateScore |
| `UC_014_ExpensePrediction.puml` | Dự đoán chi tiêu | GetHistory, AnalyzeTrend, AnalyzePattern, ApplyModel, CalculatePrediction |
| `UC_015_ManageNotifications.puml` | Quản lý thông báo | ViewNotifications, MarkRead, FilterNotifications, DeleteNotification |
| `UC_016_LinkTelegram.puml` | Liên kết Telegram | GenerateCode, ShowCode, EnterCode, ValidateCode, SaveLink |
| `UC_017_TelegramTransaction.puml` | Nhập giao dịch qua Telegram | SendMessage, IdentifyUser, FindWebUser, ProcessNLU, CreateTransaction |
| `UC_018_BudgetAlert.puml` | Cảnh báo vượt ngân sách | CalculateSpent, CalculateUsage, CheckThreshold, CreateNotification |
| `UC_019_AutoBudgetSuggestion.puml` | Đề xuất ngân sách tự động | GetIncome, AnalyzeSpending, CalculateAverage, SuggestBudget |
| `UC_020_DetailedReport.puml` | Xem báo cáo phân tích chi tiết | SelectTimeRange, QueryDatabase, CalculateMetrics, CreateChartData |

---

## 2. ACTIVITY DIAGRAMS (18+ files)

Mô tả quy trình chi tiết từng bước của các use cases:

| File | Use Case | Mô tả |
|------|----------|-------|
| `ACT_001_Register.puml` | Đăng ký tài khoản | Quy trình từ nhập form đến tạo tài khoản thành công |
| `ACT_002_Login.puml` | Đăng nhập | Quy trình xác thực và tạo session |
| `ACT_005_ViewTransactions.puml` | Xem danh sách giao dịch | Quy trình lọc, tìm kiếm, phân trang |
| `ACT_006_UpdateTransaction.puml` | Cập nhật giao dịch | Quy trình chỉnh sửa và validate |
| `ACT_007_DeleteTransaction.puml` | Xóa giao dịch | Quy trình xác nhận và xóa |
| `ACT_008_CreateGoal.puml` | Tạo mục tiêu tài chính | Quy trình tạo và validate mục tiêu |
| `ACT_009_ContributeGoal.puml` | Đóng góp vào mục tiêu | Quy trình thêm tiền và kiểm tra đạt mục tiêu |
| `ACT_010_CreateBudget.puml` | Tạo ngân sách | Quy trình tạo và validate ngân sách |
| `ACT_015_ManageNotifications.puml` | Quản lý thông báo | Quy trình xem, đánh dấu đọc, lọc |
| `ACT_016_LinkTelegram.puml` | Liên kết Telegram | Quy trình tạo code và liên kết |
| `ACT_017_TelegramTransaction.puml` | Nhập giao dịch qua Telegram | Quy trình xử lý tin nhắn và tạo giao dịch |
| `ACT_018_BudgetAlert.puml` | Cảnh báo vượt ngân sách | Quy trình tính toán và gửi cảnh báo |
| `ACT_019_AutoBudgetSuggestion.puml` | Đề xuất ngân sách tự động | Quy trình phân tích và đề xuất |
| `ACT_020_DetailedReport.puml` | Xem báo cáo phân tích chi tiết | Quy trình tính toán và hiển thị báo cáo |
| `10_Activity_TransactionEntry.puml` | Nhập giao dịch (tổng hợp) | Quy trình nhập giao dịch với NLU và thủ công |
| `11_Activity_FinancialAnalysis.puml` | Phân tích tài chính | Quy trình phân tích và cache |
| `12_Activity_BudgetManagement.puml` | Quản lý ngân sách | Quy trình quản lý và cảnh báo |
| `19_Activity_ExpensePrediction.puml` | Dự đoán chi tiêu | Quy trình dự đoán với AI |
| `20_Activity_NotificationFlow.puml` | Quy trình thông báo | Quy trình gửi thông báo qua các channels |

---

## 3. SEQUENCE DIAGRAMS (15+ files)

Mô tả luồng tương tác giữa các components:

| File | Use Case | Components |
|------|----------|------------|
| `SEQ_001_Register.puml` | Đăng ký tài khoản | User, Frontend, Backend, AuthService, Database |
| `SEQ_005_ViewTransactions.puml` | Xem danh sách giao dịch | User, Frontend, Backend, Handler, Database |
| `SEQ_006_UpdateTransaction.puml` | Cập nhật giao dịch | User, Frontend, Backend, Handler, Service, Database, Cache |
| `SEQ_007_DeleteTransaction.puml` | Xóa giao dịch | User, Frontend, Backend, Handler, Service, Database, Cache |
| `SEQ_008_CreateGoal.puml` | Tạo mục tiêu tài chính | User, Frontend, Backend, Handler, Service, Database |
| `SEQ_010_CreateBudget.puml` | Tạo ngân sách | User, Frontend, Backend, Handler, Service, Database |
| `SEQ_016_LinkTelegram.puml` | Liên kết Telegram | User, Frontend, Backend, Database, Telegram Bot |
| `SEQ_017_TelegramTransaction.puml` | Nhập giao dịch qua Telegram | User, Telegram Bot, Backend, AI Service, Service, Database |
| `SEQ_019_AutoBudgetSuggestion.puml` | Đề xuất ngân sách tự động | User, Frontend, Backend, Handler, Service, Database, AI Service |
| `06_Sequence_NLUTransaction.puml` | Nhập giao dịch bằng NLU | User, Frontend, Backend, AI Service, NLU Service, Database |
| `07_Sequence_GoalManagement.puml` | Quản lý mục tiêu | User, Frontend, Backend, Goal Service, Database, Notification |
| `08_Sequence_AnomalyDetection.puml` | Phát hiện bất thường | User, Frontend, Backend, Analytics, AI Service, Anomaly Service |
| `09_Sequence_AIChat.puml` | Chat với AI | User, Frontend, Backend, AI Handler, AI Service, Chat Service |
| `17_Sequence_Authentication.puml` | Xác thực (Register + Login) | User, Frontend, Backend, Auth Service, Database, JWT Service |
| `18_Sequence_DashboardAnalytics.puml` | Dashboard Analytics | User, Frontend, Backend, Analytics, Transaction Service, Database, Cache |

---

## 4. STATE DIAGRAMS (3 files)

Mô tả trạng thái và chuyển đổi trạng thái của các entities:

| File | Entity | States |
|------|--------|--------|
| `13_State_Transaction.puml` | Transaction | Draft, Pending, Confirmed, Updated, Cancelled, Deleted |
| `14_State_FinancialGoal.puml` | Financial Goal | Created, Active, InProgress, Paused, Achieved, Cancelled |
| `15_State_Budget.puml` | Budget | Created, Active, Monitoring, Warning, Exceeded, Paused, Completed, Cancelled |

---

## 5. ERD DIAGRAM (1 file)

| File | Mô tả |
|------|-------|
| `16_ERD_TabiMoney.puml` | Entity Relationship Diagram với 11 bảng: users, user_profiles, categories, transactions, financial_goals, budgets, notifications, ai_analysis, ai_feedback, user_sessions, telegram_accounts |

---

## 6. BIỂU ĐỒ TỔNG QUAN (5 files)

| File | Mô tả |
|------|-------|
| `01_UseCase_Overall.puml` | Use Case tổng quát của toàn hệ thống (23 use cases) |
| `02_UseCase_TransactionManagement.puml` | Use Case quản lý giao dịch |
| `03_UseCase_GoalManagement.puml` | Use Case quản lý mục tiêu |
| `04_UseCase_BudgetManagement.puml` | Use Case quản lý ngân sách |
| `05_UseCase_AIFeatures.puml` | Use Case các tính năng AI |

---

## Tổng kết

- **Tổng số file .puml**: 60+ files
- **Use Case Diagrams**: 25 files (20 chi tiết + 5 tổng quan)
- **Activity Diagrams**: 19 files
- **Sequence Diagrams**: 15 files
- **State Diagrams**: 3 files
- **ERD Diagram**: 1 file

---

## Cách sử dụng

1. **Xem biểu đồ**: Sử dụng PlantUML viewer hoặc plugin trong IDE
2. **Chỉnh sửa**: Mở file .puml và chỉnh sửa trực tiếp
3. **Export**: Sử dụng PlantUML để export sang PNG, SVG, PDF
4. **Import vào StarUML**: Có thể import file .puml vào StarUML

---

## Lưu ý

- Tất cả các biểu đồ đều được tạo dựa trên tài liệu USE_CASE_DETAILED.md
- Các biểu đồ được đánh số theo thứ tự use case (UC_001 đến UC_020)
- Một số use case có nhiều biểu đồ (Use Case, Activity, Sequence) để mô tả đầy đủ
- Các biểu đồ tổng quan (01-05) mô tả nhóm use cases theo module

