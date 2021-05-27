package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	RemoteOfficeDB *sqlx.DB
)

type SSLMode string
const (
	SSLModeEnable SSLMode = "enable"
	SSLModeDisable SSLMode = "disable"
)

// ConnectAndMigrate function connects with a given database and returns error if there is any error
func ConnectAndMigrate (host, port, databaseName, user, password string, sslMode SSLMode) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)
	DB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}
	RemoteOfficeDB = DB
	if err := migrateUp(DB); err != nil {
		return err
	}
	return nil
}

// migrate function migrate the database and handles the migration logic
func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations/",
		"postgres", driver)

	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Tx provides the transaction wrapper
func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := RemoteOfficeDB.Beginx()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to start a transaction: %v", err))
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				logrus.Errorf("failed to rollback tx: %s", err)
			}
			return
		}
		if err := tx.Commit(); err != nil {
			logrus.Errorf("failed to commit tx: %s", err)
		}
	}()
	err = fn(tx)
	return err
}
