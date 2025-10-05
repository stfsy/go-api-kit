package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stfsy/go-api-kit/server/handlers/validation"
)

func ValidatingHandler[T any](handler func(http.ResponseWriter, *http.Request, *T)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		contentLength := r.ContentLength
		// For methods that typically don't have a body, just call the handler with Nil
		if (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete) && contentLength > 0 {
			var body T
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()

			if err := decoder.Decode(&body); err != nil {
				SendBadRequest(w, nil)
				return
			}

			errors := validation.ValidateStruct(&body)
			if len(errors) != 0 {
				// convert errors to ErrorDetails
				errorDetails := make(ErrorDetails, len(errors))

				for field, errDetail := range errors {
					errorDetails[field] = ErrorDetail{Message: errDetail.Message}
				}
				SendValidationError(w, errorDetails)
				return
			}

			handler(w, r, &body)
			return
			// Optionally add validation logic here
		}

		handler(w, r, nil)
	}
}
