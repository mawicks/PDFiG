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

var registeredFilterFactoryFactories map[string] func (ProtectedDictionary) StreamFilterFactory

func RegisterFilterFactoryFactory(name string, sfff func (ProtectedDictionary) StreamFilterFactory) {
	if registeredFilterFactoryFactories == nil {
		registeredFilterFactoryFactories = make(map[string]func (ProtectedDictionary) StreamFilterFactory, 5)
	}
	registeredFilterFactoryFactories[name] = sfff
}

func FilterFactory(name string, d ProtectedDictionary) StreamFilterFactory {
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