package handlers

import (
    "net/http"

    "tabimoney/internal/models"
    "tabimoney/internal/services"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "User registration data"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req models.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		return c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req models.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Login failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token refresh failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	token := c.Get("token").(string)

	if err := h.authService.Logout(userID, token); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Logout failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Logged out successfully",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Password change data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	if err := h.authService.ChangePassword(userID, &req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Password change failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password changed successfully",
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	// Get user from database
    var user models.User
    if err := services.DB().Preload("Profile").First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user profile",
			Message: err.Error(),
		})
	}

    response := services.UserToResponse(&user)
	return c.JSON(http.StatusOK, response)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserUpdateRequest true "Profile update data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	var req models.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	// Update user profile
	var user models.User
    if err := services.DB().First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to find user",
			Message: err.Error(),
		})
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.AvatarURL = req.AvatarURL

    if err := services.DB().Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update profile",
			Message: err.Error(),
		})
	}

    response := services.UserToResponse(&user)
	return c.JSON(http.StatusOK, response)
}

// Response types
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

// IncomeResponse payload
type IncomeResponse struct {
    MonthlyIncome float64 `json:"monthly_income"`
}

// GetMonthlyIncome godoc
// @Summary Get monthly income
// @Description Get user's configured monthly income
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} IncomeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/income [get]
func (h *AuthHandler) GetMonthlyIncome(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    amount, err := h.authService.GetMonthlyIncome(userID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get income", Message: err.Error()})
    }
    return c.JSON(http.StatusOK, IncomeResponse{MonthlyIncome: amount})
}

// SetMonthlyIncome godoc
// @Summary Set monthly income
// @Description Update user's monthly income
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body IncomeResponse true "Monthly income"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/income [put]
func (h *AuthHandler) SetMonthlyIncome(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    var req IncomeResponse
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body", Message: err.Error()})
    }
    if req.MonthlyIncome < 0 {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Validation failed", Message: "monthly_income must be >= 0"})
    }
    if err := h.authService.SetMonthlyIncome(userID, req.MonthlyIncome); err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to set income", Message: err.Error()})
    }
    return c.JSON(http.StatusOK, SuccessResponse{Message: "Monthly income updated"})
}
