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

// GetBudget retrieves a specific budget
func (h *BudgetHandler) GetBudget(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid budget ID",
			Message: "Budget ID must be a valid number",
		})
	}

	budgets, err := h.budgetService.GetBudgets(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get budgets",
			Message: err.Error(),
		})
	}

	// Find the specific budget
	var budget *models.Budget
	for i := range budgets {
		if budgets[i].ID == budgetID {
			budget = &budgets[i]
			break
		}
	}

	if budget == nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Budget not found",
			Message: "The requested budget does not exist",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": budget,
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

// GetBudgetAlerts returns budgets that are approaching or exceeding their limits
func (h *BudgetHandler) GetBudgetAlerts(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	alerts, err := h.budgetService.GetBudgetAlerts(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get budget alerts",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": alerts,
	})
}
