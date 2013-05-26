package pdf

import (
	"bytes"
	"io")

type StreamFilterFactory interface {
	Name() string
	NewEncoder(io.WriteCloser) io.WriteCloser
	NewDecoder(io.Reader) io.Reader
	DecodeParms(file... File) Object
}

var registeredFilterFactoryFactories map[string] func (*Dictionary) StreamFilterFactory

func RegisterFilterFactoryFactory(name string, sfff func (*Dictionary) StreamFilterFactory) {
	if registeredFilterFactoryFactories == nil {
		registeredFilterFactoryFactories = make(map[string]func (*Dictionary) StreamFilterFactory, 5)
	}
	registeredFilterFactoryFactories[name] = sfff
}

func FilterFactory(name string, d *Dictionary) StreamFilterFactory {
	if sfff,ok := registeredFilterFactoryFactories[name]; ok {
		return sfff(d)
	}
	return nil
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