package pdf

import "bytes"
import "fmt"
import "unicode"

// PDF "String" object
// Implements:
//	pdf.Object
type String struct {
	value      string
	serializer func(t *String, w Writer)
}

// Constructor for Name object
func NewString(s string) *String {
	return &String{s, normalSerializer}
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

func normalSerializer(s *String, w Writer) {
	w.WriteByte('(')
	for _, b := range []byte(s.value) {
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

func asciiSerializer(s *String, w Writer) {
	w.WriteByte('(')
	for _, b := range []byte(s.value) {
		w.Write(stringAsciiEscapeByte(b))
	}
	w.WriteByte(')')
	return
}

func hexSerializer(s *String, w Writer) {
	w.WriteByte('<')
	for _, b := range []byte(s.value) {
		w.WriteByte(HexDigit(b / 16))
		w.WriteByte(HexDigit(b % 16))
	}
	w.WriteByte('>')
}

func (s *String) Serialize(w Writer, file ...File) {
	s.serializer(s, w)
}

func (s *String) SetNormalOutput() {
	s.serializer = normalSerializer
}

func (s *String) SetHexOutput() {
	s.serializer = hexSerializer
}

func (s *String) SetAsciiOutput() {
	s.serializer = asciiSerializer
}
