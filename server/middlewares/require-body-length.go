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

// ServeHTTP enforces a maximum request body length, responding with 413 Payload Too Large if exceeded.
func (m *RequireMaxBodyLengthMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r.Body = http.MaxBytesReader(rw, r.Body, int64(m.maxSize))
	next.ServeHTTP(rw, r)
}
