package cmd

import (
	"database/sql"
	"fmt"
	"tempest-user-service/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	// https://github.com/golang-migrate
	// file:///absolute/path
	// file://relative/path
	pathToMigration = "file://db/migrations"
)

func OpenDB(conf *config.DB) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.Username,
		conf.Password,
		conf.Name,
	)

	conn, err := sql.Open(conf.Driver, connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db, err %v", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("error sending ping to db, err %v", err)
	}

	return conn, nil
}

func MigrateDB(db *sql.DB, driverType string) error {

	driverInstance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error getting driver instance, err %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		pathToMigration,
		driverType,
		driverInstance,
	)
	if err != nil {
		return fmt.Errorf("error getting migrations, err %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrate up, err %v", err)
	}

	return nil
}
