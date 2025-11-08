package postgresstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	// Import pgx driver for database/sql usage with Postgres storage.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/config"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/storage"
)

var (
	ErrNotFound      = errors.New("event not found")
	ErrContextCancel = errors.New("operation canceled")
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.PostgresConfig) *Storage {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		panic(err) // or return error
	}
	return &Storage{db: db}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	query := `INSERT INTO events (title, description, start, "end", allday, clinic, userid, service)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query,
		event.Title, event.Description, event.Start, event.End, event.AllDay, event.Clinic, event.UserID, event.Service)
	return err
}

func (s *Storage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	var event storage.Event
	query := `SELECT id, title, description, start, "end", allday, clinic, userid, service FROM events WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Start,
		&event.End,
		&event.AllDay,
		&event.Clinic,
		&event.UserID,
		&event.Service)
	return event, err
}

func (s *Storage) ListEvents(ctx context.Context, period storage.Period) ([]storage.Event, error) {
	query := `SELECT id, title, description, start, "end", allday, clinic, userid, service FROM events`
	var whereClause string

	switch period {
	case storage.PeriodDay:
		whereClause = ` WHERE date_trunc('day', start) = date_trunc('day', CURRENT_DATE)`
	case storage.PeriodWeek:
		whereClause = ` WHERE date_trunc('week', start) = date_trunc('week', CURRENT_DATE)`
	case storage.PeriodMonth:
		whereClause = ` WHERE date_trunc('month', start) = date_trunc('month', CURRENT_DATE)`
	case storage.PeriodAll:
		// No WHERE clause needed
	default:
		// Optionally handle unknown periods (could return error instead)
		whereClause = ` WHERE 1=0` // Returns no rows for unknown periods
	}

	if whereClause != "" {
		query += whereClause
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Start,
			&event.End,
			&event.AllDay,
			&event.Clinic,
			&event.UserID,
			&event.Service); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// DeleteEvent removes an event by ID. Returns ErrNotFound if event doesn't exist.
func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	// Check context before starting operation
	select {
	case <-ctx.Done():
		return ErrContextCancel
	default:
	}

	// Execute SQL delete operation
	result, err := s.db.ExecContext(ctx,
		`DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateEvent updates an existing event by ID.
// It returns ErrNotFound if the event doesn't exist.
func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := `UPDATE events
              SET title = $1,
                  description = $2,
                  start = $3,
                  "end" = $4,
                  allday = $5,
                  clinic = $6,
                  userid = $7,
                  service = $8
              WHERE id = $9`
	result, err := s.db.ExecContext(ctx, query,
		event.Title, event.Description, event.Start, event.End,
		event.AllDay, event.Clinic, event.UserID, event.Service, event.ID)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
