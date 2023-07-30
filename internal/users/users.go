package users

import (
	"context"
	"strings"
	"time"

	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	"github.com/pkg/errors"
)

// User holds all data required to represent a user
type User struct {
	FirstName string     `json:"firstName,omitempty"`
	LastName  string     `json:"lastName,omitempty"`
	Mobile    string     `json:"mobile,omitempty"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (u *User) setDefaults() {
	now := time.Now()
	if u.CreatedAt == nil {
		u.CreatedAt = &now
	}

	if u.UpdatedAt == nil {
		u.UpdatedAt = &now
	}
}

// Sanitize is used to sanitize/cleanup the fields of User
func (u *User) Sanitize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.Mobile = strings.TrimSpace(u.Mobile)
}

// Validate is used to validate the fields of User
func (u *User) Validate() error {
	if u.Email == "" {
		return nil
	}

	err := validateEmail(u.Email)
	if err != nil {
		return err
	}

	return nil
}

func validateEmail(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email address provided")
	}

	return nil
}

type userCachestore interface {
	SetUser(ctx context.Context, email string, u *User) error
	ReadUserByEmail(ctx context.Context, email string) (*User, error)
}
type store interface {
	Create(ctx context.Context, u *User) error
	ReadByEmail(ctx context.Context, email string) (*User, error)
}

// Users struct holds all the dependencies required for the users package. And exposes all services
// provided by this package as its methods
type Users struct {
	logHandler logger.Logger
	cachestore userCachestore
	store      store
}

// CreateUser creates a new user
func (us *Users) CreateUser(ctx context.Context, u *User) (*User, error) {
	u.setDefaults()
	u.Sanitize()

	err := u.Validate()
	if err != nil {
		return nil, err
	}

	err = us.store.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ReadByEmail returns a user which matches the given email
func (us *Users) ReadByEmail(ctx context.Context, email string) (*User, error) {
	email = strings.TrimSpace(email)
	err := validateEmail(email)
	if err != nil {
		return nil, err
	}

	u, err := us.store.ReadByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = us.cachestore.SetUser(ctx, u.Email, u)
	if err != nil {
		// in case of error while storing in cache, it is only logged
		// This behaviour as well as read-through cache behaviour depends on your business logic.
		us.logHandler.Error(err.Error())
	}
	return u, nil
}

// NewService initializes the Users struct with all its dependencies and returns a new instance
// all dependencies of Users should be sent as arguments of NewService
func NewService(
	persistenceStore store,
) (*Users, error) {
	return &Users{
		store: persistenceStore,
	}, nil
}
