package services

import (
	"fmt"
	"log"
	"math"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type BudgetService struct {
	db     *gorm.DB
	config *config.Config
}

func NewBudgetService(cfg *config.Config) *BudgetService {
	return &BudgetService{
		db:     database.GetDB(),
		config: cfg,
	}
}

// CreateBudget creates a new budget
func (s *BudgetService) CreateBudget(userID uint64, req *models.BudgetCreateRequest) (*models.Budget, error) {
	// Validate basic date range
	if req.StartDate.After(req.EndDate) {
		return nil, fmt.Errorf("start_date must be before or equal to end_date")
	}

	// Prevent multiple active budgets for same category & overlapping time
	if req.CategoryID != nil {
		var count int64
		if err := s.db.Model(&models.Budget{}).
			Where("user_id = ? AND is_active = ? AND category_id = ?", userID, true, *req.CategoryID).
			Where("NOT (end_date < ? OR start_date > ?)", req.StartDate, req.EndDate).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to check existing budgets: %w", err)
		}
		if count > 0 {
			return nil, fmt.Errorf("there is already an active budget for this category in the selected period")
		}
	}

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

	// Validate basic date range
	if req.StartDate.After(req.EndDate) {
		return nil, fmt.Errorf("start_date must be before or equal to end_date")
	}

	// Prevent overlapping active budgets for same category (excluding current budget)
	if req.CategoryID != nil && req.IsActive {
		var count int64
		if err := s.db.Model(&models.Budget{}).
			Where("user_id = ? AND is_active = ? AND category_id = ? AND id <> ?", userID, true, *req.CategoryID, budgetID).
			Where("NOT (end_date < ? OR start_date > ?)", req.StartDate, req.EndDate).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to check existing budgets: %w", err)
		}
		if count > 0 {
			return nil, fmt.Errorf("another active budget for this category already exists in the selected period")
		}
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

// CheckBudgetNotifications checks and triggers budget notifications
func (s *BudgetService) CheckBudgetNotifications(userID uint64) error {
	dispatcher := NewNotificationDispatcher(s.config)

	// Chỉ kiểm tra các budget đang hoạt động
	var budgets []models.Budget
	if err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&budgets).Error; err != nil {
		return fmt.Errorf("failed to load budgets for notifications: %w", err)
	}

	now := time.Now()

	for i := range budgets {
		// Bỏ qua ngân sách không nằm trong khoảng thời gian hiện tại
		if now.Before(budgets[i].StartDate) || now.After(budgets[i].EndDate) {
			continue
		}

		// Tính lại metrics để đảm bảo số liệu mới nhất
		s.calculateBudgetMetrics(&budgets[i])

		// Check if budget needs notification
		if budgets[i].UsagePercentage >= budgets[i].AlertThreshold {
			// Check if budget is exceeded
			if budgets[i].UsagePercentage >= 100 {
				// Budget exceeded
				if err := dispatcher.TriggerBudgetExceededAlert(userID, &budgets[i]); err != nil {
					log.Printf("Failed to trigger budget exceeded alert: %v", err)
				}
			} else {
				// Budget threshold reached
				if err := dispatcher.TriggerBudgetThresholdAlert(userID, &budgets[i]); err != nil {
					log.Printf("Failed to trigger budget threshold alert: %v", err)
				}
			}
		}
	}

	return nil
}

// GetBudgetInsights computes safe-to-spend and pacing information for active budgets in current period
func (s *BudgetService) GetBudgetInsights(userID uint64) (*models.BudgetInsights, error) {
    // Load active budgets
    var budgets []models.Budget
    if err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&budgets).Error; err != nil {
        return nil, fmt.Errorf("failed to load budgets: %w", err)
    }

    now := time.Now()
    // Determine current period from budgets (default monthly)
    period := "monthly"
    if len(budgets) > 0 {
        period = budgets[0].Period
    }

    var totalRemaining float64
    var daysLeft int
    insights := &models.BudgetInsights{
        UserID: userID,
        Period: period,
        AsOf:   now,
    }

    for i := range budgets {
        s.calculateBudgetMetrics(&budgets[i])
        // Only consider budgets whose window includes now
        if now.Before(budgets[i].StartDate) || now.After(budgets[i].EndDate) {
            continue
        }
        // days left including today
        dl := int(math.Max(1, math.Ceil(budgets[i].EndDate.Sub(now).Hours()/24)))
        if daysLeft == 0 || dl < daysLeft {
            daysLeft = dl
        }
        totalRemaining += math.Max(0, budgets[i].RemainingAmount)

        // compute pacing
        totalDays := int(math.Max(1, math.Round(budgets[i].EndDate.Sub(budgets[i].StartDate).Hours()/24)))
        elapsedDays := int(math.Max(1, math.Round(now.Sub(budgets[i].StartDate).Hours()/24)))
        allowedPace := 100.0 * float64(elapsedDays) / float64(totalDays)
        actualPace := budgets[i].UsagePercentage
        bp := models.BudgetPace{
            BudgetID:        budgets[i].ID,
            Name:            budgets[i].Name,
            CategoryID:      budgets[i].CategoryID,
            Amount:          budgets[i].Amount,
            SpentAmount:     budgets[i].SpentAmount,
            RemainingAmount: budgets[i].RemainingAmount,
            UsagePercentage: budgets[i].UsagePercentage,
            AllowedPacePct:  math.Min(100, math.Max(0, allowedPace)),
            ActualPacePct:   math.Min(100, math.Max(0, actualPace)),
            IsOverPace:      actualPace > allowedPace*1.2, // 120% of allowed pace considered risky
        }
        insights.Budgets = append(insights.Budgets, bp)
        if bp.IsOverPace {
            insights.RiskBudgetIDs = append(insights.RiskBudgetIDs, bp.BudgetID)
        }
    }

    insights.TotalRemaining = totalRemaining
    insights.DaysLeft = daysLeft
    if daysLeft <= 0 {
        insights.SafeToSpendDaily = 0
        insights.SafeToSpendWeekly = 0
    } else {
        insights.SafeToSpendDaily = totalRemaining / float64(daysLeft)
        insights.SafeToSpendWeekly = insights.SafeToSpendDaily * 7
    }

    // projected end usage: average of per-budget projection weighted by amount
    var weightedUsage, sumAmount float64
    for _, b := range insights.Budgets {
        if b.Amount > 0 {
            weightedUsage += b.UsagePercentage * b.Amount
            sumAmount += b.Amount
        }
    }
    if sumAmount > 0 {
        insights.ProjectedEndUsagePct = weightedUsage / sumAmount
    }

    return insights, nil
}

// SuggestBudgets proposes category budgets based on recent spending or 50/30/20
func (s *BudgetService) SuggestBudgets(userID uint64) (*models.AutoBudgetSuggestResponse, error) {
    now := time.Now()
    startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())

    // fetch profile for income
    var profile models.UserProfile
    _ = s.db.Where("user_id = ?", userID).First(&profile).Error

    // aggregate last 3 months expenses by category
    threeMonthsAgo := startOfMonth.AddDate(0, -3, 0)
    type row struct{ CategoryID uint64; Name string; Amount float64 }
    var rows []row
    s.db.Table("transactions t").
        Select("t.category_id as category_id, c.name as name, COALESCE(SUM(t.amount),0) as amount").
        Joins("JOIN categories c ON c.id = t.category_id").
        Where("t.user_id = ? AND t.transaction_type = 'expense' AND t.transaction_date BETWEEN ? AND ?", userID, threeMonthsAgo, endOfMonth).
        Group("t.category_id, c.name").
        Scan(&rows)

    suggestions := []models.AutoBudgetSuggestion{}
    var total float64
    if len(rows) > 0 {
        for _, r := range rows {
            medianApprox := r.Amount / 3.0
            suggested := medianApprox * 0.9 // safety margin
            cid := r.CategoryID
            suggestions = append(suggestions, models.AutoBudgetSuggestion{
                CategoryID:   &cid,
                Name:         r.Name,
                SuggestedAmt: math.Max(0, math.Round(suggested*100)/100),
            })
            total += suggested
        }
    } else {
        // fallback 50/30/20 if no history
        income := profile.MonthlyIncome
        needs := income * 0.5
        wants := income * 0.3
        savings := income * 0.2

        // map needs: food 40%, transport 20%, bills 30%, healthcare 10%
        // Try to find system categories by name_en; fallback to nil
        type catRow struct{ ID uint64; Name string; NameEn string }
        var cats []catRow
        s.db.Raw("SELECT id, name, name_en FROM categories WHERE is_system = TRUE").Scan(&cats)
        findCat := func(nameEn string) *uint64 {
            for _, c := range cats {
                if c.NameEn == nameEn {
                    id := c.ID
                    return &id
                }
            }
            return nil
        }
        add := func(name string, amt float64, cid *uint64) {
            suggestions = append(suggestions, models.AutoBudgetSuggestion{CategoryID: cid, Name: name, SuggestedAmt: math.Round(amt*100) / 100})
            total += amt
        }
        add("Food & Dining", needs*0.4, findCat("Food & Dining"))
        add("Transportation", needs*0.2, findCat("Transportation"))
        add("Bills", needs*0.3, nil)
        add("Healthcare", needs*0.1, findCat("Healthcare"))
        add("Entertainment", wants*0.4, findCat("Entertainment"))
        add("Shopping", wants*0.4, findCat("Shopping"))
        add("Other", wants*0.2, findCat("Other"))
        add("Savings", savings, findCat("Savings"))
    }

    resp := &models.AutoBudgetSuggestResponse{
        UserID:         userID,
        MonthlyIncome:  profile.MonthlyIncome,
        Period:         "monthly",
        StartDate:      startOfMonth,
        EndDate:        endOfMonth,
        Suggestions:    suggestions,
        TotalSuggested: math.Round(total*100) / 100,
        Notes:          []string{"Suggestions based on recent spending or 50/30/20 fallback", "Safety margin applied where applicable"},
    }
    return resp, nil
}

// CreateBudgetsFromSuggestions creates budgets from suggestions payload
func (s *BudgetService) CreateBudgetsFromSuggestions(userID uint64, req *models.AutoBudgetCreateRequest) ([]models.Budget, error) {
    if req == nil || len(req.Budgets) == 0 {
        return nil, fmt.Errorf("no budgets provided")
    }

	// Validate basic period and date range
	if req.StartDate.After(req.EndDate) {
		return nil, fmt.Errorf("start_date must be before or equal to end_date")
	}
	if req.Period == "" {
		req.Period = "monthly"
	}

	var created []models.Budget
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, b := range req.Budgets {
			if b.SuggestedAmt <= 0 {
				continue // bỏ qua các đề xuất 0 hoặc âm
			}

			// Nếu có category thì tránh tạo budget trùng khoảng thời gian với budget đang active
			if b.CategoryID != nil {
				var count int64
				if err := tx.Model(&models.Budget{}).
					Where("user_id = ? AND is_active = ? AND category_id = ?", userID, true, *b.CategoryID).
					Where("NOT (end_date < ? OR start_date > ?)", req.StartDate, req.EndDate).
					Count(&count).Error; err != nil {
					return fmt.Errorf("failed to check existing budgets: %w", err)
				}
				if count > 0 {
					// Bỏ qua đề xuất này, vì đã có ngân sách active cùng category & thời gian
					continue
				}
			}

			name := b.Name
			if name == "" {
				name = "Budget"
			}
			budget := models.Budget{
				UserID:         userID,
				CategoryID:     b.CategoryID,
				Name:           name,
				Amount:         b.SuggestedAmt,
				Period:         req.Period,
				StartDate:      req.StartDate,
				EndDate:        req.EndDate,
				IsActive:       true,
				AlertThreshold: req.AlertThreshold,
			}
			if err := tx.Create(&budget).Error; err != nil {
				return fmt.Errorf("failed to create budget: %w", err)
			}
			// Tính metrics sau khi tạo để trả về cho FE, không cần transaction
			s.calculateBudgetMetrics(&budget)
			created = append(created, budget)
		}
		if len(created) == 0 {
			return fmt.Errorf("no valid budgets created from suggestions")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}
