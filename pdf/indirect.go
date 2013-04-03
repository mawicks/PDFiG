package pdf

import "strconv"

// Implements:
// 	pdf.Object
//	bufio.Writer
type Indirect struct {
	fileBindings map[File] ObjectNumber
	isFinal bool
}

// Constructor for Indirect object
func NewIndirect () *Indirect {
	return &Indirect{make(map[File] ObjectNumber,1000), false}
}

/*
pdf.Indirect is one of the most important and most complex
object types.  Whereas "direct" objects are rendered the same way on
every output stream, a pdf.Indirect is rendered differently
depending on which file it is associated with.  There is always an
underlying "direct object" that is written to one or more PDF files,
where it is assigned both an object and generation number.  When the
Serialize() method is invoked, a pdf.Indirect is rendered as an
indirect reference (for example, "10 1 R").  There are several
use-cases.

USE-CASE 1: A pdf.Indirect object is created, its Serialize() method
is invoked for one or more output files, then its Finalize() method is
invoked with its underlying direct object.  The PDF reference suggests
that a reference is sometimes written before the information required
for the object being referenced is available.  An example from the PDF
reference is using an indirect object reference for the length of the
stream before writing the stream, presumably so you can write the
stream dictionary at a moment when the length of the stream is
unknown.  The stream length is written as a separate object after the
stream is completed when the length is known with certainty.  This
particular example is somewhat contrived.  With the size of memory in
modern computers it's difficult to imagine a scenario where a stream
cannot be written to a buffer in its entirely and written as an atomic
stream object.  Nonetheless, pdf.IndirectS supports this programming
style where object references are written to a file before the objects
they refer to have been completely defined.  We do *not* support this
model with streams, however, and *do* require streams to be completed
in memory before any portion of the strean is written to a file.  A
pdf.Indirect can be written before the direct object it references has
been defined.  A pdf.Indirect obtains and reserves an object number
whenever it is written to a file, whether or not the object being
referenced has been specified.  Eventually, the Finalize() method must
be called passing the object being referenced.  At that moment, the
object being referenced is written to all files to which the
pdf.Indirect was written.  If the pdf.Indirect is subsequently added
to additional files, the Finalize()ed object must also written to
those files.  This, however, requires either retaining a reference in
memory indefinitely (bad) or reading it from one of the files where it
is known to exist (reading is not yet implemented; it must be read as
opposed to just copied because any indirect references contained
within it also need to be read and added to the file accordingly.  For
the time being, we elect not to keep references in memory.  Until
parsing is implemented, Serialize()ing a pdf.Indirect to a file after
calling Finalize() will generate an error.

USE-CASE 2: A pdf.Indirect is created based on a finished direct
object.  This is essentially the same as use-case 1.  An object is
constructed and immediately Finalize()'d.  Subsequent invocations of
Serialize() on files where it doesn't already exist cause it to be
added to that file.  As with USE-CASE 1, this requires either
retaining a reference in memory indefinitely (bad) or reading from one
of the files where it is known to exist (not yet implemented).  For
the time being, we elect not to keep references in memory. Until
parsing is implemented, so Serialize()ing a pdf.Indirect to a file
after calling Finalize() is an error.

USE-CASE 3: A pdf.Indirect is constructed when a token of the form "10
1 R" is read from file X.  The underlying direct object is unknown at
that moment, but it exists in X.  Since the object exists statically
in X, it is considered to have been finalized.  If the same object is
Serialize()'ed to file Y, then it should immediately added to file Y.
The objects contents are obtained from X.

POSSIBLE TEMPORARY BEHAVIOR: The referenced object is retained by the
pdf.Indirect so that it can be "dereferenced" or written to additional
files.  In future versions that are able to parse PDF files, the
object may be discarded from memory once written and read back from
disk if it is dererenced or written to another file.  Even better
would be a weak-reference to the object, but weak references are not
implemented in Go.

*/

func (i *Indirect) Serialize (w Writer, file... File) {
	if (len(file) == 0) {
		panic ("File parameter required for pdf.Indirect.Serialize()")
	}
	if (i.isFinal) {
		panic ("Serializing a finalized object is not yet allowed")
	}
	n := i.getObjectNumber (file[0])
	w.WriteString (strconv.FormatInt(int64(n.number), 10))
	w.WriteByte (' ')
	w.WriteString (strconv.FormatInt(int64(n.generation), 10))
	w.WriteString (" R")
}

func (i *Indirect) Finalize (o Object) {
	if (i.isFinal) {
		panic ("Finalize() called on a final object")
	}
	for file,objectNumber := range i.fileBindings {
		file.AddObjectAt (objectNumber, o)
	}
	i.isFinal = true
	return
}

func (i *Indirect) getObjectNumber (f File) ObjectNumber {
	result, ok := i.fileBindings[f]
	if !ok {
		result = f.ReserveObjectNumber(i)
		i.fileBindings[f] = result
	}
	return result
}
