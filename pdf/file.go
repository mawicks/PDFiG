package pdf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"github.com/mawicks/pdfdig/containers"
	"github.com/mawicks/pdfdig/readers" )

// xrefEntry type
type xrefEntry struct {
	byteOffset uint64
	generation uint16
	inUse      bool

	// "dirty" is true when the in-memory version of the object doesn't match
	// the "file" copy.
	dirty bool
}

// Write xrefEntry to output stream using Writer.
func (entry *xrefEntry) Serialize(w Writer) {
	fmt.Fprintf(w,
		"%010d %05d %c \n",
		entry.byteOffset,
		entry.generation,
		func(inuse bool) (result rune) {
			if inuse {
				result = 'n'
			} else {
				result = 'f'
			}
			return result
		}(entry.inUse))
}

type file struct {
	file *os.File
	originalSize int64
	mode int
	xrefLocation int64
	xref containers.Array
	trailerDictionary *Dictionary
	catalogIndirect *Indirect

	// "writer" is a wrapper around "file".
	// Note: Do not use "file" as a writer.  Use "writer" instead.
	// "file" must be used for low-level operations such as Seek(),
	// flush "writer" before using "file".
	writer *bufio.Writer
}

// Public methods

// Constructor for File object
func NewFile(filename string) File {
	var result *file

	success := false
	modes := [...]int {os.O_RDWR|os.O_CREATE, os.O_RDONLY}
	var f *os.File
	var err error
	var mode int
	for _,m := range modes {
		f, err = os.OpenFile(filename, m, 0777)
		if (err == nil) {
			mode = m
			success = true
			break;
		}
	}
	if !success {
		panic("Failed to open or create file")

	} else {
		result = new(file)
		result.file = f
		result.mode = mode
		result.trailerDictionary = NewDictionary()

		result.xref = containers.NewDynamicArray(1024)
		result.xref.PushBack(&xrefEntry{0, 65535, false, true})

		result.originalSize,_ = f.Seek(0,os.SEEK_END)

		if (result.originalSize != 0) {
			result.xrefLocation = getXrefLocation(f)
			getXref(f, result.xrefLocation)
		}

		result.writer = bufio.NewWriter(f)
		if (result.originalSize == 0) {
			writeHeader(result.writer)
		}
	}

	return result
}

// Read a line from a PDF file interpreting end-of-line characters
// according to the PDF specification.  In contexts where you would be
// likely to use pdf.ReadLine() are where the line consists of ASCII
// characters.  Therefore ReadLine() returns a string rather than a
// []byte.
func ReadLine(r io.ByteScanner) (result string, err error) {
	bytes := make([]byte, 0, 512)
	var byte byte
	for byte,err=r.ReadByte(); err==nil && byte!='\r' && byte!='\n'; byte,err=r.ReadByte() {
		bytes = append(bytes, byte)
	}
	// Gobble up a second end-of-line character, if present.
	// Don't gobble up two identical end-of-line-characters as
	// logically they represent two separate lines.
	if err==nil {
		secondbyte,err2:=r.ReadByte()
		if err2==nil && (secondbyte==byte || (secondbyte!='\r' && secondbyte!='\n')) {
			r.UnreadByte()
		}
	}
	if err==io.EOF {
		err = nil
	}
	result = string(bytes)
	return
}

// Parse the file for the file for the xref location, leaving the file position unchanged.
func getXrefLocation(f *os.File) (result int64) {
	save,_ := f.Seek(0,os.SEEK_END)
	regexp,_ := regexp.Compile (`\s*FOE%%\s*(\d+)(\s*ferxtrats)`)
	reader := bufio.NewReader(readers.NewReverseReader(f))
	indexes := regexp.FindReaderSubmatchIndex(reader)

	if (indexes != nil) {
		f.Seek(-int64(indexes[3]),os.SEEK_END)
		b := make([]byte,indexes[3]-indexes[2])
		_,err := f.Read(b)
		if (err == nil) {
			result,_ = strconv.ParseInt(string(b),10,64)
			fmt.Printf ("Xref location is %d\n", result)
		}
	}
	// Restore file position
	f.Seek(save,os.SEEK_SET)
	return result
}

func getXref (f* os.File, location int64) {
	fmt.Printf ("in getXref()\n")
	if _,err := f.Seek (location, os.SEEK_SET); err == nil {
		r := bufio.NewReader(f)
		if header,_ := ReadLine(r); header == "xref" {
			fmt.Printf ("Found xref!\n")
		} else {
			fmt.Printf ("Didn't find xref\n")
		}
	}
}

func (f *file) SetCatalog(catalog *Indirect) {
	f.catalogIndirect = catalog
}

func (f *file) SetInfo(i *Indirect) {
	f.trailerDictionary.Add("Info", i)
}

func (f *file) release() {
	f.xref.SetSize(0)
	f.catalogIndirect = nil
	f.trailerDictionary = nil
	f.file = nil
}

// Implements Close() in File interface
func (f *file) Close() {
	f.trailerDictionary.Add("Size", NewIntNumeric(int(f.xref.Size())))

	// The catalog appearing in the trailer must be indirect
	// object.  Create an indirect object pointed at the catalog,
	// add it to the trailer dictionary, and write out the trailer dictionary.
	if f.catalogIndirect == nil {
		panic("No document catalog has been specified.  Use File.SetCatalog() to set one.")
	}
	f.trailerDictionary.Add("Root", f.catalogIndirect)

	xrefPosition := f.Tell()
	f.writeXref()

	f.writeTrailer(xrefPosition)
	f.writer.Flush()
	f.file.Close()

	f.release()
}

func (f *file) Seek(position int64, whence int) (int64, error) {
	f.writer.Flush()
	return f.file.Seek(position, whence)
}

func (f *file) Tell() int64 {
	position, _ := f.Seek(0, os.SEEK_CUR)
	return position
}

// Implements AddObjectAt() in File interface
func (f *file) AddObjectAt(object ObjectNumber, o Object) {
	entry := (*f.xref.At(uint(object.number))).(*xrefEntry)
	if entry.byteOffset != 0 {
		panic("An object has already been written with this number")
	}
	if entry.generation != object.generation {
		panic("Generation number mismatch")
	}

	entry.byteOffset = uint64(f.Tell())

	fmt.Fprintf(f.writer, "%d %d obj\n", object.number, object.generation)
	o.Serialize(f.writer, f)
	fmt.Fprintf(f.writer, "\nendobj\n")
}

// Implements AddObject() in File interface
func (f *file) AddObject(object Object) (reference *Indirect) {
	reference = NewIndirect(f)
	reference.Finalize(object)
	return reference
}

// Implements DeleteObject() in File interface
func (f *file) DeleteObject(object ObjectNumber) {
	entry := (*f.xref.At(uint(object.number))).(*xrefEntry)
	if object.generation != entry.generation {
		panic("Generation number mismatch")
	}

	if entry.generation < 65535 {
		// Increment the generation count for the next use
		// and link into free list.
		entry.generation += 1
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
func (f *file) ReserveObjectNumber(o Object) ObjectNumber {
	var nextUnused uint
	var generation uint16

	// Find an unused node if possible.  Begin searching at
	// index=1 because first record begins free list and is always
	// marked as free.
	for nextUnused = 1; nextUnused < f.xref.Size() &&
		(*f.xref.At(nextUnused)).(*xrefEntry).generation < 65535 &&
		(*f.xref.At(nextUnused)).(*xrefEntry).inUse; nextUnused++ {
		// Empty loop
	}

	if nextUnused >= f.xref.Size() {
		// Create a new xref entry
		f.xref.PushBack(&xrefEntry{0, 0, true, true})
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
	panic("Not implemented")
}

func writeHeader(w *bufio.Writer) {
	_,err := w.WriteString("%PDF-1.4\n")
	if (err != nil) {
		panic("Unable to write PDF header")
	}
}

func nextSegment(xref containers.Array, start uint) (nextStart, length uint) {
	var i uint
	// Skip "clean" entries.
	for i = start; i < xref.Size() && !(*xref.At(i)).(*xrefEntry).dirty; i++ {
	}

	nextStart = i
	for i = nextStart; i < xref.Size() && (*xref.At(i)).(*xrefEntry).dirty; i++ {
		length += 1
	}

	return nextStart, length
}

func (f *file) writeXref() {
	f.writer.WriteString("xref\n")

	for s, l := nextSegment(f.xref, 0); s < f.xref.Size(); s, l = nextSegment(f.xref, s+l) {
		fmt.Fprintf(f.writer, "%d %d\n", s, l)
		for i := s; i < s+l; i++ {
			entry := (*f.xref.At(uint(i))).(*xrefEntry)
			if entry.byteOffset == 0 && entry.inUse {
				panic(fmt.Sprintf("Object %d reserved but never added or finalized", i))
			}
			entry.Serialize(f.writer)
		}
	}
}

func (f *file) writeTrailer(xrefPosition int64) {
	f.writer.WriteString("trailer\n")
	f.trailerDictionary.Serialize(f.writer, f)
	f.writer.WriteString("\nstartxref\n")
	fmt.Fprintf(f.writer, "%d\n", xrefPosition)
	f.writer.WriteString("%%EOF\n")
}
