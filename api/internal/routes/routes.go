package routes

import (
	"swipelearn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	flashcardHandler *handlers.FlashcardHandler,
	// deckHandler *handlers.DeckHandler,
	// userHandler *handlers.UserHandler,
	// healthHandler *handlers.HealthHandler,
) *gin.Engine {
	router := gin.New()

	// Middleware
	// router.Use(middleware.CORS())

	// Health routes (no auth needed)
	// router.GET("/health", healthHandler.Health)
	// router.GET("/ready", healthHandler.Ready)
	// router.GET("/version", healthHandler.Version)

	// API routes group (with middleware)
	apiGroup := router.Group("/api/v1")
	// apiGroup.Use(middleware.Auth()) // Apply auth to all API routes

	// Setup route groups
	SetupFlashcardRoutes(apiGroup, flashcardHandler)
	// SetupDeckRoutes(apiGroup, deckHandler)
	// SetupUserRoutes(apiGroup, userHandler)

	return router
}
