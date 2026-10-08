package middlewares

// Ported from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"net/http"
)

type noCacheHeaderEntry struct {
	canonicalKey string
	values       []string
}

var precomputedNoCacheHeaders = []noCacheHeaderEntry{
	{canonicalKey: "Cache-Control", values: []string{"no-store, no-cache, must-revalidate, proxy-revalidate"}},
	{canonicalKey: "Expires", values: []string{"0"}},
	{canonicalKey: "Pragma", values: []string{"no-cache"}},
	{canonicalKey: "Surrogate-Control", values: []string{"no-store"}},
	{canonicalKey: "X-Accel-Expires", values: []string{"0"}},
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

	// Set our NoCache headers using precomputed canonical keys and slice values to avoid allocations.
	headers := rw.Header()
	for i := range precomputedNoCacheHeaders {
		headers[precomputedNoCacheHeaders[i].canonicalKey] = precomputedNoCacheHeaders[i].values
	}

	next.ServeHTTP(rw, r)
}
