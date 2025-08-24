package jwt

import (
	"bytes"
	"fmt"
	"time"
)

// Исправь ParseToken так, чтобы из ошибки можно было достать email пользователя.

// ParseToken парсит и валидирует токен jwt, проверяя, что он подписан
// алгоритмом HMAC SHA256 с использованием ключа secret.
func ParseToken(jwt, secret []byte) (Token, error) {
	if len(jwt) == 0 {
		return Token{}, ErrEmptyJWT
	}

	parts := bytes.Split(jwt, byteDot)
	if len(parts) != 3 {
		return Token{}, ErrInvalidTokenFormat
	}

	headerData, payloadData, signData := parts[0], parts[1], parts[2]

	t, err := parsePayload(payloadData)
	if err != nil {
		return Token{}, fmt.Errorf("%w: %v", ErrInvalidPayloadEncoding, err)
	}

	email := t.Email

	h, err := parseHeader(headerData)
	if err != nil {
		return Token{}, WithEmail(fmt.Errorf("%w: %v", ErrInvalidHeaderEncoding, err), email)
	}

	if h.Typ != supportedTokenType {
		return Token{}, WithEmail(fmt.Errorf("%w: %q", ErrUnsupportedTokenType, h.Typ), email)
	}

	if err := verifySignature(
		h.Alg,
		bytes.Join([][]byte{parts[0], parts[1]}, byteDot),
		signData,
		secret,
	); err != nil {
		return Token{}, WithEmail(fmt.Errorf("verify signature: %w", err), email)
	}

	if time.Unix(t.ExpiredAt, 0).Before(time.Now()) {
		return Token{}, WithEmail(ErrExpiredToken, email)
	}

	return t, nil
}
