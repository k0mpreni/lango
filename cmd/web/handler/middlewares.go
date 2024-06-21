package handler

import (
	"context"
	"fmt"
	"lango/cmd/web/types"
	"lango/internal/auth"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		session, err := auth.SessionStore.Get(r, SessionUserKey)
		if err != nil {
			fmt.Println("err getting session", err)
			next.ServeHTTP(w, r)
			return
		}

		accessToken := session.Values["accessToken"]
		email := session.Values["email"]
		userId := uuid.UUID{}

		at, _ := accessToken.(string)
		userEmail, _ := email.(string)
		if accessToken == nil || email == nil || len(userEmail) < 1 {
			next.ServeHTTP(w, r)
			return
		}

		// Check if user is logged in
		// TODO: Check if valid user
		u := types.AuthenticatedUser{
			ID:          userId,
			Email:       userEmail,
			LoggedIn:    true,
			AccessToken: at,
		}

		ctx := context.WithValue(r.Context(), types.UserKey, u)
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
