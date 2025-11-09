package services

import (
	"fmt"
	"go-ai-eng-flashcards/db"
	"go-ai-eng-flashcards/models"
	"log/slog"
	"strings"
)

type NoteService struct {
	repo   db.NoteRepository
	logger *slog.Logger
}

func NewNoteService(repo db.NoteRepository, logger *slog.Logger) *NoteService {
	return &NoteService{repo: repo, logger: logger}
}

func (s *NoteService) CreateNote(req *models.CreateNoteRequest) (*models.Note, error) {
	s.logger.Info("Attempting to create a new note", slog.Any("content", req.Content))
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	note := &models.Note{
		Content: strings.TrimSpace(req.Content),
	}

	if err := s.repo.CreateNote(note); err != nil {
		return nil, err
	}

	s.logger.Info("Note created successfully", slog.Any("note_id", note.ID))
	return note, nil
}

func (s *NoteService) GetNoteByID(id int64) (*models.Note, error) {
	s.logger.Info("Attempting to retrieve note by ID", slog.Any("note_id", id))
	if id <= 0 {
		return nil, fmt.Errorf("invalid note ID: %d", id)
	}

	note, err := s.repo.GetNoteById(id)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Note retrieved successfully", slog.Any("note_id", note.ID))
	return note, nil
}

func (s *NoteService) GetAllNotes() ([]*models.Note, error) {
	s.logger.Info("Attempting to retrieve all notes")
	notes, err := s.repo.GetAllNotes()
	if err != nil {
		return nil, err
	}

	s.logger.Info("All notes retrieved successfully", slog.Any("count", len(notes)))
	return notes, nil
}

func (s *NoteService) UpdateNote(id int64, req *models.UpdateNoteRequest) (*models.Note, error) {
	s.logger.Info("Attempting to update note", slog.Any("note_id", id), slog.Any("updates", req))
	if id <= 0 {
		return nil, fmt.Errorf("invalid note ID: %d", id)
	}

	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if req.Content != nil {
		trimmedContent := strings.TrimSpace(*req.Content)
		if trimmedContent == "" {
			return nil, fmt.Errorf("content cannot be empty")
		}
		updates["content"] = trimmedContent
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no valid updates provided")
	}

	if err := s.repo.UpdateNote(id, updates); err != nil {
		return nil, err
	}

	s.logger.Info("Note updated successfully", slog.Any("note_id", id))
	return s.repo.GetNoteById(id)
}

func (s *NoteService) DeleteNote(id int64) error {
	s.logger.Info("Attempting to delete note", slog.Any("note_id", id))
	if id <= 0 {
		return fmt.Errorf("invalid note ID: %d", id)
	}

	err := s.repo.DeleteNote(id)
	if err != nil {
		return err
	}

	s.logger.Info("Note deleted successfully", slog.Any("note_id", id))
	return nil
}

func (s *NoteService) validateCreateRequest(req *models.CreateNoteRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return fmt.Errorf("content is required")
	}

	if len(content) > 255 {
		return fmt.Errorf("content cannot exceed 255 characters")
	}

	return nil
}

func (s *NoteService) validateUpdateRequest(req *models.UpdateNoteRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.Content == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	if req.Content != nil {
		content := strings.TrimSpace(*req.Content)
		if len(content) > 255 {
			return fmt.Errorf("content cannot exceed 255 characters")
		}
	}

	return nil
}
