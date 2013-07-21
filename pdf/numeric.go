package pdf

import "math"
import "strconv"

// PDF "Numeric" object
// Implements:
//	pdf.Object
type Numeric interface {
	Value() interface{}
	Object
}

// RealNumeric implements Numeric
type RealNumeric struct {
	value float32
}

// IntNumeric implements Numeric
type IntNumeric struct {
	value int
}

func (n *RealNumeric) Clone() Object {
	// Numerics are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

func (n *RealNumeric) Dereference() Object {
	return n
}

// Numerics are always immutable so the return value of Protected() can
// safely be cast back to Numeric or RealNumeric.
func (n *RealNumeric) Protected() Object {
	// Numerics are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

// Numerics are immutable so the return value of Unprotected() can
// safely be cast back to Numeric or RealNumeric.
func (n *RealNumeric) Unprotected() Object {
	// Numerics are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

func (n *RealNumeric) Serialize(w Writer, file ...File) {
	w.WriteString(strconv.FormatFloat(float64(n.value), 'f', -1, 32))
}

func (n *RealNumeric) Value() float32 {
	return n.value
}

func (n *IntNumeric) Clone() Object {
	// Numerics are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

func (n *IntNumeric) Dereference() Object {
	return n
}

// Numerics are always immutable so the return value of Protected() can
// safely be cast back to Numeric or IntNumeric
func (n *IntNumeric) Protected() Object {
	// Numerics are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

// Numerics are always immutable so the return value of Unprotected() can
// safely be cast back to Numeric or IntNumeric
func (n *IntNumeric) Unprotected() Object {
	// Numerics are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

func (n *IntNumeric) Serialize(w Writer, file ...File) {
	w.WriteString(strconv.FormatInt(int64(n.value), 10))
}

func (n *IntNumeric) Value() int {
	return n.value
}

func adjustRealRange(v float64) (float32Value float32) {
	switch {
	case v > math.MaxFloat32:
		float32Value = math.MaxFloat32

	case v < -math.MaxFloat32:
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

// Constructors for Numeric object
func NewIntNumeric(v int) Object {
	return &IntNumeric{v}
}

func NewRealNumeric(v float32) Object {
	return &RealNumeric{v}
}

func NewNumeric(v float64) (result Object) {
	var intValue = int(v)

	// Use IntNumeric if value can be represented as int32;
	// otherwise use RealNumeric, which uses float32 internally
	if float64(intValue) == v {
		result = &IntNumeric{intValue}
	} else {
		result = &RealNumeric{adjustRealRange(v)}
	}

	return result
}
