package pdf

import (
	"bytes"
	"io")

type StreamFilter interface {
	Name() string
	NewEncoder(io.WriteCloser) io.WriteCloser
	NewDecoder(io.Reader) io.Reader
	DecodeParms(file... File) Object
}

type BufferCloser struct {
	bytes.Buffer
}

func NewBufferCloser() *BufferCloser {
	return new(BufferCloser)
}

func (bc *BufferCloser) Close() error {
	return nil
}