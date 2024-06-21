package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type Session struct {
	UserID string
}

const (
	MaxAge = 86400 * 30
	IsProd = false
)

var SessionStore *sessions.CookieStore

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	SessionStore.Options.MaxAge = MaxAge
	SessionStore.Options.Path = "/"
	SessionStore.Options.HttpOnly = true
	SessionStore.Options.Secure = IsProd

	gothic.Store = SessionStore
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			"http://localhost:8080/auth/google/callback",
		),
	)
}
