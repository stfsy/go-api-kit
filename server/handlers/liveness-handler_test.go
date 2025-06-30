package handlers

import (
	"net/http/httptest"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestLivenessHandler(t *testing.T) {
	assert := a.New(t)

	recorder := httptest.NewRecorder()
	LivenessHandler(recorder, nil)
	res := recorder.Result()

	assert.Equal(200, res.StatusCode)
	assert.Contains(res.Header.Get("Content-Type"), "text/plain")
	assert.Equal("Ok", recorder.Body.String())
}
