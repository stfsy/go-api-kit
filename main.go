package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/stfsy/go-api-kit/server"
	"github.com/stfsy/go-api-kit/utils"
)

var logger = utils.NewLogger("main")

func main() {
	s := server.NewServer(&server.ServerConfig{
		MuxCallback: func(*http.ServeMux) {
			// add your endpoints and middlewares here
		},
		ListenCallback: func() {
			// do sth just after listen was called on the server instance and
			// just before the server starts serving requests
		},
		// port override is optional but can be used if you want to
		// define the port manually. If empty the value of env.PORT is used.
		PortOverride: "8080",
	})

	serverErr := make(chan error, 1)
	go func() {
		err := s.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		logger.Error("server startup error", "err", err)
		os.Exit(1)
	case sig := <-sigChan:
		logger.Info("shutdown signal received", "signal", sig)
		s.Stop()
		logger.Info("Graceful shutdown complete.")
	}
}
