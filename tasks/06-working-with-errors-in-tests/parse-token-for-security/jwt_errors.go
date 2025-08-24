package jwt

import "errors"

var (
	ErrEmptyJWT               = errors.New("empty jwt data")
	ErrInvalidTokenFormat     = errors.New("invalid token format: 'header.payload.signature' was expected")
	ErrInvalidHeaderEncoding  = errors.New("invalid header encoding")
	ErrUnsupportedTokenType   = errors.New("unsupported token type")
	ErrUnsupportedSigningAlgo = errors.New("unsupported the signing algo")
	ErrInvalidSignature       = errors.New("invalid signature")
	ErrInvalidPayloadEncoding = errors.New("invalid payload encoding")
	ErrExpiredToken           = errors.New("token was expired")
)

type ErrorWithEmail struct {
	cause error
	email string
}

func WithEmail(err error, email string) error {
	return &ErrorWithEmail{
		cause: err,
		email: email,
	}
}

func (e *ErrorWithEmail) Error() string {
	return e.cause.Error()
}

func (e *ErrorWithEmail) Unwrap() error {
	return e.cause
}

func (e *ErrorWithEmail) Email() string {
	return e.email
}
