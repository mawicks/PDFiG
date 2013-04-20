package pdf_test

import (
	"fmt"
	"github.com/mawicks/pdfDIG/pdf" )

func ExampleDocument() {
	doc := pdf.OpenDocument("/tmp/test-document.pdf")
	doc.SetAuthor("Mark Wicks")
	doc.SetTitle("Test Document")

	// Page 1
	page := doc.NewPage()
	f1 := pdf.NewStandardFont(pdf.Helvetica)
	name := page.AddFont(f1)
	fmt.Fprintf (page, "BT /%s 24 Tf 250 528 Td (Hello World!) Tj ET", name)

	// Page 2
	page = doc.NewPage()
	fmt.Fprintf (page, "0 0 m 612 792 l s ")
	fmt.Fprintf (page, "0 792 m 612 0 l s")

	// Page 3
	page = doc.NewPage()
	name = page.AddFont(f1)
	fmt.Fprintf (page, "BT /%s 24 Tf 250 528 Td (Goodbye World!) Tj ET", name)

	doc.Close()
}
