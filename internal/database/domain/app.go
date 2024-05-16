package domain

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type AppModel struct {
	DB *sql.DB
}

func (m *AppModel) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := m.DB.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
