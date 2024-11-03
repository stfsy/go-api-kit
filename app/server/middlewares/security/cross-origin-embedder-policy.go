package security

type HeaderKeyValue struct {
	Name  string
	Value string
}

type CrossOriginEmbedderPolicy struct{}

func NewCrossOriginEmbedderPolicy() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("Cross-Origin-Embedder-Policy", "require-corp")
}
