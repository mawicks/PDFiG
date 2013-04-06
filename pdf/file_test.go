package pdf

func ExampleFile() {
	f := NewFile("/tmp/test-file.pdf")
	o1 := NewIndirect()
	indirect1 := f.AddObject(o1)
	o1.Finalize(NewNumeric(3.14))

	indirect2 := f.AddObject(NewNumeric(2.718))

	f.AddObject(NewName("foo"))

	// Delete the *indirect reference* to the 3.14 numeric
	f.DeleteObject(indirect1.ObjectNumber(f))
	f.AddObject(NewNumeric(3))

	// Delete the 2.718 numeric object itself
	f.DeleteObject(indirect2.ObjectNumber(f))

	p := NewPage(f)
	p.SetParent(indirect1)
	p.SetMediaBox(0, 0, 612, 792)
	p.Finalize()

	catalogIndirect := NewIndirect(f)
	f.SetCatalog(catalogIndirect)

	catalog := NewDictionary()
	catalog.Add("Type", NewName("Catalog"))
	catalogIndirect.Finalize(catalog)

	f.Close()
}
