package models

import (
	"time"

	"github.com/google/uuid"
)

type Flashcard struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Front       string     `json:"front" db:"front"`
	Back        string     `json:"back" db:"back"`
	DeckID      uuid.UUID  `json:"deck_id" db:"deck_id"`
	Difficulty  float64    `json:"difficulty" db:"difficulty"`
	Interval    int        `json:"interval" db:"interval"`
	EaseFactor  float64    `json:"ease_factor" db:"ease_factor"`
	ReviewCount int        `json:"review_count" db:"review_count"`
	LastReview  *time.Time `json:"last_review" db:"last_review"`
	NextReview  *time.Time `json:"next_review" db:"next_review"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateFlashcardRequest struct {
	Front  string    `json:"front" binding:"required"`
	Back   string    `json:"back" binding:"required"`
	UserID uuid.UUID `json:"user_id" binding:"required"`
	DeckID uuid.UUID `json:"deck_id" binding:"required"`
}

type UpdateFlashcardRequest struct {
	Front      *string  `json:"front"`
	Back       *string  `json:"back"`
	Difficulty *float64 `json:"difficulty"`
	Interval   *int     `json:"interval"`
}

type ReviewFlashcardRequest struct {
	Quality int `json:"quality" binding:"required,min=0,max=5"`
}
