package pdf

import "testing"

func TestTestFile (t *testing.T) {
	f := NewFile("/tmp/foo.pdf")
	obj1 := f.AddObject (NewNumeric(3.14))
	obj2 := f.AddObject (NewNumeric(2.718))
	f.DeleteObject (obj1)
	f.AddObject (NewNumeric(3))
	f.DeleteObject (obj2)
	f.Close()
}

