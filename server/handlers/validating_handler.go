package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stfsy/go-api-kit/server/handlers/validation"
)

var EMPTY_MAP = make(map[string]map[string]string)

func ValidatingHandler[T any](handler func(http.ResponseWriter, *http.Request, *T)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			var body T
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&body); err != nil {
				SendBadRequest(w, EMPTY_MAP)
				return
			}

			errors := validation.ValidateStruct(&body)
			if len(errors) != 0 {
				// convert errors to ErrorDetails
				errorDetails := make(ErrorDetails, len(errors))

				for field, errDetail := range errors {
					fieldErrors := make(map[string]string, 1)
					fieldErrors["Message"] = errDetail.Message
					fieldErrors["Validator"] = errDetail.Validator
					errorDetails[field] = fieldErrors
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
