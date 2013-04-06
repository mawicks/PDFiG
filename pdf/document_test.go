package pdf

import "fmt"

func ExampleDocument() {
	doc := NewDocument("/tmp/test-document.pdf")

	var page *Page

	// Page 1
	page = doc.NewPage()
	f1 := NewStandardFont(Helvetica, "F1")
	page.AddFont(f1)
	fmt.Fprintf (page, "BT /F1 24 Tf 250 528 Td (Hello World!) Tj ET")

	// Page 2
	page = doc.NewPage()
	fmt.Fprintf (page, "0 0 m 612 792 l s ")
	fmt.Fprintf (page, "0 792 m 612 0 l s")

	// Page 3
	page = doc.NewPage()
	page.AddFont(f1)
	fmt.Fprintf (page, "BT /F1 24 Tf 250 528 Td (Hello World!) Tj ET")

	doc.Close()
}
