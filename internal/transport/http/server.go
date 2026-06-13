package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	httpmw "github.com/mohamedveron/go_app_template/internal/transport/http/middleware"
	serverV1 "github.com/mohamedveron/go_app_template/internal/transport/http/v1"
	"github.com/mohamedveron/go_app_template/internal/users"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Server handles HTTP server logic with versioning support
type Server struct {
	router         chi.Router
	authMiddleware *httpmw.AuthMiddleware
	workspace      string
	port           uint16
	version        string
	v1Server       *serverV1.V1Server
}

// NewServer creates a new server with versioned API support
func NewServer(authMiddleware *httpmw.AuthMiddleware, usersService *users.UsersService,
	workspace string, port uint16, version string) *Server {
	s := &Server{
		authMiddleware: authMiddleware,
		workspace:      workspace,
		port:           port,
		version:        version,
	}
	s.v1Server = serverV1.NewV1Server(authMiddleware, usersService, workspace, port)
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Use(otelhttp.NewMiddleware("go_app_template"))
	r.Use(httpmw.LoggingMiddleware())

	r.Get("/health", s.getHealth)

	s.v1Server.RegisterHandlers()
	r.Mount("/api/v1", s.v1Server.GetRouter())

	s.router = r
}

// GetRouter returns the configured router as an http.Handler
func (s *Server) GetRouter() http.Handler {
	return s.router
}

func (s *Server) getHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":                 "ok",
		"version":                s.version,
		"api_supported_versions": []string{"v1"},
	})
}
