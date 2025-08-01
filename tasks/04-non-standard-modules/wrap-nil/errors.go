package errors

import "fmt"

// Wrapf работает аналогично fmt.Errorf, только поддерживает nil-ошибки.
func Wrapf(err error, f string, v ...any) error {
	if err == nil {
		return nil
	}

	return NewError(err, f, v...)
}

func NewError(err error, f string, v ...any) *Error {
	message := fmt.Sprintf(f, v...)
	message += ": " + err.Error()

	return &Error{
		Cause:   err,
		Message: message,
	}
}

type Error struct {
	Cause   error
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}
