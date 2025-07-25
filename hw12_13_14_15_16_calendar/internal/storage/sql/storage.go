package postgresstorage

import (
	"context"
	"database/sql"

	// Import pgx driver for database/sql usage with Postgres storage.
	_ "github.com/jackc/pgx/v4/stdlib"                                              //nolint:depguard // allowed as per our webinars
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
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

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	query := `SELECT id, title, description, start, "end", allday, clinic, userid, service FROM events`
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
