package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/stfsy/go-api-kit/config"
	"github.com/stfsy/go-api-kit/server/handlers"
	"github.com/stfsy/go-api-kit/server/middlewares"
	"github.com/stfsy/go-api-kit/utils"
	cors "github.com/stfsy/go-cors"
	"github.com/urfave/negroni"
)

var logger = utils.NewLogger("server")

// ServerConfig configures the API server's endpoints, middleware, and startup behavior.
type ServerConfig struct {
	CorsConfig *cors.Options
	// CrossOriginProtection configures CSRF protection.
	CrossOriginProtection *http.CrossOriginProtection
	// MuxCallback registers endpoints and custom middlewares to the HTTP mux.
	MuxCallback func(*http.ServeMux)
	// MiddlewareCallback customizes the Negroni middleware stack before the server starts.
	MiddlewareCallback func(*negroni.Negroni) *negroni.Negroni
	// ListenCallback is called after the server starts listening, before serving requests.
	ListenCallback func()
	// PortOverride manually sets the port. If empty, uses the PORT environment variable.
	PortOverride string
}

type Server struct {
	serverConfig  *ServerConfig
	server        *http.Server
	serverContext context.Context
}

func NewServer(serverConfig *ServerConfig) *Server {
	return &Server{
		serverConfig:  serverConfig,
		server:        nil,
		serverContext: context.Background(),
	}
}

func (s *Server) Start() error {
	if s == nil || s.serverConfig == nil {
		return fmt.Errorf("server configuration is nil")
	}

	mux := http.NewServeMux()
	if s.serverConfig.MuxCallback != nil {
		s.serverConfig.MuxCallback(mux)
	}

	mux.HandleFunc("/", handlers.NotFoundHandler)

	n := createMiddlewareHandler(s.serverContext, s.serverConfig)
	if s.serverConfig.MiddlewareCallback != nil {
		n = s.serverConfig.MiddlewareCallback(n)
	}
	n.UseHandler(mux)

	csrfProtection := s.serverConfig.CrossOriginProtection
	if csrfProtection == nil {
		csrfProtection = createCrossOritinProtection()
	}

	configuration := config.Get()
	port := configuration.Port
	if s.serverConfig.PortOverride != "" {
		port = s.serverConfig.PortOverride
	}

	s.server = createServer(port, csrfProtection.Handler(n))

	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("unable to bind to %s: %w", s.server.Addr, err)
	}

	if s.serverConfig.ListenCallback != nil {
		s.serverConfig.ListenCallback()
	}

	logger.Info(fmt.Sprintf("Listening on port %s", port))
	err = s.server.Serve(ln)

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("unable to accept incoming connections: %w", err)
	}

	return nil
}

func createCrossOritinProtection() *http.CrossOriginProtection {
	return &http.CrossOriginProtection{}
}

func createServer(port string, h http.Handler) *http.Server {
	c := config.Get()

	if port == "" {
		port = "8080"
	}

	var addr string
	if runtime.GOOS == "windows" {
		addr = fmt.Sprintf("localhost:%s", port)
	} else {
		addr = fmt.Sprintf(":%s", port)
	}

	return &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(c.IdleTimeout) * time.Second,
	}
}

func createMiddlewareHandler(_ctx context.Context, sc *ServerConfig) *negroni.Negroni {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(middlewares.NewAccessLog())
	n.Use(middlewares.NewRespondWithSecurityHeadersMiddleware())
	n.Use(middlewares.NewNoCacheHeadersMiddleware())
	n.Use(middlewares.NewRequireHTTP11Middleware())
	n.Use(middlewares.NewRequireMaxBodyLengthMiddleware())
	if sc.CorsConfig != nil {
		n.Use(cors.New(*sc.CorsConfig))
	}
	n.Use(middlewares.NewRequireContentLengthOrTransferEncodingMiddleware())
	n.Use(middlewares.NewRequireContentTypeMiddleware("application/json"))
	return n
}

func (s *Server) Stop() {
	// Set a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()

	if s == nil || s.server == nil {
		logger.Info("Server not running")
		return
	}

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Info(fmt.Sprintf("Server shutdown error: %v", err))
	} else {
		logger.Info("Server stopped")
	}
}
