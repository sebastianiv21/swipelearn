package routes

import (
	"swipelearn-api/internal/handlers"
	"swipelearn-api/internal/middleware"
	"swipelearn-api/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	flashcardHandler *handlers.FlashcardHandler,
	deckHandler *handlers.DeckHandler,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	jwtService *services.JWTService,
) *gin.Engine {
	router := gin.New()

	// Middleware
	// router.Use(middleware.CORS())

	// Setup public auth routes (no JWT middleware required)
	SetupAuthRoutes(router, authHandler)

	// API routes group (with middleware)
	apiGroup := router.Group("/api/v1")
	apiGroup.Use(middleware.JWTAuth(jwtService)) // Apply JWT auth to all API routes

	// Setup route groups
	SetupFlashcardRoutes(apiGroup, flashcardHandler)
	SetupDeckRoutes(apiGroup, deckHandler)
	SetupUserRoutes(apiGroup, userHandler)

	// Protected auth routes
	authGroup := apiGroup.Group("/auth")
	{
		authGroup.POST("/logout", authHandler.Logout) // POST /api/v1/auth/logout (protected)
	}

	return router
}
