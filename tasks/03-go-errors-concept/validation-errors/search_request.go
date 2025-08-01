package requests

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const maxPageSize = 100

var (
	errIsNotRegexp     = errors.New("exp is not regexp")
	errInvalidPage     = errors.New("invalid page")
	errInvalidPageSize = errors.New("invalid page size")
)

type ValidationErrors []error

func (errs ValidationErrors) Is(err error) bool {
	for _, err := range errs {
		if err == err {
			return true
		}
	}

	return false
}

func (errs ValidationErrors) Error() string {
	builder := strings.Builder{}

	builder.WriteString("validation errors:\n")

	for _, err := range errs {
		builder.WriteString("\t" + err.Error() + "\n")
	}

	return builder.String()
}

type SearchRequest struct {
	Exp      string
	Page     int
	PageSize int
}

func (r SearchRequest) Validate() error {
	var errs ValidationErrors

	if _, err := regexp.Compile(r.Exp); err != nil {
		errs = append(errs, fmt.Errorf("%w: %v", errIsNotRegexp, err))
	}

	if r.Page <= 0 {
		errs = append(errs, fmt.Errorf("%w: %d", errInvalidPage, r.Page))
	}

	if r.PageSize <= 0 {
		errs = append(errs, fmt.Errorf("%w: %d <= 0", errInvalidPageSize, r.PageSize))
	}

	if r.PageSize > maxPageSize {
		errs = append(errs, fmt.Errorf("%w: %d > %d", errInvalidPageSize, r.PageSize, maxPageSize))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
