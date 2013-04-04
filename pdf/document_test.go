package pdf

import "fmt"
import "testing"

func TestDocument(t *testing.T) {
	f := NewDocument("/tmp/foo-document.pdf")

	// Page 1
	p := f.NewPage()
	fmt.Fprintf (p, "0 0 m 612 792 l s ")
	fmt.Fprintf (p, "0 792 m 612 0 l s")

	// Page 2
	f.NewPage()

	f.Close()
}
