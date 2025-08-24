package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func GenericTest(t *testing.T, fn func(rw http.ResponseWriter, details map[string]string), status int, title string) {
	t.Parallel()
	assert := a.New(t)

	var tests = []struct {
		title   string
		status  int
		details map[string]string
	}{
		{
			title,
			status,
			nil,
		},
		{
			title,
			status,
			map[string]string(nil),
		},
		{
			title,
			status,
			map[string]string{
				"api_key": "not allowed to be empty",
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Sets status %d and marshalls %v", tc.status, tc.details), func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			fn(recorder, tc.details)
			res := recorder.Result()

			assert.Equal(tc.status, res.StatusCode)
			assert.Equal("application/problem+json", res.Header.Get("Content-Type"))

			// Check response body
			var resp HttpError
			err := json.NewDecoder(res.Body).Decode(&resp)

			assert.Nil(err)
			assert.Equal(tc.title, resp.Title)
			assert.Equal(tc.status, resp.Status)
			assert.Equal(tc.details, resp.Details)
		})
	}
}

func TestBadRequest(t *testing.T) {
	GenericTest(t, SendBadRequest, http.StatusBadRequest, "Bad Request")
}

func TestUnauthorized(t *testing.T) {
	GenericTest(t, SendUnauthorized, http.StatusUnauthorized, "Unauthorized")
}

func TestForbidden(t *testing.T) {
	GenericTest(t, SendForbidden, http.StatusForbidden, "Forbidden")
}

func TestNotFound(t *testing.T) {
	GenericTest(t, SendNotFound, http.StatusNotFound, "Not Found")
}

func TestMethodNotAllowed(t *testing.T) {
	GenericTest(t, SendMethodNotAllowed, http.StatusMethodNotAllowed, "Method Not Allowed")
}

func TestNotAcceptable(t *testing.T) {
	GenericTest(t, SendNotAcceptable, http.StatusNotAcceptable, "Not Acceptable")
}

func TestRequestTimeout(t *testing.T) {
	GenericTest(t, SendRequestTimeout, http.StatusRequestTimeout, "Request Timeout")
}

func TestConflict(t *testing.T) {
	GenericTest(t, SendConflict, http.StatusConflict, "Conflict")
}

func TestGone(t *testing.T) {
	GenericTest(t, SendGone, http.StatusGone, "Gone")
}

func TestLengthRequired(t *testing.T) {
	GenericTest(t, SendLengthRequired, http.StatusLengthRequired, "Length Required")
}

func TestPreconditionFailed(t *testing.T) {
	GenericTest(t, SendPreconditionFailed, http.StatusPreconditionFailed, "Precondition Failed")
}

func TestPayloadTooLarge(t *testing.T) {
	GenericTest(t, SendPayloadTooLarge, http.StatusRequestEntityTooLarge, "Payload Too Large")
}

func TestURITooLong(t *testing.T) {
	GenericTest(t, SendURITooLong, http.StatusRequestURITooLong, "URI Too Long")
}

func TestUnsupportedMediaType(t *testing.T) {
	GenericTest(t, SendUnsupportedMediaType, http.StatusUnsupportedMediaType, "Unsupported Media Type")
}

func TestRangeNotSatisfiable(t *testing.T) {
	GenericTest(t, SendRangeNotSatisfiable, http.StatusRequestedRangeNotSatisfiable, "Range Not Satisfiable")
}

func TestExpectationFailed(t *testing.T) {
	GenericTest(t, SendExpectationFailed, http.StatusExpectationFailed, "Expectation Failed")
}

func TestUnprocessableEntity(t *testing.T) {
	GenericTest(t, SendUnprocessableEntity, http.StatusUnprocessableEntity, "Unprocessable Entity")
}

func TestTooManyRequests(t *testing.T) {
	GenericTest(t, SendTooManyRequests, http.StatusTooManyRequests, "Too Many Requests")
}

func TestInternalServerError(t *testing.T) {
	GenericTest(t, SendInternalServerError, http.StatusInternalServerError, "Internal Server Error")
}

func TestNotImplemented(t *testing.T) {
	GenericTest(t, SendNotImplemented, http.StatusNotImplemented, "Not Implemented")
}

func TestBadGateway(t *testing.T) {
	GenericTest(t, SendBadGateway, http.StatusBadGateway, "Bad Gateway")
}

func TestServiceUnavailable(t *testing.T) {
	GenericTest(t, SendServiceUnavailable, http.StatusServiceUnavailable, "Service Unavailable")
}

func TestGatewayTimeout(t *testing.T) {
	GenericTest(t, SendGatewayTimeout, http.StatusGatewayTimeout, "Gateway Timeout")
}

func TestHTTPVersionNotSupported(t *testing.T) {
	GenericTest(t, SendHTTPVersionNotSupported, http.StatusHTTPVersionNotSupported, "HTTP Version Not Supported")
}
