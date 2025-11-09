// Create a handler just like todoHandler.go but for notes

package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"go-ai-eng-flashcards/models"
	"go-ai-eng-flashcards/services"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	service *services.NoteService
	logger  *slog.Logger
}

func NewNoteHandler(service *services.NoteService, logger *slog.Logger) *NoteHandler {
	return &NoteHandler{service: service, logger: logger}
}

func (h *NoteHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notes", h.CreateNote).Methods("POST")
	router.HandleFunc("/notes", h.GetAllNotes).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}", h.GetNoteByID).Methods("GET")
	router.HandleFunc("/notes/{id:[0-9]+}", h.UpdateNote).Methods("PUT")
	router.HandleFunc("/notes/{id:[0-9]+}", h.DeleteNote).Methods("DELETE")
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to create a new note")
	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid JSON payload for CreateNote", slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	note, err := h.service.CreateNote(&req)
	if err != nil {
		h.logger.Error("Failed to create note", slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Note created successfully", slog.Any("note_id", note.ID))
	h.writeJSONResponse(w, http.StatusCreated, note)
}

func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to get all notes")
	notes, err := h.service.GetAllNotes()
	if err != nil {
		h.logger.Error("Failed to retrieve all notes", slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve notes")
		return
	}

	h.logger.Info("Successfully retrieved all notes", slog.Any("count", len(notes)))
	h.writeJSONResponse(w, http.StatusOK, notes)
}

func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	h.logger.Info("Received request to get note by ID", slog.String("note_id_str", idStr))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid note ID format", slog.String("note_id_str", idStr), slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := h.service.GetNoteByID(id)
	if err != nil {
		h.logger.Error("Failed to retrieve note by ID", slog.Any("note_id", id), slog.Any("error", err))
		if noteErrorContainsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve note")
		}
		return
	}

	h.logger.Info("Note retrieved successfully", slog.Any("note_id", note.ID))
	h.writeJSONResponse(w, http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	h.logger.Info("Received request to update note", slog.String("note_id_str", idStr))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid note ID format for UpdateNote", slog.String("note_id_str", idStr), slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid JSON payload for UpdateNote", slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	note, err := h.service.UpdateNote(id, &req)
	if err != nil {
		h.logger.Error("Failed to update note", slog.Any("note_id", id), slog.Any("error", err))
		if noteErrorContainsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.logger.Info("Note updated successfully", slog.Any("note_id", note.ID))
	h.writeJSONResponse(w, http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	h.logger.Info("Received request to delete note", slog.String("note_id_str", idStr))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid note ID format for DeleteNote", slog.String("note_id_str", idStr), slog.Any("error", err))
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	err = h.service.DeleteNote(id)
	if err != nil {
		h.logger.Error("Failed to delete note", slog.Any("note_id", id), slog.Any("error", err))
		if noteErrorContainsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete note")
		}
		return
	}

	h.logger.Info("Note deleted successfully", slog.Any("note_id", id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to write JSON response", slog.Any("error", err))
	}
}

func (h *NoteHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		h.logger.Error("Failed to write error response", slog.Any("error", err))
	}
}

func noteErrorContainsNotFound(message string) bool {
	return len(message) > 0 && (message[len(message)-9:] == "not found" ||
		message[:len("note with id")] == "note with id")
}
