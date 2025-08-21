package errs

import (
	"errors"
	"strings"
)

// Pretty делает цепочку ошибок более читаемой – очередная заврапленная
// ошибка будет представлена на новой строке.
func Pretty(err error) error {
	if err == nil {
		return nil
	}

	prevErr := err

	messages := make([]string, 0)

	for {
		curErr := errors.Unwrap(prevErr)
		if curErr != nil {
			messages = append(messages, cut(prevErr.Error(), curErr.Error()))
		} else {
			messages = append(messages, prevErr.Error())

			break
		}

		prevErr = curErr
	}

	return &PrettyError{
		message: strings.Join(messages, "\n"),
		cause:   err,
	}
}

func cut(str, substr string) string {
	index := strings.Index(str, substr)

	return strings.TrimSpace(str[:index])
}

type PrettyError struct {
	message string
	cause   error
}

func (e *PrettyError) Error() string {
	return e.message
}

func (e *PrettyError) Unwrap() error {
	return e.cause
}
