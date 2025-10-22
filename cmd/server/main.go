package main

import (
    "context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"
	"tabimoney/internal/handlers"
	appmw "tabimoney/internal/middleware"
	"tabimoney/internal/services"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Setup logging
	setupLogging(cfg)

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	// Run database migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// Initialize Redis
	if err := database.InitRedis(cfg); err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}
	defer database.CloseRedis()

	// Initialize services
	authService := services.NewAuthService(cfg)
    // Initialize optional services later
    txHandler := handlers.NewTransactionHandler()
    categoryHandler := handlers.NewCategoryHandler()
    aiService := services.NewAIService(cfg)
    aiHandler := handlers.NewAIHandler(aiService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup Echo server
	e := echo.New()

	// Middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())
	e.Use(appmw.CORSConfig(cfg))
	e.Use(appmw.RateLimitConfig(cfg))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		// Check database health
		if err := database.HealthCheck(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "unhealthy",
				"error":  err.Error(),
			})
		}

		// Check Redis health
		if err := database.RedisHealthCheck(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "unhealthy",
				"error":  err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := e.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
    auth.POST("/logout", authHandler.Logout, appmw.AuthMiddleware(authService))
    auth.POST("/change-password", authHandler.ChangePassword, appmw.AuthMiddleware(authService))
    auth.GET("/profile", authHandler.GetProfile, appmw.AuthMiddleware(authService))
    auth.PUT("/profile", authHandler.UpdateProfile, appmw.AuthMiddleware(authService))
    auth.GET("/income", authHandler.GetMonthlyIncome, appmw.AuthMiddleware(authService))
    auth.PUT("/income", authHandler.SetMonthlyIncome, appmw.AuthMiddleware(authService))
    
    // Telegram integration routes
    auth.POST("/telegram/generate-link-code", authHandler.GenerateTelegramLinkCode, appmw.AuthMiddleware(authService))
    auth.GET("/telegram/status", authHandler.GetTelegramStatus, appmw.AuthMiddleware(authService))
    auth.POST("/telegram/disconnect", authHandler.DisconnectTelegram, appmw.AuthMiddleware(authService))
    auth.POST("/telegram/link", authHandler.LinkTelegramAccount)

    // Transactions routes
    tx := api.Group("/transactions", appmw.AuthMiddleware(authService))
    tx.GET("", txHandler.List)
    tx.POST("", txHandler.Create)
    tx.GET("/:id", txHandler.Get)
    tx.PUT("/:id", txHandler.Update)
    tx.DELETE("/:id", txHandler.Delete)

    // Categories
    cat := api.Group("/categories", appmw.AuthMiddleware(authService))
    cat.GET("", categoryHandler.List)
    cat.POST("", categoryHandler.Create)
    cat.GET("/:id", categoryHandler.Get)
    cat.PUT("/:id", categoryHandler.Update)
    cat.DELETE("/:id", categoryHandler.Delete)

    // Goals routes
    goalHandler := handlers.NewGoalHandler()
    goals := api.Group("/goals", appmw.AuthMiddleware(authService))
    goals.GET("", goalHandler.GetGoals)
    goals.GET("/:id", goalHandler.GetGoal)
    goals.POST("", goalHandler.CreateGoal)
    goals.PUT("/:id", goalHandler.UpdateGoal)
    goals.DELETE("/:id", goalHandler.DeleteGoal)
    goals.POST("/:id/contribute", goalHandler.AddContribution)

    // Budgets routes
    budgetHandler := handlers.NewBudgetHandler()
    // Notifications routes
    notificationHandler := handlers.NewNotificationHandler()
    notifications := api.Group("/notifications", appmw.AuthMiddleware(authService))
    notifications.GET("", notificationHandler.List)
    notifications.POST(":id/read", notificationHandler.MarkRead)
    budgets := api.Group("/budgets", appmw.AuthMiddleware(authService))
    budgets.GET("", budgetHandler.GetBudgets)
    budgets.GET("/:id", budgetHandler.GetBudget)
    budgets.POST("", budgetHandler.CreateBudget)
    budgets.PUT("/:id", budgetHandler.UpdateBudget)
    budgets.DELETE("/:id", budgetHandler.DeleteBudget)
    budgets.GET("/alerts", budgetHandler.GetBudgetAlerts)

    // AI endpoints
    ai := api.Group("/ai", appmw.AuthMiddleware(authService))
    ai.POST("/suggest-category", aiHandler.SuggestCategory)
    ai.POST("/chat", aiHandler.ProcessChat)

    // Analytics routes
    analyticsHandler := handlers.NewAnalyticsHandler(cfg)
    analytics := api.Group("/analytics", appmw.AuthMiddleware(authService))
    analytics.GET("/dashboard", analyticsHandler.GetDashboardAnalytics)
    analytics.GET("/category-spending", analyticsHandler.GetCategorySpending)
    analytics.GET("/spending-patterns", analyticsHandler.GetSpendingPatterns)
    analytics.GET("/anomalies", analyticsHandler.GetAnomalies)
    analytics.GET("/predictions", analyticsHandler.GetPredictions)

	// Start server
	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      e,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Server starting on %s", cfg.GetServerAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logrus.Info("Server shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}

	logrus.Info("Server exited")
}

func setupLogging(cfg *config.Config) {
	// Set log level
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Set log format
	if cfg.Logging.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}
