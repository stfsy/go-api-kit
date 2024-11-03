package security

func NewXDownloadOptions() HeaderKeyValueProvider {
	return NewKeyValuePairProvider("X-Download-Options", "noopen")
}
