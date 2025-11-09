package services

import (
	"context"
	"fmt"
	"go-ai-eng-flashcards/models"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

const (
	systemPrompt       = "You are an expert quiz master. A user will provide you with a series of notes. Your job is to generate a single, concise question based on these notes. The question should test the user's knowledge of the provided information. Do not ask for the notes, just generate the question from the notes provided in the prompt."
	userPromptTemplate = "Here are my notes:\n\n%s\n\nPlease generate a quiz question based on these notes."
)

// QuizService handles the business logic for quiz generation.
type QuizService struct {
	llm         *googleai.GoogleAI
	noteService *NoteService
}

// NewQuizService creates a new instance of QuizService.
func NewQuizService(apiKey string, noteService *NoteService) (*QuizService, error) {
	llm, err := googleai.New(context.Background(), googleai.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini LLM: %w", err)
	}
	return &QuizService{llm: llm, noteService: noteService}, nil
}

// GenerateQuizTurn adds a new, LLM-generated assistant message to a conversation history.
func (s *QuizService) GenerateQuizTurn(currentMessages []models.Message) []models.Message {
	allNotes, err := s.noteService.GetAllNotes()
	if err != nil {
		log.Printf("Error fetching notes for quiz generation: %v", err)
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

	userPrompt := fmt.Sprintf(userPromptTemplate, noteBuilder.String())

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	ctx := context.Background()
	completion, err := s.llm.GenerateContent(ctx, messages, llms.WithTemperature(0.8))
	if err != nil {
		log.Printf("Error generating content from LLM: %v", err)
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

	return append(currentMessages, assistantMessage)
}
