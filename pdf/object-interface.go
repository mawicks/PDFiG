package pdf

import "bufio"
import "bytes"

type Writer interface {
	Write(p []byte) (nn int, err error)
	WriteByte(c byte) error
	WriteRune(r rune) (size int, err error)
	WriteString(s string) (int, error)
}

// All PDF objects implement the pdf.Object inteface
type Object interface {
	// Clone() copies the object in such a way that it can be
	// returned from or passed to a function without losing
	// control of its internal data structures.
	Clone() Object

	// Serialize() write a serial byte representation (as defined
	// by the PDF specification) of the object to the Writer.
	// Indirect references are resolved and numbered as if they
	// were being written to the optional File argument.  Having
	// separate arguments for Writer and File allows writing an
	// object to stdout, but using the indirect reference object
	// numbers as if it were contained in a specific PDF file.
	// Objects can be unserialized using Parser.Scan().
	Serialize(Writer, ...File)

	// If the target is an indirect reference, Dereference()
	// returns an object that is not an indirect reference.
	// Otherwise it returns the target.
	Dereference() Object
}

// ObjectStringDecorator adds the String() method to Object; delegating all other methods to object.
type ObjectStringDecorator struct {
	Object
}

func (o *ObjectStringDecorator) String(file ...File) string {
	var buffer bytes.Buffer
	f := bufio.NewWriter(&buffer)
	o.Serialize(f, file...)
	f.Flush()
	return buffer.String()
}
