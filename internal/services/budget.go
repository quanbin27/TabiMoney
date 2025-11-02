package services

import (
	"fmt"
	"log"

	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type BudgetService struct {
	db *gorm.DB
}

func NewBudgetService() *BudgetService {
	return &BudgetService{
		db: database.GetDB(),
	}
}

// CreateBudget creates a new budget
func (s *BudgetService) CreateBudget(userID uint64, req *models.BudgetCreateRequest) (*models.Budget, error) {
	budget := &models.Budget{
		UserID:         userID,
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Amount:         req.Amount,
		Period:         req.Period,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsActive:       true,
		AlertThreshold: req.AlertThreshold,
	}

	if err := s.db.Create(budget).Error; err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	// Calculate spent amount and remaining
	s.calculateBudgetMetrics(budget)

	return budget, nil
}

// GetBudgets retrieves user's budgets
func (s *BudgetService) GetBudgets(userID uint64) ([]models.Budget, error) {
	var budgets []models.Budget
	if err := s.db.Where("user_id = ?", userID).Preload("Category").Order("created_at DESC").Find(&budgets).Error; err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}

	// Calculate metrics for each budget
	for i := range budgets {
		s.calculateBudgetMetrics(&budgets[i])
	}

	return budgets, nil
}

// UpdateBudget updates an existing budget
func (s *BudgetService) UpdateBudget(userID uint64, budgetID uint64, req *models.BudgetUpdateRequest) (*models.Budget, error) {
	var budget models.Budget
	if err := s.db.Where("id = ? AND user_id = ?", budgetID, userID).Preload("Category").First(&budget).Error; err != nil {
		return nil, fmt.Errorf("budget not found: %w", err)
	}

	// Update fields
	budget.CategoryID = req.CategoryID
	budget.Name = req.Name
	budget.Amount = req.Amount
	budget.Period = req.Period
	budget.StartDate = req.StartDate
	budget.EndDate = req.EndDate
	budget.IsActive = req.IsActive
	budget.AlertThreshold = req.AlertThreshold

	if err := s.db.Save(&budget).Error; err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	// Calculate metrics
	s.calculateBudgetMetrics(&budget)

	return &budget, nil
}

// DeleteBudget deletes a budget
func (s *BudgetService) DeleteBudget(userID uint64, budgetID uint64) error {
	result := s.db.Where("id = ? AND user_id = ?", budgetID, userID).Delete(&models.Budget{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete budget: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("budget not found")
	}
	return nil
}

// calculateBudgetMetrics calculates spent amount, remaining amount, and usage percentage
func (s *BudgetService) calculateBudgetMetrics(budget *models.Budget) {
	// Get spent amount for this budget period
	var spentAmount float64
	query := s.db.Model(&models.Transaction{}).
		Where("user_id = ? AND transaction_type = ? AND transaction_date BETWEEN ? AND ?",
			budget.UserID, "expense", budget.StartDate, budget.EndDate)

	if budget.CategoryID != nil {
		query = query.Where("category_id = ?", *budget.CategoryID)
	}

	query.Select("COALESCE(SUM(amount), 0)").Scan(&spentAmount)

	// Calculate metrics
	budget.SpentAmount = spentAmount
	budget.RemainingAmount = budget.Amount - spentAmount
	if budget.Amount > 0 {
		budget.UsagePercentage = (spentAmount / budget.Amount) * 100
	} else {
		budget.UsagePercentage = 0
	}
}

// GetBudgetAlerts returns budgets that are approaching or exceeding their limits
func (s *BudgetService) GetBudgetAlerts(userID uint64) ([]models.Budget, error) {
	var budgets []models.Budget
	if err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Preload("Category").Find(&budgets).Error; err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}

	var alerts []models.Budget
	for _, budget := range budgets {
		s.calculateBudgetMetrics(&budget)
		
		// Check if budget is approaching or exceeding threshold
		if budget.UsagePercentage >= budget.AlertThreshold {
			alerts = append(alerts, budget)
		}
	}

	return alerts, nil
}

// CheckBudgetNotifications checks and triggers budget notifications
func (s *BudgetService) CheckBudgetNotifications(userID uint64) error {
	dispatcher := NewNotificationDispatcher()
	budgets, err := s.GetBudgets(userID)
	if err != nil {
		return err
	}

	for _, budget := range budgets {
		// Check if budget needs notification
		if budget.UsagePercentage >= budget.AlertThreshold {
			// Check if budget is exceeded
			if budget.UsagePercentage >= 100 {
				// Budget exceeded
				if err := dispatcher.TriggerBudgetExceededAlert(userID, &budget); err != nil {
					log.Printf("Failed to trigger budget exceeded alert: %v", err)
				}
			} else {
				// Budget threshold reached
				if err := dispatcher.TriggerBudgetThresholdAlert(userID, &budget); err != nil {
					log.Printf("Failed to trigger budget threshold alert: %v", err)
				}
			}
		}
	}

	return nil
}
