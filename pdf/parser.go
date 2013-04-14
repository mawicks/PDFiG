package pdf

import (
	"fmt"
	"io"
	"errors"
	"strconv" )

type Scanner interface {
	io.ByteScanner
	io.Reader
//	Peek(int) ([]byte, error)
}


var parsingError = errors.New("Parsing Error")
var expectingName = errors.New("Expecting Name")

type Parser struct {
}

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
	panic(parsingError)
}

func scanNumeric (scanner Scanner, b byte) Object {
	var buffer[]byte = make([]byte, 0, 5)
	var err error

	hasAtLeastDigit := false
	float := false

	if (b == '+' || b == '-') {
		buffer = append(buffer,b)
		b,err=scanner.ReadByte()
	}

	for ; err==nil && IsDigit(b); b,err=scanner.ReadByte() {
		buffer = append(buffer,b)
		hasAtLeastDigit = true
	}

	if (err == nil && b == '.') {
		buffer = append(buffer,b)
		b,err=scanner.ReadByte()
		float = true
	}

	for ; err==nil && IsDigit(b); b,err=scanner.ReadByte() {
		buffer = append(buffer,b)
		hasAtLeastDigit = true
	}

	if err==nil {
		scanner.UnreadByte()
	}

	if err != nil && err!=io.EOF || !hasAtLeastDigit {
		panic(parsingError)
	}

	if float {
		number,_ := strconv.ParseFloat(string(buffer),32)
		return NewFloatNumeric(float32(number))
	}
	number,_ := strconv.ParseInt(string(buffer),10,32)
	return NewIntNumeric(int(number))
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
				if b,err = scanner.ReadByte(); err != nil {
					panic(parsingError)
				}
				r = 16*r + ParseHexDigit(b)
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
	panic (parsingError)
}

func scanEscape (scanner Scanner) (b byte) {
	var err error
	if b,err =scanner.ReadByte(); err != nil {
		panic (parsingError)
	}

	if IsOctalDigit(b) {
		r := ParseOctalDigit(b)
		for i:=0; i<2; i++ {
			if b,err=scanner.ReadByte(); err != nil {
				panic(parsingError)
			}
			r = 8*r + ParseOctalDigit(b)
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
	var openCount = 1
	var buffer[]byte = make([]byte, 0, 128)
	b,err :=scanner.ReadByte()
	for ; err == nil && openCount != 0; b,err=scanner.ReadByte() {
		switch b {
		case '(':
			openCount += 1;
			buffer = append(buffer, b)
		case ')':
			openCount -= 1;
			if (openCount != 0) {
				buffer = append(buffer, b)
			}
		case '\\':
			v := scanEscape (scanner)
			buffer = append(buffer, v)
		default:
			buffer = append(buffer, b)
		}
	}
	if err != nil {
		panic (parsingError)
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
			if b,err = scanner.ReadByte(); err != nil {
				panic(parsingError)
			}
			r = 16*r + ParseHexDigit(b)
		}
		buffer = append(buffer, r)
	}
	if err != nil {
		panic(parsingError)
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
		panic(parsingError)
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
		panic(parsingError)
	}

	b,err = nextNonWhiteByte(scanner)
	if (b != '>') {
		panic(parsingError)
	}
	return d
}

func scanDictionaryOrStream (scanner Scanner) Object {
	var err error

	dictionary := scanDictionary(scanner)

	if err != nil {
		panic(parsingError)
	}

	var b byte
	b,err = nextNonWhiteByte(scanner)
	scanner.UnreadByte()

	var s string
	// Could be "stream" line.
	if b=='s' {
		s,err = ReadLine (scanner)
	}

	if (err == nil && s == "stream") {
		v := dictionary.Get("Length").(*IntNumeric)
		if v != nil {
			length := v.Value()
			contents := make([]byte, length)
			scanner.Read(contents)
		}
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
	panic(parsingError)
}

func Scan(scanner Scanner) (o Object,err error) {
	defer func() {
		if x := recover(); x!= nil {
			fmt.Printf ("An error occurred while parsing: %v\n", x)
			err = x.(error)
		}
	}()
	o = scanObject(scanner)
	return
}
