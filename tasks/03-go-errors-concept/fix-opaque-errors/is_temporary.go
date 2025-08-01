package errors

import "errors"

func IsTemporary(err error) bool {
	for {
		if isTemporary(err) {
			return true
		}

		err = errors.Unwrap(err)
		if err == nil {
			return false
		}
	}
}

func isTemporary(err error) bool {
	type t interface {
		IsTemporary() bool
	}

	if tInstance, ok := err.(t); ok {
		return tInstance.IsTemporary()
	}

	return false
}
