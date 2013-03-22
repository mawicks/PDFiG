/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "fmt"
import "io"

// PDF "Null" object
// Implements: pdf.Object
type Null struct {}

func (n *Null) Serialize (f io.Writer) {
	fmt.Fprint (f, "null")
	return
}

