package services

import (
	"fmt"
	"strings"

	"go-ai-eng-flashcards/db"
	"go-ai-eng-flashcards/models"
)

type FlashcardService struct {
	repo db.FlashcardRepository
}

func NewFlashcardService(repo db.FlashcardRepository) *FlashcardService {
	return &FlashcardService{repo: repo}
}

func (s *FlashcardService) CreateFlashcard(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	todo := &models.Flashcard{
		Content: strings.TrimSpace(req.Content),
	}

	if err := s.repo.CreateFlashcard(todo); err != nil {
		return nil, fmt.Errorf("failed to create flashcard: %w", err)
	}

	return todo, nil
}

func (s *FlashcardService) GetFlashcardByID(id int64) (*models.Flashcard, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid flashcard ID: %d", id)
	}

	todo, err := s.repo.GetFlashcardById(id)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *FlashcardService) GetAllFlashcards() ([]*models.Flashcard, error) {
	todos, err := s.repo.GetAllFlashcards()
	if err != nil {
		return nil, fmt.Errorf("failed to get flashcards: %w", err)
	}

	return todos, nil
}

func (s *FlashcardService) UpdateFlashcard(id int64, req *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid todo ID: %d", id)
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
		updates["title"] = trimmedContent
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no valid updates provided")
	}

	if err := s.repo.UpdateFlashcard(id, updates); err != nil {
		return nil, err
	}

	return s.repo.GetFlashcardById(id)
}

func (s *FlashcardService) DeleteFlashcard(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid todo ID: %d", id)
	}

	return s.repo.DeleteFlashcard(id)
}

func (s *FlashcardService) validateCreateRequest(req *models.CreateFlashcardRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	title := strings.TrimSpace(req.Content)
	if title == "" {
		return fmt.Errorf("content is required")
	}

	if len(title) > 255 {
		return fmt.Errorf("content cannot exceed 255 characters")
	}

	return nil
}

func (s *FlashcardService) validateUpdateRequest(req *models.UpdateFlashcardRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.Content == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	if req.Content != nil {
		title := strings.TrimSpace(*req.Content)
		if len(title) > 255 {
			return fmt.Errorf("content cannot exceed 255 characters")
		}
	}

	return nil
}
