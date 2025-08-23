package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidatingHandler[T any](handler func(http.ResponseWriter, *http.Request, *T)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			var body T
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&body); err != nil {
				SendBadRequest(w, map[string]string{})
				return
			}

			handler(w, r, &body)
			// Optionally add validation logic here

		default:
			// For methods that typically don't have a body, just call the handler with Nil
			handler(w, r, nil)
		}
	}
}
