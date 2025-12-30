package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SimpleAuth is a simple authentication middleware that checks for X-User-ID header
// This is a basic implementation for demonstration purposes
// In production, you should use proper JWT or OAuth authentication
func SimpleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			// For development, allow user_id as query parameter
			userIDStr = c.Query("user_id")
		}

		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID is required. Provide it in X-User-ID header or user_id query parameter",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid User ID format",
			})
			c.Abort()
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("user_id", userID)
		c.Next()
	}
}
