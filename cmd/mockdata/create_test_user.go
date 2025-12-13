package main

import (
	"fmt"
	"log"

	"tabimoney/internal/config"
	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	db := database.GetDB()

	// Check if user 15 exists
	var user models.User
	if err := db.First(&user, 15).Error; err == nil {
		fmt.Printf("✅ User ID 15 already exists: %s (%s)\n", user.Username, user.Email)
		return
	}

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123456"), bcrypt.DefaultCost)

	user = models.User{
		ID:           15,
		Email:        "test15@tabimoney.com",
		Username:     "testuser15",
		PasswordHash: string(hashedPassword),
		FirstName:    "Test",
		LastName:     "User",
		IsVerified:   true,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Create user profile
	profile := models.UserProfile{
		UserID:               15,
		MonthlyIncome:        20000000, // 20M VND
		Currency:             "VND",
		Timezone:             "Asia/Ho_Chi_Minh",
		Language:             "vi",
		NotificationSettings: "{}",
		AISettings:           "{}",
	}

	if err := db.Create(&profile).Error; err != nil {
		log.Printf("Warning: Failed to create user profile: %v", err)
	}

	fmt.Printf("✅ Created test user:\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   Username: %s\n", user.Username)
	fmt.Printf("   Email: %s\n", user.Email)
	fmt.Printf("   Password: test123456\n")
}
