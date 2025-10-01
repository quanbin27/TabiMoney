package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"

	"github.com/labstack/echo/v4"
)

// RateLimitConfig returns rate limiting middleware configuration
func RateLimitConfig(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get client IP
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			// Create rate limit key
			key := fmt.Sprintf("rate_limit:%s", clientIP)

			ctx := context.Background()

			// Check current rate limit
			current, err := database.IncrementRateLimit(ctx, key, time.Duration(cfg.RateLimit.Window)*time.Second)
			if err != nil {
				// If Redis is down, allow the request but log the error
				fmt.Printf("Rate limit check failed: %v\n", err)
				return next(c)
			}

			// Check if limit exceeded
			if current > int64(cfg.RateLimit.Requests) {
				return c.JSON(429, map[string]interface{}{
					"error": "Rate limit exceeded",
					"message": fmt.Sprintf("Too many requests. Limit: %d per %d seconds", 
						cfg.RateLimit.Requests, cfg.RateLimit.Window),
					"retry_after": cfg.RateLimit.Window,
				})
			}

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(cfg.RateLimit.Requests))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.FormatInt(int64(cfg.RateLimit.Requests)-current, 10))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Duration(cfg.RateLimit.Window)*time.Second).Unix(), 10))

			return next(c)
		}
	}
}

// UserRateLimitConfig returns user-specific rate limiting middleware
func UserRateLimitConfig(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user ID from context (set by auth middleware)
			userID, ok := c.Get("user_id").(uint64)
			if !ok {
				// If no user ID, use IP-based rate limiting
				return RateLimitConfig(cfg)(next)(c)
			}

			// Create user-specific rate limit key
			key := fmt.Sprintf("user_rate_limit:%d", userID)

			ctx := context.Background()

			// Check current rate limit
			current, err := database.IncrementRateLimit(ctx, key, time.Duration(cfg.RateLimit.Window)*time.Second)
			if err != nil {
				// If Redis is down, allow the request but log the error
				fmt.Printf("User rate limit check failed: %v\n", err)
				return next(c)
			}

			// Check if limit exceeded
			if current > int64(cfg.RateLimit.Requests) {
				return c.JSON(429, map[string]interface{}{
					"error": "Rate limit exceeded",
					"message": fmt.Sprintf("Too many requests. Limit: %d per %d seconds", 
						cfg.RateLimit.Requests, cfg.RateLimit.Window),
					"retry_after": cfg.RateLimit.Window,
				})
			}

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(cfg.RateLimit.Requests))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.FormatInt(int64(cfg.RateLimit.Requests)-current, 10))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Duration(cfg.RateLimit.Window)*time.Second).Unix(), 10))

			return next(c)
		}
	}
}
