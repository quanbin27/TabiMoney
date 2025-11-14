package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// TelegramLinkCodeRequest represents the request payload for generating link code
type TelegramLinkCodeRequest struct {
	TelegramUserID int64  `json:"telegram_user_id" validate:"required"`
	LinkCode       string `json:"link_code" validate:"required"`
}

// TelegramLinkCodeResponse represents the response payload for link code
type TelegramLinkCodeResponse struct {
	Success        bool   `json:"success"`
	LinkCode       string `json:"link_code"`
	ExpiryMinutes  int    `json:"expiry_minutes"`
	Message        string `json:"message"`
}

// TelegramStatusResponse represents the response payload for Telegram status
type TelegramStatusResponse struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}

// TelegramDisconnectResponse represents the response payload for disconnect
type TelegramDisconnectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GenerateTelegramLinkCode godoc
// @Summary Generate Telegram link code
// @Description Generate a link code for Telegram bot integration
// @Tags telegram
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} TelegramLinkCodeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/telegram/generate-link-code [post]
func (h *AuthHandler) GenerateTelegramLinkCode(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Generate link code
	linkCode, err := h.authService.GenerateTelegramLinkCode(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to generate link code",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, TelegramLinkCodeResponse{
		Success:       true,
		LinkCode:      linkCode,
		ExpiryMinutes: 10, // 10 minutes expiry
		Message:       "Link code generated successfully",
	})
}

// GetTelegramStatus godoc
// @Summary Get Telegram integration status
// @Description Check if user's account is linked with Telegram
// @Tags telegram
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} TelegramStatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/telegram/status [get]
func (h *AuthHandler) GetTelegramStatus(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Check if user has Telegram account linked
	linked, err := h.authService.IsTelegramLinked(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to check Telegram status",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, TelegramStatusResponse{
		Connected: linked,
		Message:   "Status retrieved successfully",
	})
}

// DisconnectTelegram godoc
// @Summary Disconnect Telegram account
// @Description Unlink user's Telegram account from the system
// @Tags telegram
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} TelegramDisconnectResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/telegram/disconnect [post]
func (h *AuthHandler) DisconnectTelegram(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Disconnect Telegram account
	err := h.authService.DisconnectTelegram(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to disconnect Telegram",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, TelegramDisconnectResponse{
		Success: true,
		Message: "Telegram account disconnected successfully",
	})
}

// LinkTelegramAccount godoc
// @Summary Link Telegram account with link code
// @Description Link Telegram account using the generated link code
// @Tags telegram
// @Accept json
// @Produce json
// @Param request body TelegramLinkCodeRequest true "Link code data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/telegram/link [post]
func (h *AuthHandler) LinkTelegramAccount(c echo.Context) error {
	var req TelegramLinkCodeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Validate link code and get web user ID
	webUserID, err := h.authService.ValidateTelegramLinkCode(req.LinkCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid link code",
			Message: err.Error(),
		})
	}

	// Link Telegram account
	err = h.authService.LinkTelegramAccount(req.TelegramUserID, webUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to link Telegram account",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Telegram account linked successfully",
	})
}

