package main

import (
	"fmt"
	"os"
	"github.com/mawicks/PDFiG/pdf" )

func insert_page() {
	doc := pdf.OpenDocument("/tmp/test-document.pdf", os.O_RDWR|os.O_CREATE)

	page := doc.NewPage()
	f1 := pdf.NewStandardFont(pdf.Helvetica)
	fmt.Fprintf (page, "BT /%s 24 Tf 235 528 Td (Inserted page) Tj ET", page.AddFont(f1))

	doc.Close()
}
