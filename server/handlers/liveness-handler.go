package handlers

import (
	"net/http"
)

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	SendText(w, "Ok")
}
