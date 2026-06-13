package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mohamedveron/go_app_template/internal/configs"
	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	httpmw "github.com/mohamedveron/go_app_template/internal/transport/http/middleware"
	httpserver "github.com/mohamedveron/go_app_template/internal/transport/http"
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

	server := httpserver.NewServer(authMiddleware, "", uint16(httpCfg.Port), cfg.Version)

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
