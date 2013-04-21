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
	// returned from or passed to a function without fear of
	// modifying internal data structures.

	Clone() Object
	// Serialize() writes a representation of object to Writer as if with indirect references
	// resolved using optional File.
	Serialize(Writer, ...File)
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
