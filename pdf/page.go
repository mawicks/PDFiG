package pdf

import ("errors"
	"fmt"
	"strconv")

type Page struct {
	fileList []File
	contents *Stream
	parent *Indirect

	dictionary, resources, fontResources *Dictionary

	fontMap map[Font] string
}

// There is no constructor here.  Pages are created by a PageFactory.New().

func (p *Page) Close() *Indirect {
	if (p.fontResources != nil) {
		p.resources.Add("Font", p.fontResources)
		p.fontResources = nil
	}

	p.dictionary.Add("Resources", NewIndirect(p.fileList...).Write(p.resources))
	p.resources = nil

	p.dictionary.Add("Contents", NewIndirect(p.fileList...).Write(p.contents))
	p.contents = nil

	if p.parent == nil {
		panic("No parent specified")
	}
	p.dictionary.Add("Parent", p.parent)
	p.parent = nil

	p.dictionary.Add("Type", NewName("Page"))

	indirect := NewIndirect(p.fileList...).Write(p.dictionary)
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

func (p *Page) setBox (boxname string, llx, lly, urx, ury float64) {
	if p.dictionary == nil {
		panic (errors.New(fmt.Sprintf(`Attempt to set "%s" on a closed Page`, boxname)))
	}
	p.dictionary.Add(boxname, NewRectangle(llx, lly, urx, ury))
}

func (p *Page) SetMediaBox(llx, lly, urx, ury float64) {
	p.setBox("MediaBox", llx, lly, urx, ury)
}

func (p *Page) SetCropBox(llx, lly, urx, ury float64) {
	p.setBox("CropBox", llx, lly, urx, ury)
}

func (p *Page) SetBleedBox(llx, lly, urx, ury float64) {
	p.setBox("BleedBox", llx, lly, urx, ury)
}

func (p *Page) SetTrimBox(llx, lly, urx, ury float64) {
	p.setBox("TrimBox", llx, lly, urx, ury)
}

func (p *Page) SetArtBox(llx, lly, urx, ury float64) {
	p.setBox("ArtBox", llx, lly, urx, ury)
}

func (p *Page) Write(b []byte) (int, error) {
	if p.contents == nil {
		panic (errors.New("Attempt to write to a closed Page"))
	}
	return p.contents.Write(b)
}

