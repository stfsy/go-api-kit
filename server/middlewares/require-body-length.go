package middlewares

// adapted from https://raw.githubusercontent.com/go-chi/chi/master/middleware/content_type.go

import (
	"net/http"

	"github.com/stfsy/go-api-kit/config"
)

type RequireMaxBodyLengthMiddleware struct {
	maxSize int
}

func NewRequireMaxBodyLengthMiddleware() *RequireMaxBodyLengthMiddleware {
	maxSize := config.Get().MaxBodySize
	return &RequireMaxBodyLengthMiddleware{
		maxSize: maxSize,
	}
}

// AllowContentType enforces a whitelist of request Content-Types otherwise responds
// with a 415 Unsupported Media Type status.
func (m *RequireMaxBodyLengthMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r.Body = http.MaxBytesReader(rw, r.Body, int64(m.maxSize))
	next.ServeHTTP(rw, r)
}
