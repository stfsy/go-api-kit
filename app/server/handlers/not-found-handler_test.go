package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestSendNotFound(t *testing.T) {
	assert := a.New(t)

	recorder := httptest.NewRecorder()
	NotFoundHandler(recorder, nil)
	res := recorder.Result()

	assert.Equal(404, res.StatusCode)
	assert.Equal("application/json", res.Header.Get("Content-Type"))

	// Check response body
	var payload HttpError
	err := json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		assert.Nil(err)
	}

	assert.Equal(404, payload.Status)
	assert.Equal("Not Found", payload.Title)
}
