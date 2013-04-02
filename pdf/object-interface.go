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
	Serialize(Writer,...File)		// Write a representation of object.
}

// ObjectStringDecorator adds the String() method to Object; delegating all other methods to object.
type ObjectStringDecorator struct {
	Object
}

func (o *ObjectStringDecorator) String(file...File) string {
	var buffer bytes.Buffer
	f := bufio.NewWriter(&buffer)
	o.Serialize (f, file...)
	f.Flush()
	return buffer.String()
}

