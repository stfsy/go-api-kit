package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stfsy/go-api-kit/server/handlers/validation"
)

func ValidatingHandler[T any](handler func(http.ResponseWriter, *http.Request, *T)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
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
			// Optionally add validation logic here

		default:
			// For methods that typically don't have a body, just call the handler with Nil
			handler(w, r, nil)
		}
	}
}
