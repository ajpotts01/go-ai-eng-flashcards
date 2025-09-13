package db

import (
	"database/sql"
	"fmt"
	"go-ai-eng-flashcards/models"

	_ "github.com/lib/pq"
)

type FlashcardRepository interface {
	CreateFlashcard(flashcard *models.Flashcard) error
	GetFlashcardById(id int64) (*models.Flashcard, error)
	GetAllFlashcards() ([]*models.Flashcard, error)
	UpdateFlashcard(id int64, updates map[string]any) error
	DeleteFlashcard(id int64) error
}

type PostgresFlashcardRepository struct {
	db *sql.DB
}

func NewPostgresFlashcardRepository(dbUrl string) (*PostgresFlashcardRepository, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresFlashcardRepository{db: db}, nil
}

func (r *PostgresFlashcardRepository) CreateFlashcard(flashcard *models.Flashcard) error {
	query := `
	INSERT INTO
		flashcards.flashcards (content)
	VALUES ($1)
	RETURNING id, created_at, updated_at
	`

	row := r.db.QueryRow(query, flashcard.Content)
	err := row.Scan(&flashcard.ID, &flashcard.CreatedAt, &flashcard.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create flashcard: %w", err)
	}

	return nil
}

func (r *PostgresFlashcardRepository) GetFlashcardById(id int64) (*models.Flashcard, error) {
	query := `
	SELECT 
		id, content, created_at, updated_at
	FROM
	    flashcards.flashcards
	WHERE
	    id = $1
	`

	flashcard := &models.Flashcard{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(&flashcard.ID, &flashcard.Content, &flashcard.CreatedAt, &flashcard.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flashcard with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get flashcard: %w", err)
	}

	return flashcard, nil
}

func (r *PostgresFlashcardRepository) GetAllFlashcards() ([]*models.Flashcard, error) {
	query := `
	SELECT
		id, content, created_at, updated_at
	FROM
	    flashcards.flashcards
	ORDER BY
	    created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all flashcards: %w", err)
	}
	defer rows.Close()

	flashcards := make([]*models.Flashcard, 0)
	for rows.Next() {
		flashcard := &models.Flashcard{}
		err := rows.Scan(&flashcard.ID, &flashcard.Content, &flashcard.CreatedAt, &flashcard.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan flashcard: %w", err)
		}
		flashcards = append(flashcards, flashcard)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate flashcards: %w", err)
	}

	return flashcards, nil
}

func (r *PostgresFlashcardRepository) UpdateFlashcard(id int64, updates map[string]any) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	query := "UPDATE flashcards.flashcards SET "
	args := []any{}
	argIndex := 1

	for field, value := range updates {
		if argIndex > 1 {
			query += ","
		}
		query += fmt.Sprintf("%s = $%d", field, argIndex)
		args = append(args, value)
		argIndex++
	}

	query += fmt.Sprintf(", updated_at = NOW() WHERE id = $%d", argIndex)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update flashcard: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated - flashcard with id %d not found", id)
	}

	return nil
}

func (r *PostgresFlashcardRepository) DeleteFlashcard(id int64) error {
	query := "DELETE FROM flashcards.flashcards WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete flashcard: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted - flashcard with id %d not found", id)
	}

	return nil
}

func (r *PostgresFlashcardRepository) Close() error {
	return r.db.Close()
}
