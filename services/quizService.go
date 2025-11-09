package services

import (
	"context"
	"fmt"
	"go-ai-eng-flashcards/models"
	"log/slog"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

const (
	systemPrompt       = "You are an expert quiz master. A user will provide you with a series of notes. Your job is to generate a single, concise question based on these notes. The question should test the user's knowledge of the provided information. Do not ask for the notes, just generate the question from the notes provided in the prompt."
	initialQuizPrompt  = "Generate a quiz question based on these study notes:\n\n%s"
	conversationPrompt = "Continue this quiz conversation with a follow-up question based on the notes and conversation history.\n\nNotes:\n%s\n\nConversation:\n%s\n\n"
	userPromptTemplate = "Here are my notes:\n\n%s\n\nPlease generate a quiz question based on these notes."
)

// QuizService handles the business logic for quiz generation.
type QuizService struct {
	llm         *googleai.GoogleAI
	noteService *NoteService
	logger      *slog.Logger
}

// NewQuizService creates a new instance of QuizService.
func NewQuizService(apiKey string, noteService *NoteService, logger *slog.Logger) (*QuizService, error) {
	logger.Info("Initializing QuizService")
	llm, err := googleai.New(context.Background(), googleai.WithAPIKey(apiKey))
	if err != nil {
		logger.Error("Failed to initialize Gemini LLM", slog.Any("error", err))
		return nil, fmt.Errorf("failed to initialize Gemini LLM: %w", err)
	}
	logger.Info("QuizService initialized successfully")
	return &QuizService{llm: llm, noteService: noteService, logger: logger}, nil
}

// GenerateQuizTurn adds a new, LLM-generated assistant message to a conversation history.
func (s *QuizService) GenerateQuizTurn(currentMessages []models.Message) []models.Message {
	s.logger.Info("Generating quiz turn")
	allNotes, err := s.noteService.GetAllNotes()
	if err != nil {
		s.logger.Error("Error fetching notes for quiz generation", slog.Any("error", err))
		assistantMessage := models.Message{
			Role:    "assistant",
			Content: "Sorry, I was unable to fetch the notes to generate a question.",
		}
		return append(currentMessages, assistantMessage)
	}

	var noteBuilder strings.Builder
	for _, note := range allNotes {
		noteBuilder.WriteString(note.Content)
		noteBuilder.WriteString("\n")
	}

	var convBuilder strings.Builder
	for _, m := range currentMessages {
		convBuilder.WriteString(m.Content)
		convBuilder.WriteString("\n")
	}

	var userPrompt string
	if len(currentMessages) == 0 {
		s.logger.Info("Generating initial quiz question")
		userPrompt = fmt.Sprintf(initialQuizPrompt, noteBuilder.String())
	} else {
		s.logger.Info("Generating follow-up quiz question for existing conversation")
		userPrompt = fmt.Sprintf(conversationPrompt, noteBuilder.String(), convBuilder.String())
	}

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	ctx := context.Background()
	completion, err := s.llm.GenerateContent(ctx, messages, llms.WithTemperature(0.8))
	if err != nil {
		s.logger.Error("Error generating content from LLM", slog.Any("error", err))
		// Fallback to a generic error message
		assistantMessage := models.Message{
			Role:    "assistant",
			Content: "Sorry, I was unable to generate a question at this time.",
		}
		return append(currentMessages, assistantMessage)
	}

	generatedContent := "Sorry, I couldn't generate a question."
	if len(completion.Choices) > 0 && len(completion.Choices[0].Content) > 0 {
		generatedContent = completion.Choices[0].Content
	}

	assistantMessage := models.Message{
		Role:    "assistant",
		Content: generatedContent,
	}

	s.logger.Info("Quiz turn generated successfully")
	return append(currentMessages, assistantMessage)
}
