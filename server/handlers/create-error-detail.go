package handlers

func CreateErrorDetail(key, value string) ErrorDetails {
	return ErrorDetails{
		key: ErrorDetail{
			Message: value,
		},
	}
}

func CreateMustNotBeUndefinedErrorDetail(field string) ErrorDetails {
	return CreateErrorDetail(field, "must not be undefined")
}
