package middleware

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"

	"github.com/labstack/echo/v4"
)

// RateLimitConfig returns rate limiting middleware configuration
func RateLimitConfig(cfg *config.Config) echo.MiddlewareFunc {
	// List of paths to skip rate limiting
	skipPaths := map[string]bool{
		"/health": true,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip rate limiting for health checks
			if skipPaths[c.Path()] {
				return next(c)
			}

			// Get client IP - handle Docker network properly
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			// For Docker bridge networks, use X-Forwarded-For if available
			forwardedFor := c.Request().Header.Get("X-Forwarded-For")
			if forwardedFor != "" {
				// Take the first IP from X-Forwarded-For chain
				ips := splitAndTrim(forwardedFor, ",")
				if len(ips) > 0 {
					clientIP = ips[0]
				}
			}

			// Normalize IP (remove port if present)
			if idx := indexOf(clientIP, ":"); idx > 0 {
				clientIP = clientIP[:idx]
			}

			// Skip rate limiting for localhost/internal IPs in development
			if cfg.Environment == "development" && (clientIP == "127.0.0.1" || clientIP == "::1" || 
				clientIP == "192.168.65.1" || clientIP == "172.17.0.1" || clientIP == "172.18.0.1") {
				// Still track but with more lenient limits for development
				// Allow 10000 requests per 60 seconds for localhost in dev
				devLimit := 10000
				key := fmt.Sprintf("rate_limit:dev:%s", clientIP)
				ctx := context.Background()
				current, err := database.IncrementRateLimit(ctx, key, time.Duration(cfg.RateLimit.Window)*time.Second)
				if err != nil {
					return next(c)
				}
				if current > int64(devLimit) {
					return c.JSON(429, map[string]interface{}{
						"error": "Rate limit exceeded",
						"message": fmt.Sprintf("Too many requests. Limit: %d per %d seconds", 
							devLimit, cfg.RateLimit.Window),
						"retry_after": cfg.RateLimit.Window,
					})
				}
				return next(c)
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

// Helper functions
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
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
