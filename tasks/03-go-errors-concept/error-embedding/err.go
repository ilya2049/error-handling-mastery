package errors

var (
	ErrAlreadyDone      error = &AlreadyDoneError{Err{"job is already done"}}
	ErrInconsistentData error = &InconsistentDataError{Err{"job payload is corrupted"}}
	ErrInvalidID        error = &InvalidIDError{Err{"invalid job id"}}
	ErrNotReady         error = &NotReadyError{Err{"job is not ready to be performed"}}
	ErrNotFound         error = &NotFoundError{Err{"job wasn't found"}}
)

// Реализуй тип Err и типы для ошибок выше, используя его.
type Err struct {
	Message string
}

type AlreadyDoneError struct {
	Err
}

func (e *AlreadyDoneError) Error() string {
	return e.Message
}

type InconsistentDataError struct {
	Err
}

func (e *InconsistentDataError) Error() string {
	return e.Message
}

type InvalidIDError struct {
	Err
}

func (e *InvalidIDError) Error() string {
	return e.Message
}

type NotReadyError struct {
	Err
}

func (e *NotReadyError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Err
}

func (e *NotFoundError) Error() string {
	return e.Message
}
