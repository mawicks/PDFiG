/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

// PDF "Boolean" object.  There is no "Boolean" type as such.
// Implements: pdf.Object

type TrueBoolean struct {}
type FalseBoolean struct {}

// Constructor for "Boolean" object.  There are no methods required
// for "Booleans" beyond those in pdf.Object interface, so we simply
// return a pdf.Object and never even define a Boolean type.  Since
// TrueBoolean and BooleanFalse are empty structs, returning by value
// in NewBoolean and implementing methods with value targets should be
// efficient.
func NewBoolean (v bool) Object {
	var result Object
	if (v) {
		result = TrueBoolean{}
	} else {
		result = FalseBoolean{}
	}
	return result
}

// Since TrueBoolean and FalseBoolean are empty structs, value targets
// should be efficient.
func  (b TrueBoolean) Serialize (w Writer, file... File) {
	w.WriteString("true")
}

func  (b FalseBoolean) Serialize (w Writer, file... File) {
	w.Write ([]byte("false"))
}

