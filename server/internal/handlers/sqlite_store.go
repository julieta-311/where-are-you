package handlers

import (
	"context"
	"database/sql"
	"time"

	"github.com/julieta-311/where-are-you/server/internal/models"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &SQLiteStore{db: db}
	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLiteStore) init() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS requests (
			id TEXT PRIMARY KEY,
			requester_id TEXT,
			target_id TEXT,
			status TEXT,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS locations (
			user_id TEXT PRIMARY KEY,
			latitude REAL,
			longitude REAL,
			timestamp DATETIME
		)`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLiteStore) CreateRequest(ctx context.Context, req models.LocationRequest) error {
	query := `INSERT INTO requests (id, requester_id, target_id, status, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, req.ID, req.RequesterID, req.TargetID, req.Status, req.CreatedAt)
	return err
}

func (s *SQLiteStore) GetRequest(ctx context.Context, id string) (models.LocationRequest, error) {
	query := `SELECT id, requester_id, target_id, status, created_at FROM requests WHERE id = ?`
	var req models.LocationRequest
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID, &req.RequesterID, &req.TargetID, &req.Status, &req.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return req, models.ErrNotFound
	}
	return req, err
}

func (s *SQLiteStore) GetPendingRequests(ctx context.Context, userID string) ([]models.LocationRequest, error) {
	query := `SELECT id, requester_id, target_id, status, created_at 
	          FROM requests WHERE target_id = ? AND status = 'pending'`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []models.LocationRequest
	for rows.Next() {
		var req models.LocationRequest
		if err := rows.Scan(&req.ID, &req.RequesterID, &req.TargetID, &req.Status, &req.CreatedAt); err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (s *SQLiteStore) UpdateRequestStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE requests SET status = ? WHERE id = ?`
	res, err := s.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *SQLiteStore) UpdateLocation(ctx context.Context, loc models.LocationUpdate) error {
	query := `INSERT OR REPLACE INTO locations (user_id, latitude, longitude, timestamp) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, loc.UserID, loc.Latitude, loc.Longitude, loc.Timestamp)
	return err
}

func (s *SQLiteStore) GetLocation(ctx context.Context, userID string) (models.LocationUpdate, error) {
	query := `SELECT user_id, latitude, longitude, timestamp FROM locations WHERE user_id = ?`
	var loc models.LocationUpdate
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&loc.UserID, &loc.Latitude, &loc.Longitude, &loc.Timestamp)
	if err == sql.ErrNoRows {
		return loc, models.ErrNotFound
	}
	return loc, err
}

func (s *SQLiteStore) CheckAuthorization(ctx context.Context, requesterID, targetID string) (bool, error) {
	// Authorised if there is an 'accepted' request created within the last hour.
	query := `SELECT COUNT(*) FROM requests 
	          WHERE requester_id = ? AND target_id = ? AND status = 'accepted' AND created_at > ?`
	var count int
	oneHourAgo := time.Now().Add(-time.Hour)
	err := s.db.QueryRowContext(ctx, query, requesterID, targetID, oneHourAgo).Scan(&count)
	return count > 0, err
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
