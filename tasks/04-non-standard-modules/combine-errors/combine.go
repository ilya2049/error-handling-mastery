package errors

import (
	"fmt"
)

// Combine "прицепляет" ошибки other к err так, что они начинают фигурировать при выводе
// её на экран через спецификатор `%+v`. Если err является nil, то первостепенной ошибкой
// становится первая из ошибок other.
func Combine(err error, other ...error) error {
	if err == nil {
		for _, e := range other {
			if e != nil {
				return &CombinedError{
					err: e,
				}
			}
		}

		return nil
	}

	additionalErrors := make([]error, 0)
	for _, e := range other {
		if e != nil {
			additionalErrors = append(additionalErrors, e)
		}
	}

	return &CombinedError{
		err:              err,
		additionalErrors: additionalErrors,
	}
}

type CombinedError struct {
	err              error
	additionalErrors []error
}

func (e *CombinedError) Error() string {
	return e.err.Error()
}

func (e *CombinedError) Unwrap() error {
	return e.err
}

func (e *CombinedError) Format(state fmt.State, verb rune) {
	if verb == 'v' && state.Flag('+') {
		state.Write([]byte(e.Error()))
		if len(e.additionalErrors) > 0 {
			state.Write([]byte("\n"))
		}

		for _, ad := range e.additionalErrors {
			state.Write([]byte("  - " + ad.Error()))
			state.Write([]byte("\n"))
		}

		return
	}

	state.Write([]byte(e.Error()))
}
