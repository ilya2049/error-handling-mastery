package errs

import (
	"errors"
	"fmt"
	"strings"
)

func Errorf(format string, args ...any) error {
	argNumber := -1
	errors := []error{}
	errStrings := []string{}

	formatWithoutErrors := format
	argsWithoutErrors := []any{}

	for i, r := range format {
		if r == '%' {
			argNumber++

			if format[i+1] == 'w' {
				if args[argNumber] == nil {
					errStrings = append(errStrings, "<nil>")

					continue
				}

				err, ok := args[argNumber].(error)
				if ok {
					errors = append(errors, err)
					errStrings = append(errStrings, err.Error())
				} else {
					errStrings = append(errStrings, fmt.Sprintf("%v", args[argNumber]))
				}
			} else {
				argsWithoutErrors = append(argsWithoutErrors, args[argNumber])
			}
		}
	}

	if len(errors) == 0 {
		return nil
	}

	for _, s := range errStrings {
		formatWithoutErrors = strings.Replace(formatWithoutErrors, "%w", s, 1)
	}

	message := fmt.Sprintf(formatWithoutErrors, argsWithoutErrors...)

	return Errors{
		errChain: errors,
		message:  message,
	}
}

type Errors struct {
	errChain []error
	message  string
}

func (errs Errors) Is(err error) bool {
	for _, e := range errs.errChain {
		if errors.Is(e, err) {
			return true
		}
	}

	return false
}

func (errs Errors) As(target any) bool {
	for _, e := range errs.errChain {
		if errors.As(e, target) {
			return true
		}
	}

	return false
}

func (errs Errors) Error() string {
	return errs.message
}
