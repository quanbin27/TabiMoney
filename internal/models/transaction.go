package models

import (
    "time"
)

type Transaction struct {
	ID                      uint64         `json:"id" gorm:"primaryKey"`
	UserID                  uint64         `json:"user_id" gorm:"not null"`
	CategoryID              uint64         `json:"category_id" gorm:"not null"`
	Amount                  float64        `json:"amount" gorm:"not null"`
	Description             string         `json:"description"`
	TransactionType         string         `json:"transaction_type" gorm:"type:enum('income','expense','transfer');not null"`
	TransactionDate         time.Time      `json:"transaction_date" gorm:"not null"`
	TransactionTime         *time.Time     `json:"transaction_time"`
	Location                string         `json:"location"`
	Tags                    string         `json:"tags" gorm:"type:json"`
	Metadata                string         `json:"metadata" gorm:"type:json"`
	IsRecurring             bool           `json:"is_recurring" gorm:"default:false"`
	RecurringPattern        string         `json:"recurring_pattern"`
	ParentTransactionID     *uint64        `json:"parent_transaction_id"`
	AIConfidence            float64        `json:"ai_confidence"`
	AISuggestedCategoryID   *uint64        `json:"ai_suggested_category_id"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`

	// Relations
	User                    *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category                *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	ParentTransaction       *Transaction   `json:"parent_transaction,omitempty" gorm:"foreignKey:ParentTransactionID"`
	AISuggestedCategory     *Category      `json:"ai_suggested_category,omitempty" gorm:"foreignKey:AISuggestedCategoryID"`
}

type Category struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	UserID      *uint64        `json:"user_id"`
	Name        string         `json:"name" gorm:"not null"`
	NameEn      string         `json:"name_en"`
	Description string         `json:"description"`
	ParentID    *uint64        `json:"parent_id"`
	IsSystem    bool           `json:"is_system" gorm:"default:false"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relations
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent      *Category      `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []Category     `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Transactions []Transaction `json:"-" gorm:"foreignKey:CategoryID"`
}

type FinancialGoal struct {
	ID            uint64     `json:"id" gorm:"primaryKey"`
	UserID        uint64     `json:"user_id" gorm:"not null"`
	Title         string     `json:"title" gorm:"not null"`
	Description   string     `json:"description"`
	TargetAmount  float64    `json:"target_amount" gorm:"not null"`
	CurrentAmount float64    `json:"current_amount" gorm:"default:0"`
	TargetDate    *time.Time `json:"target_date"`
	GoalType      string     `json:"goal_type" gorm:"type:enum('savings','debt_payment','investment','purchase','other');default:'savings'"`
	Priority      string     `json:"priority" gorm:"type:enum('low','medium','high','urgent');default:'medium'"`
	IsAchieved    bool       `json:"is_achieved" gorm:"default:false"`
	AchievedAt    *time.Time `json:"achieved_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Progress      float64    `json:"progress" gorm:"-"` // Calculated field, not stored in DB

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type Budget struct {
	ID              uint64     `json:"id" gorm:"primaryKey"`
	UserID          uint64     `json:"user_id" gorm:"not null"`
	CategoryID      *uint64    `json:"category_id"`
	Name            string     `json:"name" gorm:"not null"`
	Amount          float64    `json:"amount" gorm:"not null"`
	Period          string     `json:"period" gorm:"type:enum('weekly','monthly','yearly');default:'monthly'"`
	StartDate       time.Time  `json:"start_date" gorm:"not null"`
	EndDate         time.Time  `json:"end_date" gorm:"not null"`
	IsActive        bool       `json:"is_active" gorm:"default:true"`
	AlertThreshold  float64    `json:"alert_threshold" gorm:"default:80.00"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	SpentAmount     float64    `json:"spent_amount" gorm:"-"` // Calculated field
	RemainingAmount float64    `json:"remaining_amount" gorm:"-"` // Calculated field
	UsagePercentage float64    `json:"usage_percentage" gorm:"-"` // Calculated field

	// Relations
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// TransactionCreateRequest represents the request payload for creating a transaction
type TransactionCreateRequest struct {
	CategoryID      uint64    `json:"category_id" validate:"required"`
	Amount          float64   `json:"amount" validate:"required,gt=0"`
	Description     string    `json:"description" validate:"max=500"`
	TransactionType string    `json:"transaction_type" validate:"required,oneof=income expense transfer"`
	TransactionDate string    `json:"transaction_date" validate:"required"`
	TransactionTime string    `json:"transaction_time,omitempty"`
	Location        string    `json:"location" validate:"max=200"`
	Tags            []string  `json:"tags"`
	Metadata        map[string]interface{} `json:"metadata"`
	IsRecurring     bool      `json:"is_recurring"`
	RecurringPattern string   `json:"recurring_pattern" validate:"omitempty,oneof=daily weekly monthly yearly"`
}

// TransactionUpdateRequest represents the request payload for updating a transaction
type TransactionUpdateRequest struct {
	CategoryID      uint64    `json:"category_id" validate:"required"`
	Amount          float64   `json:"amount" validate:"required,gt=0"`
	Description     string    `json:"description" validate:"max=500"`
	TransactionType string    `json:"transaction_type" validate:"required,oneof=income expense transfer"`
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
	TransactionTime *time.Time `json:"transaction_time"`
	Location        string    `json:"location" validate:"max=200"`
	Tags            []string  `json:"tags"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// TransactionQueryRequest represents the request payload for querying transactions
type TransactionQueryRequest struct {
	Page           int       `json:"page" validate:"min=1"`
	Limit          int       `json:"limit" validate:"min=1,max=100"`
	CategoryID     *uint64   `json:"category_id"`
	TransactionType *string  `json:"transaction_type" validate:"omitempty,oneof=income expense transfer"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	MinAmount      *float64  `json:"min_amount"`
	MaxAmount      *float64  `json:"max_amount"`
	Search         string    `json:"search"`
	SortBy         string    `json:"sort_by" validate:"omitempty,oneof=created_at transaction_date amount"`
	SortOrder      string    `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// TransactionResponse represents the response payload for transaction data
type TransactionResponse struct {
	ID                      uint64                `json:"id"`
	UserID                  uint64                `json:"user_id"`
	CategoryID              uint64                `json:"category_id"`
	Amount                  float64               `json:"amount"`
	Description             string                `json:"description"`
	TransactionType         string                `json:"transaction_type"`
	TransactionDate         time.Time             `json:"transaction_date"`
	TransactionTime         *time.Time            `json:"transaction_time"`
	Location                string                `json:"location"`
	Tags                    []string              `json:"tags"`
	Metadata                map[string]interface{} `json:"metadata"`
	IsRecurring             bool                  `json:"is_recurring"`
	RecurringPattern        string                `json:"recurring_pattern"`
	ParentTransactionID     *uint64               `json:"parent_transaction_id"`
	AIConfidence            float64               `json:"ai_confidence"`
	AISuggestedCategoryID   *uint64               `json:"ai_suggested_category_id"`
	CreatedAt               time.Time             `json:"created_at"`
	UpdatedAt               time.Time             `json:"updated_at"`
	Category                *CategoryResponse     `json:"category,omitempty"`
	AISuggestedCategory     *CategoryResponse     `json:"ai_suggested_category,omitempty"`
}

// CategoryResponse represents the response payload for category data
type CategoryResponse struct {
	ID          uint64    `json:"id"`
	UserID      *uint64   `json:"user_id"`
	Name        string    `json:"name"`
	NameEn      string    `json:"name_en"`
	Description string    `json:"description"`
	ParentID    *uint64   `json:"parent_id"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FinancialGoalCreateRequest represents the request payload for creating a financial goal
type FinancialGoalCreateRequest struct {
	Title         string     `json:"title" validate:"required,max=200"`
	Description   string     `json:"description" validate:"max=1000"`
	TargetAmount  float64    `json:"target_amount" validate:"required,gt=0"`
	TargetDate    *time.Time `json:"target_date"`
	GoalType      string     `json:"goal_type" validate:"required,oneof=savings debt_payment investment purchase other"`
	Priority      string     `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
}

// FinancialGoalUpdateRequest represents the request payload for updating a financial goal
type FinancialGoalUpdateRequest struct {
	Title         string     `json:"title" validate:"required,max=200"`
	Description   string     `json:"description" validate:"max=1000"`
	TargetAmount  float64    `json:"target_amount" validate:"required,gt=0"`
	CurrentAmount float64    `json:"current_amount" validate:"min=0"`
	TargetDate    *time.Time `json:"target_date"`
	GoalType      string     `json:"goal_type" validate:"required,oneof=savings debt_payment investment purchase other"`
	Priority      string     `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
}

// FinancialGoalResponse represents the response payload for financial goal data
type FinancialGoalResponse struct {
	ID            uint64     `json:"id"`
	UserID        uint64     `json:"user_id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	TargetAmount  float64    `json:"target_amount"`
	CurrentAmount float64    `json:"current_amount"`
	TargetDate    *time.Time `json:"target_date"`
	GoalType      string     `json:"goal_type"`
	Priority      string     `json:"priority"`
	IsAchieved    bool       `json:"is_achieved"`
	AchievedAt    *time.Time `json:"achieved_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Progress      float64    `json:"progress"` // Calculated field
}

// BudgetCreateRequest represents the request payload for creating a budget
type BudgetCreateRequest struct {
	CategoryID     *uint64   `json:"category_id"`
	Name           string    `json:"name" validate:"required,max=200"`
	Amount         float64   `json:"amount" validate:"required,gt=0"`
	Period         string    `json:"period" validate:"required,oneof=weekly monthly yearly"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required"`
	AlertThreshold float64   `json:"alert_threshold" validate:"min=0,max=100"`
}

// BudgetUpdateRequest represents the request payload for updating a budget
type BudgetUpdateRequest struct {
	CategoryID     *uint64   `json:"category_id"`
	Name           string    `json:"name" validate:"required,max=200"`
	Amount         float64   `json:"amount" validate:"required,gt=0"`
	Period         string    `json:"period" validate:"required,oneof=weekly monthly yearly"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required"`
	IsActive       bool      `json:"is_active"`
	AlertThreshold float64   `json:"alert_threshold" validate:"min=0,max=100"`
}

// BudgetResponse represents the response payload for budget data
type BudgetResponse struct {
	ID             uint64                `json:"id"`
	UserID         uint64                `json:"user_id"`
	CategoryID     *uint64               `json:"category_id"`
	Name           string                `json:"name"`
	Amount         float64               `json:"amount"`
	Period         string                `json:"period"`
	StartDate      time.Time             `json:"start_date"`
	EndDate        time.Time             `json:"end_date"`
	IsActive       bool                  `json:"is_active"`
	AlertThreshold float64               `json:"alert_threshold"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	Category       *CategoryResponse     `json:"category,omitempty"`
	SpentAmount    float64               `json:"spent_amount"` // Calculated field
	RemainingAmount float64              `json:"remaining_amount"` // Calculated field
	UsagePercentage float64              `json:"usage_percentage"` // Calculated field
}
