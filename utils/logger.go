package utils

import (
	"log/slog"
	"os"

	"github.com/stfsy/go-api-kit/config"
)

var (
	isProduction      = config.IsProduction()
	loggerOptions     = &slog.HandlerOptions{}
	loggerJsonHandler = slog.NewJSONHandler(os.Stdout, loggerOptions)
	loggerTextHandler = slog.NewTextHandler(os.Stdout, loggerOptions)
)

func NewLogger(component string) *slog.Logger {
	var loggerHandler slog.Handler
	if isProduction {
		loggerHandler = loggerJsonHandler
	} else {
		loggerHandler = loggerTextHandler
	}
	return slog.New(loggerHandler).With(slog.String("component", component))
}
