package pdf

import "testing"

func TestTestFile(t *testing.T) {
	f := NewFile("/tmp/foo.pdf")
	o1 := NewIndirect()
	indirect1 := f.AddObject(o1)
	o1.Finalize(NewNumeric(3.14))

	indirect2 := f.AddObject(NewNumeric(2.718))

	f.AddObject(NewName("foo"))


	// Delete the *indirect reference* to the 3.14 numeric
	f.DeleteObject(indirect1.ObjectNumber(f))
	f.AddObject(NewNumeric(3))

	// Delete the 2.718 numeric itself
	f.DeleteObject(indirect2.ObjectNumber(f))
	f.Close()
}
