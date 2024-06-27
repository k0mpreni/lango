package handler

import (
	"context"
	"fmt"
	"lango/cmd/web/types"
	"lango/internal/auth"
	"lango/internal/database"
	"net/http"
	"strings"
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

		email := session.Values["email"]

		userEmail, _ := email.(string)

		fmt.Println("userEmail", userEmail)

		if email == nil || len(userEmail) < 1 {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.Background()
		user, err := database.DB.GetUserByEmail(ctx, userEmail)
		if err != nil {
			fmt.Println("Invalid user")
			next.ServeHTTP(w, r)
			return
		}

		u := types.AuthenticatedUser{
			ID:       user.ID.Bytes,
			Email:    user.Email,
			LoggedIn: true,
		}

		ctx = context.WithValue(r.Context(), types.UserKey, u)
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
