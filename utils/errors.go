package utils

type CredentialsError struct {
	Message string
}

type ConfirmEmailError struct {
	Message string
}

func (e *CredentialsError) Error() string {
	return e.Message
}

func (e *ConfirmEmailError) Error() string {
	return e.Message
}

func NewCustomError(errorType string, message string) error {
	switch errorType {
	case "credentials":
		return &CredentialsError{
			Message: message,
		}
	case "confirm_email":
		return &ConfirmEmailError{
			Message: message,
		}
	default:
		return nil
	}
}
