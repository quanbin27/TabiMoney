# Kiến trúc Hệ thống TabiMoney - AI-Powered Personal Finance Management

## Tổng quan Hệ thống

Hệ thống TabiMoney là một ứng dụng quản lý chi tiêu cá nhân thông minh, tích hợp AI Agent để cung cấp phân tích tài chính, dự đoán chi tiêu và tư vấn cá nhân hóa.

## Kiến trúc Tổng thể

```
┌─────────────────┐              ┌─────────────────┐
│   Web App       │              │  Telegram Bot   │
│   (Vue.js +     │              │   (Python)      │
│   Vuetify)      │              │                 │
└───┬─────────┬───┘              └───┬─────────┬───┘
    │         │                     │         │
    │         │                     │         │
    │         └───────────┐         │         └───────────┐
    │                     │         │                     │
    │         ┌───────────┴─────────┴───────────┐         │
    │         │      AI Service                │         │
    │         │   (Python + FastAPI)            │         │
    │         │  - NLU/Chat Processing          │         │
    │         │  - Expense Prediction           │         │
    │         │  - Anomaly Detection             │         │
    │         └───────────┬─────────────────────┘         │
    │                     │                               │
    │         ┌───────────┴───────────────────────┐       │
    │         │   Backend Service                  │       │
    │         │  (Golang + Echo Framework)          │       │
    │         │  - Authentication                   │       │
    │         │  - Transaction Management           │       │
    │         │  - Budget Management                │       │
    │         │  - Goal Tracking                    │       │
    │         │  - Analytics                        │       │
    │         │  - Gọi AI Service cho Prediction     │       │
    │         └──────┬───────────────────┬─────────┘
    │                │                   │
┌───┴────────┐       │        ┌─────────┴────────┐
│   MySQL    │       │        │     Redis        │
│  Database  │       │        │     Cache         │
└────────────┘       │        └──────────────────┘
                     │
              ┌──────┴────────┐
              │ External APIs │
              │ - Gemini API  │
              │ - Email       │
              └───────────────┘
```

## Các Thành phần Chính

### 1. Frontend Layer
- **Web App**: Vue.js 3 + Vuetify 3 (Material Design)
- **Responsive Design**: Mobile-first approach
- **Real-time Updates**: WebSocket connection
- **PWA Support**: Offline capability

### 2. Backend Layer
- **Backend Service**: Golang + Echo Framework
- **Authentication**: JWT + Refresh Token
- **Rate Limiting**: Redis-based
- **CORS & Security**: Comprehensive security headers
- **Chức năng**: Quản lý transactions, budgets, goals, analytics, và gọi AI Service cho prediction/anomaly detection

### 3. AI Service Layer
- **NLU Engine**: Google Gemini for natural language processing
- **Chat Service**: Xử lý chat (được gọi trực tiếp từ Frontend và Telegram Bot)
- **Expense Prediction**: Machine Learning models (được gọi từ Backend)
- **Anomaly Detection**: Statistical analysis + ML (được gọi từ Backend)
- **Smart Categorization**: NLP + Classification algorithms

### 4. Data Layer
- **Primary Database**: MySQL 8.0
- **Cache Layer**: Redis 7.0
- **File Storage**: Local filesystem (future: S3)

## Luồng Dữ liệu

### 1. Chat với AI (Frontend/Telegram Bot)
```
User Message → Frontend/Telegram Bot → AI Service (Chat Processing)
  → Gemini API (NLU) → Tự động tạo Transaction hoặc Query Data
  → Response về User
```

### 2. Nhập liệu Chi tiêu (Thủ công)
```
User Input (Form) → Frontend → Backend Service → Database Storage → Cache Update → Notification
```

### 3. Phân tích Tài chính
```
User Request → Frontend → Backend Service → AI Service (Prediction/Anomaly)
  → ML Processing → Backend tổng hợp → Cache → Response về User
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
