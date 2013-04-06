package pdf_test

import (
	"github.com/mawicks/goPDF/pdf" )

func ExampleFile() {
	f := pdf.NewFile("/tmp/test-file.pdf")
	o1 := pdf.NewIndirect()
	indirect1 := f.AddObject(o1)
	o1.Finalize(pdf.NewNumeric(3.14))

	indirect2 := f.AddObject(pdf.NewNumeric(2.718))

	f.AddObject(pdf.NewName("foo"))

	// Delete the *indirect reference* to the 3.14 numeric
	f.DeleteObject(indirect1.ObjectNumber(f))
	f.AddObject(pdf.NewNumeric(3))

	// Delete the 2.718 numeric object itself
	f.DeleteObject(indirect2.ObjectNumber(f))

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
