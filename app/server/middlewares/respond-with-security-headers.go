package middlewares

import (
	"net/http"

	"github.com/stfsy/go-api-kit/app/server/middlewares/security"
)

type SecurityHeadersMiddleware struct{}

var securityHeaders []security.HeaderKeyValue = []security.HeaderKeyValue{
	security.NewCrossOriginEmbedderPolicy().HeaderKeyValue,
	security.NewCrossOriginOpenerPolicy().HeaderKeyValue,
	security.NewCrossOriginResourcePolicy().HeaderKeyValue,
	security.NewOriginAgentClusterPolicy().HeaderKeyValue,
	security.NewReferrerPolicy().HeaderKeyValue,
	security.NewStrictTransportSecurityPolicy().HeaderKeyValue,
	security.NewXContentTypeOptions().HeaderKeyValue,
	security.NewXDownloadOptions().HeaderKeyValue,
	security.NewXFrameOptions().HeaderKeyValue,
	security.NewXPermittedCrossDomainOptions().HeaderKeyValue,
	security.NewXssProtection().HeaderKeyValue,
}

func NewRespondWithSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{}
}

func (m *SecurityHeadersMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	headers := rw.Header()

	for _, sh := range securityHeaders {
		headers.Set(sh.Name, sh.Value)
	}

	next(rw, r)
}
