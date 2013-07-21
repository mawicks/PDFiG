package pdf

type Name interface {
	Object
	String() string
}

type name struct {
	name string
}

// Constructor for standard Name implementation
func NewName(s string) Name {
	return &name{s}
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

func (n *name) Clone() Object {
	// Names are intended to be immutable, so return a pointer
	// to the same instance
	return n
}

func (n *name) Dereference() Object {
	return n
}

// Names are immutable so the return value of Protect() can be cast back to Name.
func (n *name) Protect() Object {
	// Names are treated as immutable, so return a pointer
	// to the same instance
	return n
}

// Names are immutable so the return value of Unprotect() can be cast back to Name.
func (n *name) Unprotect() Object {
	// Names are treated as immutable, so return a pointer
	// to the same instance
	return n
}


func (n *name) Serialize(w Writer, file ...File) {
	w.WriteByte('/')
	for _, b := range []byte(n.name) {
		w.Write(nameEscapeByte(b))
	}
	return
}

func (n *name) String() string {
	return n.name
}
