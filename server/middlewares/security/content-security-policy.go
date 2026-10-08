package security

func NewContentSecurityPolicy() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
}
