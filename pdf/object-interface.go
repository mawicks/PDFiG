/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "io"
import "bytes"

// All PDF objects implement the pdf.Object inteface
type Object interface {
	Serialize(io.Writer)	// Write a representation of object.
}

// ObjectStringDecorator adds the String() method to Object; delegating all other methods to object.
type ObjectStringDecorator struct {
	Object
}

func (o *ObjectStringDecorator) String() string {
	var buffer bytes.Buffer
	o.Serialize (&buffer)
	return buffer.String()
}

