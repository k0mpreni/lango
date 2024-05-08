package server

import (
	"encoding/json"
	"log"
	"net/http"

	"lango/cmd/web"
	"lango/cmd/web/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(handler.WithUser)

	r.Get("/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get("/", handler.Make(handler.HomeHandler))

	r.Get("/login", handler.Make(handler.LoginGetHandler))
	r.Post("/login", handler.Make(handler.LoginPostHandler))

	r.Get("/login/password-reset", handler.Make(handler.PasswordResetHandler))
	r.Post("/login/password-reset", handler.Make(handler.PasswordResetPostHandler))

	r.Get("/login/provider/google", handler.Make(handler.LoginWithGoogleHandler))
	r.Get("/login-other-methods", handler.Make(handler.LoginOtherMethodsHandler))

	r.Get("/login-magic-link", handler.Make(handler.MagicLinkGetHandler))
	r.Post("/login-magic-link", handler.Make(handler.MagicLinkCreateHandler))

	r.Get("/signup", handler.Make(handler.SignupGetHandler))
	r.Post("/signup", handler.Make(handler.SignUpCreateHandler))
	r.Post("/logout", handler.Make(handler.LogoutHandler))

	r.Get("/auth/callback", handler.Make(handler.AuthCallbackHandler))

	r.Get("/pricing", handler.Make(handler.HandlePricingIndex))

	r.Group(func(chiRouter chi.Router) {
		chiRouter.Use(handler.WithAuth)
		chiRouter.Get("/account", handler.Make(handler.AccountHandler))
		chiRouter.Put("/account", handler.Make(handler.AccountPutHandler))
		chiRouter.Post("/account/delete", handler.Make(handler.AccountDeleteHandler))
	})

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
