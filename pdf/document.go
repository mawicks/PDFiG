package pdf

import ("bufio"
	"os")

type Document struct {
	file File
	existing bool

	// Each element in the pages array is an Indirect reference to
	// a Page dictionary returned by Page.Indirect().  The "pages"
	// array is nil for pre-existing documents until NewPage() (or
	// some other method that implies pages will be generated)
	// gets called.  The "pages" array contains only *new* pages.
	pages *Array
	// procSetIndirect is nil if there are no new pages.
	procSetIndirect *Indirect
	// currentPage is nil until NextPage() is called.
	currentPage *Page

	// When a pre-existing document is opened, pageTreeRoot and
	// pageTreeRootIndirect are initialized with the pre-existing
	// dictionary.
	pageTreeRoot *Dictionary
	pageTreeRootIndirect *Indirect
	// pageCount is also initialized with the pre-existing page count.
	pageCount uint

	DocumentInfo
}

// needPageTree() initializes the structures required to write a new
// or modified page tree to this document.  It need not be called if
// the document is only being read.
func (d *Document) needPageTree() {
	if d.pages == nil {
		d.pages = NewArray()
		d.procSetIndirect = NewIndirect(d.file)

		newPageTreeRoot := NewDictionary()
		newPageTreeRootIndirect := NewIndirect(d.file)
		// If there is a pre-existing page tree insert the
		// whole thing as the first element of the pages array
		// (which will become /Kids).
		if d.existing {
			d.pages.Add(d.pageTreeRootIndirect)
			// Link the old page tree to the new one. and
			d.pageTreeRoot.Add("Parent", newPageTreeRootIndirect)
			// Write out the revised version
			d.pageTreeRootIndirect.Write(d.pageTreeRoot)
		}
		d.pageTreeRoot = newPageTreeRoot
		d.pageTreeRootIndirect = newPageTreeRootIndirect
		// SetMediaBox() must be called after d.pageTreeRoot
		// is initialized For now, this is a default to be
		// sure a box is set somewhere.  Clients can reset
		// with their own call to SetMediaBox().
		d.SetMediaBox(0, 0, 612, 792)
	}
}

// OpenDocument() constructs a document object from either a new or a pre-existing filename.
func OpenDocument(filename string, mode int) *Document {
	d := new(Document)

	d.file,d.existing,_ = OpenFile(filename, mode)

	if !d.existing {
		d.DocumentInfo = NewDocumentInfo()
		d.needPageTree()
	} else {
		d.DocumentInfo = DocumentInfo{Dictionary: d.file.Info(), dirty: false}
		oldPageTree := oldPageTree(d.file)
		d.pageTreeRoot = oldPageTree.root
		d.pageTreeRootIndirect = oldPageTree.rootReference
		d.pageCount = oldPageTree.pageCount
		out := bufio.NewWriter(os.Stdout)
		out.WriteString("Pre-existing page tree root: ")
		d.pageTreeRoot.Serialize(out,d.file)
		out.WriteString("\n")
		out.Flush()
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
		d.pages.Add(d.currentPage.Close())
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

		d.pageTreeRootIndirect.Write(d.pageTreeRoot)
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
		d.procSetIndirect.Write(procsetArray)
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
