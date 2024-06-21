package database

import (
	"context"
	"fmt"
	"lango/internal/database/domain"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
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

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("err", err)
		return err
	}
	defer conn.Close(ctx)

	queries := domain.New(conn)

	// m := domain.NewModels(db)
	DB = queries
	return nil
}
