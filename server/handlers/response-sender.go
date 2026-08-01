package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"

	ContentTypeText        = "text/plain; charset=utf-8"
	ContentTypeJson        = "application/json"
	ContentTypeProblemJson = "application/problem+json"
)

func send(rw http.ResponseWriter, response []byte, status int) bool {
	rw.WriteHeader(status)
	_, err := rw.Write(response)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to send response to stream %s", err.Error()))
		return false
	}

	return true
}

func SendText(rw http.ResponseWriter, response string) {
	rw.Header().Set(HeaderContentType, ContentTypeText)
	_ = send(rw, []byte(response), http.StatusOK)
}

func SendJson(rw http.ResponseWriter, response []byte) {
	rw.Header().Set(HeaderContentType, ContentTypeJson)
	_ = send(rw, response, http.StatusOK)
}

// SendStructAsJson streams v as JSON directly to rw, avoiding the extra buffer
// allocation and copy that json.Marshal would incur for the response body.
func SendStructAsJson(rw http.ResponseWriter, v interface{}) {
	rw.Header().Set(HeaderContentType, ContentTypeJson)
	err := json.NewEncoder(rw).Encode(v)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to encode response as JSON %s", err.Error()))
		SendInternalServerError(rw, nil)
	}
}
