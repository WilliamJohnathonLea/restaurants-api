package db

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	underlying *migrate.Migrate
}

func (m *Migrator) Close() error {
	sourceErr, dbErr := m.underlying.Close()

	return errors.Join(sourceErr, dbErr)
}

func (m *Migrator) Run() error {
	err := m.underlying.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

func NewMigrator(dbUrl, migrationsDir string) (*Migrator, error) {
	underlying, err := migrate.New(
		migrationsDir,
		dbUrl,
	)

	if err != nil {
		return nil, err
	}

	migrator := &Migrator{
		underlying: underlying,
	}

	return migrator, nil
}
