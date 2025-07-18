package postgresstorage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/config"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_calendar/internal/storage"
	"gopkg.in/yaml.v3"
)

func testConfig() (config.PostgresConfig, string) {
	configPath := filepath.Join("../../../configs/config.yaml")
	f, err := os.Open(configPath)
	if err != nil {
		panic("failed to open config.yaml: " + err.Error())
	}
	defer f.Close()

	var cfg struct {
		Storage struct {
			Postgres struct {
				DSN string `yaml:"dsn"`
			} `yaml:"postgres"`
		} `yaml:"storage"`
		MigrationsPath string `yaml:"migrations_path"`
	}
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		panic("failed to decode config.yaml: " + err.Error())
	}
	// Always join with '../../../' to ensure root-level migrations
	migrationsPath := filepath.Join("../../../", cfg.MigrationsPath)
	return config.PostgresConfig{DSN: cfg.Storage.Postgres.DSN}, migrationsPath
}

func runGooseMigrations(dsn, migrationsPath string) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, absPath)
}

func countEvents(store *Storage, ctx context.Context) (int, error) {
	var count int
	row := store.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events")
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func TestCreateAndGetEvent(t *testing.T) {
	cfg, migrationsPath := testConfig()
	if err := runGooseMigrations(cfg.DSN, migrationsPath); err != nil {
		t.Fatalf("Failed to run goose migrations: %v", err)
	}
	store := New(cfg)
	ctx := context.Background()

	countBefore, err := countEvents(store, ctx)
	if err != nil {
		t.Fatalf("Failed to count events before: %v", err)
	}
	fmt.Println("countBefore", countBefore)

	event := storage.Event{
		Title:       "Test Event",
		Description: "A test event",
		AllDay:      2, // meaning 2 hours
	}

	start := time.Now()
	end := start.Add(time.Duration(event.AllDay) * time.Hour)

	event.Start = &start
	event.End = &end

	// err := store.CreateEvent(ctx, event)
	// if err != nil {
	// 	t.Fatalf("CreateEvent failed: %v", err)
	// }

	err = store.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	var id int
	row := store.db.QueryRowContext(ctx, "SELECT id FROM events ORDER BY id DESC LIMIT 1")
	if err := row.Scan(&id); err != nil {
		t.Fatalf("Failed to get last inserted id: %v", err)
	}
	got, err := store.GetEvent(ctx, id)
	if err != nil {
		t.Fatalf("GetEvent failed: %v", err)
	}
	if got.Title != event.Title || got.Description != event.Description || got.AllDay != event.AllDay {
		t.Errorf("GetEvent returned wrong data: got %+v, want %+v", got, event)
	}

	_, err = store.db.ExecContext(ctx, "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		t.Fatalf("Failed to delete inserted event: %v", err)
	}

	countAfter, err := countEvents(store, ctx)
	if err != nil {
		t.Fatalf("Failed to count events after: %v", err)
	}
	if countBefore != countAfter {
		t.Errorf("Expected event count to be unchanged after test, before=%d after=%d", countBefore, countAfter)
	}
}
