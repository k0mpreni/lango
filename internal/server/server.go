package server

import (
	"fmt"
	"lango/internal/auth"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
}

var CurrentServer Server

func NewServer() *http.Server {
	auth.NewAuth()

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	CurrentServer = Server{
		port: port,
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
