package dbx

import (
	"database/sql"
	"embed"

	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var Migrations embed.FS

type txKeyType struct{}

var TxKey = txKeyType{}

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: Migrations,
	Root:       "migrations",
}

func MigrateUp(url string) error {
	db, err := sql.Open("postgres", url)

	applied, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to applyConditions migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")

	return nil
}

func MigrateDown(url string) error {
	db, err := sql.Open("postgres", url)

	applied, err := migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to applyConditions migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")

	return nil
}
