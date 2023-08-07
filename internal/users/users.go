package users

import (
	"context"
	"strings"

	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	"github.com/mohamedveron/go_app_template/internal/users/domain"
	"github.com/mohamedveron/go_app_template/internal/users/persistence"
)

// Users struct holds all the dependencies required for the users package. And exposes all services
// provided by this package as its methods
type Users struct {
	logHandler  logger.Logger
	persistence persistence.UsersPersistence
}

// CreateUser creates a new user
func (us *Users) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
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
func (us *Users) ReadByEmail(ctx context.Context, email string) (*domain.User, error) {
	email = strings.TrimSpace(email)

	u, err := us.persistence.ReadByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// NewService initializes the Users struct with all its dependencies and returns a new instance
// all dependencies of Users should be sent as arguments of NewService
func NewService(
	persistence persistence.UsersPersistence,
) (*Users, error) {
	return &Users{
		persistence: persistence,
	}, nil
}
