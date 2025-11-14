package handlers

import (
	"net/http"
	"strconv"

	"tabimoney/internal/models"
	"tabimoney/internal/services"

	"github.com/labstack/echo/v4"
)

type BudgetHandler struct {
	budgetService *services.BudgetService
}

func NewBudgetHandler() *BudgetHandler {
	return &BudgetHandler{
		budgetService: services.NewBudgetService(),
	}
}

// GetBudgets retrieves user's budgets
func (h *BudgetHandler) GetBudgets(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	budgets, err := h.budgetService.GetBudgets(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get budgets",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": budgets,
	})
}

// CreateBudget creates a new budget
func (h *BudgetHandler) CreateBudget(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	var req models.BudgetCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	budget, err := h.budgetService.CreateBudget(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create budget",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": budget,
	})
}

// UpdateBudget updates an existing budget
func (h *BudgetHandler) UpdateBudget(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid budget ID",
			Message: "Budget ID must be a valid number",
		})
	}

	var req models.BudgetUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	budget, err := h.budgetService.UpdateBudget(userID, budgetID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update budget",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": budget,
	})
}

// DeleteBudget deletes a budget
func (h *BudgetHandler) DeleteBudget(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid budget ID",
			Message: "Budget ID must be a valid number",
		})
	}

	if err := h.budgetService.DeleteBudget(userID, budgetID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete budget",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Budget deleted successfully",
	})
}

// GetBudgetInsights returns safe-to-spend and pacing info
func (h *BudgetHandler) GetBudgetInsights(c echo.Context) error {
    userID := c.Get("user_id").(uint64)

    insights, err := h.budgetService.GetBudgetInsights(userID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error:   "Failed to get budget insights",
            Message: err.Error(),
        })
    }
    return c.JSON(http.StatusOK, map[string]interface{}{
        "data": insights,
    })
}

// GetAutoBudgetSuggestions suggests budgets for current period
func (h *BudgetHandler) GetAutoBudgetSuggestions(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    resp, err := h.budgetService.SuggestBudgets(userID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error:   "Failed to suggest budgets",
            Message: err.Error(),
        })
    }
    return c.JSON(http.StatusOK, map[string]interface{}{
        "data": resp,
    })
}

// CreateBudgetsFromSuggestions bulk creates budgets from suggestions
func (h *BudgetHandler) CreateBudgetsFromSuggestions(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    var req models.AutoBudgetCreateRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:   "Invalid request",
            Message: err.Error(),
        })
    }
    if req.AlertThreshold == 0 {
        req.AlertThreshold = 80
    }
    created, err := h.budgetService.CreateBudgetsFromSuggestions(userID, &req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error:   "Failed to create budgets",
            Message: err.Error(),
        })
    }
    return c.JSON(http.StatusCreated, map[string]interface{}{
        "data": created,
    })
}
