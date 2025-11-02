package services

import (
	"encoding/json"
	"fmt"
	"log"

	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type NotificationPreferencesService struct {
	db *gorm.DB
}

type NotificationPreferences struct {
	// Channel preferences
	EmailEnabled     bool `json:"email_enabled"`
	TelegramEnabled  bool `json:"telegram_enabled"`
	InAppEnabled     bool `json:"in_app_enabled"`
	PushEnabled      bool `json:"push_enabled"`

	// Feature preferences
	BudgetAlerts      bool `json:"budget_alerts"`
	GoalAlerts        bool `json:"goal_alerts"`
	AIAlerts          bool `json:"ai_alerts"`
	TransactionAlerts bool `json:"transaction_alerts"`
	AnalyticsAlerts   bool `json:"analytics_alerts"`

	// Priority preferences
	UrgentNotifications bool `json:"urgent_notifications"`
	HighNotifications   bool `json:"high_notifications"`
	MediumNotifications bool `json:"medium_notifications"`
	LowNotifications    bool `json:"low_notifications"`

	// Frequency preferences
	DailyDigest     bool `json:"daily_digest"`
	WeeklyDigest    bool `json:"weekly_digest"`
	MonthlyDigest   bool `json:"monthly_digest"`
	RealTimeAlerts  bool `json:"real_time_alerts"`

	// Time preferences
	QuietHoursStart string `json:"quiet_hours_start"` // HH:MM format
	QuietHoursEnd   string `json:"quiet_hours_end"`   // HH:MM format
	Timezone        string `json:"timezone"`
}

func NewNotificationPreferencesService() *NotificationPreferencesService {
	return &NotificationPreferencesService{
		db: database.GetDB(),
	}
}

// GetUserPreferences gets user's notification preferences
func (s *NotificationPreferencesService) GetUserPreferences(userID uint64) (*NotificationPreferences, error) {
	var user models.User
	if err := s.db.Preload("Profile").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Default preferences
	preferences := &NotificationPreferences{
		EmailEnabled:         true,
		TelegramEnabled:      true,
		InAppEnabled:         true,
		PushEnabled:          false,
		BudgetAlerts:         true,
		GoalAlerts:           true,
		AIAlerts:             true,
		TransactionAlerts:    true,
		AnalyticsAlerts:      true,
		UrgentNotifications:  true,
		HighNotifications:    true,
		MediumNotifications:  true,
		LowNotifications:     false,
		DailyDigest:          false,
		WeeklyDigest:         true,
		MonthlyDigest:        true,
		RealTimeAlerts:       true,
		QuietHoursStart:      "22:00",
		QuietHoursEnd:        "08:00",
		Timezone:            "Asia/Ho_Chi_Minh",
	}

	// Load user's saved preferences
	if user.Profile != nil && user.Profile.NotificationSettings != "" {
		if err := json.Unmarshal([]byte(user.Profile.NotificationSettings), preferences); err != nil {
			log.Printf("Failed to unmarshal notification preferences for user %d: %v", userID, err)
		}
	}

	return preferences, nil
}

// UpdateUserPreferences updates user's notification preferences
func (s *NotificationPreferencesService) UpdateUserPreferences(userID uint64, preferences *NotificationPreferences) error {
	// Validate preferences
	if err := s.validatePreferences(preferences); err != nil {
		return fmt.Errorf("invalid preferences: %w", err)
	}

	// Convert to JSON
	preferencesJSON, err := json.Marshal(preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	// Update user profile
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new profile
			profile = models.UserProfile{
				UserID:               userID,
				NotificationSettings: string(preferencesJSON),
			}
			if err := s.db.Create(&profile).Error; err != nil {
				return fmt.Errorf("failed to create user profile: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get user profile: %w", err)
		}
	} else {
		// Update existing profile
		profile.NotificationSettings = string(preferencesJSON)
		if err := s.db.Save(&profile).Error; err != nil {
			return fmt.Errorf("failed to update user profile: %w", err)
		}
	}

	return nil
}

// validatePreferences validates notification preferences
func (s *NotificationPreferencesService) validatePreferences(preferences *NotificationPreferences) error {
	// Validate quiet hours format
	if preferences.QuietHoursStart != "" {
		if !s.isValidTimeFormat(preferences.QuietHoursStart) {
			return fmt.Errorf("invalid quiet hours start format, expected HH:MM")
		}
	}
	if preferences.QuietHoursEnd != "" {
		if !s.isValidTimeFormat(preferences.QuietHoursEnd) {
			return fmt.Errorf("invalid quiet hours end format, expected HH:MM")
		}
	}

	// Validate timezone
	if preferences.Timezone == "" {
		preferences.Timezone = "Asia/Ho_Chi_Minh"
	}

	return nil
}

// isValidTimeFormat checks if time string is in HH:MM format
func (s *NotificationPreferencesService) isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 {
		return false
	}
	if timeStr[2] != ':' {
		return false
	}
	hour := timeStr[:2]
	minute := timeStr[3:]
	
	// Check if hour and minute are valid numbers
	if hour < "00" || hour > "23" {
		return false
	}
	if minute < "00" || minute > "59" {
		return false
	}
	
	return true
}

// ShouldSendNotification checks if notification should be sent based on preferences
func (s *NotificationPreferencesService) ShouldSendNotification(userID uint64, notificationType, priority string) (bool, error) {
	preferences, err := s.GetUserPreferences(userID)
	if err != nil {
		return false, err
	}

	// Check priority preferences
	switch priority {
	case "urgent":
		if !preferences.UrgentNotifications {
			return false, nil
		}
	case "high":
		if !preferences.HighNotifications {
			return false, nil
		}
	case "medium":
		if !preferences.MediumNotifications {
			return false, nil
		}
	case "low":
		if !preferences.LowNotifications {
			return false, nil
		}
	}

	// Check feature preferences
	switch notificationType {
	case "warning":
		// Budget and goal warnings
		return preferences.BudgetAlerts || preferences.GoalAlerts, nil
	case "info":
		// Analytics and general info
		return preferences.AnalyticsAlerts, nil
	case "success":
		// Goal achievements
		return preferences.GoalAlerts, nil
	case "reminder":
		// Budget reminders
		return preferences.BudgetAlerts, nil
	default:
		return true, nil
	}
}

// GetEnabledChannels returns list of enabled notification channels
func (s *NotificationPreferencesService) GetEnabledChannels(userID uint64) ([]string, error) {
	preferences, err := s.GetUserPreferences(userID)
	if err != nil {
		return nil, err
	}

	var channels []string
	if preferences.EmailEnabled {
		channels = append(channels, "email")
	}
	if preferences.TelegramEnabled {
		channels = append(channels, "telegram")
	}
	if preferences.InAppEnabled {
		channels = append(channels, "in_app")
	}
	if preferences.PushEnabled {
		channels = append(channels, "push")
	}

	return channels, nil
}

// IsQuietHours checks if current time is within quiet hours
func (s *NotificationPreferencesService) IsQuietHours(userID uint64) (bool, error) {
	preferences, err := s.GetUserPreferences(userID)
	if err != nil {
		return false, err
	}

	if preferences.QuietHoursStart == "" || preferences.QuietHoursEnd == "" {
		return false, nil
	}

	// TODO: Implement timezone-aware quiet hours checking
	// For now, return false (no quiet hours)
	return false, nil
}

// GetDefaultPreferences returns default notification preferences
func (s *NotificationPreferencesService) GetDefaultPreferences() *NotificationPreferences {
	return &NotificationPreferences{
		EmailEnabled:         true,
		TelegramEnabled:      true,
		InAppEnabled:         true,
		PushEnabled:          false,
		BudgetAlerts:         true,
		GoalAlerts:           true,
		AIAlerts:             true,
		TransactionAlerts:    true,
		AnalyticsAlerts:      true,
		UrgentNotifications:  true,
		HighNotifications:    true,
		MediumNotifications:  true,
		LowNotifications:     false,
		DailyDigest:          false,
		WeeklyDigest:         true,
		MonthlyDigest:        true,
		RealTimeAlerts:       true,
		QuietHoursStart:      "22:00",
		QuietHoursEnd:        "08:00",
		Timezone:            "Asia/Ho_Chi_Minh",
	}
}

// ResetToDefaults resets user's notification preferences to defaults
func (s *NotificationPreferencesService) ResetToDefaults(userID uint64) error {
	defaultPrefs := s.GetDefaultPreferences()
	return s.UpdateUserPreferences(userID, defaultPrefs)
}

// GetPreferencesSummary returns a summary of user's notification preferences
func (s *NotificationPreferencesService) GetPreferencesSummary(userID uint64) (map[string]interface{}, error) {
	preferences, err := s.GetUserPreferences(userID)
	if err != nil {
		return nil, err
	}

	enabledChannels, err := s.GetEnabledChannels(userID)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"enabled_channels": enabledChannels,
		"budget_alerts":    preferences.BudgetAlerts,
		"goal_alerts":      preferences.GoalAlerts,
		"ai_alerts":        preferences.AIAlerts,
		"transaction_alerts": preferences.TransactionAlerts,
		"analytics_alerts": preferences.AnalyticsAlerts,
		"urgent_enabled":   preferences.UrgentNotifications,
		"high_enabled":     preferences.HighNotifications,
		"medium_enabled":   preferences.MediumNotifications,
		"low_enabled":      preferences.LowNotifications,
		"daily_digest":     preferences.DailyDigest,
		"weekly_digest":    preferences.WeeklyDigest,
		"monthly_digest":   preferences.MonthlyDigest,
		"real_time_alerts": preferences.RealTimeAlerts,
		"quiet_hours":      preferences.QuietHoursStart != "" && preferences.QuietHoursEnd != "",
		"timezone":         preferences.Timezone,
	}

	return summary, nil
}

