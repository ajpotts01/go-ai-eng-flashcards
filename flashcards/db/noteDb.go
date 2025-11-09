package db

import (
	"database/sql"
	"fmt"
	"go-ai-eng-flashcards/models"
	"log/slog"

	_ "github.com/lib/pq"
)

type NoteRepository interface {
	CreateNote(note *models.Note) error
	GetNoteById(id int64) (*models.Note, error)
	GetAllNotes() ([]*models.Note, error)
	UpdateNote(id int64, updates map[string]any) error
	DeleteNote(id int64) error
	Close() error
}

type PostgresNoteRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresNoteRepository(dbUrl string, logger *slog.Logger) (*PostgresNoteRepository, error) {
	logger.Info("Attempting to open database connection")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		logger.Error("Failed to open database", slog.Any("error", err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	logger.Info("Pinging database to verify connection")
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", slog.Any("error", err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully")
	return &PostgresNoteRepository{db: db, logger: logger}, nil
}

func (r *PostgresNoteRepository) CreateNote(note *models.Note) error {
	r.logger.Info("Attempting to create a new note", slog.Any("note_content", note.Content))
	query := `
	INSERT INTO
		flashcards.notes (content)
	VALUES ($1)
	RETURNING id, created_at, updated_at
	`

	row := r.db.QueryRow(query, note.Content)
	err := row.Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to create note", slog.Any("error", err))
		return fmt.Errorf("failed to create note: %w", err)
	}

	r.logger.Info("Note created successfully", slog.Any("note_id", note.ID))
	return nil
}

func (r *PostgresNoteRepository) GetNoteById(id int64) (*models.Note, error) {
	r.logger.Info("Attempting to retrieve note by ID", slog.Any("note_id", id))
	query := `
	SELECT 
		id, content, created_at, updated_at
	FROM
	    flashcards.notes
	WHERE
	    id = $1
	`

	note := &models.Note{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(&note.ID, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Note not found", slog.Any("note_id", id))
			return nil, fmt.Errorf("note with id %d not found", id)
		}
		r.logger.Error("Failed to get note by ID", slog.Any("note_id", id), slog.Any("error", err))
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	r.logger.Info("Note retrieved successfully", slog.Any("note_id", note.ID))
	return note, nil
}

func (r *PostgresNoteRepository) GetAllNotes() ([]*models.Note, error) {
	r.logger.Info("Attempting to retrieve all notes")
	query := `
	SELECT
		id, content, created_at, updated_at
	FROM
	    flashcards.notes
	ORDER BY
	    created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to get all notes", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get all notes: %w", err)
	}
	defer rows.Close()

	notes := make([]*models.Note, 0)
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(&note.ID, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			r.logger.Error("Failed to scan note", slog.Any("error", err))
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Failed to iterate notes", slog.Any("error", err))
		return nil, fmt.Errorf("failed to iterate notes: %w", err)
	}

	r.logger.Info("All notes retrieved successfully", slog.Any("count", len(notes)))
	return notes, nil
}

func (r *PostgresNoteRepository) UpdateNote(id int64, updates map[string]any) error {
	r.logger.Info("Attempting to update note", slog.Any("note_id", id), slog.Any("updates", updates))
	if len(updates) == 0 {
		r.logger.Warn("No updates provided for note", slog.Any("note_id", id))
		return fmt.Errorf("no updates provided")
	}

	query := "UPDATE flashcards.notes SET "
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
		r.logger.Error("Failed to update note", slog.Any("note_id", id), slog.Any("error", err))
		return fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected after update", slog.Any("note_id", id), slog.Any("error", err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Warn("No rows updated for note", slog.Any("note_id", id))
		return fmt.Errorf("no rows updated - note with id %d not found", id)
	}

	r.logger.Info("Note updated successfully", slog.Any("note_id", id))
	return nil
}

func (r *PostgresNoteRepository) DeleteNote(id int64) error {
	r.logger.Info("Attempting to delete note", slog.Any("note_id", id))
	query := "DELETE FROM flashcards.notes WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("Failed to delete note", slog.Any("note_id", id), slog.Any("error", err))
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected after delete", slog.Any("note_id", id), slog.Any("error", err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Warn("No rows deleted for note", slog.Any("note_id", id))
		return fmt.Errorf("no rows deleted - note with id %d not found", id)
	}

	r.logger.Info("Note deleted successfully", slog.Any("note_id", id))
	return nil
}

func (r *PostgresNoteRepository) Close() error {
	r.logger.Info("Closing database connection")
	err := r.db.Close()
	if err != nil {
		r.logger.Error("Failed to close database connection", slog.Any("error", err))
		return fmt.Errorf("failed to close database: %w", err)
	}
	r.logger.Info("Database connection closed successfully")
	return nil
}
