package handler

import (
	"lango/cmd/web/view"
	"lango/cmd/web/view/account"
	"lango/internal/supa"
	"log/slog"
	"net/http"
)

// func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Bad Request", http.StatusBadRequest)
// 	}

// 	name := r.FormValue("name")
// 	component := HelloPost(name)
// 	err = component.Render(r.Context(), w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
// 	}
// }

func AccountHandler(w http.ResponseWriter, r *http.Request) error {
	u := getAuthenticatedUser(r)

	return render(r, w, account.Index(account.AccountUser{Email: u.Email}))
}

func AccountPutHandler(w http.ResponseWriter, r *http.Request) error {
	params := account.AccountUser{
		Email:           r.FormValue("email"),
		CurrentPassword: r.FormValue("currentPassword"),
		NewPassword:     r.FormValue("newPassword"),
	}

	// err := checkmail.ValidateFormat(params.Email)
	// if err != nil {
	// 	return render(
	// 		r,
	// 		w,
	// 		account.AccountForm(account.AccountUser{Email: params.Email}, account.AccountErrors{
	// 			Email: "Please enter a valid email",
	// 		}),
	// 	)
	// }

	if len(params.NewPassword) < 8 {
		return render(
			r,
			w,
			account.AccountForm(account.AccountUser{Email: params.Email}, account.AccountErrors{
				NewPassword: "Password must be longer than 8 characters",
			}),
		)
	}

	u := view.AuthenticatedUser(r.Context())

	user, err := supa.Client.Auth.UpdateUser(
		r.Context(),
		u.AccessToken,
		map[string]interface{}{"email": params.Email, "password": params.NewPassword},
	)
	if err != nil {
		slog.Error("updating account error", "err", err)
		return render(
			r,
			w,
			account.AccountForm(account.AccountUser{Email: params.Email}, account.AccountErrors{
				CurrentPassword: "Impossible to update the account",
			}),
		)

	}

	return render(
		r,
		w,
		account.AccountForm(account.AccountUser{Email: user.Email}, account.AccountErrors{}),
	)
}

func AccountDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, account.AccountDeleted())
}
