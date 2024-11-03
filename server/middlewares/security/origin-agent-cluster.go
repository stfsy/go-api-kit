package security

func NewOriginAgentClusterPolicy() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("Origin-Agent-Cluster", "?1")
}
