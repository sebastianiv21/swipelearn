package routes

import (
	"swipelearn-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupDeckRoutes(apiGroup *gin.RouterGroup, deckHandler *handlers.DeckHandler) {
	// Deck routes under /api/v1/decks
	decks := apiGroup.Group("/decks")
	{
		decks.POST("", deckHandler.CreateDeck)       // POST /api/v1/decks
		decks.GET("", deckHandler.GetDecks)          // GET /api/v1/decks
		decks.GET("/:id", deckHandler.GetDeck)       // GET /api/v1/decks/:id
		decks.PUT("/:id", deckHandler.UpdateDeck)    // PUT /api/v1/decks/:id
		decks.DELETE("/:id", deckHandler.DeleteDeck) // DELETE /api/v1/decks/:id
	}
}
