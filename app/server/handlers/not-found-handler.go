package handlers

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, _r *http.Request) {
	SendNotFound(w, nil)
}
