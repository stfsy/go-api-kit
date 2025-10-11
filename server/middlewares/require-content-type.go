package middlewares

// adapted from https://raw.githubusercontent.com/go-chi/chi/master/middleware/content_type.go

import (
	"mime"
	"net/http"
	"strings"

	"github.com/stfsy/go-api-kit/server/handlers"
)

type RequireContentTypeMiddleware struct {
	AllowedContentType string
}

func NewRequireContentTypeMiddleware(allowedContentType string) *RequireContentTypeMiddleware {
	// Normalize allowed content type: parse and keep the media type only.
	allowed := strings.ToLower(strings.TrimSpace(allowedContentType))
	if i := strings.Index(allowed, ";"); i > -1 {
		allowed = allowed[0:i]
	}
	// In case caller passed parameters, attempt to parse, else keep the raw media type.
	if mt, _, err := mime.ParseMediaType(allowed); err == nil {
		allowed = strings.ToLower(strings.TrimSpace(mt))
	}

	return &RequireContentTypeMiddleware{
		AllowedContentType: allowed,
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
	// Parse the Content-Type header using mime.ParseMediaType to canonicalize comparisons.
	ctHeader := strings.TrimSpace(r.Header.Get("Content-Type"))
	if ctHeader == "" {
		handlers.SendUnsupportedMediaType(rw, nil)
		return
	}

	mediaType, _, err := mime.ParseMediaType(ctHeader)
	if err != nil {
		// If parsing fails, conservatively reject the request.
		handlers.SendUnsupportedMediaType(rw, nil)
		return
	}

	if strings.ToLower(mediaType) == m.AllowedContentType {
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
		// Consider explicit Content-Length > 0 or chunked transfer encoding.
		if r.ContentLength > 0 {
			return true
		}
		te := strings.ToLower(strings.TrimSpace(r.Header.Get("Transfer-Encoding")))
		if strings.Contains(te, "chunked") {
			return true
		}
		return false
	default:
		return false
	}
}
