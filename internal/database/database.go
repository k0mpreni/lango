package database

import (
	"database/sql"
	"fmt"
	"lango/internal/database/domain"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *domain.Models
)

var DB *domain.Models

func Init() error {
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
		return err
	}
	m := domain.NewModels(db)
	DB = &m
	return nil
}
