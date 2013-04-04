package pdf

type Document struct {
	file            File
	catalog         *Dictionary
	catalogIndirect *Indirect
	// Each element in the pages array is a Indirect reference to a Page dictionary
	// returned by Page.Indirect().
	pages *Array

	currentPage *Page
	pageCount   uint

	pageTreeRoot         *Dictionary
	pageTreeRootIndirect *Indirect
}

func NewDocument(filename string) *Document {
	d := new(Document)

	d.file = NewFile(filename)

	d.catalog = NewDictionary()
	d.catalog.Add("Type", NewName("Catalog"))

	d.catalogIndirect = NewIndirect()
	d.catalogIndirect.ObjectNumber(d.file)
	d.file.SetCatalog(d.catalogIndirect)

	d.pages = NewArray()

	d.pageTreeRootIndirect = NewIndirect()
	d.pageTreeRootIndirect.ObjectNumber(d.file)

	d.pageTreeRoot = NewDictionary()
	d.pageTreeRoot.Add("Type", NewName("Pages"))

	// For now, this is a default to be sure a box is set somewhere.
	// Clients can reset with their own call to SetMediaBox().
	d.SetMediaBox(0, 0, 612, 792)

	return d
}

func (d *Document) finishCurrentPage() {
	d.currentPage.Finalize()
	d.pages.Add(d.currentPage.Indirect())
	d.pageCount += 1
}

func (d *Document) finishPageTree() {
	d.pageTreeRoot.Add("Count", NewIntNumeric(int(d.pageCount)))
	d.pageTreeRoot.Add("Kids", d.pages)
	d.pageTreeRootIndirect.Finalize(d.pageTreeRoot)
}

func (d *Document) finishCatalog() {
	d.catalog.Add("Pages", d.pageTreeRootIndirect)
	d.catalogIndirect.Finalize(d.catalog)
}

func (d *Document) NewPage() *Page {
	if d.currentPage != nil {
		d.finishCurrentPage()
	}
	d.currentPage = NewPage()
	d.currentPage.BindToFile(d.file)
	d.currentPage.SetParent(d.pageTreeRootIndirect)
	return d.currentPage
}

func (d *Document) Close() {
	d.finishCurrentPage()
	d.finishPageTree()
	d.finishCatalog()
	d.file.Close()
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
