# Hướng dẫn Import vào Draw.io

## Các file .drawio đã tạo

Thư mục này chứa **20 file** định dạng `.drawio` (XML format) có thể import trực tiếp vào draw.io và chỉnh sửa hoàn toàn.

### Danh sách đầy đủ:

**Use Case Diagrams (5 files):**
1. **01_UseCase_Overall.drawio** - Biểu đồ Use Case tổng quát
2. **05_UseCase_TransactionManagement.drawio** - Use Case quản lý giao dịch
3. **06_UseCase_GoalManagement.drawio** - Use Case quản lý mục tiêu
4. **07_UseCase_BudgetManagement.drawio** - Use Case quản lý ngân sách
5. **08_UseCase_AIFeatures.drawio** - Use Case tính năng AI

**Sequence Diagrams (6 files):**
6. **03_Sequence_NLUTransaction.drawio** - Luồng nhập giao dịch bằng NLU
7. **09_Sequence_GoalManagement.drawio** - Luồng quản lý mục tiêu
8. **10_Sequence_AnomalyDetection.drawio** - Luồng phát hiện bất thường
9. **11_Sequence_AIChat.drawio** - Luồng chat với AI
10. **17_Sequence_Authentication.drawio** - Luồng đăng ký/đăng nhập
11. **18_Sequence_DashboardAnalytics.drawio** - Luồng Dashboard và Analytics

**Activity Diagrams (5 files):**
12. **12_Activity_TransactionEntry.drawio** - Quy trình nhập giao dịch
13. **13_Activity_FinancialAnalysis.drawio** - Quy trình phân tích tài chính
14. **14_Activity_BudgetManagement.drawio** - Quy trình quản lý ngân sách
15. **19_Activity_ExpensePrediction.drawio** - Quy trình dự đoán chi tiêu
16. **20_Activity_NotificationFlow.drawio** - Quy trình thông báo

**State Diagrams (3 files):**
17. **04_State_FinancialGoal.drawio** - Trạng thái mục tiêu tài chính
18. **15_State_Transaction.drawio** - Trạng thái giao dịch
19. **16_State_Budget.drawio** - Trạng thái ngân sách

**ERD Diagram (1 file):**
20. **02_ERD_TabiMoney.drawio** - Biểu đồ ERD (Entity Relationship Diagram)

Xem file **INDEX.md** để biết chi tiết về từng biểu đồ.

## Cách Import vào Draw.io

### Phương pháp 1: Import trực tiếp (Khuyến nghị)

1. Truy cập https://app.diagrams.net/ (hoặc https://draw.io)
2. Chọn **File** > **Open from** > **Device**
3. Chọn file `.drawio` cần import
4. File sẽ được mở và bạn có thể chỉnh sửa ngay

### Phương pháp 2: Kéo thả file

1. Truy cập https://app.diagrams.net/
2. Kéo file `.drawio` vào cửa sổ trình duyệt
3. File sẽ tự động mở

### Phương pháp 3: Sử dụng Draw.io Desktop

1. Tải Draw.io Desktop từ: https://github.com/jgraph/drawio-desktop/releases
2. Cài đặt ứng dụng
3. Mở Draw.io Desktop
4. Chọn **File** > **Open** và chọn file `.drawio`

## Chỉnh sửa trong Draw.io

Sau khi import, bạn có thể:

### 1. Thêm/Sửa/Xóa các phần tử
- Click vào phần tử để chọn
- Kéo thả để di chuyển
- Double-click để chỉnh sửa text
- Sử dụng toolbar để thêm shapes mới

### 2. Thay đổi Style
- Chọn phần tử
- Sử dụng panel bên phải để thay đổi:
  - Màu sắc (Fill)
  - Màu viền (Stroke)
  - Font chữ
  - Kích thước

### 3. Thêm Relationships
- Chọn công cụ **Connector** từ toolbar
- Click vào phần tử nguồn
- Kéo đến phần tử đích
- Chọn style mũi tên phù hợp

### 4. Layout và Alignment
- Chọn nhiều phần tử (Ctrl+Click hoặc drag)
- Sử dụng **Arrange** menu để:
  - Align (căn chỉnh)
  - Distribute (phân bố đều)
  - Group/Ungroup

## Lưu file

### Lưu trên Cloud
- **File** > **Save to** > **Device** (lưu file .drawio)
- Hoặc chọn **OneDrive**, **Google Drive**, **GitHub**, etc.

### Export sang định dạng khác
- **File** > **Export as** > Chọn định dạng:
  - PNG (hình ảnh)
  - SVG (vector)
  - PDF
  - JPG
  - XML (draw.io format)

## Tips

1. **Sử dụng Layers**: Tổ chức biểu đồ phức tạp bằng layers
2. **Templates**: Draw.io có nhiều template UML sẵn có
3. **Keyboard Shortcuts**: 
   - `Ctrl+Z`: Undo
   - `Ctrl+Y`: Redo
   - `Ctrl+C/V`: Copy/Paste
   - `Delete`: Xóa phần tử
   - `Ctrl+G`: Group
   - `Ctrl+Shift+G`: Ungroup

4. **Shape Library**: 
   - Vào **More Shapes** để thêm thư viện UML
   - Tìm "UML" trong shape library

## Tạo thêm biểu đồ mới

Nếu bạn muốn tạo thêm các biểu đồ khác từ file PlantUML:

1. Mở file `.puml` trong PlantUML online: http://www.plantuml.com/plantuml/uml/
2. Export sang PNG hoặc SVG
3. Import hình ảnh vào Draw.io
4. Vẽ lại hoặc chỉnh sửa trên nền hình ảnh

Hoặc:

1. Sử dụng Draw.io từ đầu
2. Chọn template UML phù hợp
3. Vẽ biểu đồ theo cấu trúc từ file PlantUML

## Lưu ý

- File `.drawio` là định dạng XML, có thể mở bằng text editor để chỉnh sửa trực tiếp (không khuyến nghị)
- File này tương thích với tất cả các phiên bản Draw.io
- Có thể chia sẻ file này với team để cùng chỉnh sửa

