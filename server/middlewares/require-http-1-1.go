package middlewares

import (
	"net/http"

	"github.com/stfsy/go-api-kit/server/handlers"
)

// RequireHTTP11Middleware blocks requests that are not using HTTP/1.1.
type RequireHTTP11Middleware struct{}

func NewRequireHTTP11Middleware() *RequireHTTP11Middleware {
	return &RequireHTTP11Middleware{}
}

func (m *RequireHTTP11Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	isHttp11 := r.ProtoMajor == 1 && r.ProtoMinor == 1
	if !isHttp11 {
		handlers.SendBadRequest(rw, nil)
		return
	}
	next.ServeHTTP(rw, r)
}
