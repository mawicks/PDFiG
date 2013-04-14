package pdf

import (
	"io"
	"errors" )

type Scanner interface {
	io.ByteReader
	Peek(int) ([]byte, error)
}


var parsingError = errors.New("Parsing Error")

type Parser struct {
}

// Skip white space and return the byte following the white space or error.
// If err is non-nil, the value of b is undefined.
func nextNonWhiteByte (scanner io.ByteScanner) (b byte,err error) {
	b, err = scanner.ReadByte()
	for ; err == nil && (IsWhiteSpace(b) || b == '%'); b,err=scanner.ReadByte() {
		if b == '%' {
			ReadLine(scanner)
		}
	}
	return
}

func scanName (scanner io.ByteScanner) (Object, error) {
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
	if err == io.EOF {
		err = nil
	}
	return NewName(string(buffer)),err
}

func scanEscape (scanner io.ByteScanner) (b byte,err error) {
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

func scanNormalString (scanner io.ByteScanner) (Object, error) {
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
	if err == nil {
		return NewBinaryString(buffer),nil
	}
	return nil,parsingError
}

func scanKeyword (scanner io.ByteScanner, b byte) (Object,error) {
	var buffer[]byte = make([]byte, 0, 5)
	buffer = append(buffer, b)

	b,err := scanner.ReadByte()
	for ; err == nil && IsAlpha(b); b,err=scanner.ReadByte() {
		buffer = append(buffer, b)
	}

	if err == nil {
		scanner.UnreadByte()
	}

	switch (string(buffer)) {
	case "null":
		return NewNull(),nil
	case "true":
		return NewBoolean(true),nil
	case "false":
		return NewBoolean(false),nil
	}
	return nil,parsingError
}

func Scan(scanner io.ByteScanner) (Object,error) {
	b,err := nextNonWhiteByte(scanner)
	if err == nil {
		switch  {
		case b =='/':
			return scanName (scanner)
		case IsDigit(b):
			return nil,nil
		case b=='<':
			return nil,parsingError
		case b=='(':
			return scanNormalString(scanner)
		case b=='[':
			return nil,parsingError
		case IsAlpha(b):
			return scanKeyword(scanner, b)
		}
	}
	return nil,err
}