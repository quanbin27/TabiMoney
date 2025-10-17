	package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                   uint64         `json:"id" gorm:"primaryKey"`
    Email                string         `json:"email" gorm:"size:191;uniqueIndex;not null"`
    Username             string         `json:"username" gorm:"size:191;uniqueIndex;not null"`
	PasswordHash         string         `json:"-" gorm:"not null"`
	FirstName            string         `json:"first_name"`
	LastName             string         `json:"last_name"`
	Phone                string         `json:"phone"`
	AvatarURL            string         `json:"avatar_url"`
	IsVerified           bool           `json:"is_verified" gorm:"default:false"`
	VerificationToken    string         `json:"-"`
	ResetToken           string         `json:"-"`
	ResetTokenExpiresAt  *time.Time     `json:"-"`
	LastLoginAt          *time.Time     `json:"last_login_at"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Profile      *UserProfile      `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Sessions     []UserSession     `json:"-" gorm:"foreignKey:UserID"`
	Transactions []Transaction     `json:"-" gorm:"foreignKey:UserID"`
	Categories   []Category        `json:"-" gorm:"foreignKey:UserID"`
	Goals        []FinancialGoal   `json:"-" gorm:"foreignKey:UserID"`
	Budgets      []Budget          `json:"-" gorm:"foreignKey:UserID"`
	Notifications []Notification   `json:"-" gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	ID                  uint64    `json:"id" gorm:"primaryKey"`
	UserID              uint64    `json:"user_id" gorm:"uniqueIndex;not null"`
	MonthlyIncome       float64   `json:"monthly_income" gorm:"default:0"`
	Currency            string    `json:"currency" gorm:"default:'VND'"`
	Timezone            string    `json:"timezone" gorm:"default:'Asia/Ho_Chi_Minh'"`
	Language            string    `json:"language" gorm:"default:'vi'"`
	NotificationSettings string   `json:"notification_settings" gorm:"type:json"`
	AISettings          string    `json:"ai_settings" gorm:"type:json"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type UserSession struct {
	ID                  uint64    `json:"id" gorm:"primaryKey"`
	UserID              uint64    `json:"user_id" gorm:"not null"`
    TokenHash           string    `json:"-" gorm:"size:191;not null;index"`
    RefreshTokenHash    string    `json:"-" gorm:"size:191;not null"`
	ExpiresAt           time.Time `json:"expires_at" gorm:"not null"`
	RefreshExpiresAt    time.Time `json:"refresh_expires_at" gorm:"not null"`
	UserAgent           string    `json:"user_agent"`
	IPAddress           string    `json:"ip_address"`
	IsActive            bool      `json:"is_active" gorm:"default:true"`
	CreatedAt           time.Time `json:"created_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserCreateRequest represents the request payload for creating a user
type UserCreateRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
	Phone     string `json:"phone" validate:"max=20"`
}

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateRequest represents the request payload for updating user profile
type UserUpdateRequest struct {
	FirstName string  `json:"first_name" validate:"max=100"`
	LastName  string  `json:"last_name" validate:"max=100"`
	Phone     string  `json:"phone" validate:"max=20"`
	AvatarURL string  `json:"avatar_url" validate:"url"`
}

// UserProfileUpdateRequest represents the request payload for updating user profile settings
type UserProfileUpdateRequest struct {
	MonthlyIncome       float64 `json:"monthly_income" validate:"min=0"`
	Currency            string  `json:"currency" validate:"len=3"`
	Timezone            string  `json:"timezone"`
	Language            string  `json:"language" validate:"len=2"`
	NotificationSettings string `json:"notification_settings"`
	AISettings          string `json:"ai_settings"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	AvatarURL string    `json:"avatar_url"`
	IsVerified bool     `json:"is_verified"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Profile   *UserProfileResponse `json:"profile,omitempty"`
}

// UserProfileResponse represents the response payload for user profile data
type UserProfileResponse struct {
	ID                  uint64  `json:"id"`
	UserID              uint64  `json:"user_id"`
	MonthlyIncome       float64 `json:"monthly_income"`
	Currency            string  `json:"currency"`
	Timezone            string  `json:"timezone"`
	Language            string  `json:"language"`
	NotificationSettings string `json:"notification_settings"`
	AISettings          string `json:"ai_settings"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// AuthResponse represents the response payload for authentication
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// PasswordResetRequest represents the request payload for password reset
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetConfirmRequest represents the request payload for confirming password reset
type PasswordResetConfirmRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}
