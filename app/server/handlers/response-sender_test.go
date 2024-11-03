package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestSendJson(t *testing.T) {
	assert := a.New(t)

	recorder := httptest.NewRecorder()
	SendJson(recorder, []byte(`{ "status": 200, "title": "Thank you" }`))
	res := recorder.Result()

	assert.Equal(200, res.StatusCode)
	assert.Equal("application/json", res.Header.Get("Content-Type"))

	// Check response body
	var payload HttpError
	err := json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		assert.Nil(err)
	}

	assert.Equal(200, payload.Status)
	assert.Equal("Thank you", payload.Title)
}

func TestSendText(t *testing.T) {
	assert := a.New(t)

	recorder := httptest.NewRecorder()
	SendText(recorder, "Ok")
	res := recorder.Result()

	assert.Equal(200, res.StatusCode)
	assert.Equal("text/plain; charset=utf-8", res.Header.Get("Content-Type"))

	body, _ := io.ReadAll(res.Body)
	assert.Equal([]byte("Ok"), body)
}
