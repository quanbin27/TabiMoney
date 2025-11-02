package handlers

import (
	"net/http"

	"tabimoney/internal/services"
	"github.com/labstack/echo/v4"
)

type NotificationPreferencesHandler struct {
	svc *services.NotificationPreferencesService
}

func NewNotificationPreferencesHandler() *NotificationPreferencesHandler {
	return &NotificationPreferencesHandler{
		svc: services.NewNotificationPreferencesService(),
	}
}

// GetPreferences gets user's notification preferences
func (h *NotificationPreferencesHandler) GetPreferences(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	
	preferences, err := h.svc.GetUserPreferences(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get preferences",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": preferences,
	})
}

// UpdatePreferences updates user's notification preferences
func (h *NotificationPreferencesHandler) UpdatePreferences(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	
	var preferences services.NotificationPreferences
	if err := c.Bind(&preferences); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := h.svc.UpdateUserPreferences(userID, &preferences); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update preferences",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"message": "Preferences updated successfully",
	})
}

// GetSummary gets user's notification preferences summary
func (h *NotificationPreferencesHandler) GetSummary(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	
	summary, err := h.svc.GetPreferencesSummary(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get preferences summary",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": summary,
	})
}

// ResetToDefaults resets user's notification preferences to defaults
func (h *NotificationPreferencesHandler) ResetToDefaults(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	
	if err := h.svc.ResetToDefaults(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to reset preferences",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"message": "Preferences reset to defaults",
	})
}

// GetEnabledChannels gets user's enabled notification channels
func (h *NotificationPreferencesHandler) GetEnabledChannels(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	
	channels, err := h.svc.GetEnabledChannels(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get enabled channels",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": channels,
	})
}

// TestNotification sends a test notification to user
func (h *NotificationPreferencesHandler) TestNotification(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	
	// Get channel from query parameter
	channel := c.QueryParam("channel")
	if channel == "" {
		channel = "in_app"
	}

	// Create test notification
	dispatcher := services.NewNotificationDispatcher()
	
	// Get frontend URL from config or use default
	frontendURL := c.Request().Header.Get("Origin")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	
	trigger := services.NotificationTrigger{
		UserID:          userID,
		NotificationType: "info",
		Priority:        "low",
		Title:           "🔔 Thông báo Test",
		Message:         "Đây là thông báo test để kiểm tra cài đặt thông báo của bạn có hoạt động đúng không.",
		ActionURL:       frontendURL + "/notifications",
		Metadata: map[string]interface{}{
			"test": true,
		},
	}

	if err := dispatcher.DispatchNotification(trigger); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to send test notification",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"message": "Test notification sent successfully",
	})
}

