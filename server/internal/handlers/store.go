package handlers

import (
	"context"
	"sync"
	"time"

	"github.com/julieta-311/where-are-you/server/internal/models"
)

type MemoryStore struct {
	mu        sync.RWMutex
	requests  map[string]models.LocationRequest
	locations map[string]models.LocationUpdate
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		requests:  make(map[string]models.LocationRequest),
		locations: make(map[string]models.LocationUpdate),
	}
}

func (s *MemoryStore) CreateRequest(ctx context.Context, req models.LocationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests[req.ID] = req
	return nil
}

func (s *MemoryStore) GetRequest(ctx context.Context, id string) (models.LocationRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	req, ok := s.requests[id]
	if !ok {
		return models.LocationRequest{}, models.ErrNotFound
	}
	return req, nil
}

func (s *MemoryStore) GetPendingRequests(ctx context.Context, userID string) ([]models.LocationRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var pending []models.LocationRequest
	for _, req := range s.requests {
		if req.TargetID == userID && req.Status == "pending" {
			pending = append(pending, req)
		}
	}
	return pending, nil
}

func (s *MemoryStore) UpdateRequestStatus(ctx context.Context, id string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	req, ok := s.requests[id]
	if !ok {
		return models.ErrNotFound
	}
	req.Status = status
	s.requests[id] = req
	return nil
}

func (s *MemoryStore) UpdateLocation(ctx context.Context, loc models.LocationUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.locations[loc.UserID] = loc
	return nil
}

func (s *MemoryStore) GetLocation(ctx context.Context, userID string) (models.LocationUpdate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	loc, ok := s.locations[userID]
	if !ok {
		return models.LocationUpdate{}, models.ErrNotFound
	}
	return loc, nil
}

func (s *MemoryStore) CheckAuthorization(ctx context.Context, requesterID, targetID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// For prototype: Authorised if there is an 'accepted' request within the last hour.
	for _, req := range s.requests {
		if req.RequesterID == requesterID && req.TargetID == targetID && req.Status == "accepted" {
			if time.Since(req.CreatedAt) < time.Hour {
				return true, nil
			}
		}
	}
	return false, nil
}
