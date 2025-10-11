package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
			name:         "POST HTTP/1.1 with zero Content-Length",
			method:       http.MethodPost,
			proto:        "HTTP/1.1",
			headers:      map[string]string{"Content-Length": "0"},
			expectStatus: http.StatusLengthRequired,
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

			// Prepare body reader based on headers instructions. For Content-Length
			// cases we create a body with the requested size so http.NewRequest will
			// set the Content-Length header automatically. For Transfer-Encoding
			// cases we provide a small body and set the TransferEncoding slice.
			var bodyReader io.Reader
			// Track whether we should explicitly set an empty header (special case)
			setEmptyContentLength := false

			setContentLengthHeader := false
			if v, ok := tc.headers["Content-Length"]; ok {
				if v == "" {
					// explicit empty header case — keep nil body but set header later
					setEmptyContentLength = true
				} else {
					// parse number of bytes and create a body of that size
					// if value is "0" we create an empty reader so Content-Length: 0
					// will be set by NewRequest
					// ignore parse errors — tests supply valid integers
					// use bytes.Repeat to create predictable content
					// note: strconv.Atoi is lightweight and safe here
					// to avoid importing strconv repeatedly, use fmt.Sscanf
					var n int
					_, _ = fmt.Sscanf(v, "%d", &n)
					if n == 0 {
						bodyReader = bytes.NewReader([]byte{})
						// ensure header presence for the explicit zero-length test
						setContentLengthHeader = true
					} else {
						bodyReader = bytes.NewReader(make([]byte, n))
					}
				}
			}

			// For Transfer-Encoding cases provide a small body if none set already
			if v, ok := tc.headers["Transfer-Encoding"]; ok {
				if v != "" && bodyReader == nil {
					bodyReader = bytes.NewReader([]byte("ok"))
				}
			}

			r, _ := http.NewRequest(tc.method, "/", bodyReader)
			r.Proto = tc.proto
			r.ProtoMajor = int(tc.proto[5] - '0')
			r.ProtoMinor = int(tc.proto[7] - '0')

			for k, v := range tc.headers {
				switch k {
				case "Content-Length":
					if setEmptyContentLength {
						// explicitly set empty header to simulate malformed input
						r.Header.Set("Content-Length", "")
					}
					if setContentLengthHeader {
						r.Header.Set("Content-Length", tc.headers["Content-Length"])
					}
				case "Transfer-Encoding":
					if v == "" {
						r.TransferEncoding = nil
						continue
					}
					// Set TransferEncoding slice directly (parsed form)
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
