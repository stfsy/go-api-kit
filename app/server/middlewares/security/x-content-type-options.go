package security

func NewXContentTypeOptions() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("X-Content-Type-Options", "nosniff")
}
