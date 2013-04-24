package pdf

import "bufio"
import "os"

type Document struct {
	file File
	existing bool


	// Next four pointers are nil for pre-existing documents unless NewPage() is called
	// Either all are nil or all are non-nil.
	pageTreeRoot *Dictionary
	pageTreeRootIndirect *Indirect
	procSetIndirect *Indirect
	// Each element in the pages array is an Indirect reference to a Page dictionary
	// returned by Page.Indirect().
	pages *Array

	currentPage *Page
	pageCount uint

	DocumentInfo
}

func (d *Document) needPageTree() {
	if d.pages == nil {
		d.pages = NewArray()
		d.pageTreeRoot = NewDictionary()
		d.pageTreeRootIndirect = NewIndirect(d.file)
		d.procSetIndirect = NewIndirect(d.file)

		// For now, this is a default to be sure a box is set somewhere.
		// Clients can reset with their own call to SetMediaBox().
		d.SetMediaBox(0, 0, 612, 792)
	}
}

// OpenDocument() constructs a document object from either a new or a pre-existing filename.
func OpenDocument(filename string, mode int) *Document {
	d := new(Document)

	d.file,d.existing,_ = OpenFile(filename, mode)

	d.DocumentInfo = NewDocumentInfo()

	if d.existing {
		d.DocumentInfo = DocumentInfo{d.file.Info(),false}
		d.pageTreeRoot = oldPageTree(d.file).root
		out := bufio.NewWriter(os.Stdout)
		out.WriteString("page tree root: ")
		d.pageTreeRoot.Serialize(out,d.file)
		out.WriteString("\n")
		out.Flush()
	} else {
		d.needPageTree()
	}

	// Set a default producer field.  Clients calls to SetProducer() override this.
	d.SetProducer("PDFiG")

	return d
}

func (d *Document) release() {
	d.pages = nil
	d.currentPage = nil
	d.pageTreeRoot = nil
	d.pageTreeRootIndirect = nil
	d.procSetIndirect = nil
}

func (d *Document) finishCatalog() {
	if d.pageTreeRootIndirect != nil {
		catalog := NewDictionary()
		catalog.Add("Type", NewName("Catalog"))
		catalog.Add("Pages", d.pageTreeRootIndirect)
		d.file.SetCatalog(catalog)
	}
}

func (d *Document) finishCurrentPage() {
	if d.currentPage != nil {
		d.pages.Add(d.currentPage.Finalize())
		d.pageCount += 1
	}
}

func (d *Document) finishDocumentInfo() {
	if d.DocumentInfo.IsDirty() {
		d.file.SetInfo (d.DocumentInfo)
	}
}

func (d *Document) finishPageTree() {
	if d.pageTreeRoot != nil {
		d.pageTreeRoot.Add("Type", NewName("Pages"))
		d.pageTreeRoot.Add("Count", NewIntNumeric(int(d.pageCount)))
		d.pageTreeRoot.Add("Kids", d.pages)
		d.pageTreeRootIndirect.Finalize(d.pageTreeRoot)
	}
}

func (d *Document) finishProcSet() {
	// Procset is option for PDF versions >= 1.4
	// The following full set is recommended, however, for maximal compatibility.
	if d.procSetIndirect != nil {
		procsetArray := NewArray()
		procsetArray.Add (NewName("PDF"))
		procsetArray.Add (NewName("Text"))
		procsetArray.Add (NewName("ImageB"))
		procsetArray.Add (NewName("ImageC"))
		procsetArray.Add (NewName("ImageI"))
		d.procSetIndirect.Finalize(procsetArray)
	}
}

func (d *Document) NewPage() *Page {
	d.finishCurrentPage()
	d.currentPage = NewPage(d.file)

	d.needPageTree()
	d.currentPage.SetParent(d.pageTreeRootIndirect)
	d.currentPage.setProcSet(d.procSetIndirect)

	return d.currentPage
}

func (d *Document) Close() {
	d.finishCurrentPage()
	d.finishProcSet()
	d.finishPageTree()
	d.finishCatalog()
	d.finishDocumentInfo()

	d.file.Close()

	d.release()
}

func (d *Document) SetMediaBox(llx, lly, urx, ury float64) {
	d.needPageTree()
	d.pageTreeRoot.Add("MediaBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetCropBox(llx, lly, urx, ury float64) {
	d.needPageTree()
	d.pageTreeRoot.Add("CropBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetBleedBox(llx, lly, urx, ury float64) {
	d.needPageTree()
	d.pageTreeRoot.Add("BleedBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetTrimBox(llx, lly, urx, ury float64) {
	d.needPageTree()
	d.pageTreeRoot.Add("TrimBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetArtBox(llx, lly, urx, ury float64) {
	d.needPageTree()
	d.pageTreeRoot.Add("ArtBox", NewRectangle(llx, lly, urx, ury))
}
