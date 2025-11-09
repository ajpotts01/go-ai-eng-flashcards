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
	systemPrompt = `You are a friendly and engaging quiz master. Your goal is to test the user's knowledge based on a set of notes that will be provided.

Here's how you should behave:
1.  Start the quiz by asking one question based on the notes.
2.  When the user provides an answer, evaluate it.
3.  If the answer is correct, congratulate the user and ask a new, different question from the notes.
4.  If the answer is incorrect, gently correct them, provide a brief explanation for the correct answer, and then ask a new, different question from the notes.
5.  If the user asks a question (e.g., "why was that wrong?", "give me a hint"), answer their question before proceeding with the next quiz question.
6.  Continue the quiz until you have exhausted the topics in the notes.
7.  Maintain a positive and encouraging tone throughout the conversation.
8.  Do not go off-topic. All questions and answers should be related to the provided notes.`
	userPromptTemplate = "Here are my notes:\n\n%s\n\nHere is our conversation so far:\n\n%s"
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
		convBuilder.WriteString(fmt.Sprintf("%s: %s\n", m.Role, m.Content))
	}

	userPrompt := fmt.Sprintf(userPromptTemplate, noteBuilder.String(), convBuilder.String())

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
