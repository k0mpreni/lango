package server

import (
	"encoding/json"
	"lango/cmd/web"
	"lango/cmd/web/handler"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(handler.WithUser)

	// r.Get("/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get("/", handler.Make(handler.HomeHandler))

	r.Get("/login", handler.Make(handler.LoginGetHandler))
	r.Get("/login/activate", handler.Make(handler.LoginActivateEmailHandler))
	r.Post("/login", handler.Make(handler.LoginPostHandler))

	r.Get("/login/password-reset", handler.Make(handler.PasswordResetHandler))
	r.Post("/login/password-reset", handler.Make(handler.PasswordResetPostHandler))

	// provider := chi.URLParam(req, "provider")
	//
	// req = req.WithContext(context.WithValue(context.Background(), "provider", provider))
	// // try to get the user without re-authenticating
	// if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
	// 	fmt.Println("GOTHUSER", gothUser)
	// } else {
	// 	fmt.Println("ERR", err)
	// 	gothic.BeginAuthHandler(res, req)
	// }

	r.Get("/login-other-methods", handler.Make(handler.LoginOtherMethodsHandler))

	r.Get("/login-magic-link", handler.Make(handler.MagicLinkGetHandler))
	r.Post("/login-magic-link", handler.Make(handler.MagicLinkCreateHandler))

	r.Get("/signup", handler.Make(handler.SignupGetHandler))
	r.Post("/signup", handler.Make(handler.SignUpCreateHandler))

	r.Get("/auth/provider/{provider}", handler.Make(handler.LoginWithProviderHandler))
	r.Get("/logout", handler.Make(handler.LogoutHandler))
	// r.Get("/logout/", func(res http.ResponseWriter, req *http.Request) {
	// 	gothic.Logout(res, req)
	// 	res.Header().Set("Location", "/")
	// 	res.WriteHeader(http.StatusTemporaryRedirect)
	// })

	r.Get("/auth/{provider}/callback", handler.Make(handler.AuthCallbackHandler))
	// r.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
	// 	fmt.Println("CALLBACK")
	// 	user, err := gothic.CompleteUserAuth(res, req)
	// 	if err != nil {
	// 		fmt.Fprintln(res, err)
	// 		return
	// 	}
	// 	fmt.Println("AUTH CALLBACK", user)
	// 	http.Redirect(res, req, "http://localhost:8080", http.StatusFound)
	// })

	r.Get("/pricing", handler.Make(handler.HandlePricingIndex))

	// Protected Routes
	r.Group(func(chiRouter chi.Router) {
		chiRouter.Use(handler.WithAuth)
		chiRouter.Get("/account", handler.Make(handler.AccountHandler))
		chiRouter.Put("/account", handler.Make(handler.AccountPutHandler))
		chiRouter.Delete("/account", handler.Make(handler.AccountDeleteHandler))

		chiRouter.Get("/courses", handler.Make(handler.CoursesHandler))

		chiRouter.Get("/courses/create", handler.Make(handler.CoursesCreateHandler))
		chiRouter.Post("/courses/create", handler.Make(handler.CoursesCreatePostHandler))

		chiRouter.Get("/courses/{courseId}", handler.Make(handler.CoursesJoinHandler))
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

// func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
// }
