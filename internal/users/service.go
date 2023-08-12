package users

import (
	"github.com/mohamedveron/go_app_template/internal/users/persistence"
)

// Users struct holds all the dependencies required for the users package. And exposes all services
// provided by this package as its methods
type UsersService struct {
	persistence persistence.UsersPersistence
}

// NewService initializes the Users struct with all its dependencies and returns a new instance
// all dependencies of Users should be sent as arguments of NewService
func NewService(
	persistence persistence.UsersPersistence,
) (*UsersService, error) {
	return &UsersService{
		persistence: persistence,
	}, nil
}
