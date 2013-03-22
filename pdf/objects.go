/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bytes"
import "io"
import "fmt"
import "unicode"

type Name struct {
	name string
}

func HexDigit (b byte) (result byte) {
	switch {
	case b < 10:
		result = b + '0'
	default:
		result = (b-10) + 'A'
	}
	return result
}

// Constructor for Name object
func NewName (s string) (* Name) {
	return &Name{s}
}

func nameEscapeByte (b byte) (result []byte) {
	switch {
	case b != '#' && IsRegular(b):
		result = []byte{b}
	default:
		result = []byte{'#', HexDigit(b/16), HexDigit(b%16)}
	}
	return result
}

func (n *Name) Serialize (f io.Writer) {
	f.Write ([]byte{'/'})
	for _,b := range []byte(n.name) {
		f.Write (nameEscapeByte(b))
	}
	return
}

type String struct {
	value string
	serializer func (t *String, f io.Writer)
}

// Constructor for Name object
func NewString (s string) (* String) {
	return &String{s,normalSerializer}
}

func stringMinimalEscapeByte (b byte) (result []byte) {
	switch b {
	case '(', ')', '\\':
		result = []byte{'\\', b}
	default:
		result = []byte{b}
	}
	return result
}

func normalSerializer (s *String, f io.Writer) {
	f.Write ([]byte{'('})
	for _,b := range []byte(s.value) {
		f.Write (stringMinimalEscapeByte(b))
	}
	f.Write ([]byte{')'})
	return
}

func stringAsciiEscapeByte (b byte) (result []byte) {
	switch b {
	case '(', ')', '\\':
		result = []byte{'\\', b}
	default:
		if b < 128 && unicode.IsPrint(rune(b)) {
			result = []byte{b}
		} else {
			var buffer bytes.Buffer
			buffer.Write([]byte{'\\'})
			switch b {
			case '\n':
				buffer.Write([]byte{'n'})
			case '\r':
				buffer.Write([]byte{'r'})
			case '\t':
				buffer.Write([]byte{'t'})
			case '\b':
				buffer.Write([]byte{'b'})
			case '\f':
				buffer.Write([]byte{'f'})
			default:
				fmt.Fprintf (&buffer, "%03o", b)
			}
			result = buffer.Bytes()
		}
	}
	return result
}

func asciiSerializer (s *String, f io.Writer) {
	f.Write ([]byte{'('})
	for _,b := range []byte(s.value) {
		f.Write (stringAsciiEscapeByte(b))
	}
	f.Write ([]byte{')'})
	return
}

func hexSerializer (s *String, f io.Writer) {
	f.Write ([]byte{'<'})
	for _,b := range []byte(s.value) {
		f.Write ([]byte{HexDigit(b/16), HexDigit(b%16)})
	}
	f.Write ([]byte{'>'})
}

func (s *String) Serialize (f io.Writer) {
	s.serializer(s, f)
}

func (s *String) SetNormalOutput () {
	s.serializer = normalSerializer
}

func (s *String) SetHexOutput () {
	s.serializer = hexSerializer
}

func (s *String) SetAsciiOutput () {
	s.serializer = asciiSerializer
}