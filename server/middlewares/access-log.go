package middlewares

import (
	"log"
	"os"

	"github.com/urfave/negroni"
)

// NewLogger returns a new Logger instance
func NewAccessLog() *negroni.Logger {
	l := negroni.NewLogger()
	l.ALogger = log.New(os.Stdout, "", 0)
	l.SetFormat("{{.StartTime}} \"{{.Method}} {{.Path}} {{.Request.Proto}}\" {{.Status}} {{.Duration}} {{.Request.UserAgent}} ")
	return l
}
