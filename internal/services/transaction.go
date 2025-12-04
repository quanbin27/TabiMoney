package services

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "time"

    "tabimoney/internal/config"
    "tabimoney/internal/database"
    "tabimoney/internal/models"

    "gorm.io/gorm"
)

type TransactionService struct {
	db     *gorm.DB
	config *config.Config
}

func NewTransactionService(cfg *config.Config) *TransactionService {
	return &TransactionService{
		db:     database.GetDB(),
		config: cfg,
	}
}

// CreateTransaction creates a new transaction
func (s *TransactionService) CreateTransaction(userID uint64, req *models.TransactionCreateRequest) (*models.TransactionResponse, error) {
	// Validate category exists and belongs to user or is system category
	var category models.Category
	if err := s.db.Where("id = ? AND (user_id = ? OR is_system = ?)", 
		req.CategoryID, userID, true).First(&category).Error; err != nil {
		return nil, fmt.Errorf("category not found or not accessible: %w", err)
	}

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction_date format, expected YYYY-MM-DD: %w", err)
	}

	// Parse transaction time if provided
	var transactionTime *time.Time
	if req.TransactionTime != "" {
		parsedTime, err := time.Parse("15:04", req.TransactionTime)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction_time format, expected HH:MM: %w", err)
		}
		// Combine date and time
		combinedTime := time.Date(transactionDate.Year(), transactionDate.Month(), transactionDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, transactionDate.Location())
		transactionTime = &combinedTime
	}

	// Create transaction
	transaction := &models.Transaction{
		UserID:          userID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionType: req.TransactionType,
		TransactionDate: transactionDate,
		TransactionTime: transactionTime,
		Location:        req.Location,
		Tags:            s.marshalTags(req.Tags),
		Metadata:        s.marshalMetadata(req.Metadata),
		IsRecurring:     req.IsRecurring,
		RecurringPattern: req.RecurringPattern,
	}

    if err := s.db.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Load category for response
	if err := s.db.Preload("Category").First(transaction, transaction.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load transaction with category: %w", err)
	}

	// Trigger budget threshold notifications synchronously (best-effort)
	bs := NewBudgetService(s.config)
	if err := bs.CheckBudgetNotifications(userID); err != nil {
		log.Printf("Failed to check budget notifications: %v", err)
	}

	// Check for large transaction alert (only for expense transactions)
	if transaction.TransactionType == "expense" {
		dispatcher := NewNotificationDispatcher(s.config)
		threshold := s.getLargeTransactionThreshold(userID)
		if transaction.Amount > threshold {
			if err := dispatcher.TriggerLargeTransactionAlert(userID, transaction, threshold); err != nil {
				log.Printf("Failed to trigger large transaction alert: %v", err)
			}
		}
	}

	// Clear dashboard cache
	ctx := context.Background()
	database.DeleteDashboardCache(ctx, userID)

	return s.transactionToResponse(transaction), nil
}

// GetTransactions retrieves transactions with filtering and pagination
func (s *TransactionService) GetTransactions(userID uint64, req *models.TransactionQueryRequest) ([]models.TransactionResponse, int64, error) {
	var transactions []models.Transaction
	var total int64

	// Build query
	query := s.db.Where("user_id = ?", userID)

	// Apply filters
	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}
	if req.TransactionType != nil {
		query = query.Where("transaction_type = ?", *req.TransactionType)
	}
	if req.StartDate != nil {
		query = query.Where("transaction_date >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("transaction_date <= ?", *req.EndDate)
	}
	if req.MinAmount != nil {
		query = query.Where("amount >= ?", *req.MinAmount)
	}
	if req.MaxAmount != nil {
		query = query.Where("amount <= ?", *req.MaxAmount)
	}
	if req.Search != "" {
		query = query.Where("description LIKE ? OR location LIKE ?", 
			"%"+req.Search+"%", "%"+req.Search+"%")
	}

	// Get total count
	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Apply sorting
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "transaction_date"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).
		Preload("Category").
		Preload("AISuggestedCategory").
		Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

    // Convert to response (ensure empty slice, not null)
    responses := make([]models.TransactionResponse, 0)
	for _, t := range transactions {
		responses = append(responses, *s.transactionToResponse(&t))
	}

	return responses, total, nil
}

// UpdateTransaction updates an existing transaction
func (s *TransactionService) UpdateTransaction(userID, transactionID uint64, req *models.TransactionUpdateRequest) (*models.TransactionResponse, error) {
	// Find transaction
	var transaction models.Transaction
	
	if err := s.db.Where("user_id = ? AND id = ?", userID, transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	// Validate category
	var category models.Category
	if err := s.db.Where("id = ? AND (user_id = ? OR is_system = ?)", 
		req.CategoryID, userID, true).First(&category).Error; err != nil {
		return nil, fmt.Errorf("category not found or not accessible: %w", err)
	}

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction_date format, expected YYYY-MM-DD: %w", err)
	}

	// Parse transaction time if provided
	var transactionTime *time.Time
	if req.TransactionTime != "" {
		parsedTime, err := time.Parse("15:04", req.TransactionTime)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction_time format, expected HH:MM: %w", err)
		}
		// Combine date and time
		combinedTime := time.Date(transactionDate.Year(), transactionDate.Month(), transactionDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0, transactionDate.Location())
		transactionTime = &combinedTime
	}

	// Update transaction
	transaction.CategoryID = req.CategoryID
	transaction.Amount = req.Amount
	transaction.Description = req.Description
	transaction.TransactionType = req.TransactionType
	transaction.TransactionDate = transactionDate
	transaction.TransactionTime = transactionTime
	transaction.Location = req.Location
	transaction.Tags = s.marshalTags(req.Tags)
	transaction.Metadata = s.marshalMetadata(req.Metadata)

	if err := s.db.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Load category for response
	if err := s.db.Preload("Category").First(&transaction, transaction.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load transaction with category: %w", err)
	}

	// Clear dashboard cache
    ctx := context.Background()
	database.DeleteDashboardCache(ctx, userID)

	return s.transactionToResponse(&transaction), nil
}

// DeleteTransaction deletes a transaction
func (s *TransactionService) DeleteTransaction(userID, transactionID uint64) error {
	// Find transaction
	var transaction models.Transaction
	if err := s.db.Where("user_id = ? AND id = ?", userID, transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("transaction not found")
		}
		return fmt.Errorf("failed to find transaction: %w", err)
	}

	// Delete transaction
	if err := s.db.Delete(&transaction).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

    // Trigger budget threshold notifications synchronously (best-effort)
    bs := NewBudgetService(s.config)
    if err := bs.CheckBudgetNotifications(userID); err != nil {
        log.Printf("Failed to check budget notifications: %v", err)
    }

    // Clear dashboard cache
    ctx := context.Background()
    database.DeleteDashboardCache(ctx, userID)

	return nil
}

// GetMonthlySummary retrieves monthly financial summary
func (s *TransactionService) GetMonthlySummary(userID uint64, year int, month int) (*models.DashboardAnalytics, error) {
    // Check cache first
    ctx := context.Background()
	period := fmt.Sprintf("%d-%02d", year, month)
	cacheKey := fmt.Sprintf("dashboard:%d:%s", userID, period)
	
	if cached, err := database.GetCache(ctx, cacheKey); err == nil {
		var analytics models.DashboardAnalytics
		if err := json.Unmarshal([]byte(cached), &analytics); err == nil {
			return &analytics, nil
		}
	}

	// Calculate date range
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// Get transactions
	var transactions []models.Transaction
	if err := s.db.Where("user_id = ? AND transaction_date BETWEEN ? AND ?", 
		userID, startDate, endDate).
		Preload("Category").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Calculate analytics
	analytics := s.calculateMonthlyAnalytics(userID, transactions, period)

	// Cache result
	if analyticsJSON, err := json.Marshal(analytics); err == nil {
		database.SetCache(ctx, cacheKey, analyticsJSON, 1*time.Hour)
	}

	return analytics, nil
}

// GetCategorySpending retrieves spending breakdown by category
func (s *TransactionService) GetCategorySpending(userID uint64, startDate, endDate time.Time) ([]models.CategoryAnalytics, error) {
	var results []models.CategoryAnalytics

	// Query category spending
	if err := s.db.Raw(`
		SELECT 
			c.id as category_id,
			c.name as category_name,
			SUM(t.amount) as amount,
			COUNT(t.id) as transaction_count,
			AVG(t.amount) as average_amount
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? 
			AND t.transaction_type = 'expense'
			AND t.transaction_date BETWEEN ? AND ?
		GROUP BY c.id, c.name
		ORDER BY amount DESC
	`, userID, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get category spending: %w", err)
	}

	// Calculate percentages
	var totalAmount float64
	for _, result := range results {
		totalAmount += result.Amount
	}

	for i := range results {
		if totalAmount > 0 {
			results[i].Percentage = (results[i].Amount / totalAmount) * 100
		}
	}

	return results, nil
}

// Helper methods

// getLargeTransactionThreshold gets the large transaction threshold for a user
// Returns user's custom threshold if set, otherwise returns system default (1,000,000 VND)
func (s *TransactionService) getLargeTransactionThreshold(userID uint64) float64 {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err == nil {
		if profile.LargeTransactionThreshold != nil && *profile.LargeTransactionThreshold > 0 {
			return *profile.LargeTransactionThreshold
		}
	}
	// Default threshold: 1,000,000 VND (1 million)
	return 1000000
}

func (s *TransactionService) transactionToResponse(t *models.Transaction) *models.TransactionResponse {
	response := &models.TransactionResponse{
		ID:                      t.ID,
		UserID:                  t.UserID,
		CategoryID:              t.CategoryID,
		Amount:                  t.Amount,
		Description:             t.Description,
		TransactionType:         t.TransactionType,
		TransactionDate:         t.TransactionDate,
		TransactionTime:         t.TransactionTime,
		Location:                t.Location,
		Tags:                    s.unmarshalTags(t.Tags),
		Metadata:                s.unmarshalMetadata(t.Metadata),
		IsRecurring:             t.IsRecurring,
		RecurringPattern:        t.RecurringPattern,
		ParentTransactionID:     t.ParentTransactionID,
		AIConfidence:            t.AIConfidence,
		AISuggestedCategoryID:  t.AISuggestedCategoryID,
		CreatedAt:               t.CreatedAt,
		UpdatedAt:               t.UpdatedAt,
	}

	if t.Category != nil {
		response.Category = s.categoryToResponse(t.Category)
	}

	if t.AISuggestedCategory != nil {
		response.AISuggestedCategory = s.categoryToResponse(t.AISuggestedCategory)
	}

	return response
}

func (s *TransactionService) categoryToResponse(c *models.Category) *models.CategoryResponse {
	return &models.CategoryResponse{
		ID:          c.ID,
		UserID:      c.UserID,
		Name:        c.Name,
		NameEn:      c.NameEn,
		Description: c.Description,
		ParentID:    c.ParentID,
		IsSystem:    c.IsSystem,
		IsActive:    c.IsActive,
		SortOrder:   c.SortOrder,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func (s *TransactionService) marshalTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	jsonData, _ := json.Marshal(tags)
	return string(jsonData)
}

func (s *TransactionService) unmarshalTags(tagsJSON string) []string {
	var tags []string
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &tags)
	}
	return tags
}

func (s *TransactionService) marshalMetadata(metadata map[string]interface{}) string {
	if len(metadata) == 0 {
		return "{}"
	}
	jsonData, _ := json.Marshal(metadata)
	return string(jsonData)
}

func (s *TransactionService) unmarshalMetadata(metadataJSON string) map[string]interface{} {
	var metadata map[string]interface{}
	if metadataJSON != "" {
		json.Unmarshal([]byte(metadataJSON), &metadata)
	}
	return metadata
}

func (s *TransactionService) calculateMonthlyAnalytics(userID uint64, transactions []models.Transaction, period string) *models.DashboardAnalytics {
	var totalIncome, totalExpense float64
	var incomeCount, expenseCount int

	// Calculate totals
	for _, t := range transactions {
		if t.TransactionType == "income" {
			totalIncome += t.Amount
			incomeCount++
		} else if t.TransactionType == "expense" {
			totalExpense += t.Amount
			expenseCount++
		}
	}

	netAmount := totalIncome - totalExpense

	// Calculate category breakdown
	categoryMap := make(map[uint64]*models.CategoryAnalytics)
	for _, t := range transactions {
		if t.TransactionType == "expense" {
			if analytics, exists := categoryMap[t.CategoryID]; exists {
				analytics.Amount += t.Amount
				analytics.TransactionCount++
			} else {
				analytics = &models.CategoryAnalytics{
					CategoryID:       t.CategoryID,
					CategoryName:     t.Category.Name,
					Amount:           t.Amount,
					TransactionCount: 1,
				}
				categoryMap[t.CategoryID] = analytics
			}
		}
	}

	var categoryBreakdown []models.CategoryAnalytics
	for _, analytics := range categoryMap {
		analytics.AverageAmount = analytics.Amount / float64(analytics.TransactionCount)
		analytics.Percentage = (analytics.Amount / totalExpense) * 100
		categoryBreakdown = append(categoryBreakdown, *analytics)
	}

	// Calculate financial health
	savingsRate := 0.0
	if totalIncome > 0 {
		savingsRate = (netAmount / totalIncome) * 100
	}

	financialHealth := models.FinancialHealth{
		Score:       s.calculateFinancialHealthScore(savingsRate, totalIncome, totalExpense),
		Level:       s.getFinancialHealthLevel(savingsRate),
		IncomeRatio: totalIncome,
		SavingsRate: savingsRate,
		DebtRatio:   0, // Simplified - would need debt data
		Recommendations: s.generateFinancialRecommendations(savingsRate, totalIncome, totalExpense),
	}

	return &models.DashboardAnalytics{
		UserID:             userID,
		Period:             period,
		TotalIncome:        totalIncome,
		TotalExpense:       totalExpense,
		NetAmount:          netAmount,
		TransactionCount:   len(transactions),
		CategoryBreakdown:  categoryBreakdown,
		FinancialHealth:    financialHealth,
		GeneratedAt:        time.Now(),
	}
}

func (s *TransactionService) calculateFinancialHealthScore(savingsRate, income, expense float64) float64 {
	// Simplified financial health scoring
	score := 50.0 // Base score
	
	if savingsRate > 20 {
		score += 30
	} else if savingsRate > 10 {
		score += 20
	} else if savingsRate > 0 {
		score += 10
	} else {
		score -= 20
	}

	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

func (s *TransactionService) getFinancialHealthLevel(savingsRate float64) string {
	if savingsRate >= 20 {
		return "excellent"
	} else if savingsRate >= 10 {
		return "good"
	} else if savingsRate >= 0 {
		return "fair"
	} else {
		return "poor"
	}
}

func (s *TransactionService) generateFinancialRecommendations(savingsRate, income, expense float64) []string {
	var recommendations []string

	if savingsRate < 10 {
		recommendations = append(recommendations, "Consider increasing your savings rate to at least 10%")
	}

	if expense > income {
		recommendations = append(recommendations, "Your expenses exceed your income. Review your spending habits")
	}

	if savingsRate > 20 {
		recommendations = append(recommendations, "Great job! Consider investing your excess savings")
	}

	return recommendations
}
