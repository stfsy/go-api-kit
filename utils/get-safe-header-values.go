package utils

import (
	"net/http"
	"net/textproto"
)

// GetSafeHeaderValue inspects the provided HTTP header map for the named
// header and enforces a strict byte-level policy suitable for API tokens and
// other sensitive header values. It returns the first non-empty header value
// and a boolean indicating whether the header (all values) are considered
// "safe".
//
// Safety rules:
//   - Reject ASCII control bytes (0x00..0x1F), DEL (0x7F)
//   - Reject any non-ASCII byte (>= 0x80)
//   - Reject space (0x20) because token-style headers must not contain spaces
//
// Behavior:
//   - If headerName is empty the function is a no-op and returns ("", true).
//   - Empty header values are ignored when selecting the return value but still
//     do not cause rejection.
//   - If any header value contains a disallowed byte the function returns
//     ("", false).
func GetSafeHeaderValue(headerName string, h http.Header) (string, bool) {
	if headerName == "" {
		return "", true
	}

	canonical := textproto.CanonicalMIMEHeaderKey(headerName)
	vals := h[canonical]

	// find the first non-empty value that passes GetSafeValue
	var value string
	for _, hv := range vals {
		if hv == "" {
			continue
		}
		v, ok := GetSafeValue(hv)
		if !ok {
			return "", false
		}
		if value == "" {
			value = v
		}
	}

	return value, true
}

// GetSafeHeaderValues returns all non-empty header values for the named header
// if and only if every value meets the safety rules. It preserves order.
func GetSafeHeaderValues(headerName string, h http.Header) ([]string, bool) {
	if headerName == "" {
		return nil, true
	}

	canonical := textproto.CanonicalMIMEHeaderKey(headerName)
	vals := h[canonical]

	out := make([]string, 0, len(vals))
	for _, hv := range vals {
		if hv == "" {
			continue
		}
		v, ok := GetSafeValue(hv)
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}

	return out, true
}
