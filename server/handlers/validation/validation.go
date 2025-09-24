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

// ValidateStruct validates a struct and returns a map of field errors
func ValidateStruct(s interface{}) map[string]FieldErrorDetail {
	errors := make(map[string]FieldErrorDetail)
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldMap := GetOrBuildFieldMap(t, "", "")

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
		return "must not be undefined"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters long", fe.Param())

	// Comparisons
	case "eq":
		return fmt.Sprintf("must be equal to %s", fe.Param())
	case "ne":
		return fmt.Sprintf("must not be equal to %s", fe.Param())
	case "lt":
		return fmt.Sprintf("must be less than %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "eqfield":
		return fmt.Sprintf("must be equal to %s", fe.Param())
	case "nefield":
		return fmt.Sprintf("must not be equal to %s", fe.Param())
	case "gtfield":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "gtefield":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "ltfield":
		return fmt.Sprintf("must be less than %s", fe.Param())
	case "ltefield":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", fe.Param())

	// String types
	case "alpha":
		return "must contain only alphabetic characters"
	case "alphanum":
		return "must contain only alphanumeric characters"
	case "alphanumunicode":
		return "must contain only alphanumeric characters and spaces"
	case "email":
		return "must be a valid email address"

	// Others
	case "url":
		return "must be a valid URL"
	case "uri":
		return "must be a valid URI"
	case "uuid":
		return "must be a valid UUID"
	case "uuid3":
		return "must be a valid UUIDv3"
	case "uuid4":
		return "must be a valid UUIDv4"
	case "uuid5":
		return "must be a valid UUIDv5"
	case "isbn":
		return "must be a valid ISBN"
	case "isbn10":
		return "must be a valid ISBN-10"
	case "isbn13":
		return "must be a valid ISBN-13"
	case "contains":
		return fmt.Sprintf("must contain '%s'", fe.Param())
	case "excludes":
		return fmt.Sprintf("must not contain '%s'", fe.Param())
	case "startswith":
		return fmt.Sprintf("must start with '%s'", fe.Param())
	case "endswith":
		return fmt.Sprintf("must end with '%s'", fe.Param())
	case "ip":
		return "must be a valid IP address"
	case "ipv4":
		return "must be a valid IPv4 address"
	case "ipv6":
		return "must be a valid IPv6 address"
	case "mac":
		return "must be a valid MAC address"
	case "cidr":
		return "must be a valid CIDR notation"
	case "cidrv4":
		return "must be a valid CIDR notation (IPv4)"
	case "cidrv6":
		return "must be a valid CIDR notation (IPv6)"
	case "dive":
		return "must have valid items only"

	default:
		return "is invalid"
	}
}
