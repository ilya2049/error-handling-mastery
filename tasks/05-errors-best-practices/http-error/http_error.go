package httperr

import (
	"fmt"
	"net/http"
)

// Реализуй нас.
var (
	ErrStatusOK                  = HTTPError(200)
	ErrStatusBadRequest          = HTTPError(400)
	ErrStatusNotFound            = HTTPError(404)
	ErrStatusUnprocessableEntity = HTTPError(422)
	ErrStatusInternalServerError = HTTPError(500)
)

// Реализуй меня.
type HTTPError int

func (e HTTPError) Code() int {
	return int(e)
}

func (e HTTPError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}
