package routes

import (
	"swipelearn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupFlashcardRoutes(apiGroup *gin.RouterGroup, flashcardHandler *handlers.FlashcardHandler) {
	// Flashcard routes under /api/v1/flashcards
	flashcards := apiGroup.Group("/flashcards")
	{
		flashcards.GET("", flashcardHandler.GetFlashcards)               // GET /api/v1/flashcards
		flashcards.POST("", flashcardHandler.CreateFlashcard)            // POST /api/v1/flashcards
		flashcards.PUT("/:id", flashcardHandler.UpdateFlashcard)         // PUT /api/v1/flashcards/:id
		flashcards.DELETE("/:id", flashcardHandler.DeleteFlashcard)      // DELETE /api/v1/flashcards/:id
		flashcards.POST("/:id/review", flashcardHandler.ReviewFlashcard) // POST /api/v1/flashcards/:id/review
		flashcards.GET("/due", flashcardHandler.GetDueFlashcards)        // GET /api/v1/flashcards/due
	}
}
