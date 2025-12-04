package handlers

import (
	"net/http"
	"strconv"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/models"
	"tabimoney/internal/services"

	"github.com/labstack/echo/v4"
)

type AnalyticsHandler struct {
	transactionService *services.TransactionService
	aiService         *services.AIService
}

func NewAnalyticsHandler(cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		transactionService: services.NewTransactionService(cfg),
		aiService:         services.NewAIService(cfg),
	}
}

// GetDashboardAnalytics retrieves monthly financial summary with AI predictions
func (h *AnalyticsHandler) GetDashboardAnalytics(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Parse year and month from query params
	year := time.Now().Year()
	month := int(time.Now().Month())
	
	if y := c.QueryParam("year"); y != "" {
		if parsedYear, err := strconv.Atoi(y); err == nil {
			year = parsedYear
		}
	}
	if m := c.QueryParam("month"); m != "" {
		if parsedMonth, err := strconv.Atoi(m); err == nil {
			month = parsedMonth
		}
	}

	// Get basic analytics
	analytics, err := h.transactionService.GetMonthlySummary(userID, year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get analytics",
			Message: err.Error(),
		})
	}

	// Get AI predictions for the next period
	startDate := time.Now().AddDate(0, -3, 0) // Last 3 months
	endDate := time.Now()
	
	predictionReq := &models.ExpensePredictionRequest{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	predictions, err := h.aiService.PredictExpenses(predictionReq)
	if err != nil {
		// If AI service fails, continue without predictions
		// Set default empty predictions
		predictions = &models.ExpensePredictionResponse{
			UserID:            userID,
			PredictedAmount:   0,
			ConfidenceScore:   0,
			CategoryBreakdown: []models.CategoryPrediction{},
			Trends:            []models.ExpenseTrend{},
			Recommendations:   []string{"AI Service đang khởi tạo..."},
			GeneratedAt:       time.Now(),
		}
	}

	// Combine analytics with AI predictions
	response := map[string]interface{}{
		"user_id":            analytics.UserID,
		"period":             analytics.Period,
		"total_income":       analytics.TotalIncome,
		"total_expense":      analytics.TotalExpense,
		"net_amount":         analytics.NetAmount,
		"transaction_count":  analytics.TransactionCount,
		"category_breakdown": analytics.CategoryBreakdown,
		"monthly_trends":     analytics.MonthlyTrends,
		"top_categories":     analytics.TopCategories,
		"financial_health":   analytics.FinancialHealth,
		"generated_at":       analytics.GeneratedAt,
		// Add AI predictions
		"predictions": predictions,
	}

	return c.JSON(http.StatusOK, response)
}

// GetCategorySpending retrieves spending breakdown by category
func (h *AnalyticsHandler) GetCategorySpending(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Parse date range
	startDate := time.Now().AddDate(0, -1, 0) // Default: last month
	endDate := time.Now()

	if s := c.QueryParam("start_date"); s != "" {
		if parsed, err := time.Parse("2006-01-02", s); err == nil {
			startDate = parsed
		}
	}
	if e := c.QueryParam("end_date"); e != "" {
		if parsed, err := time.Parse("2006-01-02", e); err == nil {
			endDate = parsed
		}
	}

	spending, err := h.transactionService.GetCategorySpending(userID, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get category spending",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, spending)
}

// GetSpendingPatterns analyzes spending patterns using AI
func (h *AnalyticsHandler) GetSpendingPatterns(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Parse date range
	startDate := time.Now().AddDate(0, -3, 0) // Default: last 3 months
	endDate := time.Now()

	if s := c.QueryParam("start_date"); s != "" {
		if parsed, err := time.Parse("2006-01-02", s); err == nil {
			startDate = parsed
		}
	}
	if e := c.QueryParam("end_date"); e != "" {
		if parsed, err := time.Parse("2006-01-02", e); err == nil {
			endDate = parsed
		}
	}

	granularity := c.QueryParam("granularity")
	if granularity == "" {
		granularity = "monthly"
	}

	req := &models.SpendingPatternRequest{
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
		Granularity: granularity,
	}

	patterns, err := h.aiService.AnalyzeSpendingPattern(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to analyze spending patterns",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, patterns)
}

// GetAnomalies detects spending anomalies
func (h *AnalyticsHandler) GetAnomalies(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Parse date range
	startDate := time.Now().AddDate(0, -1, 0) // Default: last month
	endDate := time.Now()

	if s := c.QueryParam("start_date"); s != "" {
		if parsed, err := time.Parse("2006-01-02", s); err == nil {
			startDate = parsed
		}
	}
	if e := c.QueryParam("end_date"); e != "" {
		if parsed, err := time.Parse("2006-01-02", e); err == nil {
			endDate = parsed
		}
	}

	threshold := 0.8 // Default threshold
	if t := c.QueryParam("threshold"); t != "" {
		if parsed, err := strconv.ParseFloat(t, 64); err == nil {
			threshold = parsed
		}
	}

	req := &models.AnomalyDetectionRequest{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Threshold: threshold,
	}

	anomalies, err := h.aiService.DetectAnomalies(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to detect anomalies",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, anomalies)
}

// GetPredictions gets expense predictions
func (h *AnalyticsHandler) GetPredictions(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Parse prediction period
	period := c.QueryParam("period")
	if period == "" {
		period = "monthly"
	}

	// Parse date range for predictions
	startDate := time.Now().AddDate(0, -3, 0) // Default: last 3 months
	endDate := time.Now()

	if s := c.QueryParam("start_date"); s != "" {
		if parsed, err := time.Parse("2006-01-02", s); err == nil {
			startDate = parsed
		}
	}
	if e := c.QueryParam("end_date"); e != "" {
		if parsed, err := time.Parse("2006-01-02", e); err == nil {
			endDate = parsed
		}
	}

	req := &models.ExpensePredictionRequest{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	predictions, err := h.aiService.PredictExpenses(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get predictions",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, predictions)
}
