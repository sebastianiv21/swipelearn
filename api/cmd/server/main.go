package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

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
	logger := logrus.New()

	// Configure Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	logger.SetLevel(logrus.InfoLevel)

	// Log server version info on startup
	logger.WithFields(logrus.Fields{
		"version":    Version,
		"git_commit": GitCommit,
		"build_time": BuildTime,
		"go_version": GoVersion,
	}).Info("Starting SwipeLearn API Server")

	// Create Gin router
	router := gin.New()

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.WithField("port", port).Info("PORT not set, using default")
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server exited")
}
