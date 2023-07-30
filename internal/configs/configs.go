package configs

import (
	"time"

	"github.com/mohamedveron/go_app_template/internal/pkg/datastore"
	"github.com/mohamedveron/go_app_template/internal/server/http"
)

// Configs struct handles all dependencies required for handling configurations
type Configs struct {
}

// HTTP returns the configuration required for HTTP package
func (cfg *Configs) HTTP() (*http.Config, error) {
	return &http.Config{
		Port:         8080,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		//DialTimeout:       time.Second * 3,
	}, nil
}

// Datastore returns datastore configuration
func (cfg *Configs) Datastore() (*datastore.Config, error) {
	return &datastore.Config{
		Host:   "postgres",
		Port:   "5432",
		Driver: "postgres",

		StoreName: "go_app",
		Username:  "root",
		Password:  "123321",

		SSLMode: "",

		ConnPoolSize: 10,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 60,
		DialTimeout:  time.Second * 10,
	}, nil
}

// New returns an instance of Config with all the required dependencies initialized
func New() (*Configs, error) {
	return &Configs{}, nil
}
