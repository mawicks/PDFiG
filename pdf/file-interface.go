/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

type ObjectNumber struct {
	number uint32
	generation uint16
}

type File interface {
	// Reserve an object number for Object in File.
	AddObjectAt (ObjectNumber, Object)
	AddObject (object Object) (objectNumber ObjectNumber)
	ReserveObjectNumber (Object) ObjectNumber
	DeleteObject (ObjectNumber)
	Close ()
}

