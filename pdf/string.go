package pdf

import "bytes"
import "fmt"
import "unicode"
import "unicode/utf16"

// PDF "String" object
// Implements:
//	pdf.Object
type ProtectedString interface {
	Object
	Bytes() []byte
}

type String interface {
	ProtectedString
	SetSerializer(func(String,Writer))
}

type stringImpl struct {
	value      []byte
	serializer func(s String, w Writer)
}

// Constructor for Name object
func NewTextString(s string) String {
	// If PDFDocEncoding works, use that
	result,ok := PDFDocEncoding ([]rune(s))
	// Otherwise use UTF16-BE
	if !ok {
		utf16result := utf16.Encode([]rune(s))
		result = make([]byte,0)
		result = append(result, 0xfe, 0xff)
		for _,w := range utf16result {
			result = append(result, byte(w>>8), byte(w&0xff))
		}
	}
	return &stringImpl{result, NormalStringSerializer}
}

func NewBinaryString(s []byte) String {
	return &stringImpl{s, NormalStringSerializer}
}

func (s *stringImpl) Bytes() (result []byte) {
	newSlice := make([]byte, len(s.value))
	copy (newSlice, s.value)
	return newSlice
}

func (s *stringImpl) Clone() Object {
	newString := *s
	return &newString
}

func (s *stringImpl) Dereference() Object {
	return s
}

// Return value of Protected() can safely be cast to ProtectedString
// but not String.
func (s *stringImpl) Protected() Object {
	return &readOnlyString{s}
}

// Return value of Unprotected() can safely be cast to String or
// ProtectedString.
func (s *stringImpl) Unprotected() Object {
	return s
}

func (s *stringImpl) Serialize(w Writer, file ...File) {
	s.serializer(s, w)
}

type readOnlyString struct {
	s String
}

func (ros readOnlyString) Bytes() []byte {
	return ros.s.Bytes()
}

func (ros readOnlyString) Clone() Object {
	return ros
}

func (ros readOnlyString) Dereference() Object {
	return ros
}

// Return value of Protected() can safely be cast to ProtectedString
// but not String.
func (ros readOnlyString) Protected() Object {
	return ros
}

// Return value of Unprotected() can safely be cast to String or
// ProtectedString.
func (ros readOnlyString) Unprotected() Object {
	return ros.s.Clone()
}

func (ros readOnlyString) Serialize(w Writer, file ...File) {
	ros.s.Serialize(w, file...)
}

func stringMinimalEscapeByte(b byte) (result []byte) {
	switch b {
	case '(', ')', '\\':
		result = []byte{'\\', b}
	default:
		result = []byte{b}
	}
	return result
}

func NormalStringSerializer(s String, w Writer) {
	w.WriteByte('(')
	for _, b := range s.Bytes() {
		w.Write(stringMinimalEscapeByte(b))
	}
	w.WriteByte(')')
	return
}

func stringAsciiEscapeByte(b byte) (result []byte) {
	switch b {
	case '(', ')', '\\':
		result = []byte{'\\', b}
	default:
		if b < 128 && unicode.IsPrint(rune(b)) {
			result = []byte{b}
		} else {
			var buffer bytes.Buffer
			buffer.WriteByte('\\')
			switch b {
			case '\n':
				buffer.WriteByte('n')
			case '\r':
				buffer.WriteByte('r')
			case '\t':
				buffer.WriteByte('t')
			case '\b':
				buffer.WriteByte('b')
			case '\f':
				buffer.WriteByte('f')
			default:
				fmt.Fprintf(&buffer, "%03o", b)
			}
			result = buffer.Bytes()
		}
	}
	return result
}

func GeneralAsciiEscapeByte(b byte) (result []byte) {
	switch b {
	case '\\':
		result = []byte{'\\', b}
	default:
		if b < 128 && unicode.IsPrint(rune(b)) {
			result = []byte{b}
		} else {
			var buffer bytes.Buffer
			buffer.WriteByte('\\')
			switch b {
			case '\n':
				buffer.WriteByte('n')
			case '\r':
				buffer.WriteByte('r')
			case '\t':
				buffer.WriteByte('t')
			case '\b':
				buffer.WriteByte('b')
			case '\f':
				buffer.WriteByte('f')
			default:
				fmt.Fprintf(&buffer, "%03o", b)
			}
			result = buffer.Bytes()
		}
	}
	return result
}

func AsciiStringSerializer(s String, w Writer) {
	w.WriteByte('(')
	for _, b := range s.Bytes() {
		w.Write(stringAsciiEscapeByte(b))
	}
	w.WriteByte(')')
	return
}

func HexStringSerializer(s String, w Writer) {
	w.WriteByte('<')
	for _, b := range s.Bytes() {
		w.WriteByte(HexDigit(b / 16))
		w.WriteByte(HexDigit(b % 16))
	}
	w.WriteByte('>')
}

func (s *stringImpl) SetSerializer (serializer func(String,Writer)) {
	s.serializer = serializer
}

