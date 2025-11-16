-- TabiMoney Database Schema
-- AI-Powered Personal Finance Management System

-- Create database
CREATE DATABASE IF NOT EXISTS tabimoney CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE tabimoney;

-- Users table
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    is_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    reset_token VARCHAR(255),
    reset_token_expires_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB;

-- User profiles for financial settings
CREATE TABLE user_profiles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    monthly_income DECIMAL(15,2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'VND',
    timezone VARCHAR(50) DEFAULT 'Asia/Ho_Chi_Minh',
    language VARCHAR(5) DEFAULT 'vi',
    notification_settings JSON,
    ai_settings JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_profile (user_id)
) ENGINE=InnoDB;

-- Categories for transactions
CREATE TABLE categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NULL, -- NULL for system categories
    name VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(7), -- HEX color
    parent_id BIGINT UNSIGNED NULL,
    is_system BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_is_system (is_system)
) ENGINE=InnoDB;

-- Financial goals
CREATE TABLE financial_goals (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    target_amount DECIMAL(15,2) NOT NULL,
    current_amount DECIMAL(15,2) DEFAULT 0.00,
    target_date DATE,
    goal_type ENUM('savings', 'debt_payment', 'investment', 'purchase', 'other') DEFAULT 'savings',
    priority ENUM('low', 'medium', 'high', 'urgent') DEFAULT 'medium',
    is_achieved BOOLEAN DEFAULT FALSE,
    achieved_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_goal_type (goal_type),
    INDEX idx_target_date (target_date)
) ENGINE=InnoDB;

-- Transactions
CREATE TABLE transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    description TEXT,
    transaction_type ENUM('income', 'expense', 'transfer') NOT NULL,
    transaction_date DATE NOT NULL,
    transaction_time TIME,
    location VARCHAR(200),
    tags JSON, -- Array of tags
    metadata JSON, -- Additional data like payment method, etc.
    is_recurring BOOLEAN DEFAULT FALSE,
    recurring_pattern VARCHAR(50), -- daily, weekly, monthly, yearly
    parent_transaction_id BIGINT UNSIGNED NULL, -- For recurring transactions
    ai_confidence DECIMAL(3,2), -- AI confidence score for categorization
    ai_suggested_category_id BIGINT UNSIGNED NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id) ON DELETE SET NULL,
    FOREIGN KEY (ai_suggested_category_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_category_id (category_id),
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_transaction_type (transaction_type),
    INDEX idx_amount (amount),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB;

-- AI Analysis and Predictions
CREATE TABLE ai_analysis (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    analysis_type ENUM('expense_prediction', 'anomaly_detection', 'category_suggestion', 'spending_pattern', 'goal_analysis') NOT NULL,
    data JSON NOT NULL, -- Analysis results
    confidence_score DECIMAL(3,2),
    model_version VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_analysis_type (analysis_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB;

-- User feedback for AI learning
CREATE TABLE ai_feedback (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    transaction_id BIGINT UNSIGNED NULL,
    feedback_type ENUM('category_correct', 'category_incorrect', 'prediction_accurate', 'prediction_inaccurate', 'suggestion_helpful', 'suggestion_not_helpful') NOT NULL,
    original_prediction JSON,
    user_correction JSON,
    feedback_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_feedback_type (feedback_type)
) ENGINE=InnoDB;

-- Notifications
CREATE TABLE notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    notification_type ENUM('info', 'warning', 'success', 'error', 'reminder') NOT NULL,
    priority ENUM('low', 'medium', 'high', 'urgent') DEFAULT 'medium',
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    action_url VARCHAR(500),
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB;

-- User sessions for JWT management
CREATE TABLE user_sessions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    refresh_expires_at TIMESTAMP NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(45),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB;

-- Budget tracking
CREATE TABLE budgets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NULL, -- NULL for total budget
    name VARCHAR(200) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    period ENUM('weekly', 'monthly', 'yearly') DEFAULT 'monthly',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    alert_threshold DECIMAL(5,2) DEFAULT 80.00, -- Alert when 80% used
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_category_id (category_id),
    INDEX idx_period (period),
    INDEX idx_start_date (start_date)
) ENGINE=InnoDB;

-- Insert default system categories
INSERT INTO categories (name, name_en, description, icon, color, is_system, is_active, sort_order) VALUES
('Ăn uống', 'Food & Dining', 'Chi phí ăn uống, nhà hàng, cà phê', 'restaurant', '#FF6B6B', TRUE, TRUE, 1),
('Giao thông', 'Transportation', 'Xăng xe, taxi, bus, vé máy bay', 'directions_car', '#4ECDC4', TRUE, TRUE, 2),
('Mua sắm', 'Shopping', 'Quần áo, đồ dùng cá nhân', 'shopping_bag', '#45B7D1', TRUE, TRUE, 3),
('Giải trí', 'Entertainment', 'Phim, game, du lịch', 'movie', '#96CEB4', TRUE, TRUE, 4),
('Y tế', 'Healthcare', 'Khám bệnh, thuốc, bảo hiểm', 'local_hospital', '#FFEAA7', TRUE, TRUE, 5),
('Học tập', 'Education', 'Khóa học, sách, tài liệu', 'school', '#DDA0DD', TRUE, TRUE, 6),
('Tiết kiệm', 'Savings', 'Tiết kiệm, đầu tư', 'savings', '#98D8C8', TRUE, TRUE, 7),
('Thu nhập', 'Income', 'Lương, thưởng, thu nhập khác', 'attach_money', '#F7DC6F', TRUE, TRUE, 8),
('Khác', 'Other', 'Chi phí khác', 'more_horiz', '#BB8FCE', TRUE, TRUE, 9);

-- Create indexes for better performance
CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date);
CREATE INDEX idx_transactions_user_type ON transactions(user_id, transaction_type);
CREATE INDEX idx_transactions_amount_range ON transactions(amount);

-- Create views for common queries
CREATE VIEW user_monthly_summary AS
SELECT 
    u.id as user_id,
    u.username,
    YEAR(t.transaction_date) as year,
    MONTH(t.transaction_date) as month,
    SUM(CASE WHEN t.transaction_type = 'income' THEN t.amount ELSE 0 END) as total_income,
    SUM(CASE WHEN t.transaction_type = 'expense' THEN t.amount ELSE 0 END) as total_expense,
    COUNT(CASE WHEN t.transaction_type = 'income' THEN 1 END) as income_count,
    COUNT(CASE WHEN t.transaction_type = 'expense' THEN 1 END) as expense_count
FROM users u
LEFT JOIN transactions t ON u.id = t.user_id
GROUP BY u.id, u.username, YEAR(t.transaction_date), MONTH(t.transaction_date);

CREATE VIEW category_spending AS
SELECT 
    t.user_id,
    c.name as category_name,
    c.icon,
    c.color,
    SUM(t.amount) as total_amount,
    COUNT(t.id) as transaction_count,
    AVG(t.amount) as avg_amount
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.transaction_type = 'expense'
GROUP BY t.user_id, c.id, c.name, c.icon, c.color;

-- Telegram Integration Tables
-- Create telegram_accounts table
CREATE TABLE IF NOT EXISTS telegram_accounts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    telegram_user_id BIGINT NOT NULL UNIQUE,
    web_user_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (web_user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_telegram_user_id (telegram_user_id),
    INDEX idx_web_user_id (web_user_id)
) ENGINE=InnoDB;

-- Create telegram_link_codes table for temporary link codes
CREATE TABLE IF NOT EXISTS telegram_link_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(16) NOT NULL UNIQUE,
    telegram_user_id BIGINT,
    web_user_id BIGINT UNSIGNED NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (web_user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_code (code),
    INDEX idx_telegram_user_id (telegram_user_id),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB;

-- Add telegram integration settings to user profiles
-- Using stored procedure to safely add columns only if they don't exist
DELIMITER $$

DROP PROCEDURE IF EXISTS AddTelegramColumnsIfNotExists$$

CREATE PROCEDURE AddTelegramColumnsIfNotExists()
BEGIN
    -- Check and add telegram_enabled column
    IF NOT EXISTS (
        SELECT * FROM INFORMATION_SCHEMA.COLUMNS 
        WHERE TABLE_SCHEMA = DATABASE() 
        AND TABLE_NAME = 'user_profiles' 
        AND COLUMN_NAME = 'telegram_enabled'
    ) THEN
        ALTER TABLE user_profiles ADD COLUMN telegram_enabled BOOLEAN DEFAULT FALSE;
    END IF;

    -- Check and add telegram_notifications column
    IF NOT EXISTS (
        SELECT * FROM INFORMATION_SCHEMA.COLUMNS 
        WHERE TABLE_SCHEMA = DATABASE() 
        AND TABLE_NAME = 'user_profiles' 
        AND COLUMN_NAME = 'telegram_notifications'
    ) THEN
        ALTER TABLE user_profiles ADD COLUMN telegram_notifications BOOLEAN DEFAULT TRUE;
    END IF;

    -- Check and add telegram_language column
    IF NOT EXISTS (
        SELECT * FROM INFORMATION_SCHEMA.COLUMNS 
        WHERE TABLE_SCHEMA = DATABASE() 
        AND TABLE_NAME = 'user_profiles' 
        AND COLUMN_NAME = 'telegram_language'
    ) THEN
        ALTER TABLE user_profiles ADD COLUMN telegram_language VARCHAR(5) DEFAULT 'vi';
    END IF;
END$$

DELIMITER ;

-- Execute the procedure to add columns safely
CALL AddTelegramColumnsIfNotExists();

-- Drop the procedure after use
DROP PROCEDURE IF EXISTS AddTelegramColumnsIfNotExists;
