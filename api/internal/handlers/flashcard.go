package handlers

import (
	"net/http"
	"strconv"
	"swipelearn-api/internal/models"
	"swipelearn-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FlashcardHandler struct {
	flashcardService *services.FlashcardService
}

func NewFlashcardHandler(fs *services.FlashcardService) *FlashcardHandler {
	return &FlashcardHandler{
		flashcardService: fs,
	}
}

// CreateFlashcard handles POST /api/v1/flashcards
func (h *FlashcardHandler) CreateFlashcard(c *gin.Context) {
	var req models.CreateFlashcardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	flashcard, err := h.flashcardService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create flashcard",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, flashcard)
}

// GetFlashcards handles GET /api/v1/flashcards
func (h *FlashcardHandler) GetFlashcards(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Parse query parameters into filters map
    filters := make(map[string]any)

    // Min difficulty filter
    if minDiffStr := c.Query("min_difficulty"); minDiffStr != "" {
        if minDiff, err := strconv.ParseFloat(minDiffStr, 64); err == nil {
            filters["min_difficulty"] = minDiff
        }
    }

    // Deck ID filter
    if deckIDStr := c.Query("deck_id"); deckIDStr != "" {
        if deckID, err := uuid.Parse(deckIDStr); err == nil {
            filters["deck_id"] = deckID
        }
    }

	flashcards, err := h.flashcardService.GetByUser(userID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve flashcards",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": flashcards,
		"count": len(flashcards),
		"filters": filters,
	})
}

// UpdateFlashcard handles PUT /api/v1/flashcards/:id
func (h *FlashcardHandler) UpdateFlashcard(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid flashcard ID",
        })
        return
    }

    var req models.UpdateFlashcardRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    flashcard, err := h.flashcardService.Update(id, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to update flashcard",
        })
        return
    }

    c.JSON(http.StatusOK, flashcard)
}

// DeleteFlashcard handles DELETE /api/v1/flashcards/:id
func (h *FlashcardHandler) DeleteFlashcard(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid flashcard ID",
        })
        return
    }

    err = h.flashcardService.Delete(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to delete flashcard",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Flashcard deleted successfully",
    })
}
