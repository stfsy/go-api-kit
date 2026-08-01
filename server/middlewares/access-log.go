package middlewares

import (
	"net/http"
	"time"

	"github.com/stfsy/go-api-kit/utils"
	"github.com/urfave/negroni/v3"
)

var logger = utils.NewLogger("access-log-middleware")

// AccessLogMiddleware logs each request via the shared structured logger instead of
// negroni's default Logger, which renders a text/template and takes a global mutex
// on every request.
type AccessLogMiddleware struct{}

// NewAccessLog returns a new AccessLogMiddleware instance
func NewAccessLog() *AccessLogMiddleware {
	return &AccessLogMiddleware{}
}

func (m *AccessLogMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	status := 0
	res, ok := rw.(negroni.ResponseWriter)
	if ok {
		status = res.Status()
	}

	logger.Info("request",
		"method", r.Method,
		"path", r.URL.Path,
		"proto", r.Proto,
		"status", status,
		"duration", time.Since(start),
		"user_agent", r.UserAgent(),
	)
}
