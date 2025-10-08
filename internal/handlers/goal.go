package handlers

import (
	"net/http"
	"strconv"

	"tabimoney/internal/models"
	"tabimoney/internal/services"

	"github.com/labstack/echo/v4"
)

type GoalHandler struct {
	goalService *services.GoalService
}

func NewGoalHandler() *GoalHandler {
	return &GoalHandler{
		goalService: services.NewGoalService(),
	}
}

// GetGoals retrieves user's financial goals
func (h *GoalHandler) GetGoals(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	goals, err := h.goalService.GetGoals(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get goals",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": goals,
	})
}

// GetGoal retrieves a specific goal
func (h *GoalHandler) GetGoal(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid goal ID",
			Message: "Goal ID must be a valid number",
		})
	}

	goals, err := h.goalService.GetGoals(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get goals",
			Message: err.Error(),
		})
	}

	// Find the specific goal
	var goal *models.FinancialGoal
	for i := range goals {
		if goals[i].ID == goalID {
			goal = &goals[i]
			break
		}
	}

	if goal == nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Goal not found",
			Message: "The requested goal does not exist",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": goal,
	})
}

// CreateGoal creates a new financial goal
func (h *GoalHandler) CreateGoal(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	var req models.FinancialGoalCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	goal, err := h.goalService.CreateGoal(userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create goal",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": goal,
	})
}

// UpdateGoal updates an existing goal
func (h *GoalHandler) UpdateGoal(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid goal ID",
			Message: "Goal ID must be a valid number",
		})
	}

	var req models.FinancialGoalUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	goal, err := h.goalService.UpdateGoal(userID, goalID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update goal",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": goal,
	})
}

// DeleteGoal deletes a goal
func (h *GoalHandler) DeleteGoal(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid goal ID",
			Message: "Goal ID must be a valid number",
		})
	}

	if err := h.goalService.DeleteGoal(userID, goalID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete goal",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Goal deleted successfully",
	})
}

// AddContribution adds money to a goal
func (h *GoalHandler) AddContribution(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	goalID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid goal ID",
			Message: "Goal ID must be a valid number",
		})
	}

	var req struct {
		Amount float64 `json:"amount"`
		Note   string  `json:"note"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
	}

	goal, err := h.goalService.AddContribution(userID, goalID, req.Amount, req.Note)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to add contribution",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": goal,
	})
}
