package migration

import (
	"database/sql"
	"embed"

	_ "github.com/golang-migrate/migrate/v4/source/file" // нужен, если где-то ещё используете file://
	_ "github.com/lib/pq"                                // PostgreSQL driver
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

type Migrator struct {
	driverName string
	dbUrl      string
	migrations migrate.EmbedFileSystemMigrationSource
}

func NewMigrator(driverName, dbUrl, root string, migrations embed.FS) *Migrator {
	return &Migrator{
		driverName: driverName,
		dbUrl:      dbUrl,
		migrations: migrate.EmbedFileSystemMigrationSource{
			FileSystem: migrations,
			Root:       root,
		},
	}
}

func (m *Migrator) RunUp() error {
	db, err := sql.Open(m.driverName, m.dbUrl)

	applied, err := migrate.Exec(db, m.driverName, m.migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")
	return nil
}

func (m *Migrator) RunDown() error {
	db, err := sql.Open(m.driverName, m.dbUrl)

	applied, err := migrate.Exec(db, m.driverName, m.migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")
	return nil
}
