package models

import (
	"time"
)

type AIAnalysis struct {
	ID              uint64    `json:"id" gorm:"primaryKey"`
	UserID          uint64    `json:"user_id" gorm:"not null"`
	AnalysisType    string    `json:"analysis_type" gorm:"type:enum('expense_prediction','anomaly_detection','category_suggestion','spending_pattern','goal_analysis');not null"`
	Data            string    `json:"data" gorm:"type:json;not null"`
	ConfidenceScore float64   `json:"confidence_score"`
	ModelVersion    string    `json:"model_version"`
	CreatedAt       time.Time `json:"created_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type AIFeedback struct {
	ID                 uint64    `json:"id" gorm:"primaryKey"`
	UserID             uint64    `json:"user_id" gorm:"not null"`
	TransactionID      *uint64   `json:"transaction_id"`
	FeedbackType       string    `json:"feedback_type" gorm:"type:enum('category_correct','category_incorrect','prediction_accurate','prediction_inaccurate','suggestion_helpful','suggestion_not_helpful');not null"`
	OriginalPrediction string    `json:"original_prediction" gorm:"type:json"`
	UserCorrection     string    `json:"user_correction" gorm:"type:json"`
	FeedbackText       string    `json:"feedback_text"`
	CreatedAt          time.Time `json:"created_at"`

	// Relations
	User        *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
}

type Notification struct {
	ID               uint64     `json:"id" gorm:"primaryKey"`
	UserID           uint64     `json:"user_id" gorm:"not null"`
	Title            string     `json:"title" gorm:"not null"`
	Message          string     `json:"message" gorm:"not null"`
	NotificationType string     `json:"notification_type" gorm:"type:enum('info','warning','success','error','reminder');not null"`
	Priority         string     `json:"priority" gorm:"type:enum('low','medium','high','urgent');default:'medium'"`
	IsRead           bool       `json:"is_read" gorm:"default:false"`
	ReadAt           *time.Time `json:"read_at"`
	ActionURL        string     `json:"action_url"`
	Metadata         string     `json:"metadata" gorm:"type:json"`
	CreatedAt        time.Time  `json:"created_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// AI Analysis Request/Response Models

type ExpensePredictionRequest struct {
	UserID     uint64    `json:"user_id" validate:"required"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
	CategoryID *uint64   `json:"category_id"`
}

type ExpensePredictionResponse struct {
	UserID            uint64               `json:"user_id"`
	PredictedAmount   float64              `json:"predicted_amount"`
	ConfidenceScore   float64              `json:"confidence_score"`
	CategoryBreakdown []CategoryPrediction `json:"category_breakdown"`
	Trends            []ExpenseTrend       `json:"trends"`
	Recommendations   []string             `json:"recommendations"`
	GeneratedAt       time.Time            `json:"generated_at"`
}

type CategoryPrediction struct {
	CategoryID      uint64  `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	PredictedAmount float64 `json:"predicted_amount"`
	ConfidenceScore float64 `json:"confidence_score"`
	Trend           string  `json:"trend"` // increasing, decreasing, stable
}

type ExpenseTrend struct {
	Period           string  `json:"period"`
	Amount           float64 `json:"amount"`
	ChangePercentage float64 `json:"change_percentage"`
	Trend            string  `json:"trend"`
}

type AnomalyDetectionRequest struct {
	UserID    uint64    `json:"user_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
	Threshold float64   `json:"threshold" validate:"min=0,max=1"`
}

type AnomalyDetectionResponse struct {
	UserID         uint64    `json:"user_id"`
	Anomalies      []Anomaly `json:"anomalies"`
	TotalAnomalies int       `json:"total_anomalies"`
	DetectionScore float64   `json:"detection_score"`
	GeneratedAt    time.Time `json:"generated_at"`
}

type Anomaly struct {
	TransactionID   uint64    `json:"transaction_id"`
	Amount          float64   `json:"amount"`
	CategoryName    string    `json:"category_name"`
	AnomalyScore    float64   `json:"anomaly_score"`
	AnomalyType     string    `json:"anomaly_type"` // amount, frequency, pattern
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
}

type CategorySuggestionRequest struct {
	UserID      uint64   `json:"user_id" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Amount      float64  `json:"amount" validate:"required,gt=0"`
	Location    string   `json:"location"`
	Tags        []string `json:"tags"`
}

type CategorySuggestionResponse struct {
	UserID          uint64               `json:"user_id"`
	Description     string               `json:"description"`
	Amount          float64              `json:"amount"`
	Suggestions     []CategorySuggestion `json:"suggestions"`
	ConfidenceScore float64              `json:"confidence_score"`
	GeneratedAt     time.Time            `json:"generated_at"`
}

type CategorySuggestion struct {
	CategoryID      uint64  `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	ConfidenceScore float64 `json:"confidence_score"`
	Reason          string  `json:"reason"`
	IsUserCategory  bool    `json:"is_user_category"`
}

type SpendingPatternRequest struct {
	UserID      uint64    `json:"user_id" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Granularity string    `json:"granularity" validate:"omitempty,oneof=daily weekly monthly"`
}

type SpendingPatternResponse struct {
	UserID          uint64            `json:"user_id"`
	Patterns        []SpendingPattern `json:"patterns"`
	Insights        []string          `json:"insights"`
	Recommendations []string          `json:"recommendations"`
	GeneratedAt     time.Time         `json:"generated_at"`
}

type SpendingPattern struct {
	CategoryID       uint64   `json:"category_id"`
	CategoryName     string   `json:"category_name"`
	TotalAmount      float64  `json:"total_amount"`
	TransactionCount int      `json:"transaction_count"`
	AverageAmount    float64  `json:"average_amount"`
	Frequency        string   `json:"frequency"` // daily, weekly, monthly
	Trend            string   `json:"trend"`
	PeakDays         []string `json:"peak_days"`
	PeakTimes        []string `json:"peak_times"`
}

type GoalAnalysisRequest struct {
	UserID uint64 `json:"user_id" validate:"required"`
	GoalID uint64 `json:"goal_id" validate:"required"`
}

type GoalAnalysisResponse struct {
	UserID          uint64     `json:"user_id"`
	GoalID          uint64     `json:"goal_id"`
	Progress        float64    `json:"progress"`
	OnTrack         bool       `json:"on_track"`
	ProjectedDate   *time.Time `json:"projected_date"`
	Recommendations []string   `json:"recommendations"`
	RiskFactors     []string   `json:"risk_factors"`
	GeneratedAt     time.Time  `json:"generated_at"`
}

// NLU (Natural Language Understanding) Models

type NLURequest struct {
	Text    string `json:"text" validate:"required"`
	UserID  uint64 `json:"user_id" validate:"required"`
	Context string `json:"context"` // conversation context
}

type NLUResponse struct {
	UserID          uint64    `json:"user_id"`
	Intent          string    `json:"intent"`
	Entities        []Entity  `json:"entities"`
	Confidence      float64   `json:"confidence"`
	SuggestedAction string    `json:"suggested_action"`
	Response        string    `json:"response"`
	GeneratedAt     time.Time `json:"generated_at"`
}

type Entity struct {
	Type       string  `json:"type"` // amount, category, date, description
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
	StartPos   int     `json:"start_pos"`
	EndPos     int     `json:"end_pos"`
}

// Dashboard Analytics Models

type DashboardAnalytics struct {
	UserID            uint64              `json:"user_id"`
	Period            string              `json:"period"`
	TotalIncome       float64             `json:"total_income"`
	TotalExpense      float64             `json:"total_expense"`
	NetAmount         float64             `json:"net_amount"`
	TransactionCount  int                 `json:"transaction_count"`
	CategoryBreakdown []CategoryAnalytics `json:"category_breakdown"`
	MonthlyTrends     []MonthlyTrend      `json:"monthly_trends"`
	TopCategories     []CategoryAnalytics `json:"top_categories"`
	FinancialHealth   FinancialHealth     `json:"financial_health"`
	GeneratedAt       time.Time           `json:"generated_at"`
}

type CategoryAnalytics struct {
	CategoryID       uint64  `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	Amount           float64 `json:"amount"`
	TransactionCount int     `json:"transaction_count"`
	Percentage       float64 `json:"percentage"`
	Trend            string  `json:"trend"`
	AverageAmount    float64 `json:"average_amount"`
}

type MonthlyTrend struct {
	Month            string  `json:"month"`
	Income           float64 `json:"income"`
	Expense          float64 `json:"expense"`
	NetAmount        float64 `json:"net_amount"`
	TransactionCount int     `json:"transaction_count"`
}

type FinancialHealth struct {
	Score           float64  `json:"score"` // 0-100
	Level           string   `json:"level"` // excellent, good, fair, poor
	IncomeRatio     float64  `json:"income_ratio"`
	SavingsRate     float64  `json:"savings_rate"`
	DebtRatio       float64  `json:"debt_ratio"`
	Recommendations []string `json:"recommendations"`
}
