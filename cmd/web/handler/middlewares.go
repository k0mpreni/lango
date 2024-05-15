package handler

import (
	"context"
	"fmt"
	"lango/internal/supa"
	"lango/internal/types"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		session, err := store.Get(r, "user")
		fmt.Println("session", session.Values["accessToken"])
		if err != nil {
			fmt.Println("err getting session", err)
			next.ServeHTTP(w, r)
			return
		}

		accessToken := session.Values["accessToken"]
		if accessToken == nil {
			next.ServeHTTP(w, r)
			return
		}

		at, ok := accessToken.(string)
		fmt.Println("auth user:", accessToken, at, ok)

		resp, err := supa.Client.Auth.User(r.Context(), at)
		if err != nil {
			fmt.Println("err auth user:", err, resp)
			store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
			session, _ := store.Get(r, "user")
			session.Values["accessToken"] = ""
			session.Save(r, w)

			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			ID:          uuid.MustParse(resp.ID),
			Email:       resp.Email,
			LoggedIn:    true,
			AccessToken: accessToken.(string),
		}

		ctx := context.WithValue(r.Context(), types.UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := getAuthenticatedUser(r)

		if !user.LoggedIn {
			path := r.URL.Path
			// TODO: handle redirection for better UX
			hxRedirect(w, r, "/login?redirect_to="+path)
			// http.Redirect(w, r, "/login?redirect_to="+path, http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
