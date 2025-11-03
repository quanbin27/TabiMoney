package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"tabimoney/internal/database"
	"tabimoney/internal/config"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type ScheduledNotificationService struct {
	dispatcher *NotificationDispatcher
	db         *gorm.DB
}

func NewScheduledNotificationService() *ScheduledNotificationService {
	return &ScheduledNotificationService{
		dispatcher: NewNotificationDispatcher(),
		db:         database.GetDB(),
	}
}

// StartScheduler starts the scheduled notification service
func (s *ScheduledNotificationService) StartScheduler(ctx context.Context) {
	log.Println("Starting scheduled notification service...")

    // Run every 7 days
    ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	// Run immediately on startup
	s.runScheduledTasks()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduled notification service stopped")
			return
		case <-ticker.C:
			s.runScheduledTasks()
		}
	}
}

// runScheduledTasks runs all scheduled notification tasks
func (s *ScheduledNotificationService) runScheduledTasks() {
	log.Println("Running scheduled notification tasks...")

	// Check budget alerts
	if err := s.checkBudgetAlerts(); err != nil {
		log.Printf("Failed to check budget alerts: %v", err)
	}

	// Check goal alerts
	if err := s.checkGoalAlerts(); err != nil {
		log.Printf("Failed to check goal alerts: %v", err)
	}

	// Check monthly reports
	if err := s.checkMonthlyReports(); err != nil {
		log.Printf("Failed to check monthly reports: %v", err)
	}

	// Check financial health alerts
	if err := s.checkFinancialHealthAlerts(); err != nil {
		log.Printf("Failed to check financial health alerts: %v", err)
	}

    // Run anomaly detection for all users (hourly)
    if err := s.RunAnomalyDetection(); err != nil {
        log.Printf("Failed to run anomaly detection: %v", err)
    }

	log.Println("Scheduled notification tasks completed")
}

// checkBudgetAlerts checks for budget alerts that need to be sent
func (s *ScheduledNotificationService) checkBudgetAlerts() error {
	var budgets []models.Budget
	if err := s.db.Where("is_active = ? AND end_date >= ?", true, time.Now()).Find(&budgets).Error; err != nil {
		return err
	}

	for _, budget := range budgets {
		// Calculate current metrics
		bs := NewBudgetService()
		bs.calculateBudgetMetrics(&budget)

		// Check if budget needs alert
		if budget.UsagePercentage >= budget.AlertThreshold {
			// Check if we already sent this alert recently
			var recentNotification models.Notification
			oneDayAgo := time.Now().Add(-24 * time.Hour)
			
			err := s.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND metadata LIKE ?", 
				budget.UserID, "warning", oneDayAgo, fmt.Sprintf("%%\"budget_id\":%d%%", budget.ID)).
				First(&recentNotification).Error
			
			if err == gorm.ErrRecordNotFound {
				// No recent alert, send one
				if budget.UsagePercentage >= 100 {
					s.dispatcher.TriggerBudgetExceededAlert(budget.UserID, &budget)
				} else {
					s.dispatcher.TriggerBudgetThresholdAlert(budget.UserID, &budget)
				}
			}
		}
	}

	return nil
}

// checkGoalAlerts checks for goal alerts that need to be sent
func (s *ScheduledNotificationService) checkGoalAlerts() error {
	var goals []models.FinancialGoal
	if err := s.db.Where("is_achieved = ? AND target_date IS NOT NULL", false).Find(&goals).Error; err != nil {
		return err
	}

	for _, goal := range goals {
		// Calculate progress
		if goal.TargetAmount > 0 {
			goal.Progress = (goal.CurrentAmount / goal.TargetAmount) * 100
		}

		// Check deadline warning (30 days before)
		if goal.TargetDate != nil {
			daysLeft := int(goal.TargetDate.Sub(time.Now()).Hours() / 24)
			if daysLeft <= 30 && daysLeft > 0 {
				// Check if we already sent this alert recently
				var recentNotification models.Notification
				oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)
				
				err := s.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND metadata LIKE ?", 
					goal.UserID, "warning", oneWeekAgo, fmt.Sprintf("%%\"goal_id\":%d%%", goal.ID)).
					First(&recentNotification).Error
				
				if err == gorm.ErrRecordNotFound {
					s.dispatcher.TriggerGoalDeadlineAlert(goal.UserID, &goal, daysLeft)
				}
			}
		}

		// Check progress milestones
		milestones := []float64{25, 50, 75, 90}
		for _, milestone := range milestones {
			if goal.Progress >= milestone && goal.Progress < milestone+5 {
				// Check if we already sent this milestone alert
				var recentNotification models.Notification
				oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)
				
				err := s.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND metadata LIKE ?", 
					goal.UserID, "info", oneWeekAgo, fmt.Sprintf("%%\"milestone\":\"%.0f%%\"%%", milestone)).
					First(&recentNotification).Error
				
				if err == gorm.ErrRecordNotFound {
					s.dispatcher.TriggerGoalProgressAlert(goal.UserID, &goal, fmt.Sprintf("%.0f%%", milestone))
				}
				break
			}
		}
	}

	return nil
}

// checkMonthlyReports checks for monthly reports that need to be sent
func (s *ScheduledNotificationService) checkMonthlyReports() error {
	// Check if it's the first day of the month
	now := time.Now()
	if now.Day() != 1 {
		return nil
	}

	// Get all active users
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		// Check if we already sent monthly report for this month
		var recentNotification models.Notification
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		
		err := s.db.Where("user_id = ? AND notification_type = ? AND created_at >= ?", 
			user.ID, "info", startOfMonth).
			Where("title LIKE ?", "%Báo cáo tài chính hàng tháng%").
			First(&recentNotification).Error
		
		if err == gorm.ErrRecordNotFound {
			// Generate and send monthly report
			ts := NewTransactionService()
			analytics, err := ts.GetMonthlySummary(user.ID, now.Year(), int(now.Month()-1))
			if err == nil {
				s.dispatcher.TriggerMonthlyReportAlert(user.ID, analytics)
			}
		}
	}

	return nil
}

// checkFinancialHealthAlerts checks for financial health alerts
func (s *ScheduledNotificationService) checkFinancialHealthAlerts() error {
	// Check if it's the first day of the month
	now := time.Now()
	if now.Day() != 1 {
		return nil
	}

	// Get all active users
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		// Generate monthly analytics
		ts := NewTransactionService()
		analytics, err := ts.GetMonthlySummary(user.ID, now.Year(), int(now.Month()-1))
		if err != nil {
			continue
		}

		// Check if financial health is poor
		if analytics.FinancialHealth.Level == "poor" {
			// Check if we already sent this alert recently
			var recentNotification models.Notification
			oneMonthAgo := time.Now().AddDate(0, -1, 0)
			
			err := s.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND title LIKE ?", 
				user.ID, "warning", oneMonthAgo, "%Sức khỏe tài chính%").
				First(&recentNotification).Error
			
			if err == gorm.ErrRecordNotFound {
				s.dispatcher.TriggerFinancialHealthAlert(user.ID, &analytics.FinancialHealth, analytics.Period)
			}
		}
	}

	return nil
}

// RunAnomalyDetection runs anomaly detection for all users
// Note: This requires AIService with proper config. 
// For now, this is a placeholder - anomaly detection should be triggered via API endpoints
func (s *ScheduledNotificationService) RunAnomalyDetection() error {
    log.Println("Running scheduled anomaly detection for all users...")
    // Load users
    var users []models.User
    if err := s.db.Find(&users).Error; err != nil {
        return err
    }

    // Init AI service
    cfg, err := config.Load()
    if err != nil {
        return err
    }
    aiSvc := NewAIService(cfg)

    // Use last 30 days window
    end := time.Now()
    start := end.AddDate(0, 0, -30)

    for _, u := range users {
        req := &models.AnomalyDetectionRequest{
            UserID:    u.ID,
            StartDate: start,
            EndDate:   end,
            Threshold: 0.9,
        }
        if _, err := aiSvc.DetectAnomalies(req); err != nil {
            log.Printf("Scheduled anomaly detection failed for user %d: %v", u.ID, err)
        }
    }
    return nil
}

// RunSpendingPrediction runs spending prediction for all users
// Note: This requires AIService with proper config.
// For now, this is a placeholder - spending prediction should be triggered via API endpoints
func (s *ScheduledNotificationService) RunSpendingPrediction() error {
    log.Println("Running scheduled spending prediction for all users...")
    // Load users
    var users []models.User
    if err := s.db.Find(&users).Error; err != nil {
        return err
    }

    // Init AI service
    cfg, err := config.Load()
    if err != nil {
        return err
    }
    aiSvc := NewAIService(cfg)

    // Predict for next month based on last 30 days
    now := time.Now()
    start := now.AddDate(0, 0, -30)
    end := now

    for _, u := range users {
        req := &models.ExpensePredictionRequest{
            UserID:    u.ID,
            StartDate: start,
            EndDate:   end,
        }
        if _, err := aiSvc.PredictExpenses(req); err != nil {
            log.Printf("Scheduled spending prediction failed for user %d: %v", u.ID, err)
        }
    }
    return nil
}
