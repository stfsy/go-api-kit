package handlers

func CreateErrorDetail(key, value string) ErrorDetails {
	return ErrorDetails{
		key: map[string]string{"message": value},
	}
}

func CreateMustNotBeUndefinedErrorDetail(field string) ErrorDetails {
	return CreateErrorDetail(field, "must not be undefined")
}
