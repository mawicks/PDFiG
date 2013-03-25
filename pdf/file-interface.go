/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

type ObjectNumber struct {
	number uint32
	generation uint16
}

type File interface {
	AssignObjectNumber (o Object) ObjectNumber
}

