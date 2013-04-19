package pdf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"github.com/mawicks/pdfDIG/containers"
	"github.com/mawicks/pdfDIG/readers" )

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
	// Note: Do not use "file.file" as a writer.  Use "file.writer" instead.
	// "file" must be used for low-level operations such as Seek(), so
	// flush "writer" before using "file".
	writer *bufio.Writer
}

// Public methods

// Constructor for File object
func NewFile(filename string) File {
	var (
		result *file
		f *os.File
		err error
		mode int
	)

	success := false
	modes := [...]int {os.O_RDWR|os.O_CREATE, os.O_RDONLY}
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
	}
	result = new(file)
	result.file = f
	result.mode = mode

	result.xref = containers.NewDynamicArray(1024)
	result.originalSize,_ = f.Seek(0, os.SEEK_END)

	if (result.originalSize == 0) {
		// There is no xref so start one
		result.xref.PushBack(&xrefEntry{0, 65535, false, true})
	} else {
		// For pre-existing files, read the xref
		result.xrefLocation = findXrefLocation(f)
		nextXref := readOneXrefSection(result, result.xrefLocation)
		for ; nextXref != 0; {
			nextXref = readOneXrefSection(result, int64(nextXref))
		}
	}

	result.trailerDictionary = NewDictionary()
	if (result.xrefLocation != 0) {
		result.trailerDictionary.Add ("Prev", NewIntNumeric(int(result.xrefLocation)))
	}

	result.writer = bufio.NewWriter(f)
	if (result.originalSize == 0) {
		writeHeader(result.writer)
	}
	f.Seek(0,os.SEEK_END)
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

// Scan the file for the xref location, returning with the original file position.
func findXrefLocation(f *os.File) (result int64) {
	save,_ := f.Seek(0,os.SEEK_END)
	regexp,_ := regexp.Compile (`\s*FOE%%\s*(\d+)(\s*ferxtrats)`)
	reader := bufio.NewReader(&io.LimitedReader{readers.NewReverseReader(f),512})
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


func readXrefSubsection(xref containers.Array, r *bufio.Reader, start, count uint) {
	var (
		position uint64
		generation uint16
		useChar rune
	)

	// Make sure xref is large enough for the section about to be read.
	if xref.Size() < start+count {
		xref.SetSize(start+count)
	}

	for i:=uint(0); i<count; i++ {
		xrefLine,_ := ReadLine(r)
		n,err := fmt.Sscanf (xrefLine, "%d %d %c", &position, &generation, &useChar)
		if err != nil || n != 3 {
			panic (fmt.Sprintf("Invalid xref line: %s", xrefLine))
		}
		fmt.Printf ("pos: %d gen: %d %c\n", position, generation, useChar)

		if useChar != 'f' && useChar != 'n' {
			panic (fmt.Sprintf("Invalid character '%c' in xref use field.", useChar))
		}
		inUse := (useChar == 'n')

		// Never overwrite a pre-existing entry.
		if *xref.At(start+i) == nil {
			*xref.At(start+i) = &xrefEntry{position, generation, inUse, false}
		}
	}
}

func readTrailer(subsectionHeader string, r *bufio.Reader, pdffile *file) *Dictionary {
	var err error
	tries := 0
	const maxTries = 4
	for tries=0; err == nil && subsectionHeader != "trailer" && tries < maxTries; tries += 1 {
		subsectionHeader,err = ReadLine(r)
	}
	if (err == nil && tries < maxTries) {
		parser := NewParser (r)
		object,err := parser.Scan(pdffile)
		if err == nil {
			if trailer,ok := object.(*Dictionary); ok {
				w := bufio.NewWriter(os.Stdout)
				trailer.Serialize(w,pdffile)
				fmt.Fprintf(w, "\n"); w.Flush()
				return trailer
			}
		}
	}
	return nil
}

func readOneXrefSection (pdffile *file, location int64) (prevXref int) {
	if _,err := pdffile.file.Seek (location, os.SEEK_SET); err != nil {
		panic ("Seeking to xref position failed")
	}
	r := bufio.NewReader(pdffile.file)
	if header,_ := ReadLine(r); header != "xref" {
		panic (`"xref" not found at expected position`)
	}
	subsectionHeader := ""
	for {
		subsectionHeader,_ = ReadLine(r)
		fmt.Printf ("subsection header: %s\n", subsectionHeader)
		start,count := uint(0),uint(0)
		n,err := fmt.Sscanf (subsectionHeader, "%d %d", &start, &count)
		fmt.Printf ("n=%d, err=%v\n", n, err)
		if (err != nil || n != 2) {
			break;
		}
		fmt.Printf ("xref subsection: start=%d, count=%d\n", start, count)
		readXrefSubsection(pdffile.xref, r, start, count)
	}

	trailer := readTrailer (subsectionHeader, r, pdffile)

	if trailer != nil {
		if prevReference,ok := trailer.Get("Prev").(*IntNumeric); ok {
			prevXref = prevReference.Value()
		}
	}
	return
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

	// The catalog appearing in the trailer must be an indirect
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
	var (
		freeIndex uint32
		generation uint16
	)

	// Find an unused node if possible taken from beginning of
	// free list.
	freeIndex = uint32((*f.xref.At(0)).(*xrefEntry).byteOffset)
	if freeIndex == 0 {
		// Create a new xref entry
		freeIndex = uint32(f.xref.Size())
		f.xref.PushBack(&xrefEntry{0, 0, true, true})
	} else {
		entry := (*f.xref.At(uint(freeIndex))).(*xrefEntry)
		// Adjust link in head of free list
		(*f.xref.At(0)).(*xrefEntry).byteOffset = entry.byteOffset
		generation = entry.generation
		entry.inUse = true
	}
	result := ObjectNumber{freeIndex, generation}
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

func dumpXref (xref containers.Array) {
	for i:=uint(0); i<xref.Size(); i++ {
		reference := xref.At(i)
		if reference == nil {
			fmt.Printf ("%d: nil\n", i)
		} else {
			entry := (*reference).(*xrefEntry)
			fmt.Printf ("%d: gen: %d inUse: %v dirty: %v\n", i, entry.generation, entry.inUse, entry.dirty)
		}
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
	dumpXref(f.xref)
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
