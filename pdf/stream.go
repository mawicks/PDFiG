package pdf

import (
	"bytes"
	"container/list"
	"io")

// Implements:
// 	pdf.Object
//	bufio.Writer

type Stream struct {
	dictionary *Dictionary
	buffer     bytes.Buffer
	// filterList is only used for writing.  Streams are fully
	// decoded when read and the in-memory stream contents are no
	// encoded in any way.  However, if a client writes an object
	// that was previously read, the client may want to use the
	// same filters.  Therefore, filters encountered while reading
	// are added to the filter list.
	filterList *list.List
}

// Constructor for Stream object
func NewStream() *Stream {
	return &Stream{NewDictionary(), bytes.Buffer{}, nil}
}

func NewStreamFromContents(dictionary *Dictionary,b []byte, filterList *list.List) *Stream {
	return &Stream{dictionary, *bytes.NewBuffer(b), filterList}
}

func (s *Stream) AddFilter(filter StreamFilter) {
	if s.filterList == nil {
		s.filterList = list.New()
	}
	s.filterList.PushBack(filter)
}

func (s *Stream) Clone() Object {
	var newFilterList *list.List
	if s.filterList != nil {
		newFilterList = list.New()
		for item:=s.filterList.Front(); item != nil; item = item.Next() {
			newFilterList.PushBack(item.Value)
		}
	}
	return NewStreamFromContents(s.dictionary,s.buffer.Bytes(), newFilterList)
}

func (s *Stream) Dereference() Object {
	return s
}

func (s *Stream) Add(key string, o Object) {
	s.dictionary.Add(key, o)
}

func (s *Stream) Reader() *bytes.Reader {
	return bytes.NewReader(s.buffer.Bytes())
}

func (s *Stream) Remove(key string) {
	s.dictionary.Remove(key)
}

func (s *Stream) Write(bytes []byte) (int, error) {
	return s.buffer.Write(bytes)
}

func (s *Stream) Serialize(w Writer, file ...File) {
	streamBuffer := NewBufferCloser()
	var streamWriter io.WriteCloser = streamBuffer

	if s.filterList != nil && s.filterList.Front() != nil {
		filters := NewArray()
		decodeParameters := NewArray()
		needDecodeParameters := false

		for item:=s.filterList.Front(); item != nil; item = item.Next() {
			streamWriter = item.Value.(StreamFilter).NewEncoder(streamWriter)
			filters.Add (NewName(item.Value.(StreamFilter).Name()))
			decodeParms := item.Value.(StreamFilter).DecodeParms(file...)
			decodeParameters.Add (decodeParms)
			if decodeParms != NewNull() {
				needDecodeParameters = true
			}
		}
		// Eliminate the arrays if they have only one element.
		if filters.Size() == 1 {
			s.dictionary.Add("Filter", filters.At(0))
			if needDecodeParameters {
				s.dictionary.Add("DecodeParms", decodeParameters.At(0))
			}
		} else {
			s.dictionary.Add("Filter", filters)
			if needDecodeParameters {
				s.dictionary.Add("DecodeParms", decodeParameters)
			}
		}
	}

	streamWriter.Write(s.buffer.Bytes())
	streamWriter.Close()

	s.dictionary.Add("Length", NewIntNumeric(streamBuffer.Len()))
	s.dictionary.Serialize(w, file...)

	w.WriteString("\nstream\n")
	w.Write(streamBuffer.Bytes())
	w.WriteString("\nendstream")
}
