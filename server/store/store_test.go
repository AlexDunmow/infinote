package store_test

import (
	"infinote/bindata"
	"fmt"
	"log"
	"testing"

	migrate_bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
)

var pool *dockertest.Pool

const TestDatabaseUser = "test_database_user"
const TestDatabasePass = "test_database_pass"
const TestDatabaseHost = "localhost"
const TestDatabaseName = "test_database_name"

func Setup(t *testing.T) (*sqlx.DB, func(), func(), func()) {
	var conn *sqlx.DB

	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "11-alpine", []string{
		"POSTGRES_USER=" + TestDatabaseUser,
		"POSTGRES_PASSWORD=" + TestDatabasePass,
		"POSTGRES_DB=" + TestDatabaseName,
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port := resource.GetPort("5432/tcp")
	if err = pool.Retry(func() error {
		conn, err = sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			TestDatabaseUser,
			TestDatabasePass,
			TestDatabaseHost,
			port,
			TestDatabaseName,
		))
		if err != nil {
			return err
		}
		err = conn.Ping()
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	conn.MustExec(`
	CREATE EXTENSION IF NOT EXISTS pg_trgm;
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	`)
	err = Migrate(conn)
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	teardown := func() {
		err := pool.Purge(resource)
		if err != nil {
			err = fmt.Errorf("could not purge resource: %w", err)
			fmt.Println(err)
		}
	}
	drop := func() {
		err := Drop(conn)
		if err != nil {
			fmt.Println("Drop database error:", err)
		}
	}
	migrate := func() {
		err := Migrate(conn)
		if err != nil {
			fmt.Println("Migrate database error:", err)
		}
	}
	return conn, drop, migrate, teardown
}
func Drop(conn *sqlx.DB) error {
	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create driver: %w", err)
	}

	s := migrate_bindata.Resource(bindata.AssetNames(),
		func(name string) ([]byte, error) {
			return bindata.Asset(name)
		})

	d, err := migrate_bindata.WithInstance(s)
	if err != nil {
		return fmt.Errorf("run dropper: %w", err)
	}
	m, err := migrate.NewWithInstance("go-bindata", d, "postgres", driver)

	if err != nil {
		return fmt.Errorf("create dropper: %w", err)
	}
	err = m.Drop()
	if err != nil {
		return fmt.Errorf("run dropper: %w", err)
	}
	return nil
}
func Migrate(conn *sqlx.DB) error {
	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create driver: %w", err)
	}

	s := migrate_bindata.Resource(bindata.AssetNames(),
		func(name string) ([]byte, error) {
			return bindata.Asset(name)
		})

	d, err := migrate_bindata.WithInstance(s)
	if err != nil {
		return fmt.Errorf("run migrator: %w", err)
	}
	m, err := migrate.NewWithInstance("go-bindata", d, "postgres", driver)
	if err != nil {
		return fmt.Errorf("run migrator: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("run migrator: %w", err)
	}
	return nil
}
