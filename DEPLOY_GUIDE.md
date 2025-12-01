# Hướng dẫn Deploy TabiMoney lên Server Hosting

## 📋 Tổng quan

Dự án TabiMoney bao gồm:
- **Backend**: Go API (port 8080)
- **Frontend**: Vue.js + Nginx (port 3000)
- **AI Service**: Python FastAPI (port 8001)
- **Telegram Bot**: Python
- **Database**: MySQL 8.0 (port 3306)
- **Cache**: Redis (port 6379)

## 🔧 Yêu cầu Server

- **OS**: Ubuntu 20.04+ hoặc Debian 11+ (khuyến nghị)
- **RAM**: Tối thiểu 2GB (khuyến nghị 4GB+)
- **Disk**: Tối thiểu 20GB
- **Docker**: Version 20.10+
- **Docker Compose**: Version 2.0+
- **Ports**: 22 (SSH), 80, 443, 3000, 8080, 8001 (có thể đóng các port này và chỉ mở 80, 443 nếu dùng Nginx reverse proxy)

## 📦 Bước 1: Chuẩn bị Server

### 1.1. Kết nối SSH vào server

```bash
ssh username@your-server-ip
# Ví dụ: ssh root@123.45.67.89
```

### 1.2. Cập nhật hệ thống

```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# Cài đặt các công cụ cần thiết
sudo apt install -y curl wget git vim ufw
```

### 1.3. Cài đặt Docker

```bash
# Xóa Docker cũ (nếu có)
sudo apt remove -y docker docker-engine docker.io containerd runc

# Cài đặt Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Thêm user vào group docker (thay 'username' bằng user của bạn)
sudo usermod -aG docker $USER
# Hoặc nếu dùng root:
# sudo usermod -aG docker root

# Khởi động lại session hoặc chạy:
newgrp docker

# Cài đặt Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Kiểm tra cài đặt
docker --version
docker-compose --version
```

### 1.4. Cấu hình Firewall (UFW)

```bash
# Cho phép SSH
sudo ufw allow 22/tcp

# Cho phép HTTP và HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Cho phép các port của ứng dụng (tùy chọn, nếu không dùng Nginx reverse proxy)
sudo ufw allow 3000/tcp
sudo ufw allow 8080/tcp
sudo ufw allow 8001/tcp

# Kích hoạt firewall
sudo ufw enable

# Kiểm tra trạng thái
sudo ufw status
```

## 📥 Bước 2: Upload Code lên Server

### Cách 1: Sử dụng Git (Khuyến nghị)

```bash
# Tạo thư mục cho dự án
mkdir -p ~/projects
cd ~/projects

# Clone repository (nếu có Git repo)
git clone <your-repository-url> TabiMoney
cd TabiMoney

# Hoặc nếu repo private, cần setup SSH key hoặc token
```

### Cách 2: Upload qua SCP từ máy local

```bash
# Từ máy local của bạn, chạy lệnh:
# scp -r /Users/quanbin27/GolandProjects/TabiMoney username@your-server-ip:~/projects/

# Sau đó trên server:
cd ~/projects/TabiMoney
```

### Cách 3: Sử dụng rsync (tốt nhất cho sync code)

```bash
# Từ máy local:
rsync -avz --exclude 'node_modules' --exclude '.git' --exclude 'venv' \
  /Users/quanbin27/GolandProjects/TabiMoney/ \
  username@your-server-ip:~/projects/TabiMoney/
```

## ⚙️ Bước 3: Cấu hình Environment

### 3.1. Tạo file .env

```bash
cd ~/projects/TabiMoney
cp config.env.example .env
nano .env
```

### 3.2. Cấu hình file .env cho Production

```env
# Database Configuration
DB_HOST=mysql
DB_PORT=3306
DB_USER=tabimoney
DB_PASSWORD=CHANGE_THIS_TO_STRONG_PASSWORD
DB_NAME=tabimoney

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration - QUAN TRỌNG: Đổi thành secret key mạnh
JWT_SECRET=CHANGE_THIS_TO_A_VERY_LONG_RANDOM_STRING_AT_LEAST_32_CHARS
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_HOURS=168

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
CORS_ORIGINS=http://YOUR_SERVER_IP:3000,http://YOUR_SERVER_IP

# AI Service URL (trong Docker network - cho backend)
AI_SERVICE_URL=http://ai-service:8001

# Frontend AI Service URL - QUAN TRỌNG: Dùng relative path để hoạt động với Nginx proxy
# Nếu dùng Nginx reverse proxy: /ai-service
# Nếu không dùng Nginx: http://YOUR_SERVER_IP:8001
VITE_AI_SERVICE_URL=/ai-service

# Gemini Configuration - BẮT BUỘC
USE_GEMINI=true
GEMINI_API_KEY=your-actual-gemini-api-key-here
GEMINI_MODEL=gemini-1.5-flash
GEMINI_MAX_TOKENS=512
GEMINI_TEMPERATURE=0.3

# Email Configuration (nếu cần)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=TabiMoney

# Telegram Bot Configuration - BẮT BUỘC nếu dùng Telegram bot
TELEGRAM_BOT_TOKEN=your-telegram-bot-token-from-botfather

# Environment
ENV=production

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

**Lưu ý quan trọng:**
- Thay `YOUR_SERVER_IP` bằng IP thực của server
- Đổi `DB_PASSWORD` thành mật khẩu mạnh
- Đổi `JWT_SECRET` thành chuỗi ngẫu nhiên dài (ít nhất 32 ký tự)
- Thêm `GEMINI_API_KEY` thực từ Google AI Studio
- Thêm `TELEGRAM_BOT_TOKEN` nếu dùng Telegram bot

### 3.3. Tạo JWT Secret mạnh

```bash
# Tạo random secret key
openssl rand -base64 32
# Copy kết quả vào JWT_SECRET trong file .env
```

## 🐳 Bước 4: Deploy với Docker Compose

### 4.1. Build và khởi động services

```bash
cd ~/projects/TabiMoney

# Build và start tất cả services
docker-compose up -d --build

# Kiểm tra trạng thái
docker-compose ps

# Xem logs
docker-compose logs -f
```

### 4.2. Kiểm tra services đã chạy

```bash
# Kiểm tra tất cả containers
docker ps

# Kiểm tra logs từng service
docker-compose logs backend
docker-compose logs frontend
docker-compose logs ai-service
docker-compose logs mysql
docker-compose logs redis
docker-compose logs telegram-bot

# Kiểm tra health
curl http://localhost:8080/health
curl http://localhost:8001/health
curl http://localhost:3000
```

### 4.3. Kiểm tra database

```bash
# Kết nối MySQL
docker exec -it tabimoney_mysql mysql -u tabimoney -p
# Nhập password từ .env

# Kiểm tra databases
SHOW DATABASES;
USE tabimoney;
SHOW TABLES;
```

## 🌐 Bước 5: Cấu hình Nginx Reverse Proxy (Khuyến nghị)

Nginx sẽ giúp:
- Truy cập qua IP mà không cần port
- Dễ dàng thêm domain sau này
- SSL/HTTPS dễ dàng hơn

### 5.1. Cài đặt Nginx

```bash
sudo apt install -y nginx
sudo systemctl enable nginx
sudo systemctl start nginx
```

### 5.2. Tạo cấu hình Nginx

```bash
sudo nano /etc/nginx/sites-available/tabimoney
```

Nội dung file (thay `YOUR_SERVER_IP` bằng IP thực):

```nginx
# Frontend - Port 80
server {
    listen 80;
    server_name YOUR_SERVER_IP;  # Hoặc _ để accept mọi request

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 10s;
        proxy_send_timeout 120s;
        proxy_read_timeout 120s;
    }

    # AI Service (nếu frontend cần gọi trực tiếp)
    location /ai-service/ {
        proxy_pass http://localhost:8001/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5.3. Kích hoạt cấu hình

```bash
# Tạo symbolic link
sudo ln -s /etc/nginx/sites-available/tabimoney /etc/nginx/sites-enabled/

# Xóa default config (tùy chọn)
sudo rm /etc/nginx/sites-enabled/default

# Test cấu hình
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

### 5.4. Cập nhật Frontend để dùng Nginx

**QUAN TRỌNG**: Frontend đã được cấu hình để dùng relative path `/ai-service` thay vì `localhost:8001`. Điều này cho phép:
- Browser tự động dùng domain/IP hiện tại
- Hoạt động với Nginx reverse proxy
- Không cần thay đổi khi có domain

Cần cập nhật file `.env` để frontend biết API URL:

```bash
nano ~/projects/TabiMoney/.env
```

Đảm bảo có:
```env
VITE_API_BASE_URL=/api/v1
VITE_AI_SERVICE_URL=/ai-service
```

Sau đó rebuild frontend (vì Vite embed env vars vào code khi build):
```bash
cd ~/projects/TabiMoney
docker-compose up -d --build frontend
```

**Lưu ý**: Nếu bạn không dùng Nginx reverse proxy và muốn truy cập trực tiếp qua port, có thể set:
```env
VITE_AI_SERVICE_URL=http://YOUR_SERVER_IP:8001
```
Nhưng khuyến nghị dùng Nginx với relative path `/ai-service` để dễ dàng thêm domain và SSL sau này.

## 🔒 Bước 6: Bảo mật Cơ bản

### 6.1. Đổi mật khẩu root MySQL

```bash
# Vào MySQL container
docker exec -it tabimoney_mysql mysql -u root -p

# Trong MySQL:
ALTER USER 'root'@'%' IDENTIFIED BY 'NEW_STRONG_PASSWORD';
FLUSH PRIVILEGES;
EXIT;
```

### 6.2. Giới hạn truy cập MySQL từ bên ngoài

Trong `docker-compose.yml`, xóa hoặc comment dòng:
```yaml
ports:
  - "3306:3306"  # Xóa dòng này để MySQL chỉ accessible trong Docker network
```

### 6.3. Tạo user non-root cho Docker (khuyến nghị)

```bash
# Tạo user mới
sudo adduser deploy
sudo usermod -aG docker deploy
sudo usermod -aG sudo deploy

# Chuyển ownership của project
sudo chown -R deploy:deploy ~/projects/TabiMoney

# Đăng nhập bằng user mới
su - deploy
```

## ✅ Bước 7: Kiểm tra và Test

### 7.1. Kiểm tra từ browser

Mở browser và truy cập:
- `http://YOUR_SERVER_IP` - Frontend
- `http://YOUR_SERVER_IP/api/v1/health` - Backend health check
- `http://YOUR_SERVER_IP/ai-service/health` - AI Service health check

### 7.2. Test đăng ký/đăng nhập

1. Truy cập frontend
2. Đăng ký tài khoản mới
3. Đăng nhập
4. Test các chức năng cơ bản

### 7.3. Kiểm tra logs

```bash
# Xem logs real-time
docker-compose logs -f

# Xem logs từng service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f ai-service
docker-compose logs -f telegram-bot
```

## 🔄 Bước 8: Tự động hóa Deploy (Tùy chọn)

### 8.1. Tạo script deploy

Tạo file `deploy.sh`:

```bash
cd ~/projects/TabiMoney
nano deploy.sh
```

Nội dung:
```bash
#!/bin/bash
set -e

echo "🚀 Starting deployment..."

# Pull latest code (nếu dùng Git)
# git pull origin main

# Backup database (tùy chọn)
# docker exec tabimoney_mysql mysqldump -u tabimoney -p$DB_PASSWORD tabimoney > backup_$(date +%Y%m%d_%H%M%S).sql

# Rebuild và restart
docker-compose down
docker-compose up -d --build

# Wait for services
sleep 10

# Check health
echo "Checking services..."
curl -f http://localhost:8080/health || echo "Backend health check failed"
curl -f http://localhost:8001/health || echo "AI Service health check failed"

echo "✅ Deployment completed!"
docker-compose ps
```

Cấp quyền thực thi:
```bash
chmod +x deploy.sh
```

Sử dụng:
```bash
./deploy.sh
```

## 📝 Bước 9: Thêm Domain sau này (Khi có domain)

### 9.1. Cập nhật DNS

Thêm A record trỏ domain về IP server:
```
A    @    123.45.67.89
A    www  123.45.67.89
```

### 9.2. Cập nhật Nginx config

```bash
sudo nano /etc/nginx/sites-available/tabimoney
```

Thay đổi:
```nginx
server_name yourdomain.com www.yourdomain.com;
```

### 9.3. Cài đặt SSL với Let's Encrypt

```bash
# Cài đặt Certbot
sudo apt install -y certbot python3-certbot-nginx

# Lấy certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal đã được setup tự động
```

### 9.4. Cập nhật .env

Cập nhật CORS_ORIGINS:
```env
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

Rebuild:
```bash
docker-compose up -d --build backend frontend
```

## 🛠️ Các lệnh hữu ích

### Xem logs
```bash
docker-compose logs -f [service-name]
```

### Restart service
```bash
docker-compose restart [service-name]
```

### Stop tất cả
```bash
docker-compose down
```

### Stop và xóa volumes (xóa dữ liệu)
```bash
docker-compose down -v
```

### Backup database
```bash
docker exec tabimoney_mysql mysqldump -u tabimoney -p$DB_PASSWORD tabimoney > backup.sql
```

### Restore database
```bash
docker exec -i tabimoney_mysql mysql -u tabimoney -p$DB_PASSWORD tabimoney < backup.sql
```

### Xem resource usage
```bash
docker stats
```

### Clean up Docker
```bash
# Xóa images không dùng
docker image prune -a

# Xóa volumes không dùng
docker volume prune
```

## 🐛 Troubleshooting

### Service không start

```bash
# Xem logs chi tiết
docker-compose logs [service-name]

# Kiểm tra port đã bị chiếm
sudo netstat -tulpn | grep :8080
sudo netstat -tulpn | grep :3000

# Restart service
docker-compose restart [service-name]
```

### Database connection failed

```bash
# Kiểm tra MySQL
docker-compose logs mysql

# Kiểm tra network
docker network ls
docker network inspect tabimoney_tabimoney_network

# Test connection
docker exec tabimoney_backend ping mysql
```

### Frontend không load

```bash
# Kiểm tra Nginx
sudo nginx -t
sudo systemctl status nginx

# Kiểm tra frontend container
docker-compose logs frontend
docker-compose ps frontend
```

### AI Service không hoạt động

```bash
# Kiểm tra API key
docker-compose exec ai-service env | grep GEMINI

# Kiểm tra logs
docker-compose logs ai-service
```

## 📞 Hỗ trợ

Nếu gặp vấn đề, kiểm tra:
1. Logs của service: `docker-compose logs [service]`
2. Health checks: `curl http://localhost:8080/health`
3. Firewall: `sudo ufw status`
4. Disk space: `df -h`
5. Memory: `free -h`

---

**Lưu ý**: 
- Đảm bảo thay đổi tất cả passwords và secrets trong file `.env`
- Backup database thường xuyên
- Monitor logs để phát hiện lỗi sớm
- Khi có domain, cài SSL ngay lập tức


