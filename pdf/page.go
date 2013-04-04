package pdf

type Page struct {
	contents *Stream
	parent   *Indirect

	dictionary *Dictionary
	resources  *Dictionary

	dictionaryIndirect *Indirect
	resourcesIndirect  *Indirect
	contentsIndirect *Indirect
}

func NewPage(file... File) *Page {
	p := new(Page)
	p.contents = NewStream()

	p.dictionary = NewDictionary()
	p.resources = NewDictionary()

	p.dictionaryIndirect = NewIndirect(file...)
	p.resourcesIndirect = NewIndirect(file...)
	p.contentsIndirect = NewIndirect(file...)

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
	p.contentsIndirect.ObjectNumber(f)
}

func (p *Page) Finalize() {
	p.dictionary.Add("Resources", p.resourcesIndirect)
	p.dictionary.Add("Type", NewName("Page"))
	p.dictionary.Add("Contents", p.contentsIndirect)
	if p.parent == nil {
		panic("No parent specified")
	}
	p.dictionary.Add("Parent", p.parent)

	p.dictionaryIndirect.Finalize(p.dictionary)
	p.contentsIndirect.Finalize(p.contents)
	p.resourcesIndirect.Finalize(p.resources)
}

func (p *Page) SetParent(i *Indirect) {
	p.parent = i
}

func (p *Page) SetProcSet(i *Indirect) {
	p.resources.Add("ProcSet", i)
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

func (p *Page) Write(b []byte) (int, error) {
	return p.contents.Write(b)
}
