package middlewares

import (
	"net/http"
	"strings"

	"github.com/stfsy/go-api-kit/server/handlers"
)

// RequireContentLengthOrTransferEncodingMiddleware blocks HTTP/1.1 POST, PATCH, and PUT requests that lack both Content-Length and Transfer-Encoding headers.
type RequireContentLengthOrTransferEncodingMiddleware struct{}

func NewRequireContentLengthOrTransferEncodingMiddleware() *RequireContentLengthOrTransferEncodingMiddleware {
	return &RequireContentLengthOrTransferEncodingMiddleware{}
}

func (m *RequireContentLengthOrTransferEncodingMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	isHttp11 := r.ProtoMajor == 1 && r.ProtoMinor == 1
	if isHttp11 {
		switch r.Method {
		case http.MethodPost, http.MethodPatch, http.MethodPut:
			// Detect presence of a Content-Length header (including Content-Length: 0 per RFC 9110) or Transfer-Encoding.
			clHeader := strings.TrimSpace(r.Header.Get("Content-Length"))
			hasContentLength := clHeader != "" || r.ContentLength > 0
			hasTransferEncoding := len(r.TransferEncoding) > 0 || strings.TrimSpace(r.Header.Get("Transfer-Encoding")) != ""
			if !hasContentLength && !hasTransferEncoding {
				handlers.SendLengthRequired(rw, nil)
				return
			}
		}
	}

	next.ServeHTTP(rw, r)
}
