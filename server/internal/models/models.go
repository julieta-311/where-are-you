package models

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type LocationRequest struct {
	ID          string    `json:"id"`
	RequesterID string    `json:"requester_id"`
	TargetID    string    `json:"target_id"`
	Status      string    `json:"status"` // Pending, accepted, denied.
	CreatedAt   time.Time `json:"created_at"`
}

type LocationUpdate struct {
	UserID    string    `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

type Store interface {
	CreateRequest(ctx context.Context, req LocationRequest) error
	GetRequest(ctx context.Context, id string) (LocationRequest, error)
	GetPendingRequests(ctx context.Context, userID string) ([]LocationRequest, error)
	UpdateRequestStatus(ctx context.Context, id string, status string) error
	UpdateLocation(ctx context.Context, loc LocationUpdate) error
	GetLocation(ctx context.Context, userID string) (LocationUpdate, error)
	// CheckAuthorization verifies if the requester is authorised to see the target's location.
	CheckAuthorization(ctx context.Context, requesterID, targetID string) (bool, error)
}
