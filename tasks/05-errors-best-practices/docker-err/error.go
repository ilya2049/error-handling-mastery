package docker

import "strings"

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) IsPullAccessDeniedError() bool {
	return strings.Contains(e.msg, "pull access denied")
}

func (e *Error) IsNoSuchContainerError() bool {
	return strings.Contains(e.msg, "No such container")
}

func (e *Error) IsContainerNotRunningError() bool {
	return strings.Contains(e.msg, "is not running")
}

func newDockerError(err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{msg: err.Error()}
}
