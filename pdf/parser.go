package pdf

import (
	"io"
	"errors"
	"github.com/mawicks/pdftig/readers"
	"strconv" )

type Scanner interface {
	io.ByteScanner
	io.Reader
//	Peek(int) ([]byte, error)
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
		return NewFloatNumeric(float32(number))
	}
	number,_ := strconv.ParseInt(string(buffer),10,32)
	return NewIntNumeric(int(number))
}

func scanDigitWithBase (scanner Scanner, parseDigit func (byte) byte) byte {
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

func scanArray (scanner Scanner) *Array {
	var array *Array = NewArray()

	b,err := nextNonWhiteByte(scanner)
	for ; err == nil && b != ']'; b,err=scanner.ReadByte() {
		scanner.UnreadByte()
		nextElement := scanObject(scanner)
		array.Add(nextElement)
	}

	if err != nil {
		panic(unexpectedEnd)
	}
	return array
}

func scanDictionary(scanner Scanner) *Dictionary {
	var d *Dictionary = NewDictionary()

	b,err := nextNonWhiteByte(scanner)
	for ; err == nil && b != '>'; b,err=scanner.ReadByte() {
		scanner.UnreadByte()
		name := scanObject(scanner).(*Name)
		if (name == nil) {
			panic(expectingName)
		}
		object := scanObject(scanner)
		d.Add(name.String(),object)
	}

	if err != nil {
		panic(unexpectedEnd)
	}

	b,err = nextNonWhiteByte(scanner)
	if (b != '>') {
		panic(expectedGreaterThan)
	}
	return d
}

func scanDictionaryOrStream (scanner Scanner) Object {
	var err error

	dictionary := scanDictionary(scanner)

	var b byte
	b,err = nextNonWhiteByte(scanner)
	scanner.UnreadByte()

	var s string
	// Could be a "stream" line.
	if b=='s' {
		s,err = ReadLine (scanner)
	}

	var stream Object
	if err == nil && s == "stream" {
		v := dictionary.Get("Length").(*IntNumeric)
		if v != nil {
			length := v.Value()
			contents := make([]byte, length)
			scanner.Read(contents)
			nextNonWhiteByte(scanner)
			scanner.UnreadByte()
			s,err = ReadLine(scanner)
			if err == nil && s == "endstream" {
				stream = NewStreamFromContents (dictionary,contents)
			}
		}
	}
	if stream != nil {
		return stream
	}
	return dictionary
}

func scanObject(scanner Scanner) Object {
	b,err := nextNonWhiteByte(scanner)
	if err == nil {
		switch  {
		case IsAlpha(b):
			return scanKeywordObject(scanner, b)
		case IsDigit(b),b=='.',b=='+',b=='-':
			return scanNumeric(scanner, b)
		case b =='/':
			return scanName (scanner)
		case b=='(':
			return scanNormalString(scanner)
		case b=='<':
			b,err = nextNonWhiteByte(scanner)
			if b == '<' {
				return scanDictionaryOrStream(scanner)
			} else {
				return scanHexString(scanner, b)
			}
		case b=='[':
			return scanArray(scanner)
		}
	}
	panic(unexpectedInput)
}

// Scan() uses Scanner to parse an arbitrary object.  If successful,
// the object is returned.  If not, err is set and context contains the
// input bytes that preceeded the error.
func Scan(scanner Scanner) (o Object,err error,context []byte) {
	historyReader := readers.NewHistoryReader(scanner,64)

	defer func() {
		if x := recover(); x!= nil {
			context = historyReader.GetHistory()
			err = x.(error)
		}
	} ()

	o = scanObject(historyReader)

	return
}
