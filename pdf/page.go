package pdf

type Page struct {
	contents *Stream
	parent   *Indirect

	dictionary *Dictionary
	resources  *Dictionary

	dictionaryIndirect *Indirect
	resourcesIndirect  *Indirect
}

func NewPage() *Page {
	p := new(Page)
	p.contents = NewStream()

	p.dictionary = NewDictionary()
	p.resources = NewDictionary()

	p.dictionaryIndirect = NewIndirect()
	p.resourcesIndirect = NewIndirect()

	return p
}

func (p *Page) Indirect() *Indirect {
	return p.dictionaryIndirect
}

func (p *Page) BindToFile(f File) {
	// The returned ObjectNumbers are ignored because they are not
	// of any use here.  We're after the side effect of binding
	// these objects to the File.
	p.dictionaryIndirect.ObjectNumber(f)
	p.resourcesIndirect.ObjectNumber(f)
}

func (p *Page) Finalize() {
	p.dictionary.Add("Resources", p.resourcesIndirect)
	p.dictionary.Add("Type", NewName("Page"))
	if p.parent == nil {
		panic("No parent specified")
	}
	p.dictionary.Add("Parent", p.parent)

	p.dictionaryIndirect.Finalize(p.dictionary)
	p.resourcesIndirect.Finalize(p.resources)
}

func (p *Page) SetParent(i *Indirect) {
	p.parent = i
}

func (p *Page) SetMediaBox(llx, lly, urx, ury float64) {
	p.dictionary.Add("MediaBox", NewRectangle(llx, lly, urx, ury))
}

func (p *Page) SetCropBox(llx, lly, urx, ury float64) {
	p.dictionary.Add("CropBox", NewRectangle(llx, lly, urx, ury))
}

func (p *Page) SetBleedBox(llx, lly, urx, ury float64) {
	p.dictionary.Add("BleedBox", NewRectangle(llx, lly, urx, ury))
}

func (p *Page) SetTrimBox(llx, lly, urx, ury float64) {
	p.dictionary.Add("TrimBox", NewRectangle(llx, lly, urx, ury))
}

func (p *Page) SetArtBox(llx, lly, urx, ury float64) {
	p.dictionary.Add("ArtBox", NewRectangle(llx, lly, urx, ury))
}
