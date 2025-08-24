package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stfsy/go-api-kit/server/handlers"
	a "github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestContentType(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name                string
		inputValue          string
		allowedContentTypes string
		want                int
	}{
		{
			"should accept requests with a matching content type",
			"application/json; charset=UTF-8",
			"application/json",
			http.StatusOK,
		},
		{
			"should accept requests with a matching content type no charset",
			"application/json",
			"application/json",
			http.StatusOK,
		},
		{
			"should accept requests with a matching content-type with extra values",
			"application/json; foo=bar; charset=UTF-8; spam=eggs",
			"application/json",
			http.StatusOK,
		},
		{
			"should accept requests with a matching content type when multiple content types are supported",
			"text/xml; charset=UTF-8",
			"text/xml",
			http.StatusOK,
		},
		{
			"should not accept requests with a mismatching content type",
			"text/plain; charset=latin-1",
			"application/json",
			http.StatusUnsupportedMediaType,
		},
		{
			"should not accept requests with a mismatching content type even if multiple content types are allowed",
			"text/plain; charset=Latin-1",
			"text/xml",
			http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			r := negroni.New()
			r.Use(NewRequireContentTypeMiddleware(tt.allowedContentTypes))

			body := []byte("This is my content. There are many like this but this one is mine")
			req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.inputValue)

			r.ServeHTTP(recorder, req)
			res := recorder.Result()

			if res.StatusCode != tt.want {
				t.Errorf("response is incorrect, got %d, want %d", recorder.Code, tt.want)
			}
		})
	}
}

func TestSendNotFound(t *testing.T) {
	assert := a.New(t)

	recorder := httptest.NewRecorder()
	r := negroni.New()
	r.Use(NewRequireContentTypeMiddleware("application/json"))

	body := []byte("This is my content. There are many like this but this one is mine")
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/xml")

	r.ServeHTTP(recorder, req)
	res := recorder.Result()

	assert.Equal(415, res.StatusCode)
	assert.Equal("application/problem+json", res.Header.Get("Content-Type"))

	// Check response body
	var payload handlers.HttpError
	err := json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		assert.Nil(err)
	}

	assert.Equal(415, payload.Status)
	assert.Equal("Unsupported Media Type", payload.Title)
}
