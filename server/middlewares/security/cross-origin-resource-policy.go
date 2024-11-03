package security

func NewCrossOriginResourcePolicy() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("Cross-Origin-Resource-Policy", "same-site")
}
