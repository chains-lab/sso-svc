package apptest

import (
	"database/sql"
	"testing"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/dbx"
)

// TEST DATABASE CONNECTION
const testDatabaseURL = "postgresql://postgres:postgres@localhost:7777/postgres?sslmode=disable"

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}

type Setup struct {
	app app.App
	Log logium.Logger
	Cfg config.Config
}

func cleanDb(t *testing.T) {
	err := dbx.MigrateDown(testDatabaseURL)
	if err != nil {
		t.Fatalf("migrate down: %v", err)
	}
	err = dbx.MigrateUp(testDatabaseURL)
	if err != nil {
		t.Fatalf("migrate up: %v", err)
	}
}

func newSetup(t *testing.T) (Setup, error) {
	cfg := config.Config{
		Database: config.DatabaseConfig{
			SQL: struct {
				URL string `mapstructure:"url"`
			}{
				URL: testDatabaseURL,
			},
		},
		JWT: config.JWTConfig{
			User: struct {
				AccessToken struct {
					SecretKey     string        `mapstructure:"secret_key"`
					TokenLifetime time.Duration `mapstructure:"token_lifetime"`
				} `mapstructure:"access_token"`
				RefreshToken struct {
					SecretKey     string        `mapstructure:"secret_key"`
					EncryptionKey string        `mapstructure:"encryption_key"`
					TokenLifetime time.Duration `mapstructure:"token_lifetime"`
				} `mapstructure:"refresh_token"`
			}{
				AccessToken: struct {
					SecretKey     string        `mapstructure:"secret_key"`
					TokenLifetime time.Duration `mapstructure:"token_lifetime"`
				}{
					SecretKey:     "UnG06MAU2i1Mvqf8", //example
					TokenLifetime: time.Minute * 15,
				},
				RefreshToken: struct {
					SecretKey     string        `mapstructure:"secret_key"`
					EncryptionKey string        `mapstructure:"encryption_key"`
					TokenLifetime time.Duration `mapstructure:"token_lifetime"`
				}{
					SecretKey:     "6DSjhhT9KIezubpR", //example
					EncryptionKey: "Zlyh20N8uojZHFdO", //example
					TokenLifetime: time.Hour * 24 * 7,
				},
			},
		},
	}

	log := logium.NewLogger("debug", "text")

	a, err := app.NewApp(cfg)
	if err != nil {
		t.Fatal(err)
	}

	return Setup{
		app: a,
		Log: logium.NewWithBase(log),
		Cfg: cfg,
	}, nil
}
