# Danh sách đầy đủ các file .drawio

Tất cả các biểu đồ UML đã được tạo ở định dạng `.drawio` (XML format) để import trực tiếp vào draw.io và chỉnh sửa.

## Use Case Diagrams (5 files)

1. **01_UseCase_Overall.drawio** - Biểu đồ Use Case tổng quát của hệ thống
   - 23 use cases chính
   - 3 actors: User, Telegram Bot, AI Service
   - Các relationships: extends, includes

2. **05_UseCase_TransactionManagement.drawio** - Use Case quản lý giao dịch
   - 10 use cases: Nhập thủ công, NLU, Xem, Lọc, Tìm kiếm, Cập nhật, Xóa, Chi tiết, Định kỳ, Phân loại tự động

3. **06_UseCase_GoalManagement.drawio** - Use Case quản lý mục tiêu tài chính
   - 8 use cases: Tạo, Xem, Cập nhật, Xóa, Đóng góp, Xem tiến độ, Đánh dấu hoàn thành, Cảnh báo

4. **07_UseCase_BudgetManagement.drawio** - Use Case quản lý ngân sách
   - 8 use cases: Tạo, Xem, Cập nhật, Xóa, Thống kê, Cảnh báo, Gợi ý tự động, Tạo từ gợi ý

5. **08_UseCase_AIFeatures.drawio** - Use Case các tính năng AI
   - 7 use cases: NLU, Phân loại, Dự đoán, Phát hiện bất thường, Chat, Phân tích, Gợi ý
   - 2 actors: User, AI Service

## Sequence Diagrams (6 files)

6. **03_Sequence_NLUTransaction.drawio** - Luồng nhập giao dịch bằng NLU
   - Actors: User, Frontend, Backend, AI Service, NLU Service, Database
   - 13 messages từ nhập câu đến hiển thị kết quả

7. **09_Sequence_GoalManagement.drawio** - Luồng quản lý mục tiêu
   - Actors: User, Frontend, Backend, Goal Service, Database, Notification Service
   - 2 scenarios: Tạo mục tiêu, Đóng góp vào mục tiêu
   - Alt block: Mục tiêu đạt được

8. **10_Sequence_AnomalyDetection.drawio** - Luồng phát hiện bất thường
   - Actors: User, Frontend, Backend, Analytics Handler, AI Service, Anomaly Service, Database, Notification Service
   - Alt block: Có bất thường nghiêm trọng

9. **11_Sequence_AIChat.drawio** - Luồng chat với AI
   - Actors: User, Frontend, Backend, AI Handler, AI Service, Chat Service, Database
   - Alt block: Câu hỏi về dữ liệu

10. **17_Sequence_Authentication.drawio** - Luồng đăng ký và đăng nhập
    - Actors: User, Frontend, Backend, Auth Service, Database, JWT Service
    - 2 sections: Đăng ký (với alt blocks), Đăng nhập (với alt blocks)

11. **18_Sequence_DashboardAnalytics.drawio** - Luồng Dashboard và Analytics
    - Actors: User, Frontend, Backend, Analytics Handler, Transaction Service, Database, Redis Cache
    - 2 scenarios: Dashboard (với cache hit/miss), Category Spending

## Activity Diagrams (5 files)

12. **12_Activity_TransactionEntry.drawio** - Quy trình nhập giao dịch
    - Decision points: Phương thức (Thủ công/NLU), User xác nhận, Dữ liệu hợp lệ, Category hợp lệ, Vượt ngân sách, Giao dịch lớn
    - 19 activities từ chọn phương thức đến hiển thị thành công

13. **13_Activity_FinancialAnalysis.drawio** - Quy trình phân tích tài chính
    - Decision points: Có cache, Có đủ dữ liệu cho AI
    - 17 activities từ yêu cầu đến trả về kết quả

14. **14_Activity_BudgetManagement.drawio** - Quy trình quản lý ngân sách
    - Decision points: Dữ liệu hợp lệ, Usage > Alert Threshold, Usage >= 100%, Chi tiêu vượt tốc độ
    - 13 activities từ tạo/cập nhật đến hiển thị kết quả

15. **19_Activity_ExpensePrediction.drawio** - Quy trình dự đoán chi tiêu
    - Decision points: Đủ dữ liệu, Confidence > threshold
    - 15 activities từ yêu cầu đến trả về kết quả

16. **20_Activity_NotificationFlow.drawio** - Quy trình thông báo
    - Decision points: User bật thông báo, Loại thông báo (Budget/Goal/Large Transaction/Scheduled), Email enabled, In-app enabled, Telegram enabled
    - 16 activities từ trigger đến gửi qua các channels

## State Diagrams (3 files)

17. **04_State_FinancialGoal.drawio** - Trạng thái mục tiêu tài chính
    - States: Created, Active, InProgress, Paused, Achieved, Cancelled
    - Composite state: InProgress (Tracking, Alert)

18. **15_State_Transaction.drawio** - Trạng thái giao dịch
    - States: Draft, Pending, Confirmed, Updated, Cancelled, Deleted
    - Composite state: Confirmed (Active, Recurring)

19. **16_State_Budget.drawio** - Trạng thái ngân sách
    - States: Created, Active, Monitoring, Warning, Exceeded, Paused, Completed, Cancelled
    - Composite state: Monitoring (Normal, Alert, Critical)

## ERD Diagram (1 file)

20. **02_ERD_TabiMoney.drawio** - Biểu đồ Entity Relationship
    - 11 bảng: users, user_profiles, categories, transactions, financial_goals, budgets, notifications, ai_analysis, ai_feedback, user_sessions, telegram_accounts
    - Tất cả relationships và foreign keys

---

## Tổng kết

- **Tổng số file**: 20 file .drawio
- **Use Case Diagrams**: 5 files
- **Sequence Diagrams**: 6 files
- **Activity Diagrams**: 5 files
- **State Diagrams**: 3 files
- **ERD Diagram**: 1 file

Tất cả các file đều tương ứng 1:1 với các file .puml trong thư mục cha và có thể import trực tiếp vào draw.io để chỉnh sửa.

