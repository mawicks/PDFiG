package pdf

import "testing"

func TestDocument(t *testing.T) {
	f := NewDocument("/tmp/foo-document.pdf")
	f.NewPage()
	f.NewPage()
	f.Close()
}
