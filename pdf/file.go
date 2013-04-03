package pdf

import "bufio"
import "fmt"
import "os"
import "maw/containers"

// TestFile is a simple file implementing the File interface for use in unit tests.
type testFile struct {
	nextObjectNumber uint32
	nextGenerationNumber uint16
}

// Constructor for testFile object
func NewTestFile (obj uint32, gen uint16) File {
	return &testFile{obj,gen}
}

// Public methods

// Implements Close() in File interface
func (f *testFile) Close() {}

// Implements AddObjectAt() in File interface
func (f *testFile) AddObjectAt (ObjectNumber, Object) {}

// Implements AddObject() in File interface
func (f *testFile) AddObject (object Object) (objectNumber ObjectNumber) {
	objectNumber = f.ReserveObjectNumber (object)
	f.AddObjectAt (objectNumber, object)
	return objectNumber
}

// Implements DeleteObject() in File interface
func (f *testFile) DeleteObject (on ObjectNumber) {}

// Implements ReserveObjectNumber() in File interface
func (f *testFile) ReserveObjectNumber (o Object) ObjectNumber {
	result := ObjectNumber{f.nextObjectNumber,f.nextGenerationNumber}
	f.nextObjectNumber += 1
	f.nextGenerationNumber += 1
	return result
}

// xrefEntry type
type xrefEntry struct {
	byteOffset uint64
	generation uint16
	inUse bool

	// "dirty" is true when the in-memory version of the object doesn't match
	// the "file" copy.
	dirty bool
}

// Write xrefEntry to output stream using Writer.
func (entry *xrefEntry) Serialize (w Writer) {
	fmt.Fprintf (w,
		"%010d %05d %c \n",
		entry.byteOffset,
		entry.generation,
		func (inuse bool) (result rune) {
			if inuse {
				result = 'n'
			} else {
				result = 'f'
			}
			return result
		} (entry.inUse))
}

type file struct {
	xref containers.Array
	trailerDictionary *Dictionary
	file *os.File

	// "writer" is a wrapper around "file".
	// Note: Do not use "file" as a writer.  Use "writer" instead.
	// "file" must be used for low-level operations such as Seek(),
	// flush "writer" before using "file".
	writer *bufio.Writer
}

// Constructor for File object
func NewFile (filename string) File {
	var result *file
	f, error := os.Create (filename)
	if error != nil {
		panic ("Failed to create file")
	} else {
		result = new(file)
		result.file = f
		result.writer = bufio.NewWriter(f)
		result.xref = containers.NewDynamicArray(1024)
		result.trailerDictionary = NewDictionary()

		result.writePdfHeader()
		result.createInitialXref()
		result.writer.Flush()
	}

	return result
}

// Public methods
// Implements Close() in File interface
func (f *file) Close () {
	xrefPosition := f.Tell()
	f.writeXref()
	f.writeTrailer(xrefPosition)
	f.writer.Flush()
	f.file.Close()
}

func (f *file) Seek(position int64, whence int) (int64,error) {
	f.writer.Flush()
	return f.file.Seek(position, whence)
}

func (f *file) Tell() (int64) {
	position,_ := f.Seek(0, 1)
	return position
}


// Implements AddObjectAt() in File interface
func (f *file) AddObjectAt (object ObjectNumber, o Object) {
	entry := (*f.xref.At(uint(object.number))).(*xrefEntry)
	if (entry.byteOffset !=  0) {
		panic ("An object has already been written with this number")
	}
	if (entry.generation != object.generation) {
		panic ("Generation number mismatch")
	}

	entry.byteOffset = uint64(f.Tell())

	fmt.Fprintf(f.writer, "%d %d obj\n", object.number, object.generation)
	o.Serialize(f.writer, f);
	fmt.Fprintf(f.writer, "\nendobj\n")
}

// Implements AddObject() in File interface
func (f *file) AddObject (object Object) (objectNumber ObjectNumber) {
	objectNumber = f.ReserveObjectNumber (object)
	f.AddObjectAt (objectNumber, object)
	return objectNumber
}

// Implements DeleteObject() in File interface
func (f *file) DeleteObject (object ObjectNumber) {
	entry := (*f.xref.At(uint(object.number))).(*xrefEntry)
	if (object.generation != entry.generation) {
		panic ("Generation number mismatch")
	}

	if (entry.generation < 65535) {
		// Increment the generation count for the next use
		// and link into free list.
		entry.generation += 1;
		entry.byteOffset = (*f.xref.At(0)).(*xrefEntry).byteOffset
		(*f.xref.At(0)).(*xrefEntry).byteOffset = uint64(object.number)
	} else {
		// Don't link into free list.  Just set byte offset to 0
		entry.byteOffset = 0
	}

	entry.inUse = false
	entry.dirty = true
}

// Implements ReserveObjectNumber() in File interface
func (f *file) ReserveObjectNumber (o Object) ObjectNumber {
	var nextUnused uint
	var generation uint16

	// Find an unused node if possible.  Begin searching at
	// index=1 because first record begins free list and is always
	// marked as free.
	for nextUnused=1;
	    nextUnused<f.xref.Size() &&
		    (*f.xref.At(nextUnused)).(*xrefEntry).generation < 65535 &&
		    (*f.xref.At(nextUnused)).(*xrefEntry).inUse;
	    nextUnused++ {
		// Empty loop
	}

	if (nextUnused >= f.xref.Size()) {
		// Create a new xref entry
		f.xref.PushBack(&xrefEntry{0,0,true,true})
	} else {
		entry := (*f.xref.At(nextUnused)).(*xrefEntry)
		// Adjust link in head of free list
		(*f.xref.At(0)).(*xrefEntry).byteOffset = entry.byteOffset
		generation = entry.generation
		entry.inUse = true
	}
	result := ObjectNumber{uint32(nextUnused), generation}
	return result
}


func (f *file) parseExistingFile() {
	panic ("Not implemented")
}

func (f *file) createInitialXref() {
	f.xref.PushBack(&xrefEntry{0,65535,false,true})
}

func (f *file) writePdfHeader () {
	f.writer.WriteString ("%PDF-1.4\n")
}

func nextSegment (xref containers.Array, start uint) (nextStart, length uint) {
	var i uint
	// Skip "clean" entries.
	for i = start; i<xref.Size() && !(*xref.At(i)).(*xrefEntry).dirty; i++ {}

	nextStart = i
	for i = nextStart; i<xref.Size() && (*xref.At(i)).(*xrefEntry).dirty; i++ {
		length += 1
	}

	return nextStart, length
}

func (f *file) writeXref() {
	f.writer.WriteString ("xref\n")

	for s,l:=nextSegment(f.xref,0); s<f.xref.Size(); s,l=nextSegment(f.xref, s+l) {
		fmt.Fprintf (f.writer, "%d %d\n", s, l)
		for i:=s; i<s+l; i++ {
			entry := (*f.xref.At(uint(i))).(*xrefEntry)
			if (entry.byteOffset == 0 && entry.inUse) {
				panic (fmt.Sprintf ("Object %d reserved but never added or finalized", i))
			}
			entry.Serialize(f.writer)
		}
	}
}

func (f *file) writeTrailer(xrefPosition int64) {
	f.writer.WriteString ("trailer\n")
	f.trailerDictionary.Serialize (f.writer)
	f.writer.WriteString ("\nstartxref\n")
	fmt.Fprintf (f.writer, "%d\n", xrefPosition)
	f.writer.WriteString ("%%EOF\n")
}
