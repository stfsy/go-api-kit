package handlers

import (
	"net/http"
)

func MethodNotAllowedHandler(w http.ResponseWriter, _r *http.Request) {
	SendMethodNotAllowed(w, nil)
}
