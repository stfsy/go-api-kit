package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/stfsy/go-api-kit/server"
	"github.com/stfsy/go-api-kit/utils"
)

var logger = utils.NewLogger("main")
var s *server.Server

func main() {
	startServerNonBlocking()
	stopServerAfterSignal()
}

func stopServerAfterSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	s.Stop()

	logger.Info("Graceful shutdown complete.")
}

func startServerNonBlocking() {
	s = server.NewServer(&server.ServerConfig{
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
	go func() {
		err := s.Start()
		if err != nil {
			panic(fmt.Errorf("unable to start server %w", err))
		}
	}()
}
