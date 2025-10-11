package handlers

func CreateErrorDetails(key, value string) ErrorDetails {
	return ErrorDetails{
		key: ErrorDetail{
			Message: value,
		},
	}
}

func CreateMustNotBeUndefinedErrorDetail(field string) ErrorDetails {
	return CreateErrorDetails(field, "must not be undefined")
}
