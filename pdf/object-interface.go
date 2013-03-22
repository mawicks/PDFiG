/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "io"

// All PDF objects implement the pdf.Object inteface
type Object interface {
	Serialize(io.Writer)	// Write a representation of object.
}

