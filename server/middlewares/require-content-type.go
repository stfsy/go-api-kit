package middlewares

// adapted from https://raw.githubusercontent.com/go-chi/chi/master/middleware/content_type.go

import (
	"net/http"
	"strings"

	"github.com/stfsy/go-api-kit/server/handlers"
)

type RequireContentTypeMiddleware struct {
	AllowedContentType string
}

func NewRequireContentTypeMiddleware(allowedContentType string) *RequireContentTypeMiddleware {
	return &RequireContentTypeMiddleware{
		AllowedContentType: allowedContentType,
	}
}

// AllowContentType enforces a whitelist of request Content-Types otherwise responds
// with a 415 Unsupported Media Type status.
func (m *RequireContentTypeMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Only enforce content-type for write requests. Non-write requests are passed through.
	if !isWriteRequest(r) {
		next.ServeHTTP(rw, r)
		return
	}

	// For write requests (POST/PUT/PATCH and DELETE with body), require a valid Content-Type.
	// This includes chunked requests where ContentLength may be -1. Use the request headers directly.
	s := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if i := strings.Index(s, ";"); i > -1 {
		s = s[0:i]
	}

	// If no Content-Type provided for a write request, reject it.
	if s == "" {
		handlers.SendUnsupportedMediaType(rw, nil)
		return
	}

	if s == m.AllowedContentType {
		next.ServeHTTP(rw, r)
		return
	}

	handlers.SendUnsupportedMediaType(rw, nil)
}

func isWriteRequest(r *http.Request) bool {
	switch r.Method {
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		return true
	case http.MethodDelete:
		// DELETE may have a body, so treat it as a write request if it does.
		return r.ContentLength != 0
	default:
		return false
	}
}
