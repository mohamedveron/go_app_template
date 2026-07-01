package users

// UsersService holds all the dependencies required for the users package and exposes all services
// provided by this package as its methods.
type UsersService struct {
	persistence usersPersistence
}

// NewService initializes UsersService with all its dependencies and returns a new instance.
func NewService(
	persistence usersPersistence,
) (*UsersService, error) {
	return &UsersService{
		persistence: persistence,
	}, nil
}
