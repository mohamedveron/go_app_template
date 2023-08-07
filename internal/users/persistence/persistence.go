package persistence

import (
	"context"

	"github.com/mohamedveron/go_app_template/internal/users/domain"
)

type UsersPersistence interface {
	Create(ctx context.Context, u *domain.User) error
	ReadByEmail(ctx context.Context, email string) (*domain.User, error)
}
