package errctx

import "errors"

type Fields map[string]any

func copyFields(fields Fields) Fields {
	fieldsCopy := make(Fields, len(fields))

	for field, value := range fields {
		fieldsCopy[field] = value
	}

	return fieldsCopy
}

type ErrorWithContext struct {
	fields Fields
	cause  error
}

func (e *ErrorWithContext) Error() string {
	return e.cause.Error()
}

func (e *ErrorWithContext) Unwrap() error {
	return e.cause
}

func AppendTo(err error, fields Fields) error {
	if err == nil {
		return nil
	}

	var errWithContext *ErrorWithContext

	if ok := errors.As(err, &errWithContext); ok {
		for field, value := range fields {
			if _, ok := errWithContext.fields[field]; !ok {
				errWithContext.fields[field] = value
			}
		}

		return err
	}

	return &ErrorWithContext{
		fields: copyFields(fields),
		cause:  err,
	}
}

func From(err error) Fields {
	if err == nil {
		return nil
	}

	var errWithContext *ErrorWithContext

	if ok := errors.As(err, &errWithContext); ok {
		return copyFields(errWithContext.fields)
	}

	return Fields{}
}
