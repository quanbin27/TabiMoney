package handlers

import (
	"net/http"

	"tabimoney/internal/models"
	"tabimoney/internal/services"

	"github.com/labstack/echo/v4"
)

type AIHandler struct {
	svc *services.AIService
}

func NewAIHandler(svc *services.AIService) *AIHandler {
	return &AIHandler{svc: svc}
}

func (h *AIHandler) SuggestCategory(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	var req models.CategorySuggestionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
	}
	req.UserID = userID
	resp, err := h.svc.SuggestCategory(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Suggestion failed", Message: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
