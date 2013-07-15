package pdf

import (
	"fmt"
	"io"
	"errors"
	"github.com/mawicks/PDFiG/readers"
	"strconv" )

type Scanner interface {
	io.ByteScanner
	io.Reader
//	Peek(int) ([]byte, error)
}

type Parser struct {
	scanner *readers.HistoryReader
	queuedObject Object
}

// NewParser constructs a new parser from the passed Scanner.
// Typically Scanner will be the pdf.File's underlying os.File, but
// this is not strictly necessary.
func NewParser(scanner Scanner) *Parser {
	return &Parser{readers.NewHistoryReader(scanner,64),nil}
}

var (
	invalidKeyword = errors.New(`Invalid keyword`)
	expectingDigit = errors.New(`Expecting digit`)
	parsingError = errors.New(`Parsing error`)
	unknownIOError = errors.New(`Unknown I/O error`)
	expectingName = errors.New(`Expecting PDF name`)
	unexpectedEnd = errors.New(`Unexpected end of input`)
	unexpectedInput = errors.New(`Unexpected character or end of input`)
	expectedGreaterThan = errors.New(`Expected ">"`)
	expectingHexDigit = errors.New(`Expecting hex digit`)
	expectingOctalDigit = errors.New(`Expecting octal digit`) )

// Skip white space and return the byte following the white space or error.
// If err is non-nil, the value of b is undefined.
func nextNonWhiteByte (scanner Scanner) (b byte,err error) {
	b, err = scanner.ReadByte()
	for ; err == nil && (IsWhiteSpace(b) || b == '%'); b,err=scanner.ReadByte() {
		if b == '%' {
			ReadLine(scanner)
		}
	}
	return
}

func scanKeyword (scanner Scanner, b byte) (string,error) {
	var buffer[]byte = make([]byte, 0, 5)
	buffer = append(buffer, b)

	b,err := scanner.ReadByte()
	for ; err == nil && IsAlpha(b); b,err=scanner.ReadByte() {
		buffer = append(buffer, b)
	}

	if err == nil {
		scanner.UnreadByte()
	} else if err == io.EOF {
		err = nil
	}

	return string(buffer),err
}

func scanKeywordObject (scanner Scanner, b byte) Object {
	keyword,err := scanKeyword(scanner, b)
	if err == nil {
		switch (keyword) {
		case "null":
			return NewNull()
		case "true":
			return NewBoolean(true)
		case "false":
			return NewBoolean(false)
		}
	}
	panic(invalidKeyword)
}

func scanNumeric (scanner Scanner, b byte) Object {
	var buffer[]byte = make([]byte, 0, 5)
	var err error

	hasAtLeastOneDigit := false
	float := false

	if (b == '+' || b == '-') {
		buffer = append(buffer,b)
		b,err=scanner.ReadByte()
	}

	for ; err==nil && IsDigit(b); b,err=scanner.ReadByte() {
		hasAtLeastOneDigit = true
		buffer = append(buffer,b)
	}

	if (err == nil && b == '.') {
		float = true
		buffer = append(buffer,b)
		b,err=scanner.ReadByte()
	}

	for ; err==nil && IsDigit(b); b,err=scanner.ReadByte() {
		hasAtLeastOneDigit = true
		buffer = append(buffer,b)
	}

	if err==nil {
		scanner.UnreadByte()
	}

	if err != nil && err!=io.EOF {
		panic(unknownIOError)
	}

	if !hasAtLeastOneDigit {
		panic(expectingDigit)
	}

	if float {
		number,_ := strconv.ParseFloat(string(buffer),32)
		return NewRealNumeric(float32(number))
	}
	number,_ := strconv.ParseInt(string(buffer),10,32)
	return NewIntNumeric(int(number))
}

func (p *Parser) scanNumericOrIndirectRef(b byte, file... File) Object {
	var n1 Object
	var err error

	if (p.queuedObject != nil) {
		p.scanner.UnreadByte()
		n1 = p.queuedObject
		p.queuedObject = nil
	} else {
		n1 = scanNumeric(p.scanner, b)
	}

	if _,ok := n1.(*IntNumeric); !ok {
		return n1
	}

	b,err = nextNonWhiteByte(p.scanner)
	if err == nil && !IsDigit(b) {
		p.scanner.UnreadByte()
	}


	if (err != nil || !IsDigit(b)) {
		return n1
	}

	n2 := scanNumeric (p.scanner, b)
	if _,ok := n2.(*IntNumeric); !ok {
		if (p.queuedObject != nil) {
			panic ("Queued object is not nil. This shouldn't happen")
		}
		p.queuedObject = n2
		return n1
	}

	b,err = nextNonWhiteByte(p.scanner)
	if b != 'R' {
		p.scanner.UnreadByte()
		p.queuedObject = n2
		return n1
	} else {
		number := uint32(n1.(*IntNumeric).Value())
		generation := uint16(n2.(*IntNumeric).Value())
		return file[0].Indirect(ObjectNumber{number,generation})
	}
	return nil
}


func scanDigitWithBase(scanner Scanner, parseDigit func (byte) byte) byte {
	b,err := scanner.ReadByte()
	if err != nil {
		panic(unexpectedEnd)
	}
	return parseDigit(b)
}

func scanHexDigit (scanner Scanner) byte {
	return scanDigitWithBase (scanner, ParseHexDigit)
}

func scanOctalDigit (scanner Scanner) byte {
	return scanDigitWithBase (scanner, ParseOctalDigit)
}

func scanName (scanner Scanner) Object {
	var buffer[]byte = make([]byte, 0, 8)
	b,err := scanner.ReadByte()
	for ; err == nil && IsRegular(b); b,err=scanner.ReadByte() {
		if (b != '#') {
			buffer = append(buffer, b)
		} else {
			r := byte(0)
			for i:=0; i<2; i++ {
				r = r<<4 + scanHexDigit(scanner)
			}
			buffer = append(buffer, r)
		}
	}
	if err == nil {
		scanner.UnreadByte()
	}
	if err == nil || err == io.EOF {
		return NewName(string(buffer))
	}
	panic (unknownIOError)
}

func scanEscape (scanner Scanner) (b byte) {
	var err error
	if b,err =scanner.ReadByte(); err != nil {
		panic (unexpectedEnd)
	}

	if IsOctalDigit(b) {
		r := ParseOctalDigit(b)
		for i:=0; i<2; i++ {
			r = r<<3 + scanOctalDigit(scanner)
		}
		return r
	}

	switch b {
	case 'n':
		b = '\n'
	case 'r':
		b = '\r'
	case 't':
		b = '\t'
	case 'b':
		b = '\b'
	case 'f':
		b = '\f'
	}
	return
}

func scanNormalString (scanner Scanner) *String {
	var openCount = 0
	var buffer[]byte = make([]byte, 0, 128)
	b,err :=scanner.ReadByte()
	for ; err == nil && (b!=')' || openCount != 0); b,err=scanner.ReadByte() {
		switch b {
		case '(':
			openCount += 1;
			buffer = append(buffer, b)
		case ')':
			openCount -= 1;
			buffer = append(buffer, b)
		case '\\':
			v := scanEscape (scanner)
			buffer = append(buffer, v)
		default:
			buffer = append(buffer, b)
		}
	}
	if err != nil {
		panic (unexpectedEnd)
	}
	return NewBinaryString(buffer)
}

func scanHexString (scanner Scanner,b byte) *String {
	var buffer[]byte = make([]byte, 0, 128)
	var err error
	for ; err == nil && b != '>'; b,err=scanner.ReadByte() {
		scanner.UnreadByte()
		r := byte(0)
		for i:=0; i<2; i++ {
			r = r<<4 + scanHexDigit(scanner)
		}
		buffer = append(buffer, r)
	}
	if err != nil {
		panic(unexpectedEnd)
	}
	return NewBinaryString(buffer)
}

func (p *Parser) scanArray (file... File) Array {
	var array Array = NewArray()

	b,err := nextNonWhiteByte(p.scanner)
	for ; p.queuedObject != nil || (err == nil && b != ']'); b,err=p.scanner.ReadByte() {
		p.scanner.UnreadByte()
		nextElement := p.scanObject(file...)
		array.Add(nextElement)
	}
	if err != nil {
		panic(unexpectedEnd)
	}
	return array
}

func (p *Parser) scanDictionary(file... File) Dictionary {
	var d Dictionary = NewDictionary()

	b,err := nextNonWhiteByte(p.scanner)
	for ; err == nil && b != '>'; b,err=p.scanner.ReadByte() {
		p.scanner.UnreadByte()
		name,ok := p.scanObject().(*Name)
		if (!ok) {
			panic(expectingName)
		}
		object := p.scanObject(file...)
		d.Add(name.String(),object)
	}

	if err != nil {
		panic(unexpectedEnd)
	}

	b,err = nextNonWhiteByte(p.scanner)
	if (b != '>') {
		panic(expectedGreaterThan)
	}
	return d
}

func (p *Parser) scanDictionaryOrStream (file... File) Object {
	var err error

	dictionary := p.scanDictionary(file...)

	var b byte
	b,err = nextNonWhiteByte(p.scanner)
	p.scanner.UnreadByte()

	var s string
	// Could be a "stream" line.
	if b=='s' {
		s,err = ReadLine (p.scanner)
	}

	var stream Object
	if err == nil && s == "stream" {
		v,ok := dictionary.Get("Length").(*IntNumeric)
		if ok {
			length := v.Value()
			contents := make([]byte, length)
			p.scanner.Read(contents)
			nextNonWhiteByte(p.scanner)
			p.scanner.UnreadByte()
			s,err = ReadLine(p.scanner)
			if err == nil && s == "endstream" {
				stream = NewStreamFromContents (dictionary,contents,nil)
			}
		}
	}
	if stream != nil {
		return stream
	}
	return dictionary
}

func (p *Parser) scanObject(file ...File) Object {
	// If there's a non-integer object left parsed during a previous
	// call, go ahead and return it.
	if p.queuedObject != nil {
		if _,ok := p.queuedObject.(*IntNumeric); !ok {
			object := p.queuedObject
			p.queuedObject = nil
			return object
		}
	}
	b,err := nextNonWhiteByte(p.scanner)
	if err == nil {
		switch  {
		case IsAlpha(b):
			return scanKeywordObject(p.scanner, b)
		case IsDigit(b),p.queuedObject != nil:
			return p.scanNumericOrIndirectRef(b, file...)
		case b=='.',b=='+',b=='-':
			return scanNumeric(p.scanner, b)
		case b =='/':
			return scanName (p.scanner)
		case b=='(':
			return scanNormalString(p.scanner)
		case b=='<':
			b,err = nextNonWhiteByte(p.scanner)
			if b == '<' {
				return p.scanDictionaryOrStream(file...)
			} else {
				return scanHexString(p.scanner, b)
			}
		case b=='[':
			return p.scanArray(file...)
		}
	}
	panic(unexpectedInput)
}

// Scan() parses an arbitrary object.  If successful, the object is
// returned.  If not, err is set and context contains the input bytes
// that preceeded the error.  The optional File argument, of which
// there should be no more than one, indicates the pdf.File to use to
// resolve indirect object references (e.g., "25 0 R").  If no File is
// supplied, the input stream may not contain any indirect object
// references.
func (p *Parser) Scan(file... File) (o Object,err error) {
	defer func() {
		if x := recover(); x!= nil {
			err,_ = x.(error)
		}
	} ()

	o = p.scanObject(file...)

	return
}

// ScanIndirect() parses an indirect object including the "%d %d obj"
// header and "endobj" trailer.  If successful the object is returned.
// It returns an error if the object number and generation do not
// match the passed ObjectNumber.  The optional File argument is as
// described in Parser.Scan().
func (p *Parser) ScanIndirect(objectNumber ObjectNumber, file... File) (object Object,err error) {
	defer func() {
		if x := recover(); x!= nil {
			err,_ = x.(error)
		}
	} ()

	header,_ := ReadLine(p.scanner)
	var (
		index uint32
		generation uint16
		obj string )

	n,err := fmt.Sscanf (header, "%d %d %s", &index, &generation, &obj)
	if err != nil || n != 3 {
		panic(errors.New(fmt.Sprintf(`Object header expected but not found at position %p`, p)))
	}
	if (objectNumber.number != index || objectNumber.generation != generation) {
		panic(errors.New(fmt.Sprintf(`Expected "%d %d obj" at location %d but found "%d %d %s"`,
			objectNumber.number, objectNumber.generation,
			index, generation, obj)))
	}
	object = p.scanObject(file...)
	nextNonWhiteByte(p.scanner)
	p.scanner.UnreadByte()

	trailer,_ := ReadLine(p.scanner)
	if trailer != "endobj" {
		panic(errors.New(fmt.Sprintf(`No "endobj" following object "%d %d obj"`,
			objectNumber.number, objectNumber.generation)))
	}
	return object,err
}


func (p *Parser) GetContext() []byte {
	return p.scanner.GetHistory()
}
