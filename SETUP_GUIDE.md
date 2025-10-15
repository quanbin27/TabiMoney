# TabiMoney - Personal Finance Management System

## Tổng quan
TabiMoney là hệ thống quản lý tài chính cá nhân với các tính năng:
- Quản lý giao dịch (thu/chi)
- Phân loại giao dịch tự động bằng AI
- Quản lý ngân sách và mục tiêu tài chính
- Phân tích và báo cáo chi tiêu
- Trợ lý AI thông minh

## Kiến trúc hệ thống
- **Frontend**: Vue.js 3 + Vuetify 3
- **Backend**: Go + Echo framework
- **AI Service**: Python + FastAPI
- **Database**: MySQL 8.0
- **Cache**: Redis
- **Container**: Docker + Docker Compose

## Yêu cầu hệ thống
- Docker và Docker Compose
- Git
- Port 3000, 8080, 8000, 3306, 6379 phải trống

## Hướng dẫn cài đặt và chạy

### 1. Chuẩn bị môi trường
```bash
# Kiểm tra Docker
docker --version
docker-compose --version

# Clone project (nếu từ Git)
git clone <repository-url>
cd TabiMoney
```

### 2. Cấu hình môi trường
```bash
# Copy file cấu hình mẫu
cp config.env.example .env

# Chỉnh sửa file .env nếu cần
nano .env
```

**Nội dung file .env:**
```env
# Database
DB_HOST=mysql
DB_PORT=3306
DB_NAME=tabimoney
DB_USER=tabimoney
DB_PASSWORD=password

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRE_HOURS=24

# AI Service
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL=gemini-2.5-flash
OLLAMA_BASE_URL=http://ollama:11434
OLLAMA_MODEL=llama3.2

# Frontend
VITE_API_BASE_URL=http://localhost:8080
VITE_AI_SERVICE_URL=http://localhost:8000
```

### 3. Chạy dự án
```bash
# Build và khởi động tất cả services
docker-compose up -d --build

# Kiểm tra trạng thái services
docker-compose ps
```

### 4. Khởi tạo dữ liệu
```bash
# Chạy script tạo dữ liệu mẫu
docker exec tabimoney_mysql mysql -u root -ppassword tabimoney -e "source /docker-entrypoint-initdb.d/01-schema.sql"

# Tạo dữ liệu mẫu (tùy chọn)
docker exec tabimoney_backend ./generate_mock_data.sh
```

### 5. Truy cập ứng dụng
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **AI Service**: http://localhost:8000
- **Database**: localhost:3306
- **Redis**: localhost:6379

## Tài khoản mặc định
- **Email**: votrungquan2002@gmail.com
- **Password**: 123456

## Các lệnh hữu ích

### Xem logs
```bash
# Xem logs tất cả services
docker-compose logs -f

# Xem logs service cụ thể
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f ai-service
```

### Restart services
```bash
# Restart tất cả
docker-compose restart

# Restart service cụ thể
docker-compose restart backend
docker-compose restart frontend
```

### Dừng và xóa
```bash
# Dừng services
docker-compose down

# Dừng và xóa volumes (xóa dữ liệu)
docker-compose down -v
```

### Backup database
```bash
# Backup
docker exec tabimoney_mysql mysqldump -u root -ppassword tabimoney > backup.sql

# Restore
docker exec -i tabimoney_mysql mysql -u root -ppassword tabimoney < backup.sql
```

## Troubleshooting

### Lỗi thường gặp

1. **Port đã được sử dụng**
```bash
# Kiểm tra port
netstat -tulpn | grep :3000
netstat -tulpn | grep :8080

# Dừng process đang sử dụng port
sudo kill -9 <PID>
```

2. **Database connection failed**
```bash
# Kiểm tra MySQL
docker-compose logs mysql

# Restart MySQL
docker-compose restart mysql
```

3. **AI Service không hoạt động**
```bash
# Kiểm tra API key
echo $GEMINI_API_KEY

# Restart AI service
docker-compose restart ai-service
```

4. **Frontend build failed**
```bash
# Xóa node_modules và rebuild
docker-compose down
docker-compose build --no-cache frontend
docker-compose up -d
```

### Kiểm tra sức khỏe hệ thống
```bash
# Kiểm tra tất cả containers
docker-compose ps

# Kiểm tra logs lỗi
docker-compose logs | grep -i error

# Kiểm tra database
docker exec tabimoney_mysql mysql -u root -ppassword -e "SHOW DATABASES;"
```

## Cấu trúc dự án
```
TabiMoney/
├── frontend/          # Vue.js frontend
├── ai-service/        # Python AI service
├── internal/          # Go backend
├── database/          # Database schema
├── docker-compose.yml # Docker configuration
├── Dockerfile.backend # Backend Dockerfile
└── README.md
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Đăng ký
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/logout` - Đăng xuất

### Transactions
- `GET /api/v1/transactions` - Lấy danh sách giao dịch
- `POST /api/v1/transactions` - Tạo giao dịch mới
- `PUT /api/v1/transactions/:id` - Cập nhật giao dịch
- `DELETE /api/v1/transactions/:id` - Xóa giao dịch

### Budgets
- `GET /api/v1/budgets` - Lấy danh sách ngân sách
- `POST /api/v1/budgets` - Tạo ngân sách mới
- `PUT /api/v1/budgets/:id` - Cập nhật ngân sách
- `DELETE /api/v1/budgets/:id` - Xóa ngân sách

### Goals
- `GET /api/v1/goals` - Lấy danh sách mục tiêu
- `POST /api/v1/goals` - Tạo mục tiêu mới
- `PUT /api/v1/goals/:id` - Cập nhật mục tiêu
- `DELETE /api/v1/goals/:id` - Xóa mục tiêu

### AI Assistant
- `POST /api/v1/ai/chat` - Chat với AI
- `POST /api/v1/ai/nlu` - Natural Language Understanding

## Liên hệ hỗ trợ
Nếu gặp vấn đề, vui lòng tạo issue hoặc liên hệ qua email.

---
**Lưu ý**: Đảm bảo có đủ quyền để chạy Docker và các port cần thiết không bị chiếm dụng.
