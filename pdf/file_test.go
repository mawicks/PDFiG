package pdf

import "testing"

func TestTestFile (t *testing.T) {
	f := NewFile("/tmp/foo.pdf")
	f.Close()
}

