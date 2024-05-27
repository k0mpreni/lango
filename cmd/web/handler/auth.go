package handler

import (
	"fmt"
	"lango/cmd/web/view/auth"
	"lango/internal/supa"
	"log/slog"
	"net/http"
	"os"

	"github.com/badoux/checkmail"
	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

const (
	sessionUserKey        = "user"
	sessionAccessTokenKey = "accesToken"
)

func SignupGetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.SignUp())
}

func SignUpCreateHandler(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	err := checkmail.ValidateFormat(params.Email)
	if err != nil {
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

	user, err := supa.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}

	return render(r, w, auth.EmailSent(user.Email))
}

func LoginGetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func PasswordResetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.ResetPassword(auth.LoginErrors{}))
}

func PasswordResetPostHandler(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email: r.FormValue("email"),
	}

	err := supa.Client.Auth.ResetPasswordForEmail(r.Context(), credentials.Email)
	if err != nil {
		fmt.Println(err)
		slog.Error("error sending reset password link", "err", err)
		return render(
			r,
			w,
			auth.ResetPassword(
				auth.LoginErrors{Email: "Something went wrong, please check the email address."},
			),
		)
	}

	return render(r, w, auth.EmailSent(credentials.Email))
}

func LoginOtherMethodsHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.LoginForm(supabase.UserCredentials{}, auth.LoginErrors{}))
}

func MagicLinkGetHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.LoginMagicLinkForm(supabase.UserCredentials{}, auth.LoginErrors{}))
}

func MagicLinkCreateHandler(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email: r.FormValue("email"),
	}

	err := supa.Client.Auth.SendMagicLink(r.Context(), credentials.Email)
	if err != nil {
		fmt.Println(err)
		slog.Error("error sending magic link", "err", err)
	}

	return render(r, w, auth.EmailSent(credentials.Email))
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

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
	resp, err := supa.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "Please enter valid credentials",
		}))

	}

	if err := setAuthSession(w, r, resp.AccessToken); err != nil {
		fmt.Println("login post handler", err)
		return err
	}

	hxRedirect(w, r, "/")
	return nil
}

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	fmt.Println(accessToken)

	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}

	if err := setAuthSession(w, r, accessToken); err != nil {
		return err
	}

	hxRedirect(w, r, "/")

	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values["accessToken"] = accessToken
	err := session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, sessionUserKey)

	err := supa.Client.Auth.SignOut(r.Context(), session.Values["accessToken"].(string))
	if err != nil {
		fmt.Println("Couldn't sign out", err)
	}

	session.Values["accessToken"] = ""
	session.Save(r, w)

	hxRedirect(w, r, "/")

	return nil
}

func LoginWithGoogleHandler(w http.ResponseWriter, r *http.Request) error {
	appUrl := os.Getenv("APP_URL")

	resp, err := supa.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: appUrl + "/auth/callback",
	})
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	hxRedirect(w, r, resp.URL)
	return nil
}
