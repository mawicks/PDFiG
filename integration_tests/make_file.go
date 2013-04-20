package main

import (
	"github.com/mawicks/PDFiG/pdf" )

// make_file() produces a file using low-level methods of the pdf.File
// type.  It does not work at the document layer and it does *not*
// produce a PDF document that a viewer will understand.
func make_file() {
	f := pdf.OpenFile("/tmp/test-file.pdf")
	o1 := pdf.NewIndirect()
	indirect1 := f.AddObject(o1)
	o1.Finalize(pdf.NewNumeric(3.14))

	indirect2 := f.AddObject(pdf.NewNumeric(2.718))

	f.AddObject(pdf.NewName("foo"))

	// Delete the *indirect reference* to the 3.14 numeric
	f.DeleteObject(indirect1)
	f.AddObject(pdf.NewNumeric(3))

	// Delete the 2.718 numeric object itself
	f.DeleteObject(indirect2)

	p := pdf.NewPage(f)
	p.SetParent(indirect1)
	p.SetMediaBox(0, 0, 612, 792)
	p.Finalize()

	catalogIndirect := pdf.NewIndirect(f)
	f.SetCatalog(catalogIndirect)

	catalog := pdf.NewDictionary()
	catalog.Add("Type", pdf.NewName("Catalog"))
	catalogIndirect.Finalize(catalog)

	f.Close()
}
