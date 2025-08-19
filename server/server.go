package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/stfsy/go-api-kit/config"
	"github.com/stfsy/go-api-kit/server/handlers"
	"github.com/stfsy/go-api-kit/server/middlewares"
	"github.com/stfsy/go-api-kit/utils"
	"github.com/urfave/negroni"
)

var logger = utils.NewLogger("server")

type ServerConfig struct {
	MuxCallback    func(*http.ServeMux)
	ListenCallback func()
	PortOverride   string
}

type Server struct {
	serverConfig *ServerConfig
	server       *http.Server
}

func NewServer(serverConfig *ServerConfig) *Server {
	return &Server{
		serverConfig: serverConfig,
		server:       nil,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	if s.serverConfig.MuxCallback != nil {
		s.serverConfig.MuxCallback(mux)
	}

	mux.HandleFunc("/", handlers.NotFoundHandler)

	n := createMiddlewareHandler()
	n.UseHandler(mux)

	configuration := config.Get()
	port := configuration.Port
	if s.serverConfig.PortOverride != "" {
		port = s.serverConfig.PortOverride
	}

	s.server = createServer(port, n)

	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("unable to bind to ip %w", err)
	}

	if s.serverConfig.ListenCallback != nil {
		s.serverConfig.ListenCallback()
	}

	logger.Info(fmt.Sprintf("Listening on port %s", port))
	err = s.server.Serve(ln)

	if err != nil && err.Error() != "http: Server closed" {
		return fmt.Errorf("unable to accept incoming connections %w", err)
	}

	return nil
}

func createServer(port string, n *negroni.Negroni) *http.Server {
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
		Addr:    addr,
		Handler: n,
	}
}

func createMiddlewareHandler() *negroni.Negroni {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(middlewares.NewRequireMaxBodyLengthMiddleware())
	n.Use(middlewares.NewAccessLog())
	n.Use(middlewares.NewRequireContentTypeMiddleware("application/json"))
	n.Use(middlewares.NewRespondWithSecurityHeadersMiddleware())
	n.Use(middlewares.NewNoCacheHeadersMiddleware())
	return n
}

func (s *Server) Stop() {
	// Set a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	} else {
		logger.Info("Server stopped")
	}
}
