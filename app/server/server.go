package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/stfsy/go-api-kit/app/config"
	"github.com/stfsy/go-api-kit/app/server/handlers"
	"github.com/stfsy/go-api-kit/app/server/middlewares"
	"github.com/stfsy/go-api-kit/app/utils"
	"github.com/urfave/negroni"
)

var logger = utils.NewLogger("server")

var server *http.Server

func Start(muxCallback func(*http.ServeMux), listenCallback func()) error {
	mux := http.NewServeMux()
	if muxCallback != nil {
		muxCallback(mux)
	}

	mux.HandleFunc("/", handlers.NotFoundHandler)

	n := createMiddlewareHandler()
	n.UseHandler(mux)

	configuration := config.Get()
	port := configuration.Port

	server = createServer(port, n)

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("unable to bind to ip %w", err)
	}

	if listenCallback != nil {
		listenCallback()
	}

	logger.Info(fmt.Sprintf("Listening on port %s", port))
	err = server.Serve(ln)

	if err != nil && err.Error() != "http: Server closed" {
		return fmt.Errorf("unable to accept incoming connections %w", err)
	}

	return nil
}

func createServer(port string, n *negroni.Negroni) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: n,
	}
}

func createMiddlewareHandler() *negroni.Negroni {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(middlewares.NewAccessLog())
	n.Use(middlewares.NewRequireContentTypeMiddleware("application/json"))
	n.Use(middlewares.NewRespondWithSecurityHeadersMiddleware())
	n.Use(middlewares.NewNoCacheHeadersMiddleware())
	return n
}

func Stop() {
	// Set a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	logger.Info("Server stopped")
}
