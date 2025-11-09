package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"go-ai-eng-flashcards/config"
	"go-ai-eng-flashcards/db"
	"go-ai-eng-flashcards/handlers"
	"go-ai-eng-flashcards/services"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()
	logger := config.NewLogger()

	if cfg.DatabaseURL == "" {
		logger.Error("DB_URL environment variable is required")
		return
	}

	todoRepo, err := db.NewPostgresTodoRepository(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to initialize database", slog.Any("error", err))
		return
	}
	defer todoRepo.Close()

	noteRepo, err := db.NewPostgresNoteRepository(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to initialize note database", slog.Any("error", err))
		return
	}
	defer noteRepo.Close()

	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	noteService := services.NewNoteService(noteRepo, logger)
	noteHandler := handlers.NewNoteHandler(noteService, logger)

	quizService, err := services.NewQuizService(cfg.GeminiAPIKey, noteService, logger)
	if err != nil {
		logger.Error("Failed to initialize quiz service", slog.Any("error", err))
		return
	}
	quizHandler := handlers.NewQuizHandler(quizService, logger)

	router := mux.NewRouter()

	router.Use(corsMiddleware)
	router.Use(jsonMiddleware)

	todoHandler.RegisterRoutes(router)
	noteHandler.RegisterRoutes(router)
	quizHandler.RegisterRoutes(router)

	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	addr := ":" + cfg.Port
	fmt.Printf("Server starting on port %s\n", cfg.Port)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("Server failed to start", slog.Any("error", err))
		return
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}
