# Kiến trúc Hệ thống TabiMoney - AI-Powered Personal Finance Management

## Tổng quan Hệ thống

Hệ thống TabiMoney là một ứng dụng quản lý chi tiêu cá nhân thông minh, tích hợp AI Agent để cung cấp phân tích tài chính, dự đoán chi tiêu và tư vấn cá nhân hóa.

## Kiến trúc Tổng thể

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web App       │    │  Telegram Bot   │    │   Mobile App    │
│   (Vue.js +     │    │                 │    │   (Future)      │
│   Vuetify)      │    │                 │    │                 │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │     API Gateway           │
                    │   (Golang + Echo)         │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │      AI Agent Layer       │
                    │  - NLU Processing         │
                    │  - Expense Prediction     │
                    │  - Anomaly Detection      │
                    │  - Smart Categorization   │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │      Business Logic       │
                    │  - User Management        │
                    │  - Transaction Processing │
                    │  - Financial Analytics    │
                    │  - Goal Tracking          │
                    └─────────────┬─────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
┌─────────┴─────────┐    ┌─────────┴─────────┐    ┌─────────┴─────────┐
│   MySQL Database  │    │   Redis Cache     │    │   External APIs   │
│   - Users         │    │   - Session       │    │   - Gemini API     │
│   - Transactions  │    │   - Dashboard     │    │   - Email Service  │
│   - Categories    │    │   - Real-time     │    │   - Notification   │
│   - Goals         │    │   - AI Cache      │    │   - Analytics      │
│   - Analytics     │    │                   │    │                   │
└───────────────────┘    └───────────────────┘    └───────────────────┘
```

## Các Thành phần Chính

### 1. Frontend Layer
- **Web App**: Vue.js 3 + Vuetify 3 (Material Design)
- **Responsive Design**: Mobile-first approach
- **Real-time Updates**: WebSocket connection
- **PWA Support**: Offline capability

### 2. Backend Layer
- **API Gateway**: Golang + Echo Framework
- **Authentication**: JWT + Refresh Token
- **Rate Limiting**: Redis-based
- **CORS & Security**: Comprehensive security headers

### 3. AI Agent Layer
- **NLU Engine**: Google Gemini for natural language processing
- **Expense Prediction**: Machine Learning models
- **Anomaly Detection**: Statistical analysis + ML
- **Smart Categorization**: NLP + Classification algorithms

### 4. Data Layer
- **Primary Database**: MySQL 8.0
- **Cache Layer**: Redis 7.0
- **File Storage**: Local filesystem (future: S3)

## Luồng Dữ liệu

### 1. Nhập liệu Chi tiêu
```
User Input → NLU Processing → AI Categorization → Database Storage → Cache Update → Real-time Notification
```

### 2. Phân tích Tài chính
```
Historical Data → AI Analysis → Pattern Recognition → Prediction → Dashboard Update → User Notification
```

### 3. Chatbot Interaction
```
User Message → NLU Processing → Intent Recognition → Database Query → AI Response → User Reply
```

## Công nghệ Sử dụng

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
- **Real-time**: WebSocket

### AI & ML
- **NLU**: Google Gemini
- **Prediction**: Scikit-learn, TensorFlow
- **Anomaly Detection**: Isolation Forest, LSTM
- **Categorization**: NLP + Classification

## Bảo mật

- **Authentication**: JWT with refresh tokens
- **Authorization**: Role-based access control
- **Data Encryption**: AES-256 for sensitive data
- **API Security**: Rate limiting, CORS, CSRF protection
- **Input Validation**: Comprehensive validation and sanitization

## Scalability

- **Horizontal Scaling**: Stateless API design
- **Database**: Read replicas, connection pooling
- **Cache**: Redis clustering
- **Load Balancing**: Nginx reverse proxy
- **Monitoring**: Prometheus + Grafana

## Deployment

- **Containerization**: Docker
- **Orchestration**: Docker Compose (development)
- **Environment**: Development, Staging, Production
- **CI/CD**: GitHub Actions
