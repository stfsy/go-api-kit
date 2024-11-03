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

	if !hasContentLength(r) || !isWriteRequest(r) {
		next.ServeHTTP(rw, r)
		return
	}

	s := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if i := strings.Index(s, ";"); i > -1 {
		s = s[0:i]
	}

	if s == m.AllowedContentType {
		next.ServeHTTP(rw, r)
		return
	}

	handlers.SendUnsupportedMediaType(rw, nil)
}

func hasContentLength(r *http.Request) bool {
	return r.ContentLength >= 0
}

func isWriteRequest(r *http.Request) bool {
	if r.Method == http.MethodGet {
		return false
	} else if r.Method == http.MethodPatch || r.Method == http.MethodPost || r.Method == http.MethodPut {
		return true
	} else if r.Method == http.MethodDelete {
		return r.Body != nil
	}
	return false
}
