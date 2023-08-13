package main

import (
	"fmt"

	"github.com/mohamedveron/go_app_template/cmd/server/http"
	"github.com/mohamedveron/go_app_template/internal/api"
	"github.com/mohamedveron/go_app_template/internal/configs"
	"github.com/mohamedveron/go_app_template/internal/pkg/datastore"
	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	"github.com/mohamedveron/go_app_template/internal/users"
	"github.com/mohamedveron/go_app_template/internal/users/persistence"
)

func main() {

	// load configuration
	cfg, err := configs.New()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	dscfg, err := cfg.Datastore()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	pqdriver, err := datastore.NewPostgresService(dscfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}
	userStore, err := persistence.NewUserPostgresPersistence(pqdriver)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}
	us, err := users.NewService(userStore)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	a, err := api.NewService(us)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	httpCfg, err := cfg.HTTP()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	server := http.New(a, httpCfg)
	server.Start()

}
