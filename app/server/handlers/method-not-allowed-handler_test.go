package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestMethodNotAllowedHandler(t *testing.T) {
	assert := a.New(t)

	recorder := httptest.NewRecorder()
	MethodNotAllowedHandler(recorder, nil)
	res := recorder.Result()

	assert.Equal(405, res.StatusCode)
	assert.Equal("application/json", res.Header.Get("Content-Type"))

	// Check response body
	var payload HttpError
	err := json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		assert.Nil(err)
	}

	assert.Equal(405, payload.Status)
	assert.Equal("Method Not Allowed", payload.Title)
}
