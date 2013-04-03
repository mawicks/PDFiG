package pdf

import "testing"

func TestTestFile(t *testing.T) {
	f := NewFile("/tmp/foo.pdf")
	o1 := NewIndirect()
	direct1 := f.AddObject(o1)

	o2 := NewNumeric(3.14)
	o1.Finalize(o2)

	o3 := NewNumeric(2.718)
	direct2 := f.AddObject(o3)

	f.DeleteObject(direct1)
	f.AddObject(NewNumeric(3))
	f.DeleteObject(direct2)
	f.Close()
}
