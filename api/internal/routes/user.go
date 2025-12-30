package routes

import (
	"swipelearn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(apiGroup *gin.RouterGroup, userHandler *handlers.UserHandler) {
	// User routes under /api/v1/users
	users := apiGroup.Group("/users")
	{
		users.POST("", userHandler.CreateUser)                    // POST /api/v1/users
		users.GET("", userHandler.GetUsers)                       // GET /api/v1/users
		users.GET("/:id", userHandler.GetUser)                    // GET /api/v1/users/:id
		users.PUT("/:id", userHandler.UpdateUser)                 // PUT /api/v1/users/:id
		users.DELETE("/:id", userHandler.DeleteUser)              // DELETE /api/v1/users/:id
		users.GET("/by-email/:email", userHandler.GetUserByEmail) // GET /api/v1/users/by-email/:email
	}
}
