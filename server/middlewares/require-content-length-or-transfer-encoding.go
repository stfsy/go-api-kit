package middlewares

import (
	"net/http"

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
			// detect presence of the Content-Length header (including zero)
			hasContentLength := r.Header.Get("Content-Length") != ""
			hasTransferEncoding := len(r.TransferEncoding) > 0
			if !hasContentLength && !hasTransferEncoding {
				handlers.SendLengthRequired(rw, nil)
				return
			}
		}
	}

	next.ServeHTTP(rw, r)
}
