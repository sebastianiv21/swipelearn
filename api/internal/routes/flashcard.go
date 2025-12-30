package routes

import (
	"swipelearn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupFlashcardRoutes(apiGroup *gin.RouterGroup, flashcardHandler *handlers.FlashcardHandler) {
	// Flashcard routes under /api/v1/flashcards
	apiGroup.GET("", flashcardHandler.GetFlashcards)
	apiGroup.POST("", flashcardHandler.CreateFlashcard)
	apiGroup.PUT("/:id", flashcardHandler.UpdateFlashcard)
	apiGroup.DELETE("/:id", flashcardHandler.DeleteFlashcard)
	apiGroup.POST("/:id/review", flashcardHandler.ReviewFlashcard) // POST /api/v1/flashcards/:id/review
	apiGroup.GET("/due", flashcardHandler.GetDueFlashcards)        // GET /api/v1/flashcards/due
}
