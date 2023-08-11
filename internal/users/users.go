package users

import (
	"context"
	"strings"

	"github.com/mohamedveron/go_app_template/internal/users/domain"
)

// CreateUser creates a new user
func (us *UsersService) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	u.SetDefaults()
	u.Sanitize()

	err := u.Validate()
	if err != nil {
		return nil, err
	}

	err = us.persistence.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ReadByEmail returns a user which matches the given email
func (us *UsersService) ReadByEmail(ctx context.Context, email string) (*domain.User, error) {
	email = strings.TrimSpace(email)

	u, err := us.persistence.ReadByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}
