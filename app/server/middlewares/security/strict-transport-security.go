package security

func NewStrictTransportSecurityPolicy() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
}
