package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type TelegramService struct {
	botToken string
	apiURL   string
	db       *gorm.DB
}

type TelegramMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type TelegramKeyboard struct {
	InlineKeyboard [][]TelegramInlineButton `json:"inline_keyboard"`
}

type TelegramInlineButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
	Result      struct {
		MessageID int64 `json:"message_id"`
	} `json:"result,omitempty"`
}

func NewTelegramService() *TelegramService {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s", botToken)

	return &TelegramService{
		botToken: botToken,
		apiURL:   apiURL,
		db:       database.GetDB(),
	}
}

// SendNotificationMessage sends a notification message to user's Telegram
func (s *TelegramService) SendNotificationMessage(userID uint64, notification *models.Notification, data map[string]interface{}) error {
	if s.botToken == "" {
		log.Printf("Telegram bot token not configured, skipping Telegram notification for user %d", userID)
		return nil
	}

	// Get user's Telegram chat ID
	chatID, err := s.getUserTelegramChatID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user telegram chat ID: %w", err)
	}

	if chatID == 0 {
		log.Printf("User %d has no Telegram account linked", userID)
		return nil
	}

	// Format message based on notification type
	message := s.formatNotificationMessage(notification, data)

	// Send message
	return s.sendMessage(chatID, message, notification)
}

// getUserTelegramChatID gets user's Telegram chat ID from database
func (s *TelegramService) getUserTelegramChatID(userID uint64) (int64, error) {
	var telegramAccount models.TelegramAccount
	if err := s.db.Where("web_user_id = ?", userID).First(&telegramAccount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // User has no Telegram account
		}
		return 0, err
	}
	return telegramAccount.TelegramUserID, nil
}

// formatNotificationMessage formats notification message for Telegram
func (s *TelegramService) formatNotificationMessage(notification *models.Notification, data map[string]interface{}) string {
	var message string

	// Add emoji based on notification type
	switch notification.NotificationType {
	case "warning":
		if notification.Priority == "urgent" || notification.Priority == "high" {
			message += "🚨 *CẢNH BÁO KHẨN CẤP*\n\n"
		} else {
			message += "⚠️ *CẢNH BÁO*\n\n"
		}
	case "error":
		message += "❌ *LỖI*\n\n"
	case "success":
		message += "✅ *THÀNH CÔNG*\n\n"
	case "reminder":
		message += "🔔 *NHẮC NHỞ*\n\n"
	default:
		message += "📊 *THÔNG BÁO*\n\n"
	}

	// Add title
	message += fmt.Sprintf("*%s*\n\n", notification.Title)

	// Add message content
	message += notification.Message + "\n\n"

	// Add specific data based on notification type
	if amount, ok := data["amount"].(float64); ok && amount > 0 {
		message += fmt.Sprintf("💰 Số tiền: *%.0f VND*\n", amount)
	}

	if categoryName, ok := data["category_name"].(string); ok && categoryName != "" {
		message += fmt.Sprintf("📂 Danh mục: *%s*\n", categoryName)
	}

	if budgetName, ok := data["budget_name"].(string); ok && budgetName != "" {
		message += fmt.Sprintf("📊 Ngân sách: *%s*\n", budgetName)
	}

	if goalName, ok := data["goal_name"].(string); ok && goalName != "" {
		message += fmt.Sprintf("🎯 Mục tiêu: *%s*\n", goalName)
	}

	if progress, ok := data["progress"].(float64); ok && progress > 0 {
		message += fmt.Sprintf("📈 Tiến độ: *%.1f%%*\n", progress)
	}

	if usagePercentage, ok := data["usage_percentage"].(float64); ok && usagePercentage > 0 {
		message += fmt.Sprintf("📊 Sử dụng: *%.1f%%*\n", usagePercentage)
	}

	// Add timestamp
	message += fmt.Sprintf("\n🕐 %s", notification.CreatedAt.Format("02/01/2006 15:04"))

	return message
}

// sendMessage sends message to Telegram
func (s *TelegramService) sendMessage(chatID int64, text string, notification *models.Notification) error {
	// Prepare message payload
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}

	// Send HTTP request
	resp, err := http.Post(s.apiURL+"/sendMessage", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	log.Printf("Telegram message sent successfully to chat %d", chatID)
	return nil
}

// SendBudgetAlert sends budget alert to Telegram
func (s *TelegramService) SendBudgetAlert(userID uint64, budget *models.Budget, alertType string) error {
	notification := &models.Notification{
		Title:            s.getBudgetAlertTitle(alertType),
		Message:          s.getBudgetAlertMessage(budget, alertType),
		NotificationType: "warning",
		Priority:         s.getBudgetAlertPriority(alertType),
		CreatedAt:        time.Now(),
	}

	data := map[string]interface{}{
		"budget_name":      budget.Name,
		"amount":           budget.Amount,
		"usage_percentage": budget.UsagePercentage,
		"remaining_amount": budget.RemainingAmount,
	}

	return s.SendNotificationMessage(userID, notification, data)
}

// SendGoalAlert sends goal alert to Telegram
func (s *TelegramService) SendGoalAlert(userID uint64, goal *models.FinancialGoal, alertType string) error {
	notification := &models.Notification{
		Title:            s.getGoalAlertTitle(alertType),
		Message:          s.getGoalAlertMessage(goal, alertType),
		NotificationType: s.getGoalAlertType(alertType),
		Priority:         s.getGoalAlertPriority(alertType),
		CreatedAt:        time.Now(),
	}

	data := map[string]interface{}{
		"goal_name": goal.Title,
		"amount":    goal.TargetAmount,
		"progress":  goal.Progress,
	}

	return s.SendNotificationMessage(userID, notification, data)
}

// SendAnomalyAlert sends anomaly detection alert to Telegram
func (s *TelegramService) SendAnomalyAlert(userID uint64, anomaly *models.Anomaly) error {
	notification := &models.Notification{
		Title:            "Phát hiện giao dịch bất thường",
		Message:          fmt.Sprintf("Giao dịch %.0f VND tại %s có vẻ bất thường", anomaly.Amount, anomaly.CategoryName),
		NotificationType: "warning",
		Priority:         "high",
		CreatedAt:        time.Now(),
	}

	data := map[string]interface{}{
		"amount":        anomaly.Amount,
		"category_name": anomaly.CategoryName,
		"anomaly_score": anomaly.AnomalyScore,
	}

	return s.SendNotificationMessage(userID, notification, data)
}

// Helper methods for budget alerts
func (s *TelegramService) getBudgetAlertTitle(alertType string) string {
	switch alertType {
	case "threshold_reached":
		return "Ngân sách đạt ngưỡng cảnh báo"
	case "budget_exceeded":
		return "Ngân sách đã vượt quá"
	case "budget_achievement":
		return "Hoàn thành tiết kiệm ngân sách"
	default:
		return "Cập nhật ngân sách"
	}
}

func (s *TelegramService) getBudgetAlertMessage(budget *models.Budget, alertType string) string {
	switch alertType {
	case "threshold_reached":
		return fmt.Sprintf("Ngân sách '%s' đã đạt %.1f%% ngưỡng cảnh báo", budget.Name, budget.UsagePercentage)
	case "budget_exceeded":
		return fmt.Sprintf("Ngân sách '%s' đã vượt quá %.1f%%", budget.Name, budget.UsagePercentage)
	case "budget_achievement":
		return fmt.Sprintf("Chúc mừng! Bạn đã hoàn thành tiết kiệm ngân sách '%s'", budget.Name)
	default:
		return fmt.Sprintf("Cập nhật về ngân sách '%s'", budget.Name)
	}
}

func (s *TelegramService) getBudgetAlertPriority(alertType string) string {
	switch alertType {
	case "budget_exceeded":
		return "urgent"
	case "threshold_reached":
		return "high"
	default:
		return "medium"
	}
}

// Helper methods for goal alerts
func (s *TelegramService) getGoalAlertTitle(alertType string) string {
	switch alertType {
	case "progress_update":
		return "Cập nhật tiến độ mục tiêu"
	case "deadline_warning":
		return "Cảnh báo hạn chót mục tiêu"
	case "goal_achieved":
		return "Chúc mừng hoàn thành mục tiêu!"
	case "behind_schedule":
		return "Mục tiêu chậm tiến độ"
	default:
		return "Cập nhật mục tiêu"
	}
}

func (s *TelegramService) getGoalAlertMessage(goal *models.FinancialGoal, alertType string) string {
	switch alertType {
	case "progress_update":
		return fmt.Sprintf("Mục tiêu '%s' đã đạt %.1f%%", goal.Title, goal.Progress)
	case "deadline_warning":
		return fmt.Sprintf("Mục tiêu '%s' sắp đến hạn", goal.Title)
	case "goal_achieved":
		return fmt.Sprintf("Chúc mừng! Bạn đã hoàn thành mục tiêu '%s'", goal.Title)
	case "behind_schedule":
		return fmt.Sprintf("Mục tiêu '%s' đang chậm tiến độ", goal.Title)
	default:
		return fmt.Sprintf("Cập nhật về mục tiêu '%s'", goal.Title)
	}
}

func (s *TelegramService) getGoalAlertType(alertType string) string {
	switch alertType {
	case "goal_achieved":
		return "success"
	case "deadline_warning", "behind_schedule":
		return "warning"
	default:
		return "info"
	}
}

func (s *TelegramService) getGoalAlertPriority(alertType string) string {
	switch alertType {
	case "deadline_warning", "behind_schedule":
		return "high"
	case "goal_achieved":
		return "medium"
	default:
		return "low"
	}
}

// SendMonthlyReport sends monthly financial report to Telegram
func (s *TelegramService) SendMonthlyReport(userID uint64, report *models.DashboardAnalytics) error {
	message := fmt.Sprintf(`📊 *BÁO CÁO THÁNG %s*

💰 Tổng thu nhập: *%.0f VND*
💸 Tổng chi tiêu: *%.0f VND*
📈 Chênh lệch: *%.0f VND*

🏥 Sức khỏe tài chính: *%s* (%.1f/100)

📂 Top danh mục chi tiêu:
`, report.Period, report.TotalIncome, report.TotalExpense, report.NetAmount,
		report.FinancialHealth.Level, report.FinancialHealth.Score)

	// Add top categories
	for i, category := range report.CategoryBreakdown {
		if i >= 5 { // Limit to top 5
			break
		}
		message += fmt.Sprintf("%d. %s: *%.0f VND* (%.1f%%)\n",
			i+1, category.CategoryName, category.Amount, category.Percentage)
	}

	message += fmt.Sprintf("\n🕐 %s", report.GeneratedAt.Format("02/01/2006 15:04"))

	notification := &models.Notification{
		Title:            "Báo cáo tài chính hàng tháng",
		Message:          message,
		NotificationType: "info",
		Priority:         "low",
		CreatedAt:        time.Now(),
	}

	return s.SendNotificationMessage(userID, notification, nil)
}

// SendLargeTransactionAlert sends alert for large transactions
func (s *TelegramService) SendLargeTransactionAlert(userID uint64, transaction *models.Transaction, threshold float64) error {
	notification := &models.Notification{
		Title:            "Giao dịch lớn được phát hiện",
		Message:          fmt.Sprintf("Giao dịch %.0f VND tại %s vượt quá ngưỡng %.0f VND", transaction.Amount, transaction.Category.Name, threshold),
		NotificationType: "warning",
		Priority:         "medium",
		CreatedAt:        time.Now(),
	}

	data := map[string]interface{}{
		"amount":        transaction.Amount,
		"category_name": transaction.Category.Name,
		"description":   transaction.Description,
	}

	return s.SendNotificationMessage(userID, notification, data)
}
