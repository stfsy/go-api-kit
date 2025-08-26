package middlewares

import (
	"net/http"

	"github.com/stfsy/go-api-kit/server/handlers"
)

// RequireContentLengthOrTransferEncodingMiddleware blocks HTTP/1.1 POST requests that lack both Content-Length and Transfer-Encoding headers.
type RequireContentLengthOrTransferEncodingMiddleware struct{}

func NewRequireContentLengthOrTransferEncodingMiddleware() *RequireContentLengthOrTransferEncodingMiddleware {
	return &RequireContentLengthOrTransferEncodingMiddleware{}
}

func (m *RequireContentLengthOrTransferEncodingMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == http.MethodPost && r.ProtoMajor == 1 && r.ProtoMinor == 1 {
		_, hasContentLength := r.Header["Content-Length"]
		_, hasTransferEncoding := r.Header["Transfer-Encoding"]
		if !hasContentLength && !hasTransferEncoding {
			handlers.SendLengthRequired(rw, nil)
			return
		}
	}
	next.ServeHTTP(rw, r)
}
