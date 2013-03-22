/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "fmt"
import "io"

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

