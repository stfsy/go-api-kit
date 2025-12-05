package middlewares

import (
	"bytes"
	"net/http/httptest"
	"testing"

	a "github.com/stretchr/testify/assert"
	"github.com/urfave/negroni/v3"
)

func TestNoCacheResponseHeaders(t *testing.T) {
	assert := a.New(t)

	recorder := httptest.NewRecorder()

	r := negroni.New()
	r.Use(NewNoCacheHeadersMiddleware())

	body := []byte("This is my content. There are many like this but this one is mine")
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	r.ServeHTTP(recorder, req)
	res := recorder.Result()

	assert.NotEmpty(res.Header.Get("Cache-Control"))
	assert.NotEmpty(res.Header.Get("Expires"))
	assert.NotEmpty(res.Header.Get("Pragma"))
	assert.NotEmpty(res.Header.Get("Surrogate-Control"))
	assert.NotEmpty(res.Header.Get("X-Accel-Expires"))
}
