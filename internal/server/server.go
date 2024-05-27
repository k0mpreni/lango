package server

import (
	"fmt"
	"lango/internal/database"
	"lango/internal/database/domain"
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

	db *domain.Models
	sb *supabase.Client
}

var CurrentServer Server

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if err := supa.Init(); err != nil {
		fmt.Println("error init supabase client", err)
	}
	if err := database.Init(); err != nil {
		fmt.Println("error init database", err)
	}

	CurrentServer = Server{
		port: port,
		db:   database.DB,
		sb:   supa.Client,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", CurrentServer.port),
		Handler:      CurrentServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
