package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type MigrationDB struct {
	DB *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *MigrationDB
)

func New() *MigrationDB {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		database,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("err", err)
	}
	dbInstance = &MigrationDB{DB: db}
	return dbInstance
}

func migrateDB(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Point to your migration files. Here we're using local files, but it could be other sources.
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrate/migrations", // source URL
		"postgres", // database name
		driver,     // database instance
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		fmt.Println("Running Migrations")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}

func main() {
	db := New()

	migrateDB(db.DB)
}
