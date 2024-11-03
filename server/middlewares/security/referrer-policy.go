package security

func NewReferrerPolicy() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("Referrer-Policy", "same-origin")
}
