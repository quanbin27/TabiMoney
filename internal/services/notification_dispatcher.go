package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type NotificationDispatcher struct {
	db              *gorm.DB
	config          *config.Config
	notificationSvc *NotificationService
	emailSvc        *EmailService
	telegramSvc     *TelegramService
}

type NotificationTrigger struct {
	UserID           uint64
	NotificationType string
	Priority         string
	Title            string
	Message          string
	Metadata         map[string]interface{}
}

func NewNotificationDispatcher(cfg *config.Config) *NotificationDispatcher {
	return &NotificationDispatcher{
		db:              database.GetDB(),
		config:          cfg,
		notificationSvc: NewNotificationService(),
		emailSvc:        NewEmailService(),
		telegramSvc:     NewTelegramService(),
	}
}

// DispatchNotification sends notification through all enabled channels
func (d *NotificationDispatcher) DispatchNotification(trigger NotificationTrigger) error {
	// Get user
	var user models.User
	if err := d.db.Preload("Profile").First(&user, trigger.UserID).Error; err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check user notification preferences
	preferences := d.getUserNotificationPreferences(&user)
	if !d.shouldSendNotification(preferences, trigger.NotificationType, trigger.Priority) {
		log.Printf("Notification skipped for user %d based on preferences", trigger.UserID)
		return nil
	}

	// Create notification record
	metadataJSON := "{}"
	if trigger.Metadata != nil {
		if data, err := json.Marshal(trigger.Metadata); err == nil {
			metadataJSON = string(data)
		}
	}

	notification, err := d.notificationSvc.Create(
		trigger.UserID,
		trigger.Title,
		trigger.Message,
		trigger.NotificationType,
		trigger.Priority,
		metadataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Send through different channels
	go d.sendEmailNotification(&user, notification, trigger.Metadata)
	go d.sendTelegramNotification(trigger.UserID, notification, trigger.Metadata)

	return nil
}

// getUserNotificationPreferences gets user's notification preferences
func (d *NotificationDispatcher) getUserNotificationPreferences(user *models.User) map[string]interface{} {
	preferences := map[string]interface{}{
		"email_enabled":      true,
		"telegram_enabled":   true,
		"in_app_enabled":     true,
		"budget_alerts":      true,
		"goal_alerts":        true,
		"ai_alerts":          true,
		"transaction_alerts": true,
		"analytics_alerts":   true,
	}

	if user.Profile != nil && user.Profile.NotificationSettings != "" {
		var userPrefs map[string]interface{}
		if err := json.Unmarshal([]byte(user.Profile.NotificationSettings), &userPrefs); err == nil {
			for k, v := range userPrefs {
				preferences[k] = v
			}
		}
	}

	return preferences
}

// shouldSendNotification checks if notification should be sent based on preferences
func (d *NotificationDispatcher) shouldSendNotification(preferences map[string]interface{}, notificationType, priority string) bool {
	// Check if notification type is enabled
	switch notificationType {
	case "warning":
		if priority == "urgent" || priority == "high" {
			return true // Always send urgent/high priority warnings
		}
		return preferences["budget_alerts"].(bool)
	case "info":
		return preferences["analytics_alerts"].(bool)
	case "success":
		return preferences["goal_alerts"].(bool)
	case "reminder":
		return preferences["budget_alerts"].(bool)
	default:
		return true
	}
}

// sendEmailNotification sends email notification
func (d *NotificationDispatcher) sendEmailNotification(user *models.User, notification *models.Notification, metadata map[string]interface{}) {
	preferences := d.getUserNotificationPreferences(user)
	if !preferences["email_enabled"].(bool) {
		return
	}

	emailData := EmailData{
		Title:            notification.Title,
		Message:          notification.Message,
		Priority:         notification.Priority,
		NotificationType: notification.NotificationType,
		Date:             notification.CreatedAt.Format("02/01/2006"),
		Time:             notification.CreatedAt.Format("15:04"),
	}

	// Add metadata to email data
	if metadata != nil {
		if amount, ok := metadata["amount"].(float64); ok {
			emailData.Amount = amount
		}
		if categoryName, ok := metadata["category_name"].(string); ok {
			emailData.CategoryName = categoryName
		}
		if budgetName, ok := metadata["budget_name"].(string); ok {
			emailData.BudgetName = budgetName
		}
		if goalName, ok := metadata["goal_name"].(string); ok {
			emailData.GoalName = goalName
		}
		if progress, ok := metadata["progress"].(float64); ok {
			emailData.Progress = progress
		}
	}

	if err := d.emailSvc.SendNotificationEmail(user, notification, emailData); err != nil {
		log.Printf("Failed to send email notification to user %d: %v", user.ID, err)
	}
}

// sendTelegramNotification sends telegram notification
func (d *NotificationDispatcher) sendTelegramNotification(userID uint64, notification *models.Notification, metadata map[string]interface{}) {
	preferences := d.getUserNotificationPreferences(&models.User{ID: userID})
	if !preferences["telegram_enabled"].(bool) {
		return
	}

	if err := d.telegramSvc.SendNotificationMessage(userID, notification, metadata); err != nil {
		log.Printf("Failed to send telegram notification to user %d: %v", userID, err)
	}
}

// Budget Notification Triggers

// TriggerBudgetThresholdAlert triggers budget threshold alert
func (d *NotificationDispatcher) TriggerBudgetThresholdAlert(userID uint64, budget *models.Budget) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "warning",
		Priority:         "high",
		Title:            "Ngân sách đạt ngưỡng cảnh báo",
		Message:          fmt.Sprintf("Ngân sách '%s' đã đạt %.1f%% ngưỡng cảnh báo (%.0f%%).", budget.Name, budget.UsagePercentage, budget.AlertThreshold),
		Metadata: map[string]interface{}{
			"budget_id":        budget.ID,
			"budget_name":      budget.Name,
			"amount":           budget.Amount,
			"usage_percentage": budget.UsagePercentage,
			"remaining_amount": budget.RemainingAmount,
			"alert_threshold":  budget.AlertThreshold,
		},
	}

	return d.DispatchNotification(trigger)
}

// TriggerBudgetExceededAlert triggers budget exceeded alert
func (d *NotificationDispatcher) TriggerBudgetExceededAlert(userID uint64, budget *models.Budget) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "warning",
		Priority:         "urgent",
		Title:            "Ngân sách đã vượt quá",
		Message:          fmt.Sprintf("Ngân sách '%s' đã vượt quá %.1f%%!", budget.Name, budget.UsagePercentage),
		Metadata: map[string]interface{}{
			"budget_id":        budget.ID,
			"budget_name":      budget.Name,
			"amount":           budget.Amount,
			"usage_percentage": budget.UsagePercentage,
			"exceeded_amount":  budget.SpentAmount - budget.Amount,
		},
	}

	return d.DispatchNotification(trigger)
}

// TriggerBudgetPacingAlert warns when spending pace is too fast
func (d *NotificationDispatcher) TriggerBudgetPacingAlert(userID uint64, budget *models.Budget, allowedPacePct, actualPacePct float64, daysLeft int) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "warning",
		Priority:         "medium",
		Title:            "Tốc độ chi vượt pace ngân sách",
		Message:          fmt.Sprintf("Ngân sách '%s' đang chi %.1f%% so với pace cho phép (%.1f%%). Còn %d ngày trong kỳ.", budget.Name, actualPacePct, allowedPacePct, daysLeft),
		Metadata: map[string]interface{}{
			"budget_id":        budget.ID,
			"budget_name":      budget.Name,
			"allowed_pace_pct": allowedPacePct,
			"actual_pace_pct":  actualPacePct,
			"days_left":        daysLeft,
		},
	}
	return d.DispatchNotification(trigger)
}

// TriggerBudgetAchievementAlert triggers budget achievement alert
func (d *NotificationDispatcher) TriggerBudgetAchievementAlert(userID uint64, budget *models.Budget) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "success",
		Priority:         "medium",
		Title:            "Hoàn thành tiết kiệm ngân sách",
		Message:          fmt.Sprintf("Chúc mừng! Bạn đã hoàn thành tiết kiệm ngân sách '%s'.", budget.Name),
		Metadata: map[string]interface{}{
			"budget_id":   budget.ID,
			"budget_name": budget.Name,
			"amount":      budget.Amount,
		},
	}

	return d.DispatchNotification(trigger)
}

// Goal Notification Triggers

// TriggerGoalProgressAlert triggers goal progress alert
func (d *NotificationDispatcher) TriggerGoalProgressAlert(userID uint64, goal *models.FinancialGoal, milestone string) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "info",
		Priority:         "medium",
		Title:            fmt.Sprintf("Mục tiêu đạt %s", milestone),
		Message:          fmt.Sprintf("Mục tiêu '%s' đã đạt %.1f%%!", goal.Title, goal.Progress),
		Metadata: map[string]interface{}{
			"goal_id":   goal.ID,
			"goal_name": goal.Title,
			"amount":    goal.TargetAmount,
			"progress":  goal.Progress,
			"milestone": milestone,
		},
	}

	return d.DispatchNotification(trigger)
}

// TriggerGoalDeadlineAlert triggers goal deadline alert
func (d *NotificationDispatcher) TriggerGoalDeadlineAlert(userID uint64, goal *models.FinancialGoal, daysLeft int) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "warning",
		Priority:         "high",
		Title:            "Cảnh báo hạn chót mục tiêu",
		Message:          fmt.Sprintf("Mục tiêu '%s' còn %d ngày nữa đến hạn!", goal.Title, daysLeft),
		Metadata: map[string]interface{}{
			"goal_id":   goal.ID,
			"goal_name": goal.Title,
			"amount":    goal.TargetAmount,
			"progress":  goal.Progress,
			"days_left": daysLeft,
		},
	}

	return d.DispatchNotification(trigger)
}

// TriggerGoalAchievedAlert triggers goal achieved alert
func (d *NotificationDispatcher) TriggerGoalAchievedAlert(userID uint64, goal *models.FinancialGoal) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "success",
		Priority:         "high",
		Title:            "Chúc mừng hoàn thành mục tiêu!",
		Message:          fmt.Sprintf("Chúc mừng! Bạn đã hoàn thành mục tiêu '%s'!", goal.Title),
		Metadata: map[string]interface{}{
			"goal_id":   goal.ID,
			"goal_name": goal.Title,
			"amount":    goal.TargetAmount,
			"progress":  goal.Progress,
		},
	}

	return d.DispatchNotification(trigger)
}

// AI Notification Triggers

// TriggerAnomalyAlert triggers anomaly detection alert
func (d *NotificationDispatcher) TriggerAnomalyAlert(userID uint64, anomaly *models.Anomaly) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "warning",
		Priority:         "high",
		Title:            "Phát hiện giao dịch bất thường",
		Message:          fmt.Sprintf("Giao dịch %.0f VND tại %s có vẻ bất thường (điểm số: %.2f).", anomaly.Amount, anomaly.CategoryName, anomaly.AnomalyScore),
		Metadata: map[string]interface{}{
			"transaction_id": anomaly.TransactionID,
			"amount":         anomaly.Amount,
			"category_name":  anomaly.CategoryName,
			"anomaly_score":  anomaly.AnomalyScore,
			"anomaly_type":   anomaly.AnomalyType,
		},
	}

	return d.DispatchNotification(trigger)
}

// TriggerSpendingPredictionAlert triggers spending prediction alert
func (d *NotificationDispatcher) TriggerSpendingPredictionAlert(userID uint64, prediction *models.ExpensePredictionResponse) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "info",
		Priority:         "medium",
		Title:            "Dự đoán chi tiêu tháng tới",
		Message:          fmt.Sprintf("Dự đoán chi tiêu tháng tới: %.0f VND (độ tin cậy: %.1f%%)", prediction.PredictedAmount, prediction.ConfidenceScore*100),
		Metadata: map[string]interface{}{
			"predicted_amount": prediction.PredictedAmount,
			"confidence_score": prediction.ConfidenceScore,
			"recommendations":  prediction.Recommendations,
		},
	}

	return d.DispatchNotification(trigger)
}

// Transaction Notification Triggers

// TriggerLargeTransactionAlert triggers large transaction alert
func (d *NotificationDispatcher) TriggerLargeTransactionAlert(userID uint64, transaction *models.Transaction, threshold float64) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "warning",
		Priority:         "medium",
		Title:            "Giao dịch lớn được phát hiện",
		Message:          fmt.Sprintf("Giao dịch %.0f VND tại %s vượt quá ngưỡng %.0f VND", transaction.Amount, transaction.Category.Name, threshold),
		Metadata: map[string]interface{}{
			"transaction_id": transaction.ID,
			"amount":         transaction.Amount,
			"category_name":  transaction.Category.Name,
			"description":    transaction.Description,
			"threshold":      threshold,
		},
	}

	return d.DispatchNotification(trigger)
}

// Analytics Notification Triggers

// TriggerMonthlyReportAlert triggers monthly report alert
func (d *NotificationDispatcher) TriggerMonthlyReportAlert(userID uint64, analytics *models.DashboardAnalytics) error {
	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: "info",
		Priority:         "low",
		Title:            "Báo cáo tài chính hàng tháng",
		Message:          fmt.Sprintf("Báo cáo tháng %s: Thu %.0f VND, Chi %.0f VND, Chênh lệch %.0f VND", analytics.Period, analytics.TotalIncome, analytics.TotalExpense, analytics.NetAmount),
		Metadata: map[string]interface{}{
			"period":        analytics.Period,
			"total_income":  analytics.TotalIncome,
			"total_expense": analytics.TotalExpense,
			"net_amount":    analytics.NetAmount,
			"health_score":  analytics.FinancialHealth.Score,
		},
	}

	return d.DispatchNotification(trigger)
}

// TriggerFinancialHealthAlert triggers financial health alert
func (d *NotificationDispatcher) TriggerFinancialHealthAlert(userID uint64, health *models.FinancialHealth, period string) error {
	var priority string
	var notificationType string

	switch health.Level {
	case "excellent":
		priority = "low"
		notificationType = "success"
	case "good":
		priority = "low"
		notificationType = "info"
	case "fair":
		priority = "medium"
		notificationType = "warning"
	case "poor":
		priority = "high"
		notificationType = "warning"
	default:
		priority = "medium"
		notificationType = "info"
	}

	trigger := NotificationTrigger{
		UserID:           userID,
		NotificationType: notificationType,
		Priority:         priority,
		Title:            fmt.Sprintf("Sức khỏe tài chính: %s", health.Level),
		Message:          fmt.Sprintf("Sức khỏe tài chính tháng %s: %.1f/100 điểm. Tỷ lệ tiết kiệm: %.1f%%", period, health.Score, health.SavingsRate),
		Metadata: map[string]interface{}{
			"period":          period,
			"health_score":    health.Score,
			"health_level":    health.Level,
			"savings_rate":    health.SavingsRate,
			"recommendations": health.Recommendations,
		},
	}

	return d.DispatchNotification(trigger)
}

// Scheduled Notification Triggers

// CheckAndTriggerScheduledNotifications checks for scheduled notifications
func (d *NotificationDispatcher) CheckAndTriggerScheduledNotifications() error {
	// Check budget alerts
	if err := d.checkBudgetAlerts(); err != nil {
		log.Printf("Failed to check budget alerts: %v", err)
	}

	// Check goal alerts
	if err := d.checkGoalAlerts(); err != nil {
		log.Printf("Failed to check goal alerts: %v", err)
	}

	// Check monthly reports
	if err := d.checkMonthlyReports(); err != nil {
		log.Printf("Failed to check monthly reports: %v", err)
	}

	return nil
}

// checkBudgetAlerts checks for budget alerts that need to be sent
func (d *NotificationDispatcher) checkBudgetAlerts() error {
	now := time.Now()
	var budgets []models.Budget
	// Chỉ lấy ngân sách đang hoạt động và đã bắt đầu (tránh báo sớm cho budget tương lai)
	if err := d.db.Where("is_active = ? AND start_date <= ? AND end_date >= ?", true, now, now).Find(&budgets).Error; err != nil {
		return err
	}

	for _, budget := range budgets {
		// Calculate current metrics
		bs := NewBudgetService(d.config)
		bs.calculateBudgetMetrics(&budget)

		// Check if budget needs alert
		if budget.UsagePercentage >= budget.AlertThreshold {
			// Check if we already sent this alert recently
			var recentNotification models.Notification
			oneDayAgo := time.Now().Add(-24 * time.Hour)

			err := d.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND metadata LIKE ?",
				budget.UserID, "warning", oneDayAgo, fmt.Sprintf("%%\"budget_id\":%d%%", budget.ID)).
				First(&recentNotification).Error

			if err == gorm.ErrRecordNotFound {
				// No recent alert, send one
				if budget.UsagePercentage >= 100 {
					d.TriggerBudgetExceededAlert(budget.UserID, &budget)
				} else {
					d.TriggerBudgetThresholdAlert(budget.UserID, &budget)
				}
			}
		}
	}

	return nil
}

// checkGoalAlerts checks for goal alerts that need to be sent
func (d *NotificationDispatcher) checkGoalAlerts() error {
	var goals []models.FinancialGoal
	if err := d.db.Where("is_achieved = ? AND target_date IS NOT NULL", false).Find(&goals).Error; err != nil {
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

				err := d.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND metadata LIKE ?",
					goal.UserID, "warning", oneWeekAgo, fmt.Sprintf("%%\"goal_id\":%d%%", goal.ID)).
					First(&recentNotification).Error

				if err == gorm.ErrRecordNotFound {
					d.TriggerGoalDeadlineAlert(goal.UserID, &goal, daysLeft)
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

				err := d.db.Where("user_id = ? AND notification_type = ? AND created_at > ? AND metadata LIKE ?",
					goal.UserID, "info", oneWeekAgo, fmt.Sprintf("%%\"milestone\":\"%.0f%%\"%%", milestone)).
					First(&recentNotification).Error

				if err == gorm.ErrRecordNotFound {
					d.TriggerGoalProgressAlert(goal.UserID, &goal, fmt.Sprintf("%.0f%%", milestone))
				}
				break
			}
		}
	}

	return nil
}

// checkMonthlyReports checks for monthly reports that need to be sent
func (d *NotificationDispatcher) checkMonthlyReports() error {
	// Check if it's the first day of the month
	now := time.Now()
	if now.Day() != 1 {
		return nil
	}

	// Get all active users
	var users []models.User
	if err := d.db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		// Check if we already sent monthly report for this month
		var recentNotification models.Notification
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

		err := d.db.Where("user_id = ? AND notification_type = ? AND created_at >= ?",
			user.ID, "info", startOfMonth).
			Where("title LIKE ?", "%Báo cáo tài chính hàng tháng%").
			First(&recentNotification).Error

		if err == gorm.ErrRecordNotFound {
			// Generate and send monthly report
			ts := NewTransactionService(d.config)
			analytics, err := ts.GetMonthlySummary(user.ID, now.Year(), int(now.Month()-1))
			if err == nil {
				d.TriggerMonthlyReportAlert(user.ID, analytics)
			}
		}
	}

	return nil
}
