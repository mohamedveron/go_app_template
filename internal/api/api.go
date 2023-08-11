package api

import (
	"time"

	"github.com/mohamedveron/go_app_template/internal/users"
)

var (
	now = time.Now()
)

// API holds all the dependencies required to expose APIs. And each API is a function with *API as its receiver
type API struct {
	users *users.UsersService
}

// Health returns the health of the app along with other info like version
func (a *API) Health() (map[string]interface{}, error) {
	return map[string]interface{}{
		"env":        "testing",
		"version":    "v0.1.0",
		"commit":     "<git commit hash>",
		"status":     "all systems up and running",
		"startedAt":  now.String(),
		"releasedOn": now.String(),
	}, nil

}

// NewService returns a new instance of API with all the dependencies initialized
func NewService(us *users.UsersService) (*API, error) {
	return &API{
		users: us,
	}, nil
}
