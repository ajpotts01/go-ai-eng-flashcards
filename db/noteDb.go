package db

import (
	"database/sql"
	"fmt"
	"go-ai-eng-flashcards/models"

	_ "github.com/lib/pq"
)

type NoteRepository interface {
	CreateNote(note *models.Note) error
	GetNoteById(id int64) (*models.Note, error)
	GetAllNotes() ([]*models.Note, error)
	UpdateNote(id int64, updates map[string]any) error
	DeleteNote(id int64) error
}

type PostgresNoteRepository struct {
	db *sql.DB
}

func NewPostgresNoteRepository(dbUrl string) (*PostgresNoteRepository, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresNoteRepository{db: db}, nil
}

func (r *PostgresNoteRepository) CreateNote(note *models.Note) error {
	query := `
	INSERT INTO
		notes.notes (content)
	VALUES ($1)
	RETURNING id, created_at, updated_at
	`

	row := r.db.QueryRow(query, note.Content)
	err := row.Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	return nil
}

func (r *PostgresNoteRepository) GetNoteById(id int64) (*models.Note, error) {
	query := `
	SELECT 
		id, content, created_at, updated_at
	FROM
	    notes.notes
	WHERE
	    id = $1
	`

	note := &models.Note{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(&note.ID, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return note, nil
}

func (r *PostgresNoteRepository) GetAllNotes() ([]*models.Note, error) {
	query := `
	SELECT
		id, content, created_at, updated_at
	FROM
	    notes.notes
	ORDER BY
	    created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all notes: %w", err)
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(&note.ID, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notes: %w", err)
	}

	return notes, nil
}

func (r *PostgresNoteRepository) UpdateNote(id int64, updates map[string]any) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	query := "UPDATE notes.notes SET "
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
		return fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated - note with id %d not found", id)
	}

	return nil
}

func (r *PostgresNoteRepository) DeleteNote(id int64) error {
	query := "DELETE FROM notes.notes WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted - note with id %d not found", id)
	}

	return nil
}

func (r *PostgresNoteRepository) Close() error {
	return r.db.Close()
}