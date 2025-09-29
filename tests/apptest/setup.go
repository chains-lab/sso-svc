package apptest

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/cmd/migrations"
	"github.com/chains-lab/sso-svc/internal"
	"github.com/chains-lab/sso-svc/internal/data"
	"github.com/chains-lab/sso-svc/internal/data/pgdb"
	"github.com/chains-lab/sso-svc/internal/domain"
	"github.com/chains-lab/sso-svc/internal/domain/services/session"
	"github.com/chains-lab/sso-svc/internal/domain/services/user"
	"github.com/chains-lab/sso-svc/internal/infra/jwtmanager"
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
	app domain.Core

	Log logium.Logger
	Cfg internal.Config
}

func cleanDb(t *testing.T) {
	err := migrations.MigrateDown(testDatabaseURL)
	if err != nil {
		t.Fatalf("migrate down: %v", err)
	}
	err = migrations.MigrateUp(testDatabaseURL)
	if err != nil {
		t.Fatalf("migrate up: %v", err)
	}
}

func newSetup(t *testing.T) (Setup, error) {
	cfg := internal.Config{
		Database: internal.DatabaseConfig{
			SQL: struct {
				URL string `mapstructure:"url"`
			}{
				URL: testDatabaseURL,
			},
		},
		JWT: internal.JWTConfig{
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

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.NewDatabase(
		pgdb.NewUsers(pg),
		pgdb.NewSessions(pg),
	)

	jwtTokenManager := jwtmanager.NewManager(jwtmanager.Config{
		AccessSK:   cfg.JWT.User.AccessToken.SecretKey,
		RefreshSK:  cfg.JWT.User.RefreshToken.SecretKey,
		AccessTTL:  cfg.JWT.User.AccessToken.TokenLifetime,
		RefreshTTL: cfg.JWT.User.RefreshToken.TokenLifetime,
		Iss:        cfg.Service.Name,
	})

	logic := domain.NewCore(
		user.New(database),
		session.New(database, jwtTokenManager),
	)

	return Setup{
		app: logic,
	}, nil
}
