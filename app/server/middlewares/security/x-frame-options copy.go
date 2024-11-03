package security

func NewXFrameOptions() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("X-Frame-Options", "DENY")
}
