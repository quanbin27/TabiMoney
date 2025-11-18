package database

import (
	"fmt"
	"log"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) error {
	dsn := cfg.GetDatabaseDSN()
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}

func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Auto migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.UserSession{},
		&models.Category{},
		&models.Transaction{},
		&models.FinancialGoal{},
		&models.Budget{},
		&models.AIAnalysis{},
		&models.Notification{},
		&models.ChatMessage{},
		&models.TelegramAccount{},
		&models.TelegramLinkCode{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func CloseDatabase() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Transaction executes a function within a database transaction
func Transaction(fn func(*gorm.DB) error) error {
	return DB.Transaction(fn)
}

// Health check for database connection
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
