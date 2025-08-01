package errors

// TrimStackTrace режет все стектрейсы в цепочке ошибок err.
func TrimStackTrace(err error) error {
	if err != nil {
		return &Error{
			cause: err,
		}
	}

	return nil
}

type Error struct {
	cause error
}

func (e *Error) Error() string {
	return e.cause.Error()
}

func (e *Error) Unwrap() error {
	return e.cause
}
