package security

func NewCrossOriginOpenerPolicy() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("Cross-Origin-Opener-Policy", "same-origin")
}
