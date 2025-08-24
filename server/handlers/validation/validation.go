package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// FieldErrorDetail represents the tag and message for a field error
type FieldErrorDetail struct {
	Validator string `json:"validator"`
	Message   string `json:"message"`
}

// ValidationErrorResponse represents the response for validation errors (grouped by field)
type ValidationErrorResponse struct {
	Title  string                      `json:"title"`
	Status int                         `json:"status"`
	Errors map[string]FieldErrorDetail `json:"errors"`
}

// ValidateStruct validates a struct and returns a map of field errors
func ValidateStruct(s interface{}) map[string]FieldErrorDetail {
	errors := make(map[string]FieldErrorDetail)
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fieldMap := buildJSONFieldMap(t, "", "")

	err := validate.Struct(s)
	if err != nil {
		for _, ferr := range err.(validator.ValidationErrors) {
			// Always use StructNamespace for lookup, which is dot-separated path
			ns := ferr.StructNamespace()              // e.g. Address.City
			ns = strings.TrimPrefix(ns, t.Name()+".") // Remove root struct name
			key := fieldMap[ns]
			// If not found, fallback to just the field name (for root fields)
			if key == "" {
				key = strings.ToLower(ferr.Field())
			}
			// Normalize to dot notation (e.g. address.city)
			key = strings.ReplaceAll(key, ".", ".")
			errors[key] = FieldErrorDetail{
				Validator: ferr.Tag(),
				Message:   getErrorMessage(ferr),
			}
		}
	}

	return errors
}

// getErrorMessage returns a human-readable error message for validation errors
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	// Required and length
	case "required":
		return fmt.Sprintf("%s must not be undefined", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", fe.Field(), fe.Param())

	// Comparisons
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	case "nefield":
		return fmt.Sprintf("%s must not be equal to %s", fe.Field(), fe.Param())
	case "gtfield":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "gtefield":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "ltfield":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "ltefield":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fe.Field(), fe.Param())

	// String types
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", fe.Field())
	case "alphanumunicode":
		return fmt.Sprintf("%s must contain only alphanumeric characters and spaces", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())

	// Others
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", fe.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "uuid3":
		return fmt.Sprintf("%s must be a valid UUIDv3", fe.Field())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUIDv4", fe.Field())
	case "uuid5":
		return fmt.Sprintf("%s must be a valid UUIDv5", fe.Field())
	case "isbn":
		return fmt.Sprintf("%s must be a valid ISBN", fe.Field())
	case "isbn10":
		return fmt.Sprintf("%s must be a valid ISBN-10", fe.Field())
	case "isbn13":
		return fmt.Sprintf("%s must be a valid ISBN-13", fe.Field())
	case "contains":
		return fmt.Sprintf("%s must contain '%s'", fe.Field(), fe.Param())
	case "excludes":
		return fmt.Sprintf("%s must not contain '%s'", fe.Field(), fe.Param())
	case "startswith":
		return fmt.Sprintf("%s must start with '%s'", fe.Field(), fe.Param())
	case "endswith":
		return fmt.Sprintf("%s must end with '%s'", fe.Field(), fe.Param())
	case "ip":
		return fmt.Sprintf("%s must be a valid IP address", fe.Field())
	case "ipv4":
		return fmt.Sprintf("%s must be a valid IPv4 address", fe.Field())
	case "ipv6":
		return fmt.Sprintf("%s must be a valid IPv6 address", fe.Field())
	case "mac":
		return fmt.Sprintf("%s must be a valid MAC address", fe.Field())
	case "cidr":
		return fmt.Sprintf("%s must be a valid CIDR notation", fe.Field())
	case "cidrv4":
		return fmt.Sprintf("%s must be a valid CIDR notation (IPv4)", fe.Field())
	case "cidrv6":
		return fmt.Sprintf("%s must be a valid CIDR notation (IPv6)", fe.Field())
	case "dive":
		return fmt.Sprintf("%s must have valid items only", fe.Field())

	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
