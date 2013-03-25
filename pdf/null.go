/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bufio"

// PDF "Null" object
// Implements: pdf.Object
type Null struct {}

var nullSingleInstance Null

func NewNull() Object {
	return &nullSingleInstance
}

func (n *Null) Serialize (f *bufio.Writer, file... File) {
	f.WriteString("null")
	return
}

