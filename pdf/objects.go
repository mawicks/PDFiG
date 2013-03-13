/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "io"
import "fmt"
import "math"

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

func NewName (s string) (* Name) {
	return &Name{s}
}

func (n *Name) Serialize (f io.Writer) {
	// All "string" types, including n.name, are assumed to be
	// encoded in UTF-8, so printing with "%s" format just works.
	fmt.Fprintf (f, "/")
	for _,b := range []byte(n.name) {
		f.Write ([]byte{b})
	}
	return
}
