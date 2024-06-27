package database

import (
	"context"
	"fmt"
	"lango/internal/database/domain"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

var DB *domain.Queries

var pgPool *pgxpool.Pool

type CancelConnection func() error

func Init() error {
	ctx := context.Background()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		database,
	)

	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal("err", err)
		return err
	}

	queries := domain.New(dbpool)

	DB = queries
	return nil
}

func Close() {
	pgPool.Close()
}
