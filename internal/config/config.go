package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
	OpenAI   OpenAIConfig
	Email    EmailConfig
	Upload   UploadConfig
	RateLimit RateLimitConfig
	Logging  LoggingConfig
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret            string
	ExpireHours       int
	RefreshExpireHours int
}

type ServerConfig struct {
	Port        string
	Host        string
	CORSOrigins []string
}

type OpenAIConfig struct {
	APIKey     string
	Model      string
	MaxTokens  int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFromName string
}

type UploadConfig struct {
	MaxSize       int64
	AllowedTypes  []string
}

type RateLimitConfig struct {
	Requests int
	Window   int
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "tabimoney"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
			ExpireHours:       getEnvAsInt("JWT_EXPIRE_HOURS", 24),
			RefreshExpireHours: getEnvAsInt("JWT_REFRESH_EXPIRE_HOURS", 168),
		},
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", "localhost"),
			CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:8080"), ","),
		},
		OpenAI: OpenAIConfig{
			APIKey:    getEnv("OPENAI_API_KEY", ""),
			Model:     getEnv("OPENAI_MODEL", "gpt-4"),
			MaxTokens: getEnvAsInt("OPENAI_MAX_TOKENS", 1000),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			SMTPFromName: getEnv("SMTP_FROM_NAME", "TabiMoney"),
		},
		Upload: UploadConfig{
			MaxSize:      getEnvAsInt64("UPLOAD_MAX_SIZE", 10485760), // 10MB
			AllowedTypes: strings.Split(getEnv("UPLOAD_ALLOWED_TYPES", "image/jpeg,image/png,image/gif"), ","),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", 1000),
			Window:   getEnvAsInt("RATE_LIMIT_WINDOW", 60),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Environment: getEnv("ENV", "development"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *Config) GetJWTExpiration() time.Duration {
	return time.Duration(c.JWT.ExpireHours) * time.Hour
}

func (c *Config) GetJWTRefreshExpiration() time.Duration {
	return time.Duration(c.JWT.RefreshExpireHours) * time.Hour
}

func (c *Config) GetDatabaseDSN() string {
	return c.Database.User + ":" + c.Database.Password + "@tcp(" + c.Database.Host + ":" + strconv.Itoa(c.Database.Port) + ")/" + c.Database.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port)
}

func (c *Config) GetServerAddr() string {
	return c.Server.Host + ":" + c.Server.Port
}
