package pdf

func ParseHexDigit(b byte) (byte) {
	switch {
	case b>='0' && b<='9':
		return b-'0'
	case b>='a' && b<='f':
		return b-'a'+10
	case b>='A' && b<='F':
		return b-'A'+10
	}
	panic (expectingHexDigit)
}

func ParseOctalDigit(b byte) (byte) {
	switch {
	case b>='0' && b<='7':
		return b-'0'
	}
	panic (expectingOctalDigit)
}

func AsciiFromBytes (b []byte) string {
	escaped := make([]byte,0,len(b))
	for i:=0; i<len(b); i++ {
		escaped = append(escaped, generalAsciiEscapeByte(b[i])...)
	}
	return string(escaped)
}

