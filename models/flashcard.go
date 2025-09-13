package models

import "time"

type Flashcard struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateFlashcardRequest struct {
	Content string `json:"content"`
}

type UpdateFlashcardRequest struct {
	Content *string `json:"content,omitempty"` // Why did tutorial's Claude Code use a pointer?
}
