"""
Database migration for Telegram bot integration
"""

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
);

-- Create telegram_link_codes table for temporary link codes
CREATE TABLE IF NOT EXISTS telegram_link_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(16) NOT NULL UNIQUE,
    telegram_user_id BIGINT NOT NULL,
    web_user_id BIGINT UNSIGNED,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (web_user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_code (code),
    INDEX idx_telegram_user_id (telegram_user_id),
    INDEX idx_expires_at (expires_at)
);

-- Add telegram integration settings to user profiles
ALTER TABLE user_profiles 
ADD COLUMN telegram_enabled BOOLEAN DEFAULT FALSE,
ADD COLUMN telegram_notifications BOOLEAN DEFAULT TRUE,
ADD COLUMN telegram_language VARCHAR(5) DEFAULT 'vi';
