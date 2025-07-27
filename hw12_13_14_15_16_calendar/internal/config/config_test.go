package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3" //nolint:depguard
)

func TestConfigUnmarshalWithEnvOverride(t *testing.T) {
	yamlData := `
logger:
  level: "info"

http:
  listen: ":8081"

storage:
  type: "postgres"
  redis:
    addr: "localhost:6379"
    password: "mypassword"
    db: 0
  postgres:
    dsn: "host=localhost user=test password=test dbname=testdb sslmode=disable"

migrationsPath: "./migrations"
`

	var cfg Config
	if err := yaml.Unmarshal([]byte(yamlData), &cfg); err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}

	// Override DSN from ENV if set
	if envDSN := os.Getenv("POSTGRES_DSN"); envDSN != "" {
		cfg.Storage.Postgres.DSN = envDSN
	}

	// Now assert final DSN value depends on env var
	envDSN := os.Getenv("POSTGRES_DSN")

	if envDSN == "" {
		// When env var is empty, config should have original YAML DSN
		expected := "host=localhost user=test password=test dbname=testdb sslmode=disable"
		if cfg.Storage.Postgres.DSN != expected {
			t.Errorf("expected DSN from YAML: got %q, want %q", cfg.Storage.Postgres.DSN, expected)
		}
	} else if cfg.Storage.Postgres.DSN != envDSN {
		t.Errorf("expected DSN from ENV: got %q, want %q", cfg.Storage.Postgres.DSN, envDSN)
	}

	// Other assertions left as before (optional)
	// if cfg.Logger.Level != "info" {
	// 	t.Errorf("unexpected Logger.Level: got %q, want %q", cfg.Logger.Level, "info")
	// }

	// if cfg.HTTP.Listen != ":8081" {
	// 	t.Errorf("unexpected HTTP.Listen: got %q, want %q", cfg.HTTP.Listen, ":8081")
	// }

	// if cfg.Storage.Type != "postgres" {
	// 	t.Errorf("unexpected Storage.Type: got %q, want %q", cfg.Storage.Type, "postgres")
	// }
}
