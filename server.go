package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedveron/go_app_template/internal/transport/http/middleware"
	serverV1 "github.com/mohamedveron/go_app_template/internal/transport/http/v1"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Server handles HTTP server logic with versioning support
type Server struct {
	router         *gin.Engine
	authMiddleware *middleware.AuthMiddleware
	workspace      string
	port           uint16
	version        string
	// Versioned servers
	v1Server *serverV1.V1Server
}

// NewServer creates a new server with versioned API support
func NewServer(authMiddleware *middleware.AuthMiddleware,
	workspace string, port uint16, version string) *Server {
	server := &Server{
		authMiddleware: authMiddleware,
		workspace:      workspace,
		port:           port,
		version:        version,
	}

	// Create versioned servers
	server.v1Server = serverV1.NewV1Server(authMiddleware, workspace, port)

	server.setupRoutes()
	return server
}

// setupRoutes configures all the HTTP routes with versioning
func (s *Server) setupRoutes() {
	// Create gin router without default middleware (no default logging)
	r := gin.New()

	// Add recovery middleware
	r.Use(gin.Recovery())

	// Add OpenTelemetry middleware
	r.Use(otelgin.Middleware("go_app_template"))

	// Add structured logging middleware
	r.Use(middleware.LoggingMiddleware())

	// Register versioned routes
	s.registerVersionedRoutes(r)

	s.router = r
}

// registerVersionedRoutes sets up all versioned API routes
func (s *Server) registerVersionedRoutes(router *gin.Engine) {
	// Health endpoint (unversioned)
	router.GET("/health", s.GetHealth)

	// Register v1 API routes
	s.v1Server.RegisterHandlers()

	// Mount v1 routes directly
	v1Router := s.v1Server.GetRouter()

	// Register all v1 routes under /api/v1 prefix
	v1APIGroup := router.Group("/api/v1")
	{
		// Forward all requests to the v1 router
		v1APIGroup.Any("/*path", func(c *gin.Context) {
			v1Router.ServeHTTP(c.Writer, c.Request)
		})
	}

	// WebSocket routes under /ws/v1 prefix
	v1WSGroup := router.Group("/ws/v1")
	{
		v1WSGroup.Any("/*path", func(c *gin.Context) {
			v1Router.ServeHTTP(c.Writer, c.Request)
		})
	}
}

// GetRouter returns the configured Gin router
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetHealth provides the main health endpoint
func (s *Server) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":                 "ok",
		"version":                s.version,
		"api_supported_versions": []string{"v1"},
	})
}
