package supa

import (
	"errors"
	"os"

	"github.com/nedpals/supabase-go"
)

var (
	Client   *supabase.Client
	sbUrl    string
	sbSecret string
)

func Init() error {
	sbUrl = os.Getenv("SUPABASE_URL")
	if sbUrl == "" {
		return errors.New("supabase url is required")
	}
	sbSecret = os.Getenv("SUPABASE_SECRET")
	if sbSecret == "" {
		return errors.New("supabase secret is required")
	}
	Client = supabase.CreateClient(sbUrl, sbSecret)

	return nil
}
