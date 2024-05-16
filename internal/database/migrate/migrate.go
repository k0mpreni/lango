package main

import (
	"lango/internal/database"
	"lango/internal/database/domain"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func createDB() (*domain.Models, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return database.New(), nil
}

func main() {
	db := database.New()

	// Create migration instance
	driver, err := postgres.WithInstance(db.App.DB, &postgres.Config{})
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
