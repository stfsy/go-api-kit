package utils

import (
	"log/slog"
	"os"

	"github.com/stfsy/go-api-kit/config"
)

func NewLogger(component string) *slog.Logger {
	options := &slog.HandlerOptions{}
	var handler slog.Handler
	if config.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}
	return slog.New(handler).With(slog.String("component", component))
}
