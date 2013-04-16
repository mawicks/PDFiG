package pdf

type Document struct {
	file File

	// Each element in the pages array is a Indirect reference to a Page dictionary
	// returned by Page.Indirect().
	pages *Array

	currentPage *Page
	pageCount uint

	pageTreeRoot *Dictionary
	pageTreeRootIndirect *Indirect

	procSetIndirect *Indirect

	DocumentInfo
}

func NewDocument(filename string) *Document {
	d := new(Document)

	d.file = NewFile(filename)

	d.pages = NewArray()

	d.pageTreeRoot = NewDictionary()
	d.pageTreeRootIndirect = NewIndirect(d.file)

	d.procSetIndirect = NewIndirect(d.file)

	// For now, this is a default to be sure a box is set somewhere.
	// Clients can reset with their own call to SetMediaBox().
	d.SetMediaBox(0, 0, 612, 792)

	d.DocumentInfo = NewDocumentInfo()
	// Set a default producer field.  Clients calls to SetProducer() override this.
	d.SetProducer("pdfdig")

	return d
}

func (d *Document) release() {
	d.pages = nil
	d.currentPage = nil
	d.pageTreeRoot = nil
	d.pageTreeRootIndirect = nil
	d.procSetIndirect = nil
}

func (d *Document) finishCurrentPage() {
	d.currentPage.Finalize()
	d.pages.Add(d.currentPage.Indirect())
	d.pageCount += 1
}

func (d *Document) finishPageTree() {
	d.pageTreeRoot.Add("Type", NewName("Pages"))
	d.pageTreeRoot.Add("Count", NewIntNumeric(int(d.pageCount)))
	d.pageTreeRoot.Add("Kids", d.pages)
	d.pageTreeRootIndirect.Finalize(d.pageTreeRoot)
}

func (d *Document) finishProcSet() {
	// Procset is option for PDF versions >= 1.4
	// The following full set is recommended, however, for maximal compatibility.
	procsetArray := NewArray()
	procsetArray.Add (NewName("PDF"))
	procsetArray.Add (NewName("Text"))
	procsetArray.Add (NewName("ImageB"))
	procsetArray.Add (NewName("ImageC"))
	procsetArray.Add (NewName("ImageI"))
	d.procSetIndirect.Finalize(procsetArray)
}

func (d *Document) finishCatalog() {
	catalog := NewDictionary()
	catalog.Add("Type", NewName("Catalog"))
	catalogIndirect := NewIndirect(d.file)
	d.file.SetCatalog(catalogIndirect)
	catalog.Add("Pages", d.pageTreeRootIndirect)
	catalogIndirect.Finalize(catalog)
}

func (d *Document) finishDocumentInfo() {
	if d.DocumentInfo.Size() != 0 {
		documentInfoIndirect := NewIndirect(d.file)
		d.file.SetInfo (documentInfoIndirect)
		documentInfoIndirect.Finalize(d.DocumentInfo)
	}
}

func (d *Document) NewPage() *Page {
	if d.currentPage != nil {
		d.finishCurrentPage()
	}
	d.currentPage = NewPage(d.file)
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
	d.pageTreeRoot.Add("MediaBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetCropBox(llx, lly, urx, ury float64) {
	d.pageTreeRoot.Add("CropBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetBleedBox(llx, lly, urx, ury float64) {
	d.pageTreeRoot.Add("BleedBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetTrimBox(llx, lly, urx, ury float64) {
	d.pageTreeRoot.Add("TrimBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetArtBox(llx, lly, urx, ury float64) {
	d.pageTreeRoot.Add("ArtBox", NewRectangle(llx, lly, urx, ury))
}
