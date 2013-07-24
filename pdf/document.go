package pdf

import ("bufio"
	"fmt"
	"os")

type Document struct {
	file File
	// existing is true if the xref and trailer were read from an
	// existing document when the document was opened with
	// OpenDocument().
	existing bool

	// Pre-existing documents are opened as if they are read-only
	// until there's a reason to write something to them.  Adding
	// pages to a pre-existing document requires modifying the
	// page tree and rewriting it.  readyForNewPages is false for
	// pre-existing documents until NewPage() or some other method
	// is called that implies that new pages will be generated.
	// readyForNewPage is false if and only if pages and
	// procSetIndirect are both nil.
	readyForNewPages bool

	// When pages are added to a pre-existing document, the
	// existing page tree root is inserted as the first element in
	// the pages array.  Other elements in the pages array are
	// Indirect references to a Page dictionary returned by
	// page.Close() on a page obtained with NewPage()
	pages Array

	// procSetIndirect is nil if there are no new pages.
	procSetIndirect Indirect

	pageFactory *PageFactory
	streamFactory *StreamFactory

	// currentPage is nil until NextPage() is called.
	currentPage *Page

	// When a pre-existing document is opened, pageTreeRoot and
	// pageTreeRootIndirect are initialized with the pre-existing
	// dictionary.  Both are reset to a newly generated page tree
	// root if pages are added to an existing document (or one of
	// the page boxes is set).  Both are initialized using a newly
	// generated page tree if a new document is opened.  They are
	// not nil.
	pageTreeRoot *IndirectDictionary

	// pageCount is initialized with the pre-existing page count.
	pageCount uint

	// DocumentInfo is initialized from a pre-existing documents
	// document info dictionary.  Otherwise it is initialized to
	// an empty dictionary.  It is not nil.
	DocumentInfo
}

var (
	defaultStreamFactory *StreamFactory)

func init() {
	defaultStreamFactory = NewStreamFactory()
	ff := new(FlateFilter)
	ff.SetCompressionLevel(9)
	defaultStreamFactory.AddFilter(ff)
}

// makeNewPageTree() initializes the structures required to write a
// new or modified page tree to this document.  It need not be called
// if the document is only being read.
func (d *Document) makeNewPageTree() {
	d.readyForNewPages = true
	d.pages = NewArray()
	d.procSetIndirect = NewIndirect(d.file)

	newPageTreeRoot := NewIndirectDictionary(d.file)
	newPageTreeRoot.Add("Type", NewName("Pages"))
	newPageTreeRoot.Add("Count", NewIntNumeric(int(d.pageCount)))
	newPageTreeRoot.Add("Kids", d.pages)

	// If there is a pre-existing page tree insert the
	// whole thing as the first element of the pages array,
	// which is the /Kids array.
	if d.existing {
		d.pages.Add(d.pageTreeRoot)
		// Link the old page tree to the new one. and
		d.pageTreeRoot.Add("Parent", newPageTreeRoot)
		// Write out the revised version
		d.pageTreeRoot.Write()
	}
	d.pageTreeRoot = newPageTreeRoot

	// SetMediaBox() must be called after d.pageTreeRoot
	// is initialized For now, this is a default to be
	// sure a box is set somewhere.  Clients can reset
	// with their own call to SetMediaBox().
	d.SetMediaBox(0, 0, 612, 792)
}

// OpenDocument() constructs a document object from either a new or a pre-existing filename.
func OpenDocument(filename string, mode int) *Document {
	d := new(Document)

	d.file,d.existing,_ = OpenFile(filename, mode)

	if !d.existing {
		d.DocumentInfo = NewDocumentInfo()
		d.makeNewPageTree()
	} else {
		existingInfo := d.file.Info();
		if existingInfo == nil {
			d.DocumentInfo = NewDocumentInfo()
		} else {
			d.DocumentInfo = DocumentInfo{Dictionary: existingInfo, dirty: false}
		}
		

		existingPageTree := existingPageTree(d.file)
		d.pageTreeRoot = existingPageTree.root
		d.pageCount = existingPageTree.pageCount
		out := bufio.NewWriter(os.Stdout)
		out.WriteString("Pre-existing page tree root: ")
		d.pageTreeRoot.Serialize(out,d.file)
		out.WriteString("\n")
		out.Flush()
	}

	d.streamFactory = defaultStreamFactory
	d.pageFactory = NewPageFactory()
	d.pageFactory.SetStreamFactory(d.streamFactory)

	// Set a default producer field.  Clients calls to SetProducer() override this.
	d.SetProducer("PDFiG")

	return d
}

func (d *Document) release() {
	d.pages = nil
	d.currentPage = nil
	d.pageTreeRoot = nil
	d.procSetIndirect = nil
}

func (d *Document) finishCatalog() {
	if d.pageTreeRoot != nil {
		catalog := NewDictionary()
		catalog.Add("Type", NewName("Catalog"))
// TODO TODO 
		catalog.Add("Pages", d.pageTreeRoot)
		d.file.SetCatalog(catalog)
	}
}

func (d *Document) finishCurrentPage() {
	if d.currentPage != nil {
		d.pages.Add(d.currentPage.Finish())
		d.pageCount += 1
		d.pageTreeRoot.Add("Count", NewIntNumeric(int(d.pageCount)))
	}
}

func (d *Document) finishDocumentInfo() {
	if d.DocumentInfo.IsDirty() {
		d.file.SetInfo (d.DocumentInfo)
	}
}

func (d *Document) finishPageTree() {
	if d.pageTreeRoot != nil {
		d.pageTreeRoot.Write()
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

// NewPage() returns a Page reference.  A Page reference implements
// the io.writer interface which may be used to write raw PDF streams
// to the page's contents stream.  Pages created with
// Document.NewPage() are closed by the next call to
// Document.NewPage() or the call to Document.Close().
func (d *Document) NewPage() *Page {
	d.finishCurrentPage()
	d.currentPage = d.pageFactory.New(d.file)

	if !d.readyForNewPages {
		d.makeNewPageTree()
	}
	d.currentPage.SetParent(d.pageTreeRoot)
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

// Page(n) returns the ExistingPage (which contains a PageDictionary
// and an Indirect object) associated with page "n" of the document.
// The first page is numbered 0.  Any inheritable attributes found
// while descending the page tree are copied into the dictionary, so
// the dictionary may not exactly match the one in the file.
func (d *Document) Page(n uint) *ExistingPage {
	writer := bufio.NewWriter(os.Stdout)

	fmt.Fprintf (writer, "Page(%d) called with page root: ", n)
	d.pageTreeRoot.Serialize(writer, d.file)
	writer.WriteString("\n")
	writer.Flush()

	return pageFromTree(d.pageTreeRoot, n)
}

// SetStreamFactory() sets the StreamFactory used by the document for
// constructing page stream.  The client may call NewStreamFactory(),
// add filters, etc., and tell the document to use that factory.  The
// default factory uses LZW encoded streams.
func (d *Document) SetStreamFactory(sf *StreamFactory) {
	d.streamFactory = sf
	d.pageFactory.SetStreamFactory(sf)
}

func (d *Document) SetMediaBox(llx, lly, urx, ury float64) {
	if !d.readyForNewPages {
		d.makeNewPageTree()
	}
	d.pageTreeRoot.Add("MediaBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetCropBox(llx, lly, urx, ury float64) {
	if !d.readyForNewPages {
		d.makeNewPageTree()
	}
	d.pageTreeRoot.Add("CropBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetBleedBox(llx, lly, urx, ury float64) {
	if !d.readyForNewPages {
		d.makeNewPageTree()
	}
	d.pageTreeRoot.Add("BleedBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetTrimBox(llx, lly, urx, ury float64) {
	if !d.readyForNewPages {
		d.makeNewPageTree()
	}
	d.pageTreeRoot.Add("TrimBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) SetArtBox(llx, lly, urx, ury float64) {
	if !d.readyForNewPages {
		d.makeNewPageTree()
	}
	d.pageTreeRoot.Add("ArtBox", NewRectangle(llx, lly, urx, ury))
}

func (d *Document) WriteObject(object Object) Indirect {
	return NewIndirect(d.file).Write(object)
}
