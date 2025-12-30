package handlers

import (
	"net/http"
	"swipelearn-api/internal/models"
	"swipelearn-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeckHandler struct {
	deckService *services.DeckService
}

func NewDeckHandler(ds *services.DeckService) *DeckHandler {
	return &DeckHandler{
		deckService: ds,
	}
}

// CreateDeck handles POST /api/v1/decks
func (h *DeckHandler) CreateDeck(c *gin.Context) {
	var req models.CreateDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	deck, err := h.deckService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create deck",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, deck)
}

// GetDecks handles GET /api/v1/decks
func (h *DeckHandler) GetDecks(c *gin.Context) {
	decks, err := h.deckService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve decks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  decks,
		"count": len(decks),
	})
}

// GetDeck handles GET /api/v1/decks/:id
func (h *DeckHandler) GetDeck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID",
		})
		return
	}

	deck, err := h.deckService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Deck not found",
		})
		return
	}

	c.JSON(http.StatusOK, deck)
}

// UpdateDeck handles PUT /api/v1/decks/:id
func (h *DeckHandler) UpdateDeck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID",
		})
		return
	}

	var req models.UpdateDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	deck, err := h.deckService.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update deck",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deck)
}

// DeleteDeck handles DELETE /api/v1/decks/:id
func (h *DeckHandler) DeleteDeck(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID",
		})
		return
	}

	err = h.deckService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete deck",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Deck deleted successfully",
	})
}
