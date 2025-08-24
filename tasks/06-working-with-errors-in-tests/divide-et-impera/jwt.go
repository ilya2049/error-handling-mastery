package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

func parseHeader(data []byte) (h Header, err error) {
	base64Data, err := io.ReadAll(base64.NewDecoder(base64.RawURLEncoding, bytes.NewReader(data)))
	if err != nil {
		return Header{}, fmt.Errorf("%w: %w", ErrInvalidHeaderEncoding, ErrInvalidBase64)
	}

	d := json.NewDecoder(bytes.NewReader(base64Data))

	var header Header

	if err := d.Decode(&header); err != nil {
		return Header{}, fmt.Errorf("%w: %w", ErrInvalidHeaderEncoding, ErrInvalidJSON)
	}

	return header, nil
}

func parsePayload(data []byte) (t Token, err error) {
	base64Data, err := io.ReadAll(base64.NewDecoder(base64.RawURLEncoding, bytes.NewReader(data)))
	if err != nil {
		return Token{}, fmt.Errorf("%w: %w", ErrInvalidPayloadEncoding, ErrInvalidBase64)
	}

	d := json.NewDecoder(bytes.NewReader(base64Data))

	var token Token

	if err := d.Decode(&token); err != nil {
		return Token{}, fmt.Errorf("%w: %w", ErrInvalidPayloadEncoding, ErrInvalidJSON)
	}

	return token, nil
}
