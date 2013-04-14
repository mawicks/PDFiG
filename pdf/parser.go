package pdf

import (
	"io"
	"errors"
	"strconv" )

type Scanner interface {
	io.ByteScanner
	io.Reader
//	Peek(int) ([]byte, error)
}


var parsingError = errors.New("Parsing Error")

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

func scanKeywordObject (scanner Scanner, b byte) (Object,error) {
	keyword,err := scanKeyword(scanner, b)
	if err == nil {
		switch (keyword) {
		case "null":
			return NewNull(),nil
		case "true":
			return NewBoolean(true),nil
		case "false":
			return NewBoolean(false),nil
		}
	}
	return nil,parsingError
}

func scanNumeric (scanner Scanner, b byte) (Object,error) {
	var buffer[]byte = make([]byte, 0, 5)
	var err error

	atLeastOneDigit := false
	float := false

	if (b == '+' || b == '-') {
		buffer = append(buffer,b)
		b,err=scanner.ReadByte()
	}

	for ; err==nil && IsDigit(b); b,err=scanner.ReadByte() {
		buffer = append(buffer,b)
		atLeastOneDigit = true
	}

	if (err == nil && b == '.') {
		buffer = append(buffer,b)
		b,err=scanner.ReadByte()
		float = true
	}

	for ; err==nil && IsDigit(b); b,err=scanner.ReadByte() {
		buffer = append(buffer,b)
		atLeastOneDigit = true
	}

	if err==nil {
		scanner.UnreadByte()
	}

	if err != nil && err!=io.EOF || !atLeastOneDigit {
		return nil,parsingError
	}

	if float {
		number,_ := strconv.ParseFloat(string(buffer),32)
		return NewFloatNumeric(float32(number)),nil
	}
	number,_ := strconv.ParseInt(string(buffer),10,32)
	return NewIntNumeric(int(number)),nil
}

func scanName (scanner Scanner) (Object, error) {
	var buffer[]byte = make([]byte, 0, 8)
	b,err := scanner.ReadByte()
	for ; err == nil && IsRegular(b); b,err=scanner.ReadByte() {
		if (b != '#') {
			buffer = append(buffer, b)
		} else {
			r := byte(0)
			var digit byte
			for i:=0; i<2; i++ {
				if b,err = scanner.ReadByte(); err != nil {
					return nil,parsingError
				}
				if digit, err = ParseHexDigit(b); err != nil {
					return nil,parsingError
				}
				r = 16*r + digit
			}
			buffer = append(buffer, r)
		}
	}
	if err == nil {
		scanner.UnreadByte()
	}
	if err == nil || err == io.EOF {
		return NewName(string(buffer)),nil
	}
	return nil,parsingError
}

func scanEscape (scanner Scanner) (b byte,err error) {
	if b,err =scanner.ReadByte(); err != nil {
		return b,err
	}

	if IsOctalDigit(b) {
		r,err := ParseOctalDigit(b)
		var digit byte
		for i:=0; i<2; i++ {
			if b,err=scanner.ReadByte(); err != nil {
				return 0,parsingError
			}
			if digit,err=ParseOctalDigit(b); err != nil {
				return 0,parsingError
			}
			r = 8*r + digit
		}
		return r,nil
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

func scanNormalString (scanner Scanner) (*String, error) {
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
			v,err := scanEscape (scanner)
			if err != nil {
				return nil,err
			}
			buffer = append(buffer, v)
		default:
			buffer = append(buffer, b)
		}
	}
	if err != nil {
		return nil,parsingError
	}
	return NewBinaryString(buffer),nil
}

func scanHexString (scanner Scanner,b byte) (Object, error) {
	var buffer[]byte = make([]byte, 0, 128)
	var err error
	for ; err == nil && b != '>'; b,err=scanner.ReadByte() {
		scanner.UnreadByte()
		var digit byte
		r := byte(0)
		for i:=0; i<2; i++ {
			if b,err = scanner.ReadByte(); err != nil {
				return nil,parsingError
			}
			if digit, err = ParseHexDigit(b); err != nil {
				return nil,parsingError
			}
			r = 16*r + digit
		}
		buffer = append(buffer, r)
	}
	if err == nil {
		return NewBinaryString(buffer),nil
	}
	return nil,parsingError
}

func scanArray (scanner Scanner) (*Array, error) {
	var array *Array = NewArray()

	b,err := nextNonWhiteByte(scanner)
	for ; err == nil && b != ']'; b,err=scanner.ReadByte() {
		scanner.UnreadByte()
		nextElement,err := scanObject(scanner)
		if err != nil {
			return nil,err
		}
		array.Add(nextElement)
	}

	if err != nil {
		return nil,parsingError
	}
	return array,nil
}

func scanDictionary(scanner Scanner) (*Dictionary, error) {
	var d *Dictionary = NewDictionary()

	b,err := nextNonWhiteByte(scanner)
	for ; err == nil && b != '>'; b,err=scanner.ReadByte() {
		scanner.UnreadByte()
		var name *Name
		object,err := scanObject(scanner)
		if err != nil {
			return nil,err
		}
		name = object.(*Name)
		object,err = scanObject(scanner)
		if err != nil {
			return nil,err
		}
		d.Add(name.String(),object)
	}

	if err != nil {
		return nil,parsingError
	}

	b,err = nextNonWhiteByte(scanner)
	if (b != '>') {
		return d,parsingError
	}
	return d,nil
}

func scanDictionaryOrStream (scanner Scanner) (Object, error) {
	dictionary,err := scanDictionary(scanner)

	if err != nil {
		return nil,err
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
	return dictionary,nil
}

func scanObject(scanner Scanner) (Object,error) {
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
	return nil,err
}

func Scan(scanner Scanner) (Object,error) {
	return scanObject (scanner)
}
