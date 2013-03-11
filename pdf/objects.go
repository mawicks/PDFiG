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
type Null struct {}

func (n Null) Serialize (f io.Writer) {
	fmt.Fprint (f, "null")
	return
}

// PDF "Boolean" object
type Boolean struct {
	value bool
}

// Constructor for Boolean object
func NewBoolean (v bool) Boolean {
	return Boolean{v}
}

func  (b Boolean) Serialize (f io.Writer) {
	if b.value {
		fmt.Fprint (f, "true")
	} else {
		fmt.Fprint (f, "false")
	}
}

// PDF "Numeric" object

type FloatNumeric struct {
	value float32
}

type IntNumeric struct {
	value int32
}

func  (n FloatNumeric) Serialize (f io.Writer) {
	fmt.Fprint (f, n.value);
}

func  (n IntNumeric) Serialize (f io.Writer) {
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

	// If value can be stored as int32, use int32;
	// otherwise use float32
	if float64(intValue) == v {
		result = IntNumeric{intValue}
	} else {
		result = FloatNumeric{adjustFloatRange(v)}
	}

	return result
}

	