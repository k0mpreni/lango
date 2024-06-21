package handler

import (
	"context"
	"database/sql"
	"fmt"
	"lango/cmd/web/types"
	"lango/cmd/web/view/auth"
	internalAuth "lango/internal/auth"
	"lango/internal/database"
	"lango/internal/database/domain"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
	"github.com/nedpals/supabase-go"
)

const (
	SessionUserKey        = "user"
	sessionAccessTokenKey = "accesToken"
)

func SignupGetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.SignUpForm(auth.SignupParams{}, auth.SignUpErrors{}))
}

func SignUpCreateHandler(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	err := checkmail.ValidateFormat(params.Email)
	if err != nil {
		fmt.Println(err)
		return render(r, w, auth.SignUpForm(params, auth.SignUpErrors{
			Email: "Please enter a valid email",
		}))
	}

	// TODO: Check for an email validation util or create one
	if len(params.Password) < 8 {
		return render(r, w, auth.SignUpForm(params, auth.SignUpErrors{
			Password: "Password must be longer than 8 characters",
		}))
	}

	if params.Password != params.ConfirmPassword {
		return render(r, w, auth.SignUpForm(params, auth.SignUpErrors{
			Password:        "Passwords do not match",
			ConfirmPassword: "Passwords do not match",
		}))
	}

	// _, err = supa.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
	// 	Email:    params.Email,
	// 	Password: params.Password,
	// })
	// if err != nil {
	// 	switch err.Error() {
	// 	case "User already registered":
	// 		return render(r, w, auth.EmailSent(params.Email, auth.LoginErrors{Email: err.Error()}))
	// 	default:
	// 		return err
	// 	}
	// }

	return render(r, w, auth.EmailSent(params.Email, auth.LoginErrors{}))
}

func LoginGetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.LoginForm(supabase.UserCredentials{}, auth.LoginErrors{}))
}

func PasswordResetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.ResetPassword(auth.LoginErrors{}))
}

func PasswordResetPostHandler(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email: r.FormValue("email"),
	}

	// err := supa.Client.Auth.ResetPasswordForEmail(r.Context(), credentials.Email)
	// if err != nil {
	// 	fmt.Println(err)
	// 	slog.Error("error sending reset password link", "err", err)
	// 	return render(
	// 		r,
	// 		w,
	// 		auth.ResetPassword(
	// 			auth.LoginErrors{Email: "Something went wrong, please check the email address."},
	// 		),
	// 	)
	// }

	return render(r, w, auth.EmailSent(credentials.Email, auth.LoginErrors{}))
}

func LoginOtherMethodsHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.LoginForm(supabase.UserCredentials{}, auth.LoginErrors{}))
}

func MagicLinkGetHandler(w http.ResponseWriter, r *http.Request) error {
	// TODO: Render full page
	return render(r, w, auth.LoginMagicLinkForm(supabase.UserCredentials{}, auth.LoginErrors{}))
}

func MagicLinkCreateHandler(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email: r.FormValue("email"),
	}

	// err := supa.Client.Auth.SendMagicLink(r.Context(), credentials.Email)
	// if err != nil {
	// 	fmt.Println(err)
	// 	slog.Error("error sending magic link", "err", err)
	// }

	return render(r, w, auth.EmailSent(credentials.Email, auth.LoginErrors{}))
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) error {
	// credentials := supabase.UserCredentials{
	// 	Email:    r.FormValue("email"),
	// 	Password: r.FormValue("password"),
	// }

	// err := checkmail.ValidateFormat(credentials.Email)
	// if err != nil {
	// 	return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
	// 		Email: "Please enter a valid email",
	// 	}))
	// }

	// TODO: Check for an email validation util or create one
	// if len(credentials.Password) < 8 {
	// 	return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
	// 		Password: "Password must be longer than 8 characters",
	// 	}))
	// }

	// call supabase
	// resp, err := supa.Client.Auth.SignIn(r.Context(), credentials)
	// if err != nil {
	// 	slog.Error("login error", "err", err)
	// 	return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
	// 		InvalidCredentials: "Please enter valid credentials",
	// 	}))
	//
	// }
	//
	// if err := setAuthSession(w, r, resp.AccessToken); err != nil {
	// 	fmt.Println("login post handler", err)
	// 	return err
	// }

	hxRedirect(w, r, "/")
	return nil
}

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) error {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return err
	}

	ctx := context.Background()
	u, err := database.DB.GetUserByEmail(ctx, user.Email)
	if err != nil {
		fmt.Println("ERRRR", err.Error(), u, err == sql.ErrNoRows)
	}

	newUser := domain.CreateUserParams{
		Email:        u.Email,
		Provider:     u.Provider,
		ProviderID:   u.ProviderID,
		PasswordHash: []byte{},
	}

	createdUser, err := database.DB.CreateUser(ctx, newUser)
	if err != nil {
		fmt.Println("ERRRR creating user", err)
	}
	fmt.Println(createdUser)

	if err := setAuthSession(w, r, createdUser); err != nil {
		return err
	}

	hxRedirect(w, r, "/")

	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, user domain.User) error {
	session, err := internalAuth.SessionStore.Get(r, SessionUserKey)
	session.Values["email"] = user.Email
	session.Values["userId"] = user.ID

	if err != nil {
		return err
	}

	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	u := types.AuthenticatedUser{
		LoggedIn: false,
	}

	session, _ := internalAuth.SessionStore.Get(r, SessionUserKey)
	ctx := context.WithValue(r.Context(), types.UserKey, u)
	session.Values["accessToken"] = ""
	session.Values["email"] = ""
	session.Values["userId"] = ""
	session.Save(r, w)

	gothic.Logout(w, r)

	// hxRedirect(w, r, "/")
	http.Redirect(w, r.WithContext(ctx), "/", http.StatusSeeOther)

	return nil
}

func LoginWithProviderHandler(w http.ResponseWriter, r *http.Request) error {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(w, r)

	return nil
}
