package routes

import (
	"swipelearn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	// Authentication routes (public)
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)    // POST /api/v1/auth/register
		auth.POST("/login", authHandler.Login)          // POST /api/v1/auth/login
		auth.POST("/refresh", authHandler.RefreshToken) // POST /api/v1/auth/refresh
		auth.POST("/logout", authHandler.Logout)        // POST /api/v1/auth/logout
	}
}
