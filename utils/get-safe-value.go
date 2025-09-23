package utils

// GetSafeValue validates a single string value using the same byte-level
// policy used by header validators. It is independent of HTTP types so it
// can be reused by any package.
//
// Rules:
//  - Empty string is allowed and returns ("", true).
//  - Reject control bytes (0x00..0x1F), DEL (0x7F), non-ASCII bytes (>=0x80)
//  - Reject space (0x20) — intended for token-style values
//  - Reject values with length >= 4096 bytes
func GetSafeValue(s string) (string, bool) {
	if s == "" {
		return "", true
	}

	const maxBytes = 4096
	if len(s) >= maxBytes {
		return "", false
	}

	for i := 0; i < len(s); i++ {
		b := s[i]
		if b <= 31 || b == 127 || b >= 128 {
			return "", false
		}
	}

	return s, true
}
