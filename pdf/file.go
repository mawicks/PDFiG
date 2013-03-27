/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bufio"
import "fmt"
import "os"
import "maw/containers"

// TestFile is a simple file implementing the File interface for use in unit tests.
type TestFile struct {
	nextObjectNumber uint32
	nextGenerationNumber uint16
}

// Constructor for Stream object
func NewTestFile (obj uint32, gen uint16) File {
	return &TestFile{obj,gen}
}

func (f *TestFile) AssignObjectNumber (o Object) ObjectNumber {
	result := ObjectNumber{f.nextObjectNumber,f.nextGenerationNumber}
	f.nextObjectNumber += 1
	f.nextGenerationNumber += 1
	return result
}

func writePdfHeader (w Writer) {
	w.WriteString ("%PDF-1.4\n")
}

// xrefEntry type
type xrefEntry struct {
	byteOffset uint64
	generation uint16
	inUse bool
	dirty bool
}

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

	// Note: Do not use "file" as a writer.  Use "writer" instead.
	// Flush "writer" before using "file".

	writer *bufio.Writer
	objectCount uint32
}

func (f *file) AssignObjectNumber (o Object) ObjectNumber {
	result := ObjectNumber{f.objectCount,0}
	f.objectCount += 1
	return result
}

func (f *file) parseExistingFile() {
panic ("Not implemented")
}

func (f *file) createInitialXref() {
	f.xref.PushBack(&xrefEntry{0,65535,false,true})
}

func (f *file) writeXref() {
	f.writer.WriteString ("xref\n")
	for i:= uint(0); i<f.xref.Size(); i++ {
		(*f.xref.At(uint(i))).(*xrefEntry).Serialize(f.writer)
	}
}

func NewFile (filename string) File {
	var result *file
	f, error := os.Create (filename)
	if error != nil {
		panic ("Failed to create file")
	} else {
		result = new(file)
		result.xref = containers.NewDynamicArray(1024)
		result.trailerDictionary = NewDictionary()
		result.file = f
		result.writer = bufio.NewWriter(f)
		result.objectCount = 0

		writePdfHeader (result.writer)
		result.createInitialXref()
		result.writeXref()
		result.writer.Flush()
	}
	return result
}
