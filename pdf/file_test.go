package pdf_test

import (
	"os"
	"strings"
	"testing"
	"github.com/mawicks/PDFiG/pdf" )

func ExampleFile_creation() {
	f,_,_ := pdf.OpenFile("/tmp/test-file.pdf", os.O_RDWR|os.O_CREATE)

	o1 := pdf.NewIndirect()
	indirect1 := f.AddObject(o1)
	o1.Write(pdf.NewNumeric(3.14))

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
	p.Close()

	catalog := pdf.NewDictionary()
	catalog.Add("Type", pdf.NewName("Catalog"))
	f.SetCatalog(catalog)

	f.Close()
}

func ExampleFile_modification() {
	f,_,_ := pdf.OpenFile("/tmp/pre-existing-file.pdf", os.O_RDWR|os.O_CREATE)
	documentInfo := pdf.NewDocumentInfo()
	documentInfo.SetTitle("Rewritten Title")
	documentInfo.SetAuthor("Nobody")
	documentInfo.SetCreator("Nothing")

	f.SetInfo (documentInfo)

	f.Close()
}

func TestPDFReadLine (t *testing.T) {
	teststring := "abc\ndef\rghi\r\njkl\n\r\n\r123\n\r\r\n456\n\n789"
	lines := [...]string{
		"abc", "def", "ghi", "jkl", "", "123", "", "456", "", "789"}
	reader := strings.NewReader(teststring)
	for _,line := range lines {
		s,err := pdf.ReadLine(reader)
		if err != nil || s != line {
			t.Errorf (`Got "%s"; expected "%s" (err=%v)`, s, line, err)
		}
	}
}

