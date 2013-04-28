package main

import (
	"bufio"
	"os"
	"github.com/mawicks/PDFiG/pdf" )

// make_file() produces a file using low-level methods of the pdf.File
// type.  It does not work at the document layer and it does *not*
// produce a PDF document that a viewer will understand.
func modify_file() {
	f,_,_ := pdf.OpenFile("/tmp/test-document.pdf", os.O_RDWR|os.O_CREATE)
	documentInfo := pdf.NewDocumentInfo()
	documentInfo.SetTitle("Rewritten Title")
	documentInfo.SetAuthor("Nobody")
	documentInfo.SetCreator("Nothing")

	// Verify that we can retrieve an arbitrary object
	writer := bufio.NewWriter(os.Stdout)
	o,_ := f.Object(pdf.NewObjectNumber(10,0))
	writer.WriteString("Object number 10: ")
	o.Serialize(writer, f)
	writer.WriteString("\n")
	writer.Flush()

	if info := f.Info(); info != nil {
		writer.WriteString("Pre-existing document info: ")
		info.Serialize(writer, f)
		writer.WriteString("\n")
		writer.Flush()
	}

	if catalog := f.Catalog(); catalog != nil {
		writer.WriteString("Pre-existing document catalog: ")
		catalog.Serialize(writer, f)
		writer.WriteString("\n")
		writer.Flush()
	}

	f.SetInfo (documentInfo)

	f.Close()
}
