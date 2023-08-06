package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/mohamedveron/go_app_template/internal/pkg/datastore"
	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	"github.com/mohamedveron/go_app_template/internal/server/http"
)

// Configs struct handles all dependencies required for handling configurations
type Configs struct {
}

// HTTP returns the configuration required for HTTP package
func (cfg *Configs) HTTP() (*http.Config, error) {
	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "9090"
	}

	port, err := strconv.Atoi(envPort)
	if err != nil {
		logger.Error("wrong port provided", envPort)
	}

	return &http.Config{
		Port:         port,
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
