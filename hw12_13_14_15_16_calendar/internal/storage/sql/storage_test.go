package postgresstorage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/config"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/storage"
	"gopkg.in/yaml.v3"
)

func testConfig() (config.PostgresConfig, string) {
	cwd, _ := os.Getwd()
	fmt.Println("!!! Current working directory:", cwd)
	configPath := filepath.Join("../../../configs/config.yaml") //nolint:gocritic //allowed for test files
	f, err := os.Open(configPath)
	if err != nil {
		panic("Test Events failed to open config.yaml: " + err.Error())
	}
	defer f.Close()

	var cfg struct {
		Storage struct {
			Postgres struct {
				DSN string `yaml:"dsn"`
			} `yaml:"postgres"`
		} `yaml:"storage"`
		MigrationsPath string `yaml:"migrationsPath"`
	}
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		panic("Test Events failed to decode config.yaml: " + err.Error())
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
		errorMsg := err.Error()
		// Sanitize the error to remove the DSN
		// Remove DSN details from error message
		re := regexp.MustCompile(`host=[^\s]+`)
		errorMsg = re.ReplaceAllString(errorMsg, "host=***")

		re = regexp.MustCompile(`user=[^\s]+`)
		errorMsg = re.ReplaceAllString(errorMsg, "user=***")

		re = regexp.MustCompile(`password=[^\s]+`)
		errorMsg = re.ReplaceAllString(errorMsg, "password=***")

		re = regexp.MustCompile(`dbname=[^\s]+`)
		errorMsg = re.ReplaceAllString(errorMsg, "dbname=***")

		return fmt.Errorf("failed to open database connection: %s", errorMsg)
	}

	defer db.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, absPath)
}

//nolint:revive // temporary
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
	dsn := os.Getenv("POSTGRES_DSN")
	cfg.DSN = dsn
	fmt.Println("ENV POSTGRES_DSN: ", dsn)
	if err := runGooseMigrations(cfg.DSN, migrationsPath); err != nil {
		//t.Fatalf("Failed to run goose migrations: please check DSN privatelly")
		if os.Getenv("CI") == "" { // only show details locally
			t.Skipf("Skipping SQL tests: could not run migrations (%v)", err)
		} else {
			t.Skip("Skipping PSQL tests: could not run migrations (details hidden in CI)")
		}
	}
	store := New(cfg)
	ctx := context.Background()

	countBefore, err := countEvents(store, ctx)
	if err != nil {
		t.Fatalf("Failed to count events before: %v", err)
	}
	fmt.Println("countBefore: ", countBefore)

	event := storage.Event{
		Title:       "Test Event",
		Description: "A test event",
		AllDay:      2, // meaning 2 hours
	}

	start := time.Now()
	end := start.Add(time.Duration(event.AllDay) * time.Hour)

	event.Start = &start
	event.End = &end

	err = store.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	countAfterCreate, err := countEvents(store, ctx)
	if err != nil {
		t.Fatalf("Failed to count events after inserting 1 evnt: %v", err)
	}
	if countBefore != countAfterCreate-1 {
		t.Errorf("Expected event count to be equal after test, before=%d after=%d", countBefore, countAfterCreate-1)
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

	assertEventCountUnchanged(ctx, t, store, countBefore)
}

func assertEventCountUnchanged(ctx context.Context, t *testing.T, store *Storage, countBefore int) {
	t.Helper()
	countAfter, err := countEvents(store, ctx)
	if err != nil {
		t.Fatalf("Failed to count events after: %v", err)
	}
	if countBefore != countAfter {
		t.Errorf("Expected event count to be unchanged after test, before=%d after=%d", countBefore, countAfter)
	}
}
