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

var registeredFilterFactoryFactories map[string] func (ReadOnlyDictionary) StreamFilterFactory

func RegisterFilterFactoryFactory(name string, sfff func (ReadOnlyDictionary) StreamFilterFactory) {
	if registeredFilterFactoryFactories == nil {
		registeredFilterFactoryFactories = make(map[string]func (ReadOnlyDictionary) StreamFilterFactory, 5)
	}
	registeredFilterFactoryFactories[name] = sfff
}

func FilterFactory(name string, d ReadOnlyDictionary) StreamFilterFactory {
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