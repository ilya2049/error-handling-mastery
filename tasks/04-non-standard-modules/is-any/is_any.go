package errors

import "errors"

func IsAny(err error, references ...error) bool {
	for _, referenceErr := range references {
		if errors.Is(err, referenceErr) {
			return true
		}
	}

	return false
}
