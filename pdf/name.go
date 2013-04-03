package pdf

func HexDigit(b byte) (result byte) {
	switch {
	case b < 10:
		result = b + '0'
	default:
		result = (b - 10) + 'A'
	}
	return result
}

type Name struct {
	name string
}

// Constructor for Name object
func NewName(s string) Name {
	return Name{s}
}

func nameEscapeByte(b byte) (result []byte) {
	switch {
	case b != '#' && IsRegular(b):
		result = []byte{b}
	default:
		result = []byte{'#', HexDigit(b / 16), HexDigit(b % 16)}
	}
	return result
}

func (n Name) Serialize(w Writer, file ...File) {
	w.WriteByte('/')
	for _, b := range []byte(n.name) {
		w.Write(nameEscapeByte(b))
	}
	return
}
