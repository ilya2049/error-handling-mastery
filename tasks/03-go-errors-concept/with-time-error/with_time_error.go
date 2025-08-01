package errs

import (
	"fmt"
	"time"
)

type WithTimeError struct {
	occurredAt time.Time
	err        error
}

func NewWithTimeError(err error) error {
	return newWithTimeError(err, time.Now)
}

func newWithTimeError(err error, timeFunc func() time.Time) error {
	if timeFunc == nil {
		panic("invalid usage of newWithTimeError")
	}

	return &WithTimeError{
		occurredAt: timeFunc(),
		err:        err,
	}
}

func (e *WithTimeError) Error() string {
	return fmt.Sprintf("%v, occurred at: %s", e.err, e.occurredAt)
}

func (e *WithTimeError) Unwrap() error {
	return e.err
}

func (e *WithTimeError) Time() time.Time {
	return e.occurredAt
}
