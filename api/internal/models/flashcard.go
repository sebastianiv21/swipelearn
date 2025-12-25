package models

import (
	"time"
)

type Flashcard struct {
	ID          int64      `json:"id" db:"id"`
	Front       string     `json:"front" db:"front"`
	Back        string     `json:"back" db:"back"`
	DeckID      int64      `json:"deck_id" db:"deck_id"`
	Difficulty  float64    `json:"difficulty" db:"difficulty"`
	Interval    int        `json:"interval" db:"interval"`
	EaseFactor  float64    `json:"ease_factor" db:"ease_factor"`
	ReviewCount int        `json:"review_count" db:"review_count"`
	LastReview  *time.Time `json:"last_review" db:"last_review"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type Deck struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateFlashcardRequest struct {
	Front  string `json:"front" binding:"required"`
	Back   string `json:"back" binding:"required"`
	DeckID int64  `json:"deck_id" binding:"required"`
}

type UpdateFlashcardRequest struct {
	Front      *string  `json:"front"`
	Back       *string  `json:"back"`
	Difficulty *float64 `json:"difficulty"`
}

type CreateDeckRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
