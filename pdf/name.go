package pdf

type Name struct {
	name string
}

// Constructor for Name object
func NewName(s string) *Name {
	return &Name{s}
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

func (n *Name) Clone() Object {
	// Names are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

func (n *Name) Dereference() Object {
	return n
}

func (n *Name) Serialize(w Writer, file ...File) {
	w.WriteByte('/')
	for _, b := range []byte(n.name) {
		w.Write(nameEscapeByte(b))
	}
	return
}

func (n *Name) String() string {
	return n.name
}
