package types

import "github.com/google/uuid"

// TODO: move to the web
type AuthenticatedUser struct {
	ID       uuid.UUID
	Email    string
	LoggedIn bool
	Provider string
}

type UserContextKey string

const UserKey = UserContextKey("user")
