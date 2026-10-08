package middlewares

import (
	"net/http"

	"github.com/stfsy/go-api-kit/server/middlewares/security"
)

type SecurityHeadersMiddleware struct{}

type securityHeaderEntry struct {
	canonicalKey string
	values       []string
}

var precomputedSecurityHeaders = func() []securityHeaderEntry {
	providers := []security.HeaderKeyValueProvider{
		security.NewContentSecurityPolicy(),
		security.NewCrossOriginEmbedderPolicy(),
		security.NewCrossOriginOpenerPolicy(),
		security.NewCrossOriginResourcePolicy(),
		security.NewOriginAgentClusterPolicy(),
		security.NewReferrerPolicy(),
		security.NewStrictTransportSecurityPolicy(),
		security.NewXContentTypeOptions(),
		security.NewXDownloadOptions(),
		security.NewXFrameOptions(),
		security.NewXPermittedCrossDomainOptions(),
		security.NewXssProtection(),
	}
	entries := make([]securityHeaderEntry, len(providers))
	for i, p := range providers {
		entries[i] = securityHeaderEntry{
			canonicalKey: http.CanonicalHeaderKey(p.Name),
			values:       []string{p.Value},
		}
	}
	return entries
}()

func NewRespondWithSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{}
}

func (m *SecurityHeadersMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	headers := rw.Header()
	for i := range precomputedSecurityHeaders {
		headers[precomputedSecurityHeaders[i].canonicalKey] = precomputedSecurityHeaders[i].values
	}
	next(rw, r)
}
