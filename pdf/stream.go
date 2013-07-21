package pdf

import (
	"bytes"
	"container/list"
	"io")

// Implements:
// 	pdf.Object

type ProtectedStream interface {
	Object
	Reader() (result io.Reader)
}

// Implements:
//	io.Writer

type Stream interface {
	ProtectedStream
	io.Writer
	AddFilter(filter StreamFilterFactory)
	Remove(key string)
}

type stream struct {
	dictionary Dictionary
	buffer     bytes.Buffer
	// filterList is only used for writing.  Streams are fully
	// decoded when read and the in-memory stream contents are no
	// encoded in any way.  However, if a client writes an object
	// that was previously read, the client may want to use the
	// same filters.  Therefore, filters encountered while reading
	// are added to the filter list.
	filterList *list.List
}

// Constructor for standard implementation of Stream.
func NewStream() Stream {
	return &stream{NewDictionary(), bytes.Buffer{}, nil}
}

func NewStreamFromContents(dictionary Dictionary,b []byte, filterList *list.List) Stream {
	return &stream{dictionary, *bytes.NewBuffer(b), filterList}
}

func (s *stream) AddFilter(filter StreamFilterFactory) {
	if s.filterList == nil {
		s.filterList = list.New()
	}
	s.filterList.PushBack(filter)
}

func (s *stream) Clone() Object {
	var newFilterList *list.List
	if s.filterList != nil {
		newFilterList = list.New()
		for item:=s.filterList.Front(); item != nil; item = item.Next() {
			newFilterList.PushBack(item.Value)
		}
	}
	return NewStreamFromContents(s.dictionary,s.buffer.Bytes(), newFilterList)
}

func (s *stream) Dereference() Object {
	return s
}

// Return value of Protect can safely be cast to ProtectedStream.
func (s *stream) Protect() Object {
	return protectedStream{s}
}

// Return value of Protect can safely be cast to Stream
func (s *stream) Unprotect() Object {
	return s
}

func (s *stream) Reader() (result io.Reader) {
	result = bytes.NewReader(s.buffer.Bytes())
	if filters := s.dictionary.GetArray("Filter"); filters != nil {
		parms := s.dictionary.GetArray("DecodeParms")
		for i:=0; i<filters.Size(); i++ {
			if n,ok := filters.At(i).(Name); ok {
				var d Dictionary
				if parms != nil && i < parms.Size() {
					d,_ = parms.At(i).(Dictionary)
				}
				if sff := FilterFactory(n.String(),d); sff != nil {
					result = sff.NewDecoder(result)
				} else {
					return nil
				}
			}
		}
	} else if n,ok := s.dictionary.GetName("Filter"); ok {
		d := s.dictionary.GetDictionary("DecodeParms")
		if sff := FilterFactory(n,d); sff != nil {
			result = sff.NewDecoder(result)
		} else {
			return nil
		}
	}
	return result
}

func (s *stream) Remove(key string) {
	s.dictionary.Remove(key)
}

func (s *stream) Write(bytes []byte) (int, error) {
	return s.buffer.Write(bytes)
}

func (s *stream) Serialize(w Writer, file ...File) {
	streamBuffer := NewBufferCloser()
	dictionary := s.dictionary.Clone().(Dictionary)

	var streamWriter io.WriteCloser = streamBuffer

	if s.filterList != nil && s.filterList.Front() != nil {
		filters := NewArray()
		decodeParameters := NewArray()
		needDecodeParameters := false

		for item:=s.filterList.Front(); item != nil; item = item.Next() {
			streamWriter = item.Value.(StreamFilterFactory).NewEncoder(streamWriter)
			filters.Add (NewName(item.Value.(StreamFilterFactory).Name()))
			decodeParms := item.Value.(StreamFilterFactory).DecodeParms(file...)
			decodeParameters.Add (decodeParms)
			if decodeParms != NewNull() {
				needDecodeParameters = true
			}
		}

		if f := s.dictionary.GetArray("Filter"); f != nil {
			filters.Append(f.(Array))
			if d := s.dictionary.GetArray("DecodeParms"); d != nil {
				decodeParameters.Append(d)
				needDecodeParameters = true
			} else if needDecodeParameters {
				for i := 0; i<f.Size(); i++ {
					decodeParameters.Add (NewNull())
				}
			}
		}

		if n,ok := s.dictionary.GetName("Filter"); ok {
			filters.Add(NewName(n))
			if d,ok := s.dictionary.GetName("DecodeParms"); ok {
				decodeParameters.Add(NewName(d))
			} else if needDecodeParameters {
				decodeParameters.Add (NewNull())
			}
		}

		// Eliminate the arrays if they have only one element.
		if filters.Size() == 1 {
			dictionary.Add("Filter", filters.At(0))
			if needDecodeParameters {
				dictionary.Add("DecodeParms", decodeParameters.At(0))
			}
		} else {
			dictionary.Add("Filter", filters)
			if needDecodeParameters {
				dictionary.Add("DecodeParms", decodeParameters)
			}
		}
	}

	streamWriter.Write(s.buffer.Bytes())
	streamWriter.Close()

	dictionary.Add("Length", NewIntNumeric(streamBuffer.Len()))
	dictionary.Serialize(w, file...)

	w.WriteString("\nstream\n")
	w.Write(streamBuffer.Bytes())
	w.WriteString("\nendstream")
}

type protectedStream struct {
	s Stream
}

// Return value of Clone() can safely be cast to ProtectedStream.
func (ps protectedStream) Clone() Object {
	return ps
}

// Return value of Dereference() can safely be cast to ProtectedStream.
func (ps protectedStream) Dereference() Object {
	return ps
}

// Return value of Protect() can safely be cast to ProtectedStream.
func (ps protectedStream) Protect() Object {
	return ps
}

// Return value of Unprotect() can safely be cast to Stream or
// ProtectedStream.
func (ps protectedStream) Unprotect() Object {
	return ps.s.Clone()
}

func (ps protectedStream) Reader() io.Reader {
	return ps.s.Reader()
}

func (ps protectedStream) Serialize(w Writer, file ...File) {
	ps.s.Serialize(w, file...)
}

func (ps protectedStream) Write(bytes []byte) (int,error) {
	return ps.s.Write(bytes)
}














