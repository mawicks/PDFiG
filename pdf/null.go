package pdf

// PDF "Null" object
// Implements: pdf.Object
type Null struct{}

var nullSingleInstance Null

func NewNull() Object {
	return &nullSingleInstance
}

func (n *Null) Serialize(w Writer, file ...File) {
	w.WriteString("null")
	return
}

func (n *Null) Clone() Object {
	// All copies of null point to the same instance.
	return n
}
