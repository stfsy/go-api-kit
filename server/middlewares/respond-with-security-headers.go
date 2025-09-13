package middlewares

import (
	"net/http"

	"github.com/stfsy/go-api-kit/server/middlewares/security"
)

type SecurityHeadersMiddleware struct{}

var securityHeadersMap = func() map[string]string {
	m := make(map[string]string, 11)
	m[security.NewCrossOriginEmbedderPolicy().Name] = security.NewCrossOriginEmbedderPolicy().Value
	m[security.NewCrossOriginOpenerPolicy().Name] = security.NewCrossOriginOpenerPolicy().Value
	m[security.NewCrossOriginResourcePolicy().Name] = security.NewCrossOriginResourcePolicy().Value
	m[security.NewOriginAgentClusterPolicy().Name] = security.NewOriginAgentClusterPolicy().Value
	m[security.NewReferrerPolicy().Name] = security.NewReferrerPolicy().Value
	m[security.NewStrictTransportSecurityPolicy().Name] = security.NewStrictTransportSecurityPolicy().Value
	m[security.NewXContentTypeOptions().Name] = security.NewXContentTypeOptions().Value
	m[security.NewXDownloadOptions().Name] = security.NewXDownloadOptions().Value
	m[security.NewXFrameOptions().Name] = security.NewXFrameOptions().Value
	m[security.NewXPermittedCrossDomainOptions().Name] = security.NewXPermittedCrossDomainOptions().Value
	m[security.NewXssProtection().Name] = security.NewXssProtection().Value
	return m
}()

func NewRespondWithSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{}
}

func (m *SecurityHeadersMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	headers := rw.Header()
	for k, v := range securityHeadersMap {
		headers.Set(k, v)
	}
	next(rw, r)
}
