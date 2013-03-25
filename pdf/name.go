/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bufio"

func HexDigit (b byte) (result byte) {
	switch {
	case b < 10:
		result = b + '0'
	default:
		result = (b-10) + 'A'
	}
	return result
}

type Name struct {
	name string
}

// Constructor for Name object
func NewName (s string) Name {
	return Name{s}
}

func nameEscapeByte (b byte) (result []byte) {
	switch {
	case b != '#' && IsRegular(b):
		result = []byte{b}
	default:
		result = []byte{'#', HexDigit(b/16), HexDigit(b%16)}
	}
	return result
}

func (n Name) Serialize (f *bufio.Writer, file... File) {
	f.WriteByte('/')
	for _,b := range []byte(n.name) {
		f.Write (nameEscapeByte(b))
	}
	return
}

