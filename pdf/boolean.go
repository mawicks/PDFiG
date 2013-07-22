package pdf

// PDF "Boolean" object.  There is no "Boolean" type as such.
// Implements: pdf.Object

type Boolean interface {
	Object
	Value() bool
}

type TrueBoolean struct{}
type FalseBoolean struct{}

// Constructor for "Boolean" object.  Since TrueBoolean and
// BooleanFalse are empty structs, returning by value in NewBoolean
// and implementing methods with value targets should be efficient.
func NewBoolean(v bool) Object {
	var result Object
	if v {
		result = TrueBoolean{}
	} else {
		result = FalseBoolean{}
	}
	return result
}

func (b TrueBoolean) Clone() Object {
	return b
}

func (b TrueBoolean) Dereference() Object {
	return b
}

func (b TrueBoolean) Protect() Object {
	return b
}

func (b TrueBoolean) Unprotect() Object {
	return b
}

func (b TrueBoolean) Serialize(w Writer, file ...File) {
	w.WriteString("true")
}

func (b TrueBoolean) Value() bool {
	return true
}

func (b FalseBoolean) Clone() Object {
	return b
}

func (b FalseBoolean) Dereference() Object {
	return b
}

func (b FalseBoolean) Protect() Object {
	return b
}

func (b FalseBoolean) Unprotect() Object {
	return b
}

func (b FalseBoolean) Serialize(w Writer, file ...File) {
	w.Write([]byte("false"))
}

func (b FalseBoolean) Value() bool {
	return false
}

