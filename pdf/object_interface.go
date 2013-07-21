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
	// control of its internal data structures.  The returned
	// value can safely be cast to the interface type whose
	// Clone() method was invoked.
	Clone() Object

	// If the target is an indirect reference, Dereference()
	// returns an object that is not an indirect reference.
	// Otherwise it returns the target.
	Dereference() Object

	// Most PDF Objects are expected to have two interfaces: a
	// protected and an unprotected interfaces.  This is for cases
	// where a class returns an object but wants to be certain the
	// client cannot alter the object.  Protected() returns a
	// protected interface instance. The protected interface
	// implements a subset of the methods available for the
	// unprotected version of the interface.  The return value can
	// safely be cast to the documented protected version of the
	// interface whose Protected() method was invoked.
	Protected() Object

	// Unprotected() returns an unprotected interface instance.
	// Unprotected interface implementations typically return
	// themselves.  Protected interface implementations will
	// return a Clone() of the protected instance.  The return
	// value can safely be cast to the documented unprotected
	// version of the interface whose Unprotected() method was
	// invoked.
	Unprotected() Object

	// Serialize() write a serial byte representation (as defined
	// by the PDF specification) of the object to the Writer.
	// Indirect references are resolved and numbered as if they
	// were being written to the optional File argument.  Having
	// separate arguments for Writer and File allows writing an
	// object to stdout, but using the indirect reference object
	// numbers as if it were contained in a specific PDF file.
	// Objects can be unserialized using Parser.Scan().
	Serialize(Writer, ...File)

}

// ObjectStringDecorator adds the String() method to Object; delegating all other methods to object.
type ObjectStringDecorator struct {
	Object
}

func (o ObjectStringDecorator) String(file ...File) string {
	var buffer bytes.Buffer
	f := bufio.NewWriter(&buffer)
	o.Serialize(f, file...)
	f.Flush()
	return buffer.String()
}
