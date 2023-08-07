package api

import (
	"context"

	"github.com/mohamedveron/go_app_template/internal/users/domain"
)

// CreateUser is the API to create/signup a new user
func (a *API) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	u, err := a.users.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ReadUserByEmail is the API to read an existing user by their email
func (a *API) ReadUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := a.users.ReadByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}
