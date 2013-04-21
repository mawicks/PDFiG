package pdf

import "strconv"

type Page struct {
	fileList []File
	contents *Stream
	parent *Indirect

	dictionary, resources, fontResources *Dictionary

	fontMap map[Font] string
}

func NewPage(file... File) *Page {
	p := new(Page)
	p.fileList = file
	p.contents = NewStream()
	p.parent = nil

	p.dictionary = NewDictionary()
	p.resources = NewDictionary()
	p.fontResources = nil

	p.fontMap = make(map[Font]string, 15)

	return p
}

func (p *Page) Finalize() *Indirect {
	if (p.fontResources != nil) {
		p.resources.Add("Font", p.fontResources)
		p.fontResources = nil
	}

	p.dictionary.Add("Resources", NewIndirect(p.fileList...).Finalize(p.resources))
	p.resources = nil

	p.dictionary.Add("Contents", NewIndirect(p.fileList...).Finalize(p.contents))
	p.contents = nil

	if p.parent == nil {
		panic("No parent specified")
	}
	p.dictionary.Add("Parent", p.parent)
	p.parent = nil

	p.dictionary.Add("Type", NewName("Page"))

	indirect := NewIndirect(p.fileList...).Finalize(p.dictionary)
	p.dictionary = nil

	return indirect
}

func (p *Page) AddFont (font Font) string {
	fontCount := len(p.fontMap)

	if fontCount >= (1<<20) {
		panic("Too many fonts on one page")
	}

	if (p.fontResources == nil) {
		p.fontResources = NewDictionary()
	}

	name,exists := p.fontMap[font]

	if (!exists) {
		name = "F" + strconv.Itoa(fontCount + 1)
		for _,file := range p.fileList {
			p.fontResources.Add(name, font.Indirect(file))
		}
		p.fontMap[font] = name
	}

	return name
}

func (p *Page) SetParent(i *Indirect) {
	p.parent = i
}

func (p *Page) setProcSet(i *Indirect) {
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
