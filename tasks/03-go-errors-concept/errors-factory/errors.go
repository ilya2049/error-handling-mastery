package errors

// NewError возвращает новое значение-ошибку, текст которой является msg.
// Две ошибки с одинаковым текстом, созданные через NewError, не равны между собой:
//
//	NewError("end of file") != NewError("end of file")
func NewError(msg string) error {
	return &Error{Message: msg}
}

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
