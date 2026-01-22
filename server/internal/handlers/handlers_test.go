package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/julieta-311/where-are-you/server/internal/models"
)

func TestCreateRequest(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := NewMemoryStore()
	h := NewHandler(store, logger)

	body, _ := json.Marshal(map[string]string{
		"requester_id": "user1",
		"target_id":    "user2",
	})
	req, _ := http.NewRequest("POST", "/requests", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateRequest(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response models.LocationRequest
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.RequesterID != "user1" || response.TargetID != "user2" {
		t.Errorf("handler returned unexpected body: %+v", response)
	}
}

func TestFlow(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := NewMemoryStore()
	h := NewHandler(store, logger)
	router := mux.NewRouter()
	router.HandleFunc("/requests/{requestId}/respond", h.RespondToRequest)
	router.HandleFunc("/location/{userId}", h.GetLocation)

	reqID := "req123"
	if err := store.CreateRequest(context.TODO(), models.LocationRequest{
		ID:          reqID,
		RequesterID: "requester",
		TargetID:    "target",
		Status:      "pending",
		CreatedAt:   time.Now(),
	}); err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	respBody, _ := json.Marshal(map[string]string{"status": "accepted"})
	respReq, _ := http.NewRequest("POST", "/requests/"+reqID+"/respond", bytes.NewBuffer(respBody))
	respRR := httptest.NewRecorder()
	router.ServeHTTP(respRR, respReq)

	if respRR.Code != http.StatusOK {
		t.Errorf("failed to accept request: %d", respRR.Code)
	}

	if err := store.UpdateLocation(context.TODO(), models.LocationUpdate{
		UserID:    "target",
		Latitude:  1.23,
		Longitude: 4.56,
	}); err != nil {
		t.Fatalf("failed to update location: %v", err)
	}

	getReq, _ := http.NewRequest("GET", "/location/target?requester_id=requester", nil)
	getRR := httptest.NewRecorder()
	router.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Errorf("failed to get location: %d", getRR.Code)
	}

	var loc models.LocationUpdate
	if err := json.NewDecoder(getRR.Body).Decode(&loc); err != nil {
		t.Fatalf("failed to decode location: %v", err)
	}
	if loc.Latitude != 1.23 {
		t.Errorf("wrong latitude: %f", loc.Latitude)
	}
}
