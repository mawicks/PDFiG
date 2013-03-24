/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bufio"
import "bytes"

// All PDF objects implement the pdf.Object inteface
type Object interface {
	Serialize(*bufio.Writer)			// Write a representation of object.
}

// ObjectStringDecorator adds the String() method to Object; delegating all other methods to object.
type ObjectStringDecorator struct {
	Object
}

func (o *ObjectStringDecorator) String() string {
	var buffer bytes.Buffer
	f := bufio.NewWriter(&buffer)
	o.Serialize (f)
	f.Flush()
	return buffer.String()
}

