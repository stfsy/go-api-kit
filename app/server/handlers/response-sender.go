package handlers

import (
	"fmt"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"

	ContentTypeText = "text/plain"
	ContentTypeJson = "application/json"
)

func send(rw http.ResponseWriter, response []byte) bool {
	_, err := rw.Write(response)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to send response to stream %s", err.Error()))
		return false
	}

	return true
}

func SendText(rw http.ResponseWriter, response string) {
	_ = send(rw, []byte(response))
}

func SendJson(rw http.ResponseWriter, response []byte) {
	rw.Header().Set(HeaderContentType, ContentTypeJson)
	_ = send(rw, response)
}
