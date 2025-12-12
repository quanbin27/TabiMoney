# TabiMoney - AI-Powered Personal Finance Management

TabiMoney là một hệ thống quản lý chi tiêu cá nhân thông minh, tích hợp AI Agent để cung cấp phân tích tài chính, dự đoán chi tiêu và tư vấn cá nhân hóa.

🌐 Sản phẩm: https://tabimoney.site

## 🚀 Tính năng chính

### 📱 Chức năng đầy đủ của Ứng dụng Quản lý Chi tiêu Cá nhân hóa

1. **Nhập liệu chi tiêu**
   - Form thủ công: chọn danh mục, nhập số tiền, mô tả, ngày tháng
   - Ngôn ngữ tự nhiên (NLU): "tôi vừa ăn bún bò 50k" → hệ thống tự động nhận diện danh mục, số tiền, ngày

2. **Quản lý thu nhập & mục tiêu**
   - Khai báo thu nhập hàng tháng
   - Đặt mục tiêu tài chính (vd: tiết kiệm mua xe 240 triệu trong 12 tháng)
   - Theo dõi tiến độ tiết kiệm theo thời gian
   - Cảnh báo nếu chi tiêu vượt mức, ảnh hưởng đến mục tiêu
   - Gợi ý điều chỉnh chi tiêu để đạt mục tiêu

3. **Quản lý chi tiêu & giao dịch**
   - Lưu trữ chi tiêu chi tiết
   - Lọc/tìm kiếm theo ngày, danh mục, số tiền
   - Gom nhóm chi tiêu theo danh mục
   - So sánh chi tiêu tháng này với tháng trước

4. **Dashboard & phân tích tài chính**
   - Biểu đồ chi tiêu theo danh mục, theo thời gian
   - So sánh xu hướng chi tiêu qua các tháng
   - Hiển thị "sức khỏe tài chính" (tỷ lệ chi/thu, % tiết kiệm)

5. **Tư vấn cá nhân hóa (AI - Machine Learning)**
   - Dự báo chi tiêu tháng tới dựa vào dữ liệu lịch sử
   - Phát hiện bất thường: giao dịch quá lớn, sai lệch thói quen
   - Phân loại chi tiêu thông minh: học từ feedback của user
   - Đưa gợi ý cá nhân hóa: giảm/giữ/cân bằng chi tiêu theo thói quen riêng

6. **Chatbot & hỏi đáp ngôn ngữ tự nhiên**
   - Ghi chi tiêu bằng chat: "tôi mua cà phê 40k"
   - Hỏi đáp tài chính: "Tháng này tôi tiêu bao nhiêu cho ăn uống?"
   - Nhận phản hồi/gợi ý ngay trong Web App

7. **Thông báo & cảnh báo**
   - Realtime notification khi chi vượt ngưỡng
   - Nhắc nhở định kỳ về tiến độ mục tiêu tài chính

## 🏗️ Kiến trúc hệ thống

### Backend
- **Language**: Golang 1.21+
- **Framework**: Echo v4
- **Database**: MySQL 8.0, Redis 7.0
- **ORM**: GORM v2
- **Authentication**: JWT, bcrypt
- **AI Integration**: Google Gemini API, Custom ML models

### Frontend
- **Framework**: Vue.js 3
- **UI Library**: Vuetify 3
- **State Management**: Pinia
- **HTTP Client**: Axios
- **Charts**: Chart.js

### AI & ML
- **NLU**: Google Gemini
- **Prediction**: Scikit-learn, TensorFlow
- **Anomaly Detection**: Isolation Forest, LSTM
- **Categorization**: NLP + Classification

## 🛠️ Cài đặt và chạy

### Yêu cầu hệ thống
- Docker & Docker Compose
- Node.js 18+ (cho development)
- Go 1.21+ (cho development)

### Chạy với Docker Compose

1. **Clone repository**
```bash
git clone <repository-url>
cd TabiMoney
```

2. **Cấu hình environment**
```bash
cp config.env.example .env
# Chỉnh sửa file .env với thông tin của bạn
```

3. **Chạy hệ thống**
```bash
docker-compose up -d
```

4. **Truy cập ứng dụng**
- Sản phẩm: https://tabimoney.site
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Documentation: http://localhost:8080/docs

### Development

#### Backend
```bash
# Cài đặt dependencies
go mod download

# Chạy database migrations
go run cmd/server/main.go migrate

# Chạy server
go run cmd/server/main.go
```

#### Frontend
```bash
cd frontend

# Cài đặt dependencies
npm install

# Chạy development server
npm run dev
```

## 📊 Cơ sở dữ liệu

### MySQL Schema
- **Users**: Thông tin người dùng
- **Transactions**: Giao dịch thu/chi
- **Categories**: Danh mục chi tiêu
- **Financial Goals**: Mục tiêu tài chính
- **Budgets**: Ngân sách
- **AI Analysis**: Phân tích AI
- **Notifications**: Thông báo

### Redis Cache
- Session management
- Dashboard cache
- Rate limiting
- Real-time notifications

## 🤖 AI Features

### Natural Language Understanding (NLU)
- Xử lý ngôn ngữ tự nhiên để nhập giao dịch
- Phân tích intent và entities
- Gợi ý danh mục tự động

### Expense Prediction
- Dự đoán chi tiêu tháng tới
- Phân tích xu hướng chi tiêu
- Gợi ý tối ưu hóa ngân sách

### Anomaly Detection
- Phát hiện giao dịch bất thường
- Cảnh báo chi tiêu vượt mức
- Phân tích pattern chi tiêu

### Smart Categorization
- Tự động phân loại giao dịch
- Học từ feedback người dùng
- Cải thiện độ chính xác theo thời gian

## 🔐 Bảo mật

- **Authentication**: JWT với refresh tokens
- **Authorization**: Role-based access control
- **Data Encryption**: AES-256 cho dữ liệu nhạy cảm
- **API Security**: Rate limiting, CORS, CSRF protection
- **Input Validation**: Validation và sanitization toàn diện

## 📈 Monitoring & Logging

- **Logging**: Structured logging với logrus
- **Health Checks**: Database và Redis health monitoring
- **Metrics**: Performance metrics và error tracking
- **Alerts**: Real-time notifications cho hệ thống

## 🚀 Deployment

### Production
```bash
# Build và deploy với Docker
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=tabimoney
DB_PASSWORD=your_password
DB_NAME=tabimoney

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRE_HOURS=24

# Gemini (Required)
USE_GEMINI=true
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL=gemini-1.5-flash

# Server
SERVER_PORT=8080
CORS_ORIGINS=http://localhost:3000
```

## 📝 API Documentation

### Authentication Endpoints
- `POST /api/v1/auth/register` - Đăng ký
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/refresh` - Làm mới token
- `POST /api/v1/auth/logout` - Đăng xuất

### Transaction Endpoints
- `GET /api/v1/transactions` - Lấy danh sách giao dịch
- `POST /api/v1/transactions` - Tạo giao dịch mới
- `PUT /api/v1/transactions/:id` - Cập nhật giao dịch
- `DELETE /api/v1/transactions/:id` - Xóa giao dịch

### AI Endpoints
- `POST /api/v1/ai/nlu` - Xử lý ngôn ngữ tự nhiên
- `POST /api/v1/ai/predict-expenses` - Dự đoán chi tiêu
- `POST /api/v1/ai/detect-anomalies` - Phát hiện bất thường
- `POST {AI_SERVICE_URL}/api/v1/chat/process` - Chat với AI (Frontend & Telegram Bot gọi trực tiếp AI Service; backend không còn expose proxy)

## 🤝 Contributing

1. Fork repository
2. Tạo feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Tạo Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📞 Contact

- **Email**: contact@tabimoney.com
- **Website**: https://tabimoney.com
- **GitHub**: https://github.com/tabimoney

## 🙏 Acknowledgments

- Google Gemini cho AI API
- Vue.js và Vuetify cho frontend framework
- Echo framework cho Golang backend
- Chart.js cho data visualization
- Material Design Icons cho icon set
