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
	// Clone() copies the full state of the object so that future
	// changes to the clone do not affect the original object and
	// vice-versa.  For immutable objects Clone() may return the
	// original reference.  However, for container objects, a
	// clone is a deep copy.  As such it can be a relatively
	// expensive operation on a large container.  Protect() and
	// Unprotect() are less expensive and achieve a similar
	// effect.  If an interface has a protected and an unprotected
	// version, the valued returned from Clone() is the
	// unprotected version.  The returned value from Clone() can
	// safely be cast to the unprotected version of the interface
	// type whose Clone() method was invoked.
	Clone() Object

	// If the target is an indirect reference, Dereference()
	// returns an object that is not an indirect reference.
	// Otherwise it returns the target.
	Dereference() Object

	// Most PDF Objects have two versions of interfaces: a
	// protected one and an unprotected one.  The protected
	// interface is used when a method returns an object but wants
	// to be certain that the client cannot alter the object.
	// Protect() returns a protected interface instance. The
	// protected interface implements a subset of the methods
	// available for the unprotected version of the interface.
	// The returned value can safely be cast to the documented
	// protected version of the interface whose Protect() method
	// was invoked.  Although semantically protection is deep
	// (i.e., protecting a container protects all objects within
	// the container), the implementation is relatively efficient.
	Protect() Object

	// Unprotect() returns an unprotected interface instance.
	// When invoked on a protected interface, the semantics are
	// essentially copy-on-write.  When invoked on an unprotected
	// interface, it is a no-op.  Unprotecting is shallow.  It
	// applies to the container itself and not to the objects
	// within the container.  To modify a protected object in an
	// unprotected container, the protected object must be
	// unprotected, modified, and (because the unprotected object
	// is a copy) rewritten to the container.  The returned Object
	// may be cast to the documented unprotected version of the
	// interface whose Unprotect() method was invoked.
	Unprotect() Object

	// Serialize() write a serial byte representation (as defined
	// by the PDF specification) of the object to the Writer.
	// Indirect references are resolved and numbered as if they
	// were being written to the optional File argument.  Having
	// separate arguments for Writer and File allows writing an
	// object to any Writer interface (e.g., stdout for debugging
	// or a string writer for internal formatting), but using the
	// indirect reference object numbers as if it were contained
	// in a specific PDF file.  Objects can be unserialized using
	// Parser.Scan().
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
