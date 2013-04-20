package main

import (
	"github.com/mawicks/pdfDIG/pdf" )

// make_file() produces a file using low-level methods of the pdf.File
// type.  It does not work at the document layer and it does *not*
// produce a PDF document that a viewer will understand.
func modify_file() {
	f := pdf.NewFile("/tmp/test-document.pdf")
	documentInfoIndirect := pdf.NewIndirect(f)
	f.SetInfo (documentInfoIndirect)
	documentInfo := pdf.NewDocumentInfo()
	documentInfo.SetTitle("Rewritten Title")
	documentInfo.SetAuthor("Nobody")
	documentInfo.SetCreator("Nothing")
	documentInfoIndirect.Finalize(documentInfo)
	f.Close()
}
