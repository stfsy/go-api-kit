package security

func NewXPermittedCrossDomainOptions() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("X-Permitted-Cross-Domain-Policies", "none")
}
