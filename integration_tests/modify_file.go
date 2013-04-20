package main

import (
	"bufio"
	"os"
	"github.com/mawicks/PDFiG/pdf" )

// make_file() produces a file using low-level methods of the pdf.File
// type.  It does not work at the document layer and it does *not*
// produce a PDF document that a viewer will understand.
func modify_file() {
	f := pdf.OpenFile("/tmp/test-document.pdf")
	documentInfo := pdf.NewDocumentInfo()
	documentInfo.SetTitle("Rewritten Title")
	documentInfo.SetAuthor("Nobody")
	documentInfo.SetCreator("Nothing")

	// Verify that we can retrieve an arbitrary object
	writer := bufio.NewWriter(os.Stdout)
	o,_ := f.Object(pdf.NewObjectNumber(10,0))
	o.Serialize(writer, f)
	writer.WriteString("\n")
	writer.Flush()

	if info := f.Info(); info != nil {
		info.Serialize(writer, f)
		writer.WriteString("\n")
		writer.Flush()
	}

	documentInfoIndirect := pdf.NewIndirect(f)
	documentInfoIndirect.Finalize(documentInfo)
	f.SetInfo (documentInfoIndirect)

	f.Close()
}
