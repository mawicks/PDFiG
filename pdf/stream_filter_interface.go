package pdf

import "bytes"

type StreamFilter interface {
	Encode([]byte) []byte
	Decode([]byte) (decoded []byte,ok bool)
	Name() string
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