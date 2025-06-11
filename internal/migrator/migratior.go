package migrator

import (
	"database/sql"
	"embed"

	"github.com/chains-lab/chains-auth/internal/config"
	_ "github.com/golang-migrate/migrate/v4/source/file" // нужен, если где-то ещё используете file://
	_ "github.com/lib/pq"                                // PostgreSQL driver
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var Migrations embed.FS

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: Migrations,
	Root:       "migrations",
}

func RunUp(cfg config.Config) error {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)

	applied, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")
	return nil
}

func RunDown(cfg config.Config) error {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)

	applied, err := migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")
	return nil
}
