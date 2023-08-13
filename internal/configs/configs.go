package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mohamedveron/go_app_template/cmd/server/http"
	"github.com/mohamedveron/go_app_template/internal/pkg/datastore"
	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
)

// Configs struct handles all dependencies required for handling configurations
type Configs struct {
	AppName              string `yaml:"appName" env:"APP_NAME" envDefault:"dsp"`
	Version              string `yaml:"version" env:"APP_VERSION" envDefault:"v0.0.0"`
	ServiceAccountBase64 string `yaml:"serviceAccountBase64" env:"SERVICE_ACCOUNT_BASE64"`
	Environment          string `yaml:"goenv" env:"GOENV"`

	MongoDB struct {
		URI                    string        `env:"CONNECTION_STRING"`
		PingTimeout            time.Duration `env:"PING_TIMEOUT" envDefault:"3s"`
		ConnectTimeout         time.Duration `env:"CONNECT_TIMEOUT" envDefault:"5s"`
		HeartbeatInterval      time.Duration `env:"HEARTBEAT_INTERVAL" envDefault:"10s"`
		LocalThreshold         time.Duration `env:"LOCAL_THRESHOLD" envDefault:"15ms"`
		MaxConnIdleTime        time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"60s"`
		MaxPoolSize            uint64        `env:"MAX_POOL_SIZE" envDefault:"100"`
		MinPoolSize            uint64        `env:"MIN_POOL_SIZE" envDefault:"1"`
		ServerSelectionTimeout time.Duration `env:"SERVER_SELECTION_TIMEOUT" envDefault:"30s"`
	}
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

func (cfg *Configs) AppFullname() string {
	return fmt.Sprintf("%s%s", cfg.AppName, cfg.Version)
}

// New returns an instance of Config with all the required dependencies initialized
func New() (*Configs, error) {
	return &Configs{}, nil
}
