# Mock Data Generator

Script để tạo dữ liệu giả (mock data) cho việc test và thử nghiệm các tính năng AI như prediction và anomaly detection.

## Cách sử dụng

### 1. Build script

```bash
cd cmd/mockdata
go build -o mockdata main.go
```

Hoặc chạy trực tiếp:

```bash
go run cmd/mockdata/main.go [flags]
```

### 2. Chạy với các tùy chọn

```bash
# Tạo 200 transactions cho user ID 15, trải đều trong 6 tháng, bao gồm anomalies
go run cmd/mockdata/main.go -user=15 -count=200 -months=6 -anomalies=true

# Tạo nhiều transactions hơn cho testing prediction (cần nhiều dữ liệu)
go run cmd/mockdata/main.go -user=15 -count=500 -months=12 -anomalies=true

# Tạo dữ liệu không có anomalies (chỉ dữ liệu bình thường)
go run cmd/mockdata/main.go -user=15 -count=200 -months=6 -anomalies=false

# Sử dụng seed cố định để tạo dữ liệu reproducible
go run cmd/mockdata/main.go -user=15 -count=200 -months=6 -seed=12345

# Hoặc chạy với user ID mặc định (15)
go run cmd/mockdata/main.go -count=200 -months=6
```

## Các tham số

- `-user`: User ID để tạo dữ liệu (mặc định: 15)
- `-count`: Số lượng transactions cần tạo (mặc định: 200)
- `-months`: Số tháng để trải đều transactions (mặc định: 6)
- `-anomalies`: Có bao gồm anomaly transactions không (mặc định: true)
- `-seed`: Random seed để tạo dữ liệu reproducible (mặc định: 42)

## Dữ liệu được tạo

### Income Transactions
- Lương hàng tháng: 15-20 triệu VND
- Tự động tạo mỗi tháng trong khoảng thời gian chỉ định

### Expense Transactions
- Phân bổ theo các category hệ thống:
  - **Ăn uống** (35%): 30k - 300k VND
  - **Giao thông** (20%): 20k - 500k VND
  - **Mua sắm** (15%): 100k - 2M VND
  - **Giải trí** (10%): 50k - 500k VND
  - **Y tế** (5%): 100k - 5M VND
  - **Học tập** (5%): 50k - 1M VND
  - **Tiết kiệm** (5%): 500k - 5M VND
  - **Khác** (5%): 20k - 500k VND

### Anomaly Transactions (10% nếu enabled)
- **High amount anomalies**: Giao dịch có số tiền cao bất thường (3-10x mức bình thường)
- **Low amount anomalies**: Giao dịch có số tiền thấp bất thường
- **Odd timing**: Giao dịch vào giờ lạ (2-5 AM)
- Được đánh dấu với tag "anomaly" và prefix "[ANOMALY]" trong description

## Test các tính năng AI

### Cách 1: Sử dụng Python Test Script (Khuyến nghị) ⭐

Script này test giống như frontend gọi API qua Go backend:

```bash
# Test với user 15 qua frontend proxy (port 3000) - giống frontend thực tế
python3 cmd/mockdata/test_ai_features.py 15 http://localhost:3000

# Test với user 15 qua backend trực tiếp (port 8080) - nhanh hơn
python3 cmd/mockdata/test_ai_features.py 15 http://localhost:8080

# Test với user 1
python3 cmd/mockdata/test_ai_features.py 1 http://localhost:8080
```

**Lưu ý:**
- Script tự động login với test user 15 (test15@tabimoney.com / test123456)
- Với user khác, bạn cần cung cấp credentials hoặc token
- Frontend proxy (port 3000) giống như cách frontend thực tế gọi
- Backend trực tiếp (port 8080) nhanh hơn, không cần frontend chạy

### Cách 2: Sử dụng curl với Bearer Token

**Bước 1: Login để lấy token**
```bash
# Login với test user 15
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test15@tabimoney.com",
    "password": "test123456"
  }'
```

**Bước 2: Copy `access_token` từ response, sau đó:**

```bash
# Test Prediction (qua Go backend - giống frontend)
curl -X GET "http://localhost:8080/api/v1/analytics/predictions?start_date=2025-06-01&end_date=2025-12-31" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE"

# Test Anomaly Detection (qua Go backend - giống frontend)
curl -X GET "http://localhost:8080/api/v1/analytics/anomalies?start_date=2025-06-01&end_date=2025-12-31&threshold=0.6" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE"
```

### Cách 3: Test qua Frontend UI

1. Mở frontend tại `http://localhost:3000`
2. Login với test user 15 (test15@tabimoney.com / test123456) hoặc user 1
3. Vào trang **Analytics** để xem predictions và anomalies
4. Mở DevTools > Network để xem API calls thực tế

## Lưu ý

1. **User phải tồn tại**: Đảm bảo user ID bạn chỉ định đã tồn tại trong database
2. **Categories phải có**: Script cần các system categories đã được tạo (chạy migrations)
3. **Dữ liệu sẽ được thêm**: Script sẽ thêm transactions mới vào database, không xóa dữ liệu cũ
4. **Prediction cần nhiều dữ liệu**: Để test prediction tốt, nên tạo ít nhất 200-500 transactions trải đều trong 6-12 tháng

## Ví dụ workflow hoàn chỉnh

```bash
# 1. Tạo dữ liệu mock cho user 15
go run cmd/mockdata/main.go -user=15 -count=500 -months=12 -anomalies=true

# 2. Kiểm tra dữ liệu đã tạo (xem summary output)

# 3. Test prediction và anomaly detection bằng Python script
python3 cmd/mockdata/test_ai_features.py 15 http://localhost:3000

# Hoặc test qua curl:
# 3a. Login để lấy token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test15@tabimoney.com","password":"test123456"}' \
  | jq -r '.access_token')

# 3b. Test prediction
curl -X GET "http://localhost:8080/api/v1/analytics/predictions?start_date=2025-06-01&end_date=2025-12-31" \
  -H "Authorization: Bearer $TOKEN"

# 3c. Test anomaly detection
curl -X GET "http://localhost:8080/api/v1/analytics/anomalies?start_date=2025-06-01&end_date=2025-12-31&threshold=0.6" \
  -H "Authorization: Bearer $TOKEN"
```

