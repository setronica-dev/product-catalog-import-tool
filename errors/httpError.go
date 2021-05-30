package errors

func NewHttpError(code int, message string) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e HttpError) Error() string {
	return e.Message
}
