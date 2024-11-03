// Package main is the entry point for the k6 CLI application. It assembles all the crucial components for the running.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/stfsy/go-api-kit/app/server"
	"github.com/stfsy/go-api-kit/app/utils"
)

var logger = utils.NewLogger("main")

func main() {
	startServerNonBlocking()
	stopServerAfterSignal()
}

func stopServerAfterSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	server.Stop()

	logger.Info("Graceful shutdown complete.")
}

func startServerNonBlocking() {
	go func() {
		err := server.Start(nil, nil)
		if err != nil {
			panic(fmt.Errorf("unable to start server %w", err))
		}
	}()
}
