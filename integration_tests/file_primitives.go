package main

import (
//	"fmt"
	"bufio"
	"os"
	"github.com/mawicks/PDFiG/pdf" )

func file_primitives() {
	logger := bufio.NewWriter(os.Stderr)

	catalog := pdf.NewDictionary()
	catalog.Add("Type", pdf.NewName("Catalog"))

	fileName1 := "/tmp/test-file1.pdf"
	os.Remove(fileName1)
	f1,_,_ := pdf.OpenFile(fileName1, os.O_RDWR|os.O_CREATE)
	d1 := pdf.NewDictionary()
	i1 := pdf.NewIndirect(f1)	// Object number 1
	i2 := pdf.NewIndirect(f1)	// Object number 2
	// Object 1 is a name
	i1.Write(pdf.NewName("bar"))
	// Object 2 is a dictionary containing the name.
	d1.Add("foo", i1)
	i2.Write(d1)

	// Every file must have a catalog.
	f1.SetCatalog(catalog)
	f1.Close()

	fileName2 := "/tmp/test-file2.pdf"
	os.Remove(fileName2)
	f1,_,_ = pdf.OpenFile(fileName1, os.O_RDONLY)
	f2,_,_ := pdf.OpenFile(fileName2, os.O_RDWR|os.O_CREATE)

	// Retrieve object 2 (the dictionary) from file 1
	o,_ := f1.Object(pdf.NewObjectNumber(2,0))

	// Read back an object immediately to force a read from the
	// serialization that is cached while writing
	x,_ := f2.Object(f2.WriteObject(pdf.NewNumeric(3.14)).ObjectNumber(f2))
	logger.WriteString("Object: ")
	x.Serialize(logger, f2)
	logger.WriteString("\n")

	// Explicitly add it to file 2; object 1 will be added
	// automatically.  The objects are renumbered automatically in
	// the new file and are intentionally reversed (2 becomes 1; 1
	// becomes 2)
	pdf.NewIndirect(f2).Write(o)

	f2.SetCatalog(catalog)

	f1.Close()
	f2.Close()

	logger.Flush()
}
