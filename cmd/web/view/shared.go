package view

import (
	"context"
	"fmt"
	"lango/cmd/web/types"
)

// This is to get display data like usernames, themes... not auth data
func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserKey).(types.AuthenticatedUser)
	fmt.Println("Authenticated", user.ID, user.Email, ok)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}
