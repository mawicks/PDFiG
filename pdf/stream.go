/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bytes"

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

func (s *Stream) Serialize (w Writer, file... File) {
	s.dictionary.Add ("Length", NewIntNumeric(s.buffer.Len()))
	s.dictionary.Serialize(w)
	w.WriteString ("\nstream\n")
	w.Write (s.buffer.Bytes())
	w.WriteString ("\nendstream\n")
}
