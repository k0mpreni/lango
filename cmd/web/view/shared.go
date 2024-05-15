package view

import (
	"context"
	"fmt"
	"lango/internal/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserKey).(types.AuthenticatedUser)
	fmt.Println(user, ok)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}
