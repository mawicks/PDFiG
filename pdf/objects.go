/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bytes"
import "io"
import "fmt"
import "math"
import "unicode"

// All PDF objects implement the pdf.Object inteface
type Object interface {
	Serialize(io.Writer)	// Write a representation of object.
}


// PDF "Null" object
// Implements: pdf.Object
type Null struct {}

func (n *Null) Serialize (f io.Writer) {
	fmt.Fprint (f, "null")
	return
}

// PDF "Boolean" object
// Implements: pdf.Object
type Boolean struct {
	value bool
}

// Constructor for Boolean object
func NewBoolean (v bool) *Boolean {
	return &Boolean{v}
}

func  (b *Boolean) Serialize (f io.Writer) {
	if b.value {
		fmt.Fprint (f, "true")
	} else {
		fmt.Fprint (f, "false")
	}
}

// PDF "Numeric" object
// Implements: pdf.Object
type FloatNumeric struct {
	value float32
}

type IntNumeric struct {
	value int32
}

func  (n *FloatNumeric) Serialize (f io.Writer) {
	fmt.Fprint (f, n.value);
}

func  (n *IntNumeric) Serialize (f io.Writer) {
	fmt.Fprint (f, n.value);
}

func adjustFloatRange (v float64) (float32Value float32) {
	switch {
	case v > math.MaxFloat32:
		float32Value = math.MaxFloat32
		
	case v < - math.MaxFloat32:
		float32Value = -math.MaxFloat32
		
	case math.Abs(v) < 1.175494351e-38:
		// Smallest 32-bit floating point number without losing precision.
		// PDF spec says set values below 1.175e-38 to 0 for readers using
		// 32-bit floats
		float32Value = 0.0
		
	default:
		float32Value = float32(v)
	}
	return float32Value
}

// Constructor for Numeric object
func NewNumeric (v float64) (result Object) {
	var intValue = int32(v)

	// Use IntNumeric if value can be represented as int32;
	// otherwise use FloatNumeric, which uses float32 internally
	if float64(intValue) == v {
		result = &IntNumeric{intValue}
	} else {
		result = &FloatNumeric{adjustFloatRange(v)}
	}

	return result
}

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
		result = []byte{'#', hexDigit(b/16), hexDigit(b%16)}
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
		f.Write ([]byte{hexDigit(b/16), hexDigit(b%16)})
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