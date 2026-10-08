package security

func NewXssProtection() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("X-XSS-Protection", "0")
}
