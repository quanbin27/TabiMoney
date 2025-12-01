#!/bin/bash

# TabiMoney Deployment Script
# Usage: ./deploy.sh [options]
# Options:
#   --build      Force rebuild all images
#   --pull       Pull latest code from git (if using git)
#   --backup     Backup database before deployment
#   --restart   Restart all services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="$PROJECT_DIR/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    log_info "Docker and Docker Compose are installed"
}

check_env_file() {
    if [ ! -f "$PROJECT_DIR/.env" ]; then
        log_warn ".env file not found. Creating from config.env.example..."
        if [ -f "$PROJECT_DIR/config.env.example" ]; then
            cp "$PROJECT_DIR/config.env.example" "$PROJECT_DIR/.env"
            log_warn "Please edit .env file with your configuration before continuing!"
            exit 1
        else
            log_error "config.env.example not found. Cannot create .env file."
            exit 1
        fi
    fi
    log_info ".env file exists"
}

backup_database() {
    if [ ! -d "$BACKUP_DIR" ]; then
        mkdir -p "$BACKUP_DIR"
    fi
    
    log_info "Backing up database..."
    
    # Load DB credentials from .env
    source "$PROJECT_DIR/.env"
    
    DB_USER=${DB_USER:-tabimoney}
    DB_PASSWORD=${DB_PASSWORD:-password}
    DB_NAME=${DB_NAME:-tabimoney}
    
    BACKUP_FILE="$BACKUP_DIR/backup_${DATE}.sql"
    
    if docker exec tabimoney_mysql mysqldump -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" > "$BACKUP_FILE" 2>/dev/null; then
        log_info "Database backed up to $BACKUP_FILE"
        
        # Compress backup
        gzip "$BACKUP_FILE"
        log_info "Backup compressed to ${BACKUP_FILE}.gz"
        
        # Keep only last 10 backups
        ls -t "$BACKUP_DIR"/backup_*.sql.gz | tail -n +11 | xargs -r rm
        log_info "Old backups cleaned (keeping last 10)"
    else
        log_warn "Database backup failed (this is OK if database is not running)"
    fi
}

pull_latest_code() {
    if [ -d "$PROJECT_DIR/.git" ]; then
        log_info "Pulling latest code from git..."
        cd "$PROJECT_DIR"
        git pull origin main || git pull origin master || log_warn "Git pull failed or no remote configured"
    else
        log_warn "Not a git repository, skipping git pull"
    fi
}

build_images() {
    log_info "Building Docker images..."
    cd "$PROJECT_DIR"
    docker-compose build --no-cache
    log_info "Images built successfully"
}

start_services() {
    log_info "Starting services..."
    cd "$PROJECT_DIR"
    docker-compose up -d
    log_info "Services started"
}

wait_for_services() {
    log_info "Waiting for services to be healthy..."
    sleep 5
    
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -f http://localhost:8080/health > /dev/null 2>&1; then
            log_info "Backend is healthy"
            break
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 2
    done
    echo
    
    if [ $attempt -eq $max_attempts ]; then
        log_warn "Backend health check timeout (this might be OK if service is still starting)"
    fi
}

check_services() {
    log_info "Checking service status..."
    cd "$PROJECT_DIR"
    docker-compose ps
    
    log_info "Checking service health..."
    
    # Check backend
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_info "✓ Backend is healthy"
    else
        log_warn "✗ Backend health check failed"
    fi
    
    # Check AI service
    if curl -f http://localhost:8001/health > /dev/null 2>&1; then
        log_info "✓ AI Service is healthy"
    else
        log_warn "✗ AI Service health check failed"
    fi
    
    # Check frontend
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        log_info "✓ Frontend is accessible"
    else
        log_warn "✗ Frontend check failed"
    fi
}

show_logs() {
    log_info "Recent logs (last 20 lines per service):"
    cd "$PROJECT_DIR"
    docker-compose logs --tail=20
}

restart_services() {
    log_info "Restarting all services..."
    cd "$PROJECT_DIR"
    docker-compose restart
    log_info "Services restarted"
}

# Main deployment function
deploy() {
    log_info "🚀 Starting TabiMoney deployment..."
    
    check_docker
    check_env_file
    
    cd "$PROJECT_DIR"
    
    # Parse arguments
    BUILD=false
    PULL=false
    BACKUP=false
    RESTART=false
    
    for arg in "$@"; do
        case $arg in
            --build)
                BUILD=true
                ;;
            --pull)
                PULL=true
                ;;
            --backup)
                BACKUP=true
                ;;
            --restart)
                RESTART=true
                ;;
        esac
    done
    
    # Execute steps
    if [ "$BACKUP" = true ]; then
        backup_database
    fi
    
    if [ "$PULL" = true ]; then
        pull_latest_code
    fi
    
    if [ "$RESTART" = true ]; then
        restart_services
        wait_for_services
        check_services
        return
    fi
    
    if [ "$BUILD" = true ]; then
        build_images
    fi
    
    start_services
    wait_for_services
    check_services
    
    log_info "✅ Deployment completed!"
    
    echo ""
    log_info "Useful commands:"
    echo "  View logs:        docker-compose logs -f"
    echo "  Restart service:  docker-compose restart [service-name]"
    echo "  Stop all:         docker-compose down"
    echo "  View status:      docker-compose ps"
}

# Run deployment
deploy "$@"


