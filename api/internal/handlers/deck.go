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

	// Get user_id from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID type in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	deck, err := h.deckService.Create(&req, userID)
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
	// Get user_id from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID type in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	decks, err := h.deckService.GetByUser(userID)
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

	// Get user_id from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID type in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	deck, err := h.deckService.GetByIDWithOwnership(id, userID)
	if err != nil {
		if err.Error() == "unauthorized: deck does not belong to user" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to access this deck",
			})
			return
		}
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

	// Get user_id from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID type in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
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

	deck, err := h.deckService.UpdateWithOwnership(id, userID, &req)
	if err != nil {
		if err.Error() == "unauthorized: deck does not belong to user" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to update this deck",
			})
			return
		}
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

	// Get user_id from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID type in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.deckService.DeleteWithOwnership(id, userID)
	if err != nil {
		if err.Error() == "unauthorized: deck does not belong to user" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You are not authorized to delete this deck",
			})
			return
		}
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
