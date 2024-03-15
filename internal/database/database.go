package database

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/synthao/orders/internal/config"
)

const dbName = "postgres"

func NewConnection(cnf *config.DB) *sqlx.DB {
	return sqlx.MustConnect(dbName, cnf.GetDSN())
}

func ApplyMigrations(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", dbName, driver)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
