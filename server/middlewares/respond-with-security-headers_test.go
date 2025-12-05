package middlewares

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni/v3"
)

func TestSecurityHeaders(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		key   string
		value string
	}{
		{
			"Cross-Origin-Embedder-Policy",
			"require-corp",
		},
		{
			"Cross-Origin-Opener-Policy",
			"same-origin",
		},
		{
			"Cross-Origin-Resource-Policy",
			"same-site",
		},
		{
			"Origin-Agent-Cluster",
			"?1",
		},
		{
			"Referrer-Policy",
			"same-origin",
		},
		{
			"Strict-Transport-Security",
			"max-age=31536000; includeSubDomains; preload",
		},
		{
			"X-Content-Type-Options",
			"nosniff",
		},
		{
			"X-Download-Options",
			"noopen",
		},
		{
			"X-Frame-Options",
			"DENY",
		},
		{
			"X-Permitted-Cross-Domain-Policies",
			"none",
		},
		{
			"X-XSS-Protection",
			"1; mode=block",
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Has header %s with value %s", tc.key, tc.value), func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			r := negroni.New()
			r.Use(NewRespondWithSecurityHeadersMiddleware())

			body := []byte("")
			req := httptest.NewRequest("POST", "/", bytes.NewReader(body))

			r.ServeHTTP(recorder, req)
			res := recorder.Result()

			if res.Header.Get(tc.key) != tc.value {
				t.Errorf("expected header %s to have value %s but got %s", tc.key, tc.value, res.Header.Get(tc.key))
			}
		})
	}
}
