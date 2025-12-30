package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"swipelearn-api/internal/db"
	"swipelearn-api/internal/handlers"
	"swipelearn-api/internal/repositories"
	"swipelearn-api/internal/routes"
	"swipelearn-api/internal/services"
	"swipelearn-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Version information - will be set at build time
var (
	Version   = "dev"     // Set during build: -ldflags "-X main.Version=1.0.0"
	GitCommit = "unknown" // Set during build: -ldflags "-X main.GitCommit=abc123"
	BuildTime = "unknown" // Set during build: -ldflags "-X main.BuildTime=2024-12-24T10:30:00Z"
	GoVersion = runtime.Version()
)

func main() {
	// Initialize structured logger
	logger := utils.SetupLogger()

	// Initialize database
	database, err := db.NewDatabase(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}
	defer database.Close()

	// Initialize layers (Dependency Injection)
	flashcardRepo := repositories.NewFlashcardRepository(database.DB, logger)
	flashcardService := services.NewFlashcardService(flashcardRepo, logger)
	flashcardHandler := handlers.NewFlashcardHandler(flashcardService)

	userRepo := repositories.NewUserRepository(database.DB, logger)
	userService := services.NewUserService(userRepo, logger)
	userHandler := handlers.NewUserHandler(userService)

	deckRepo := repositories.NewDeckRepository(database.DB, logger)
	deckService := services.NewDeckService(deckRepo, logger)
	deckHandler := handlers.NewDeckHandler(deckService)

	// JWT and Auth services
	jwtService := services.NewJWTService(logger)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(database.DB, logger)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, jwtService, logger)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup routes
	router := routes.SetupRouter(
		flashcardHandler,
		deckHandler,
		userHandler,
		authHandler,
		jwtService,
	)

	// Setup auth routes (public)
	routes.SetupAuthRoutes(router, authHandler)

	// Add logging middleware
	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"status_code": params.StatusCode,
			"latency":     params.Latency,
			"client_ip":   params.ClientIP,
			"method":      params.Method,
			"path":        params.Path,
			"user_agent":  params.Request.UserAgent(),
			"error":       params.ErrorMessage,
		}).Info("HTTP Request")

		return ""
	}))

	// Add recovery middleware with custom logging
	router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		logger.WithFields(logrus.Fields{
			"error":     err,
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"client_ip": c.ClientIP(),
		}).Error("Panic recovered")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Internal Server Error",
			"status": "error",
		})
	}))

	// Version endpoint - returns detailed version info
	router.GET("/version", func(c *gin.Context) {
		logger.Info("Version endpoint requested")
		c.JSON(http.StatusOK, gin.H{
			"version":    Version,
			"git_commit": GitCommit,
			"build_time": BuildTime,
			"go_version": GoVersion,
			"service":    "swipelearn-api",
			"timestamp":  time.Now().UTC(),
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		logger.Info("Health check requested")
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	})

	// Ready endpoint
	router.GET("/ready", func(c *gin.Context) {
		logger.Info("Ready check requested")
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"time":   time.Now().UTC(),
		})
	})

	// Configure server with proper timeouts
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "5050"
		logger.WithField("port", port).Info("PORT not set, using default")
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       utils.GetEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      utils.GetEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		ReadHeaderTimeout: utils.GetEnvAsDuration("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
		IdleTimeout:       utils.GetEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
	}

	// Start server in a goroutine
	go func() {
		logger.WithField("port", port).Info("Starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), utils.GetEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second))
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server exited")
}
