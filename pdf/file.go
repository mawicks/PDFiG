package pdf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"github.com/mawicks/PDFiG/containers"
	"github.com/mawicks/PDFiG/readers" )

// xrefEntry type
type xrefEntry struct {
	byteOffset uint64
	generation uint16
	inUse      bool

	// "dirty" is true when the in-memory version of the xref entry doesn't match
	// the "file" copy.
	dirty bool

	// serialization is used to hold a serialized version of the
	// object while the object is being written to disk.  It is
	// used by file.Object() to retrieve a requested object by
	// number that has not yet been written to disk.
	serialization []byte
}

type writeQueueEntry struct {
	index uint32
	xrefEntry *xrefEntry
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

func (entry *xrefEntry) clear (nextFree uint64) {
	if entry.inUse && entry.generation < 65535 {
		entry.generation += 1
	}
	entry.byteOffset = nextFree
	entry.inUse = false
	entry.dirty = true
}

func (entry *xrefEntry) setInUse (location uint64) {
	entry.byteOffset = location
	entry.inUse = true
	entry.dirty = true
}

type file struct {
	pdfVersion uint
	file *os.File
	mode int
	originalSize int64
	// Location of xref for pre-existing files.
	xrefLocation int64
	xref containers.Array

	// trailerDictionary is never nil
	// It is initialized from a pre-existing trailer
	// or is initialized to an empty dictionary
	trailerDictionary *Dictionary

	// catalogIndirect is never nil
	catalogIndirect *Indirect

	// "dirty" is true iff this PDF file requires an update (new
	// xref, new trailer, etc.) when it is closed.
	dirty bool

	// "writer" is a wrapper around "file".
	// Note: Do not use "file.file" as a writer.  Use "file.writer" instead.
	// "file" must be used for low-level operations such as Seek(), so
	// flush "writer" before using "file".
	writer *bufio.Writer
	lastWritePosition int64

	writeQueue chan writeQueueEntry
	writingFinished chan bool

	// semaphore protects access to "file" and to "serialization"
	// member of xref entries
	semaphore chan bool
	closed bool
}

// OpenFile() construct a File object from either a new or a pre-existing filename.
func OpenFile(filename string, mode int) (result *file,exists bool,err error) {
	var f *os.File
	f,err = os.OpenFile(filename, mode, 0666)
	if err != nil {
		return
	}

	result = new(file)
	result.file = f
	result.mode = mode

	result.xref = containers.NewDynamicArray(1024)
	result.originalSize,_ = f.Seek(0, os.SEEK_END)

	if (result.originalSize == 0) {
		// There is no xref so start one
		result.xref.PushBack(&xrefEntry{0, 65535, false, true, nil})
		result.dirty = true
	} else {
		exists = true
		// For pre-existing files, read the xref
		result.xrefLocation = findXrefLocation(f)
		var nextXref int
		nextXref,result.trailerDictionary = readOneXrefSection(result, result.xrefLocation)
		for ; nextXref != 0; {
			nextXref,_ = readOneXrefSection(result, int64(nextXref))
		}
	}
	// If no pre-existing trailer was parsed, create a new dictionary.
	if result.trailerDictionary == nil {
		result.trailerDictionary = NewDictionary()
	}

	// Link the new current trailer to the most recent pre-existing xref.
	if (result.xrefLocation != 0) {
		result.trailerDictionary.Add ("Prev", NewIntNumeric(int(result.xrefLocation)))
	}

	result.writer = bufio.NewWriter(f)
	if (result.originalSize == 0) {
		writeHeader(result.writer)
	}

	result.writeQueue = make(chan writeQueueEntry)
	result.writingFinished = make(chan bool)
	result.semaphore = make(chan bool, 1)
	result.semaphore <- true

	result.lastWritePosition,_ = result.Seek(0,os.SEEK_END)

	go result.gowriter()

	return
}

// Implements WriteObject() in File interface
func (f *file) WriteObject(object Object) (reference *Indirect) {
	returnValue := NewIndirect(f).Write(object)
	return returnValue
}

// Implements DeleteObject() in File interface
func (f *file) DeleteObject(indirect *Indirect) {
	objectNumber := indirect.ObjectNumber(f)
	entry := (*f.xref.At(uint(objectNumber.number))).(*xrefEntry)
	if objectNumber.generation != entry.generation {
		panic("Generation number mismatch")
	}

	if entry.generation < 65535 {
		// Increment the generation count for the next use
		// and link into free list.
		freeHead := (*f.xref.At(0)).(*xrefEntry)
		entry.clear(freeHead.byteOffset)
		freeHead.clear(uint64(objectNumber.number))
	} else {
		// Don't link into free list.  Just set byte offset to 0
		entry.clear(0)
	}

	f.dirty = true
}

// Object() retrieves an object that already exists (or is in the
// process of being written to) a PDF file.  Each call causes a new
// object to be unserialized from the file or a buffer so the caller
// has exclusive ownership of the returned object.
func (f *file) Object(o ObjectNumber) (Object,error) {
	entry := (*f.xref.At(uint(o.number))).(*xrefEntry)
	var r Scanner

	<-f.semaphore
	if entry.serialization == nil {
		f.Seek(int64(entry.byteOffset),os.SEEK_SET)
		r = bufio.NewReader(f.file)
	} else {
		r = bytes.NewReader(entry.serialization)
	}
	object,err := NewParser(r).ScanIndirect(o, f)
	f.semaphore<-true

	return object,err
}

// Implements ReserveObjectNumber() in File interface
func (f *file) ReserveObjectNumber(o Object) ObjectNumber {
	var (
		newNumber uint32
		generation uint16
	)

	// Find an unused node if possible taken from beginning of
	// free list.
	newNumber = uint32((*f.xref.At(0)).(*xrefEntry).byteOffset)
	if newNumber == 0 {
		// Create a new xref entry
		newNumber = uint32(f.xref.Size())
		f.xref.PushBack(&xrefEntry{0, 0, false, true,nil})
	} else {
		// Adjust link in head of free list
		freeHead := (*f.xref.At(0)).(*xrefEntry)
		entry := (*f.xref.At(uint(newNumber))).(*xrefEntry)
		freeHead.clear(entry.byteOffset)

		entry.clear(0)
		generation = entry.generation
	}
	f.dirty = true
	result := ObjectNumber{newNumber, generation}
	return result
}

// Implements Close() in File interface
func (f *file) Close() {
	close(f.writeQueue)
	<- f.writingFinished
	if f.dirty {
		// If client specified a catalog, use it.  Otherwise
		// re-use use pre-existing catalog if it exists.
		if f.catalogIndirect != nil {
			f.trailerDictionary.Add("Root", f.catalogIndirect)
		}

		if f.trailerDictionary.Get("Root") == nil {
			panic("No document catalog has been specified.  Use File.SetCatalog() to set one.")
		}

		f.trailerDictionary.Add("Size", NewIntNumeric(int(f.xref.Size())))

		xrefPosition := f.Tell()

//	 	dumpXref(f.xref)

		f.writeXref()

		f.writeTrailer(xrefPosition)
	}

	f.writer.Flush()
	f.file.Close()

	f.release()
}

func (f *file) Closed() bool {
	return f.closed
}

// ReadLine() reads a line from a PDF file interpreting end-of-line
// characters according to the PDF specification.  In contexts where
// you would be likely to use pdf.ReadLine() are where the line
// consists of ASCII characters.  Therefore ReadLine() returns a
// string rather than a []byte.
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


func (f *file) dictionaryFromTrailer(name string) *Dictionary {
	if infoValue := f.trailerDictionary.Get(name); infoValue != nil {
		indirect := infoValue.(*Indirect)
		if direct,_ := f.Object(indirect.ObjectNumber(f)); direct != nil {
			if info,ok := direct.(*Dictionary); ok {
				return info
			}
		}
	}
	return nil
}

func (f *file) dictionaryToTrailer(name string, d *Dictionary) {
	f.trailerDictionary.Add(name,NewIndirect(f).Write(d))
}

// Catalog() returns the current document catalog of nil if one doesn't
// exist (either from a pre-existing file or from file.SetCatalog())
func (f *file) Catalog() *Dictionary {
	return f.dictionaryFromTrailer("Root")
}

func (f *file) SetCatalog(catalog *Dictionary) {
	f.dictionaryToTrailer("Root",catalog)
}

// Info() returns the current document info dictionary or nil if one
// doesn't exist (either from a pre-existing file or from
// file.SetInfo())
func (f *file) Info() *Dictionary {
	return f.dictionaryFromTrailer("Info")
}

func (f *file) SetInfo(info DocumentInfo) {
	f.dictionaryToTrailer("Info", info.Dictionary)
}

// Trailer() returns the current trailer, which is never nil
func (f *file) Trailer() *Dictionary {
	// Return a clone so nobody can alter the real dictionary
	return f.trailerDictionary.Clone().(*Dictionary)
}

// Using pdf.file.Seek() rather than calling pdf.file.file.Seek()
// directly provides a measure of safety by making sure the internal
// writer is flushed before the file position is moved.
func (f *file) Seek(position int64, whence int) (int64, error) {
	f.writer.Flush()
	return f.file.Seek(position, whence)
}

func (f *file) Tell() int64 {
	// Make sure to use the flushing version of Seek() here...
	position, _ := f.Seek(0, os.SEEK_CUR)
	return position
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

		if useChar != 'f' && useChar != 'n' {
			panic (fmt.Sprintf("Invalid character '%c' in xref use field.", useChar))
		}
		inUse := (useChar == 'n')

		// Never overwrite a pre-existing entry.
		if *xref.At(start+i) == nil {
			*xref.At(start+i) = &xrefEntry{position, generation, inUse, false, nil}
		}
	}
}

func readTrailer(subsectionHeader string, r *bufio.Reader, f *file) *Dictionary {
	var err error
	tries := 0
	const maxTries = 4
	for tries=0; err == nil && subsectionHeader != "trailer" && tries < maxTries; tries += 1 {
		subsectionHeader,err = ReadLine(r)
	}
	if (err == nil && tries < maxTries) {
		parser := NewParser (r)
		object,err := parser.Scan(f)
		if err == nil {
			if trailer,ok := object.(*Dictionary); ok {
				return trailer
			}
		}
	}
	return nil
}

func readOneXrefSection (f *file, location int64) (prevXref int, trailer *Dictionary) {

	if _,err := f.file.Seek (location, os.SEEK_SET); err != nil {
		panic ("Seeking to xref position failed")
	}

	r := bufio.NewReader(f.file)
 	if header,_ := ReadLine(r); header != "xref" {
		panic (`"xref" not found at expected position`)
	}

	subsectionHeader := ""
	for {
		subsectionHeader,_ = ReadLine(r)
		start,count := uint(0),uint(0)
		n,err := fmt.Sscanf (subsectionHeader, "%d %d", &start, &count)
		if (err != nil || n != 2) {
			break;
		}
		readXrefSubsection(f.xref, r, start, count)
	}

	trailer = readTrailer (subsectionHeader, r, f)
	if trailer == nil {
		panic ("Expected trailer not found")
	} else if prevReference,ok := trailer.Get("Prev").(*IntNumeric); ok {
		prevXref = prevReference.Value()
	}
	return
}

func (f *file) release() {
	f.file = nil
	f.xref.SetSize(0)
	f.xref = nil
	f.trailerDictionary = nil
	f.catalogIndirect = nil
	f.catalogIndirect = nil
	f.writer = nil
	f.writeQueue = nil
	f.writingFinished = nil
	f.semaphore = nil
	f.closed = true
}

func (f* file) gowriter () {
	for entry := range f.writeQueue {
		entry.xrefEntry.setInUse(uint64(f.lastWritePosition))

		<-f.semaphore
		f.Seek(f.lastWritePosition, os.SEEK_SET)
		fmt.Fprintf(f.writer, "%d %d obj\n", entry.index, entry.xrefEntry.generation)

		_,err := f.writer.Write(entry.xrefEntry.serialization)
		if err != nil {
			panic(errors.New("Unable to write serialized object in file.writeObject()"))
		}
		fmt.Fprintf(f.writer, "\nendobj\n")
		f.lastWritePosition = f.Tell()
		entry.xrefEntry.serialization = nil
		f.semaphore<-true

		f.dirty = true
	}
	f.writingFinished <- true
}

// Implements WriteObjectAt() in File interface
func (f *file) WriteObjectAt(objectNumber ObjectNumber, object Object) {
	xrefEntry := (*f.xref.At(uint(objectNumber.number))).(*xrefEntry)
	if xrefEntry.generation != objectNumber.generation {
		panic(fmt.Sprintf("Generation number mismatch: object %d current generation is %d but attempted to write %d",
			objectNumber.number, xrefEntry.generation, objectNumber.generation))
	}
	buffer := new(bytes.Buffer)
	object.Serialize(buffer, f)
	xrefEntry.serialization = buffer.Bytes()
	f.writeQueue<-writeQueueEntry{objectNumber.number,xrefEntry}
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
