package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestRequireMaxBodyLengthMiddleware(t *testing.T) {
	const maxSize = 1024

	// Create middleware with 1MB limit
	mw := &RequireMaxBodyLengthMiddleware{maxSize: maxSize}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body) // Read the full request body

		if err != nil {
			w.WriteHeader(413)
		}

		w.WriteHeader(http.StatusOK)
	})

	n := negroni.New()
	n.Use(mw)
	n.UseHandler(handler)

	t.Run("should accept request under limit", func(t *testing.T) {
		body := bytes.Repeat([]byte("a"), maxSize-1)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		n.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should reject request over limit", func(t *testing.T) {
		body := bytes.Repeat([]byte("a"), maxSize+1)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		n.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	})
}
