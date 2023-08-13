// Package http
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mohamedveron/go_app_template/internal/api"
	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	"github.com/pkg/errors"
)

const (
	errorLogHTTPStatusCodeThreshold = 499
)

type Config struct {
	Host              string
	Port              int
	Environment       string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	JwkURL            string
	AllowedOrigins    []string
}

type HTTP struct {
	lock   *sync.Mutex
	server *http.Server
	// apis has all the APIs, and respective HTTP handlers will call using this
	apis                      *api.API
	shutdownInitiated         bool
	serverStartTime           time.Time
	liveHealthResponse        map[string]string
	shutdownInitiatedResponse []byte
}

func (ht *HTTP) Start() error {
	ht.serverStartTime = time.Now()
	ht.AppendHealthResponse(
		"http",
		fmt.Sprintf("OK: %s", ht.serverStartTime.Format(time.RFC3339Nano)),
	)
	err := ht.server.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "failed to start http server")
	}
	return nil
}

func (ht *HTTP) ResetHealthResponse() {
	ht.lock.Lock()
	ht.liveHealthResponse = map[string]string{}
	ht.lock.Unlock()
}

func (ht *HTTP) AppendHealthResponse(key, value string) {
	ht.lock.Lock()
	ht.liveHealthResponse[key] = value
	ht.lock.Unlock()
}
func (ht *HTTP) healthResponse() (message []byte, status int) {
	ht.lock.Lock()
	if ht.shutdownInitiated {
		message = ht.shutdownInitiatedResponse
		status = http.StatusServiceUnavailable
	} else {
		status = http.StatusOK
		message, _ = json.Marshal(ht.liveHealthResponse)
	}
	ht.lock.Unlock()

	return message, status
}

type HandlerFuncErr func(w http.ResponseWriter, req *http.Request) error

func (ht *HTTP) HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	/*status, message, _ := HTTPStatusCodeMessage(err)

	response := Error{
		Code:    int32(status),
		Message: message,
	}*/

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	_ = json.NewEncoder(w).Encode(err)

	// log the full error here for troubleshooting
	// maybe we just need internal errors to be logged
	/*if status > errorLogHTTPStatusCodeThreshold {
		logger.ErrWithStacktrace(err)
	}*/
}
func (ht *HTTP) ErrorHandler(fn HandlerFuncErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ht.HandleError(w, fn(w, r))
	}
}
func (ht *HTTP) Shutdown(ctx context.Context) error {
	ht.lock.Lock()
	ht.shutdownInitiated = true
	ht.shutdownInitiatedResponse = []byte(fmt.Sprintf("server is shutting down | %s", time.Now().Format(time.RFC3339Nano)))
	ht.lock.Unlock()

	err := ht.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to shutdown")
	}
	return err
}

func (ht *HTTP) Health(w http.ResponseWriter, _ *http.Request) {
	msg, status := ht.healthResponse()
	if status == http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)
	_, _ = w.Write(msg)
}

func New(apis *api.API, cfg *Config) *HTTP {
	ht := &HTTP{
		lock: &sync.Mutex{},
		apis: apis,
	}
	ht.ResetHealthResponse()
	router := chi.NewRouter()
	/*if cfg.Environment == config.EnvironmentLocal {
		router.Use(middleware.Logger)
	}*/
	router.Get("/-/health", ht.Health)
	v1Router := chi.NewRouter()
	v1Router.Use(middleware.Recoverer)
	v1Router.Use(
		cors.Handler(
			cors.Options{
				AllowCredentials: true,
				AllowedOrigins:   cfg.AllowedOrigins,
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "OPTIONS"},
				AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
			},
		),
	)
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info("address of the app= ", address)
	HandlerFromMux(ht, v1Router)
	router.Mount("/api/v1", v1Router)
	ht.server = &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}
	return ht
}
