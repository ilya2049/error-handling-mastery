package reader

import (
	"errors"
	"io"
)

var ErrInvalidChunkSize error

func ReadByChunk(r io.Reader, chunkSize int) ([][]byte, error) {
	if r == nil {
		return nil, nil
	}

	if chunkSize < 1 {
		return nil, ErrInvalidChunkSize
	}

	var (
		chunks    = NewChunks(chunkSize)
		err       error
		readBytes int
	)

	for err == nil {
		buf := make([]byte, chunkSize)
		readBytes, err = r.Read(buf)
		chunks.Append(buf[:readBytes])
	}

	if errors.Is(err, io.EOF) {
		return chunks.Flush(), nil
	}

	return nil, err
}

type Chunks struct {
	chunkSize int
	chunks    [][]byte
	tmpChunk  []byte
}

func NewChunks(chunkSize int) *Chunks {
	return &Chunks{
		chunkSize: chunkSize,
		chunks:    make([][]byte, 0),
		tmpChunk:  make([]byte, 0, chunkSize),
	}
}

func (ch *Chunks) Append(bytes []byte) {
	for _, b := range bytes {
		if len(ch.tmpChunk) == ch.chunkSize {
			ch.chunks = append(ch.chunks, ch.tmpChunk)
			ch.tmpChunk = make([]byte, 0, ch.chunkSize)
		}

		ch.tmpChunk = append(ch.tmpChunk, b)
	}
}

func (ch *Chunks) Flush() [][]byte {
	ch.chunks = append(ch.chunks, ch.tmpChunk)

	return ch.chunks
}
