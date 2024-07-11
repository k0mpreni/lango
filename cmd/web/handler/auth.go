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
	"lango/internal/mail"
	"log/slog"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth/gothic"
	"github.com/nedpals/supabase-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionUserKey = "user"
	ProviderKey    = "provider"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

	passHash, err := HashPassword(params.Password)
	if err != nil {
		return render(r, w, auth.SignUpForm(params, auth.SignUpErrors{
			Password: "A problem happened, please try again",
		}))
	}

	ctx := context.Background()
	user, err := database.DB.CreateUser(ctx, domain.CreateUserParams{
		Activated:    false,
		Email:        params.Email,
		PasswordHash: []byte(passHash),
		Provider:     "email",
		ProviderID:   pgtype.Text{},
	})
	if err != nil {
		return render(r, w, auth.SignUpForm(params, auth.SignUpErrors{
			Password: "A problem happened, please try again",
		}))
	}

	mail.SendEmail(params.Email)

	return render(r, w, auth.EmailSent(user.Email, auth.LoginErrors{}))
}

func LoginGetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.LoginForm(auth.LoginCredentials{}, auth.LoginErrors{}))
}

func LoginActivateEmailHandler(w http.ResponseWriter, r *http.Request) error {
	token := chi.URLParam(r, "token")
	fmt.Println("TOKEN ACTIVATION", token)

	ctx := context.Background()

	user, err := database.DB.ActivateUser(ctx, "alex.labonne@pm.me")
	if err != nil {
		return render(
			r,
			w,
			auth.LoginForm(
				auth.LoginCredentials{},
				auth.LoginErrors{InvalidCredentials: "Something went wrong"},
			),
		)
	}

	if err := setAuthSession(w, r, user); err != nil {
		fmt.Println("Error activating user", err)
		return err
	}

	hxRedirect(w, r, "/")
	return nil
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
	return render(r, w, auth.LoginForm(auth.LoginCredentials{}, auth.LoginErrors{}))
}

func MagicLinkGetHandler(w http.ResponseWriter, r *http.Request) error {
	// TODO: Render full page
	return render(r, w, auth.LoginMagicLinkForm(auth.LoginCredentials{}, auth.LoginErrors{}))
}

func MagicLinkCreateHandler(w http.ResponseWriter, r *http.Request) error {
	credentials := auth.LoginCredentials{Email: r.FormValue("email")}

	// HANDLE MAGIC LINK
	// if err != nil {
	// 	fmt.Println(err)
	// 	slog.Error("error sending magic link", "err", err)
	// }

	return render(r, w, auth.EmailSent(credentials.Email, auth.LoginErrors{}))
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	credentials := auth.LoginCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	err := checkmail.ValidateFormat(credentials.Email)
	if err != nil {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			Email: "Please enter a valid email",
		}))
	}

	// TODO: Check for an email validation util or create one
	if len(credentials.Password) < 8 {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			Password: "Password must be longer than 8 characters",
		}))
	}

	user, err := database.DB.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "Please enter valid credentials",
		}))
	}

	if isValid := CheckPasswordHash(credentials.Password, string(user.PasswordHash)); !isValid {
		slog.Error("login error invalid password", "email", user.Email)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "Please enter valid credentials",
		}))

	}

	if err := setAuthSession(w, r, user); err != nil {
		fmt.Println("login post handler", err)
		return err
	}

	hxRedirect(w, r, "/")
	return nil
}

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) error {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), ProviderKey, provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return err
	}

	var newUser domain.CreateUserParams
	var u domain.User

	ctx := context.Background()
	u, err = database.DB.GetUserByEmail(ctx, user.Email)
	if err != nil {
		fmt.Println("ERROR", err.Error(), u, err == sql.ErrNoRows)
	}

	if len(u.Email) == 0 {
		newUser = domain.CreateUserParams{
			Email:        user.Email,
			Provider:     user.Provider,
			ProviderID:   pgtype.Text{String: user.UserID, Valid: true},
			PasswordHash: []byte{},
			Activated:    true,
		}

		u, err = database.DB.CreateUser(ctx, newUser)
		if err != nil {
			fmt.Println("ERROR creating user", err)
		}
	}

	if err := setAuthSession(w, r, u); err != nil {
		return err
	}

	hxRedirect(w, r, "/")

	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, user domain.User) error {
	session, err := internalAuth.SessionStore.Get(r, SessionUserKey)
	session.Values["email"] = user.Email
	session.Values["userId"] = string(user.ID.Bytes[:])

	if err != nil {
		fmt.Println("ERR SETTING SESSION", err)
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
	r = r.WithContext(context.WithValue(context.Background(), ProviderKey, provider))

	u := types.AuthenticatedUser{
		LoggedIn: false,
	}

	session, _ := internalAuth.SessionStore.Get(r, SessionUserKey)
	ctx := context.WithValue(r.Context(), types.UserKey, u)
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
	r = r.WithContext(context.WithValue(context.Background(), ProviderKey, provider))
	gothic.BeginAuthHandler(w, r)

	return nil
}
