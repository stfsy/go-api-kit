package security

func NewXssProtection() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("X-XSS-Protection", "1; mode=block")
}
