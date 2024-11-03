package security

type HeaderKeyValueProvider struct {
	HeaderKeyValue
}

func NewKeyValuePairProvider(name string, value string) HeaderKeyValueProvider {
	return HeaderKeyValueProvider{
		HeaderKeyValue: HeaderKeyValue{
			Name:  name,
			Value: value,
		},
	}
}
