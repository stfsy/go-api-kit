package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireHTTP11Middleware(t *testing.T) {
	mw := NewRequireHTTP11Middleware()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	cases := []struct {
		name         string
		proto        string
		expectStatus int
	}{
		{
			name:         "HTTP/1.1 allowed",
			proto:        "HTTP/1.1",
			expectStatus: http.StatusOK,
		},
		{
			name:         "HTTP/1.0 blocked",
			proto:        "HTTP/1.0",
			expectStatus: http.StatusBadRequest,
		},
		{
			name:         "HTTP/2.0 blocked",
			proto:        "HTTP/2.0",
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.Proto = tc.proto
			r.ProtoMajor = int(tc.proto[5] - '0')
			r.ProtoMinor = int(tc.proto[7] - '0')
			mw.ServeHTTP(rec, r, handler.ServeHTTP)
			assert.Equal(t, tc.expectStatus, rec.Result().StatusCode)
		})
	}
}
