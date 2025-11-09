package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-ai-eng-flashcards/models"
	"go-ai-eng-flashcards/services"
	"log/slog"
	"net/http"
)

// quizRequest is the expected structure of the request body for the /quiz endpoint.
// It is defined locally within the handler package.
type quizRequest struct {
	Messages []models.Message `json:"messages"`
}

// quizResponse is the structure of the response body for the /quiz endpoint.
// It is defined locally within the handler package.
type quizResponse struct {
	Messages []models.Message `json:"messages"`
}

// QuizHandler manages HTTP requests for the /quiz endpoint.
type QuizHandler struct {
	service *services.QuizService
	logger  *slog.Logger
}

// NewQuizHandler creates a new instance of QuizHandler.
func NewQuizHandler(service *services.QuizService, logger *slog.Logger) *QuizHandler {
	return &QuizHandler{service: service, logger: logger}
}

// GenerateQuizHandler handles a single request to generate a quiz turn.
func (h *QuizHandler) GenerateQuizHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to generate a quiz turn")
	var req quizRequest
	// Decode the incoming JSON payload into the local request struct.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body for GenerateQuizHandler", slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Call the service to get the updated message list.
	updatedMessages := h.service.GenerateQuizTurn(req.Messages)

	// Prepare the response using the local response struct.
	res := quizResponse{
		Messages: updatedMessages,
	}

	h.logger.Info("Quiz turn generated successfully")
	h.writeJSONResponse(w, http.StatusOK, res)
}

func (h *QuizHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/quiz", h.GenerateQuizHandler).Methods("POST")
}

func (h *QuizHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to write JSON response", slog.Any("error", err))
	}
}

func (h *QuizHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		h.logger.Error("Failed to write error response", slog.Any("error", err))
	}
}
