package main

import (
	"fmt"
	"github.com/mawicks/pdfDIG/pdf" )

func make_document() {
	doc := pdf.NewDocument("/tmp/test-document.pdf")
	doc.SetAuthor("Mark Wicks")
	doc.SetTitle("Test Document")
	// Following is to test string encoding
	doc.SetKeywords("Résumé")

	// Page 1
	page := doc.NewPage()

	for f:=pdf.TimesRoman; f<=pdf.CourierBoldOblique; f++ {
		font := pdf.NewStandardFont(f)
		fmt.Fprintf (page, "BT /%s 24 Tf 230 %d Td (Hello World!) Tj ET ",
			page.AddFont(font), 760-24*int(f))
		// Use the same font again to test whether the font
		// dictionary gets repeated in the PDF file.  Same
		// font should not be added twice.  Here "same" means
		// the same object, not equivalent fonts.
		fmt.Fprintf (page, "BT /%s 24 Tf 230 %d Td (Hello World!) Tj ET ",
			page.AddFont(font), 360-24*int(f))
	}

	// Page 2
	page = doc.NewPage()
	fmt.Fprintf (page, "0 0 m 612 792 l s ")
	fmt.Fprintf (page, "0 792 m 612 0 l s")

	// Page 3
	page = doc.NewPage()
	f1 := pdf.NewStandardFont(pdf.Helvetica)
	fmt.Fprintf (page, "BT /%s 24 Tf 235 528 Td (Goodbye World!) Tj ET", page.AddFont(f1))

	doc.Close()
}
