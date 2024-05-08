package types

import "github.com/google/uuid"

type AuthenticatedUser struct {
	ID          uuid.UUID
	Email       string
	LoggedIn    bool
	AccessToken string
}

type UserContextKey string

const UserKey = UserContextKey("user")
