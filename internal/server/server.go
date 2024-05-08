package server

import (
	"fmt"
	"lango/internal/database"
	"lango/internal/supa"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nedpals/supabase-go"
)

type Server struct {
	port int

	db database.Service
	sb *supabase.Client
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if err := supa.Init(); err != nil {
		fmt.Println("error init supabase client", err)
	}
	NewServer := &Server{
		port: port,
		db:   database.New(),
		sb:   supa.Client,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	supa.Init()

	return server
}
