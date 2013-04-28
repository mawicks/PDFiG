package main

import (
	"fmt"
	"os"
	"github.com/mawicks/PDFiG/pdf" )

func modify_document() {
	fmt.Printf ("\n\nMODIFY DOCUMENT\n")
	doc := pdf.OpenDocument("/tmp/test-document.pdf", os.O_RDWR|os.O_CREATE)

	doc.Page(1)

	page := doc.NewPage()
	f1 := pdf.NewStandardFont(pdf.Helvetica)
	fmt.Fprintf (page, "BT /%s 24 Tf 235 528 Td (Inserted page) Tj ET", page.AddFont(f1))

	doc.Close()
	fmt.Printf ("END MODIFY DOCUMENT\n\n")
}
