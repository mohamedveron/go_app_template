package serverV1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	api_gen "github.com/mohamedveron/go_app_template/internal/transport/http/api_server_gen/v1"
	httpmw "github.com/mohamedveron/go_app_template/internal/transport/http/middleware"
	"github.com/mohamedveron/go_app_template/internal/users"
)

// HTTP implements api_gen.ServerInterface
type HTTP struct {
	authMiddleware *httpmw.AuthMiddleware
	users          *users.UsersService
	workspace      string
	port           uint16
}

// V1Server handles HTTP server logic for API v1
type V1Server struct {
	router         chi.Router
	authMiddleware *httpmw.AuthMiddleware
	usersService   *users.UsersService
	workspace      string
	port           uint16
}

// NewV1Server creates a new v1 server with the given dependencies
func NewV1Server(authMiddleware *httpmw.AuthMiddleware, usersService *users.UsersService, workspace string, port uint16) *V1Server {
	server := &V1Server{
		authMiddleware: authMiddleware,
		usersService:   usersService,
		workspace:      workspace,
		port:           port,
	}
	server.setupRoutes()
	return server
}

func (s *V1Server) setupRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(httpmw.CorsMiddleware())
	r.Use(s.authMiddleware.RequireAuth())
	s.router = r
}

// RegisterHandlers mounts the generated chi handler onto the v1 router
func (s *V1Server) RegisterHandlers() {
	h := &HTTP{
		authMiddleware: s.authMiddleware,
		users:          s.usersService,
		workspace:      s.workspace,
		port:           s.port,
	}
	api_gen.HandlerFromMux(h, s.router)
}

// GetRouter returns the configured chi router as an http.Handler
func (s *V1Server) GetRouter() http.Handler {
	return s.router
}
