package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mohamedveron/go_app_template/internal/configs"
	"github.com/mohamedveron/go_app_template/internal/pkg/datastore"
	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	httpserver "github.com/mohamedveron/go_app_template/internal/transport/http"
	httpmw "github.com/mohamedveron/go_app_template/internal/transport/http/middleware"
	"github.com/mohamedveron/go_app_template/internal/users"
	"github.com/mohamedveron/go_app_template/internal/users/persistence"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	httpCfg, err := cfg.HTTP()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	authMiddleware, err := httpmw.NewAuthMiddleware(httpmw.AuthConfig{})
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	dsCfg, err := cfg.Datastore()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	pgPool, err := datastore.NewPostgresService(dsCfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}
	defer pgPool.Close()

	usersPersistence, err := persistence.NewUserPostgresPersistence(pgPool)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	usersService, err := users.NewService(usersPersistence)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v", err))
		return
	}

	server := httpserver.NewServer(authMiddleware, usersService, "", uint16(httpCfg.Port), cfg.Version)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", httpCfg.Port),
		Handler:      server.GetRouter(),
		ReadTimeout:  httpCfg.ReadTimeout,
		WriteTimeout: httpCfg.WriteTimeout,
		IdleTimeout:  time.Second * 60,
	}

	logger.Info(fmt.Sprintf("starting server on :%d", httpCfg.Port))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(fmt.Sprintf("server error: %+v", err))
	}
}
