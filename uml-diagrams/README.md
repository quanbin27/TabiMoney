# TabiMoney - UML Diagrams

Thư mục này chứa tất cả các biểu đồ UML của dự án TabiMoney ở định dạng PlantUML (.puml), có thể import vào StarUML hoặc các công cụ UML khác.

## Danh sách Biểu đồ

### Use Case Diagrams
1. **01_UseCase_Overall.puml** - Biểu đồ Use Case tổng quát của hệ thống
2. **02_UseCase_TransactionManagement.puml** - Use Case quản lý giao dịch
3. **03_UseCase_GoalManagement.puml** - Use Case quản lý mục tiêu tài chính
4. **04_UseCase_BudgetManagement.puml** - Use Case quản lý ngân sách
5. **05_UseCase_AIFeatures.puml** - Use Case các tính năng AI

### Sequence Diagrams
6. **06_Sequence_NLUTransaction.puml** - Luồng nhập giao dịch bằng NLU
7. **07_Sequence_GoalManagement.puml** - Luồng quản lý mục tiêu
8. **08_Sequence_AnomalyDetection.puml** - Luồng phát hiện bất thường
9. **09_Sequence_AIChat.puml** - Luồng chat với AI
10. **17_Sequence_Authentication.puml** - Luồng đăng ký và đăng nhập
11. **18_Sequence_DashboardAnalytics.puml** - Luồng Dashboard và Analytics

### Activity Diagrams
10. **10_Activity_TransactionEntry.puml** - Quy trình nhập giao dịch
11. **11_Activity_FinancialAnalysis.puml** - Quy trình phân tích tài chính
12. **12_Activity_BudgetManagement.puml** - Quy trình quản lý ngân sách
13. **19_Activity_ExpensePrediction.puml** - Quy trình dự đoán chi tiêu
14. **20_Activity_NotificationFlow.puml** - Quy trình thông báo

### State Diagrams
13. **13_State_Transaction.puml** - Trạng thái giao dịch
14. **14_State_FinancialGoal.puml** - Trạng thái mục tiêu tài chính
15. **15_State_Budget.puml** - Trạng thái ngân sách

### ERD Diagram
16. **16_ERD_TabiMoney.puml** - Biểu đồ Entity Relationship của database

## Cách Import vào StarUML

### Phương pháp 1: Sử dụng PlantUML Plugin trong StarUML

1. **Cài đặt PlantUML Plugin:**
   - Mở StarUML
   - Vào `Tools` > `Extension Manager`
   - Tìm và cài đặt plugin "PlantUML"

2. **Import file .puml:**
   - Vào `File` > `Import` > `PlantUML`
   - Chọn file .puml cần import
   - Biểu đồ sẽ được tự động chuyển đổi

### Phương pháp 2: Sử dụng PlantUML Online Editor

1. Truy cập: http://www.plantuml.com/plantuml/uml/
2. Copy nội dung file .puml vào editor
3. Export sang định dạng PNG, SVG, hoặc PDF
4. Import hình ảnh vào StarUML nếu cần

### Phương pháp 3: Sử dụng VS Code Extension

1. Cài đặt extension "PlantUML" trong VS Code
2. Mở file .puml
3. Nhấn `Alt+D` để preview
4. Export sang định dạng khác nếu cần

## Cách Import vào các công cụ khác

### Visual Paradigm
- Hỗ trợ import PlantUML trực tiếp
- Vào `File` > `Import` > `PlantUML`

### Draw.io / diagrams.net
- Sử dụng PlantUML syntax trong draw.io
- Hoặc import qua PlantUML online editor trước

### Enterprise Architect
- Sử dụng plugin PlantUML hoặc chuyển đổi qua format trung gian

## Cấu trúc Use Cases Chính

### 1. Quản lý Giao dịch
- Nhập giao dịch thủ công
- Nhập giao dịch bằng NLU (Natural Language Understanding)
- Xem, cập nhật, xóa giao dịch
- Phân loại tự động

### 2. Quản lý Mục tiêu Tài chính
- Tạo và quản lý mục tiêu
- Theo dõi tiến độ
- Đóng góp vào mục tiêu
- Cảnh báo mục tiêu

### 3. Quản lý Ngân sách
- Tạo và quản lý ngân sách
- Theo dõi chi tiêu theo ngân sách
- Cảnh báo vượt ngân sách
- Gợi ý ngân sách tự động

### 4. Tính năng AI
- Xử lý ngôn ngữ tự nhiên (NLU)
- Dự đoán chi tiêu
- Phát hiện bất thường
- Chat với AI
- Phân tích chi tiêu thông minh

### 5. Analytics & Dashboard
- Xem tổng quan tài chính
- Phân tích chi tiêu theo danh mục
- Phân tích xu hướng
- Financial health score

### 6. Thông báo
- Thông báo real-time
- Cảnh báo ngân sách
- Cảnh báo mục tiêu
- Thông báo giao dịch lớn

## Database Schema (ERD)

Biểu đồ ERD bao gồm các bảng chính:
- **users**: Thông tin người dùng
- **user_profiles**: Hồ sơ tài chính người dùng
- **categories**: Danh mục chi tiêu
- **transactions**: Giao dịch thu/chi
- **financial_goals**: Mục tiêu tài chính
- **budgets**: Ngân sách
- **notifications**: Thông báo
- **ai_analysis**: Phân tích AI
- **ai_feedback**: Phản hồi AI
- **user_sessions**: Phiên đăng nhập
- **telegram_accounts**: Tài khoản Telegram

## Lưu ý

- Tất cả các file sử dụng định dạng PlantUML chuẩn
- Có thể chỉnh sửa trực tiếp trong file .puml
- Màu sắc và style có thể tùy chỉnh trong file
- Các biểu đồ có thể được render online tại: http://www.plantuml.com/plantuml/uml/

## Tài liệu tham khảo

- PlantUML Documentation: https://plantuml.com/
- StarUML Documentation: http://staruml.io/
- UML 2.5 Specification: https://www.omg.org/spec/UML/

