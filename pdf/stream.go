/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bytes"
import "bufio"

// Implements:
// 	pdf.Object
//	bufio.Writer
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

func (s *Stream) Serialize (f *bufio.Writer, file... File) {
	s.dictionary.Add ("Length", NewIntNumeric(s.buffer.Len()))
	s.dictionary.Serialize(f)
	f.WriteString ("\nstream\n")
	f.Write (s.buffer.Bytes())
	f.WriteString ("\nendstream\n")
}
