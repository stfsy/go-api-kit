package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireContentLengthOrTransferEncodingMiddleware(t *testing.T) {
	mw := NewRequireContentLengthOrTransferEncodingMiddleware()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	cases := []struct {
		name         string
		method       string
		proto        string
		headers      map[string]string
		expectStatus int
	}{
		{
			name:         "POST HTTP/1.1 missing both headers",
			method:       http.MethodPost,
			proto:        "HTTP/1.1",
			headers:      map[string]string{},
			expectStatus: http.StatusLengthRequired,
		},
		{
			name:         "POST HTTP/1.1 with Content-Length",
			method:       http.MethodPost,
			proto:        "HTTP/1.1",
			headers:      map[string]string{"Content-Length": "10"},
			expectStatus: http.StatusOK,
		},
		{
			name:         "POST HTTP/1.1 with Transfer-Encoding",
			method:       http.MethodPost,
			proto:        "HTTP/1.1",
			headers:      map[string]string{"Transfer-Encoding": "chunked"},
			expectStatus: http.StatusOK,
		},
		{
			name:         "POST HTTP/1.1 with empty Content-Length and Transfer-Encoding",
			method:       http.MethodPost,
			proto:        "HTTP/1.1",
			headers:      map[string]string{"Content-Length": "", "Transfer-Encoding": ""},
			expectStatus: http.StatusLengthRequired,
		},
		{
			name:         "PATCH HTTP/1.1 missing both headers",
			method:       http.MethodPatch,
			proto:        "HTTP/1.1",
			headers:      map[string]string{},
			expectStatus: http.StatusLengthRequired,
		},
		{
			name:         "PATCH HTTP/1.1 with Content-Length",
			method:       http.MethodPatch,
			proto:        "HTTP/1.1",
			headers:      map[string]string{"Content-Length": "5"},
			expectStatus: http.StatusOK,
		},
		{
			name:         "PUT HTTP/1.1 missing both headers",
			method:       http.MethodPut,
			proto:        "HTTP/1.1",
			headers:      map[string]string{},
			expectStatus: http.StatusLengthRequired,
		},
		{
			name:         "PUT HTTP/1.1 with Transfer-Encoding",
			method:       http.MethodPut,
			proto:        "HTTP/1.1",
			headers:      map[string]string{"Transfer-Encoding": "chunked"},
			expectStatus: http.StatusOK,
		},
		{
			name:         "POST HTTP/1.0 missing headers (should allow)",
			method:       http.MethodPost,
			proto:        "HTTP/1.0",
			headers:      map[string]string{},
			expectStatus: http.StatusOK,
		},
		{
			name:         "GET HTTP/1.1 missing headers (should allow)",
			method:       http.MethodGet,
			proto:        "HTTP/1.1",
			headers:      map[string]string{},
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r, _ := http.NewRequest(tc.method, "/", nil)
			r.Proto = tc.proto
			r.ProtoMajor = int(tc.proto[5] - '0')
			r.ProtoMinor = int(tc.proto[7] - '0')
			for k, v := range tc.headers {
				switch k {
				case "Content-Length":
					if v == "" {
						// leave default (0)
						continue
					}
					if n, err := strconv.ParseInt(v, 10, 64); err == nil {
						r.ContentLength = n
					}
				case "Transfer-Encoding":
					if v == "" {
						r.TransferEncoding = nil
						continue
					}
					r.TransferEncoding = []string{v}
				default:
					r.Header.Set(k, v)
				}
			}
			mw.ServeHTTP(rec, r, handler.ServeHTTP)
			assert.Equal(t, tc.expectStatus, rec.Result().StatusCode)
		})
	}
}
