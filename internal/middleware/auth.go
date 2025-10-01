package middleware

import (
	"strings"

	"tabimoney/internal/services"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]string{
					"error": "Authorization header required",
				})
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(401, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate token
			userID, err := authService.ValidateToken(token)
			if err != nil {
				return c.JSON(401, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Set user context
			c.Set("user_id", userID)
			c.Set("token", token)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware validates JWT token if present but doesn't require it
func OptionalAuthMiddleware(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if userID, err := authService.ValidateToken(token); err == nil {
					c.Set("user_id", userID)
					c.Set("token", token)
				}
			}

			return next(c)
		}
	}
}
