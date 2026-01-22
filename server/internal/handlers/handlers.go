package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/julieta-311/where-are-you/server/internal/models"
)

type Handler struct {
	store  models.Store
	logger *slog.Logger
}

func NewHandler(store models.Store, logger *slog.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
	}
}

func (h *Handler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RequesterID string `json:"requester_id"`
		TargetID    string `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	req := models.LocationRequest{
		ID:          uuid.New().String(),
		RequesterID: input.RequesterID,
		TargetID:    input.TargetID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := h.store.CreateRequest(r.Context(), req); err != nil {
		h.logger.Error("failed to create request", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(req); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *Handler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	reqs, err := h.store.GetPendingRequests(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get pending requests", "error", err, "user_id", userID)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reqs); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *Handler) RespondToRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["requestId"]

	var input struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if input.Status != "accepted" && input.Status != "denied" {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateRequestStatus(r.Context(), id, input.Status); err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			h.logger.Error("failed to update request", "error", err, "id", id)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	var loc models.LocationUpdate
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	loc.Timestamp = time.Now()

	if err := h.store.UpdateLocation(r.Context(), loc); err != nil {
		h.logger.Error("failed to update location", "error", err, "user_id", loc.UserID)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetID := vars["userId"]
	requesterID := r.URL.Query().Get("requester_id")

	allowed, err := h.store.CheckAuthorization(r.Context(), requesterID, targetID)
	if err != nil {
		h.logger.Error("auth check failed", "error", err, "requester", requesterID, "target", targetID)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	loc, err := h.store.GetLocation(r.Context(), targetID)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(loc); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}
