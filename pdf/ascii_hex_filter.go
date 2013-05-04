package pdf

type AsciiHexFilter struct {
}

func (fileter AsciiHexFilter) Name() string {
	return "ASCIIHexDecode"
}

func (filter AsciiHexFilter) Encode(buffer []byte) []byte {
	length := len(buffer)
	result := make([]byte, 0, 2*length + length/40 + 1)
	for i:=0; i<length; i++ {
		result = append(result,
			HexDigit(buffer[i]/16),
			HexDigit(buffer[i]%16))
		if i != 0 && i%40 == 0 {
			result = append(result,'\n')
		}
	}
	result = append(result,'>')
	return result
}

func (filter AsciiHexFilter) Decode(buffer []byte) (result []byte, ok bool) {
	length := len(buffer)
	result = make([]byte, 0, length/2)
	count := 0
	nextByte := byte(0);
	for i:= 0; i<length; i++ {
		switch  {
		case IsHexDigit(buffer[i]):
			nextByte = nextByte*16 + ParseHexDigit(buffer[i])
			count += 1
			if (count%2 == 0) {
				result = append(result, nextByte)
				nextByte = 0
			}
		case buffer[i] == '>':
			if (count%2 == 1) {
				nextByte = nextByte*16
				result = append(result, nextByte)
			}
			break;
		case IsWhiteSpace(buffer[i]):
			// Do nothing
		default:
			return nil,false
		}
	}
	return result,true
}
