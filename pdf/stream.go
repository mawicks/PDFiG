/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bytes"
import "io"

// Implements:
// 	pdf.Object
//	io.Writer
type Stream struct {
	dictionary *Dictionary
	buffer bytes.Buffer
}

// Constructor for Stream object
func NewStream () *Stream {
	return &Stream{NewDictionary(),bytes.Buffer{}}
}

func (s *Stream) Add (key string, o Object) {
	s.dictionary.Add(key, o)
}

func (s *Stream) Remove (key string) {
	s.dictionary.Remove(key)
}

func (s *Stream) Write(bytes []byte) (int, error) {
	return s.buffer.Write(bytes)
}

func (s *Stream) Serialize (f io.Writer) {
	s.dictionary.Add ("Length", NewIntNumeric(s.buffer.Len()))
	s.dictionary.Serialize(f)
	f.Write ([]byte("\nstream\n"))
	f.Write (s.buffer.Bytes())
	f.Write ([]byte("\nendstream\n"))
}
