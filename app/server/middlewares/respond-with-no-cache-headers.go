package middlewares

// Ported from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"net/http"
)

var noCacheHeaders = map[string]string{
	"Cache-Control":     "no-store, no-cache, must-revalidate, proxy-revalidate",
	"Expires":           "0",
	"Pragma":            "no-cache",
	"Surrogate-Control": "no-store",
	"X-Accel-Expires":   "0",
}

type NoCacheHeadersMiddleware struct{}

func NewNoCacheHeadersMiddleware() *NoCacheHeadersMiddleware {
	return &NoCacheHeadersMiddleware{}
}

// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
// a router (or subrouter) from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//
//	Expires: Thu, 01 Jan 1970 00:00:00 UTC
//	Cache-Control: no-cache, private, max-age=0
//	X-Accel-Expires: 0
//	Pragma: no-cache (for HTTP/1.0 proxies/clients)
func (m *NoCacheHeadersMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Set our NoCache headers
	for k, v := range noCacheHeaders {
		rw.Header().Set(k, v)
	}

	next.ServeHTTP(rw, r)
}
