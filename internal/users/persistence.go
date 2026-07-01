package users

import (
	"context"

	"github.com/mohamedveron/go_app_template/internal/users/domain"
)

type usersPersistence interface {
	Create(ctx context.Context, u *domain.User) error
	ReadByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, cursor domain.Cursor) (*domain.UserPage, error)
}
