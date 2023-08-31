package db

import (
	"fmt"
	"social-network/pkg/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sqlx.DB
}

func Open() (*Database, error) {
	path := utils.GetEnv("DBPATH")
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) RunMigrations() error {
	driver, err := sqlite3.WithInstance(d.DB.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./pkg/db/migrations/",
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migration failed... %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying up migrations... %v", err)
	}
	return nil
}

func (d *Database) MigrateToVersion(version uint) error {
	driver, err := sqlite3.WithInstance(d.DB.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./pkg/db/migrations/",
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migration failed... %v", err)
	}

	err = m.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrating to version %d... %v", version, err)
	}
	return nil
}
