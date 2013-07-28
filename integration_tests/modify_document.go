package main

import (
	"bufio"
	"fmt"
	"os"
	"github.com/mawicks/PDFiG/pdf" )

func modify_document() {
	fmt.Printf ("\nMODIFY DOCUMENT\n")
	doc := pdf.OpenDocument(OutputDirectory  + "/test-document.pdf", os.O_RDWR|os.O_CREATE)

	// Verify that we can retrieve an arbitrary page
	oldPage := doc.Page(1)	// Page 2
	// Try to read the contents on the page.
	if r := oldPage.Reader(); r != nil {
		c := []byte{0}
		fmt.Printf ("Contents: ")
		for n,_ := r.Read(c); n != 0; n,_ = r.Read(c) {
			os.Stdout.Write(pdf.GeneralAsciiEscapeByte(c[0]))
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Creating background contents\n")
	backgroundStream := pdf.NewStreamFactory().New()
	backgroundStream.Write([]byte("q 0.9 g 0.9 G 0 768 612 24 re b Q"))

	fmt.Printf("Writing contents stream to file\n")
	id := doc.WriteObject(backgroundStream)

	fmt.Printf("Prepending reference to object to page stream\n")
	oldPage.PrependContents(id)

	fmt.Printf("Rewriting page\n")
	oldPage.Rewrite()

	fmt.Printf("Displaying page(1) contents\n")
	writer := bufio.NewWriter(os.Stdout)
	if oldPage == nil {
		fmt.Fprintf (writer, "Page(1) returned nil\n")
	} else {
		fmt.Fprintf (writer, "Page(1) returned: ")
		oldPage.Clone().Serialize(writer, nil)
		writer.WriteString("\n")
	}
	writer.Flush()

	page := doc.NewPage()
	f1 := pdf.NewStandardFont(pdf.Helvetica)
	fmt.Fprintf (page, "BT /%s 24 Tf 235 528 Td (Inserted page) Tj ET", page.AddFont(f1))

	doc.Close()
	fmt.Printf ("END MODIFY DOCUMENT\n\n")
}
