package pdf

// import (
// 	"fmt" )

type PageFactory struct {
	*StreamFactory
}

func NewPageFactory() *PageFactory {
	result := new(PageFactory)
	result.StreamFactory = NewStreamFactory()
	return result
}

func (pf *PageFactory) New (file... File) *Page {
	p := new(Page)

	p.fileList = file
	p.contents = pf.StreamFactory.New()

	p.parent = nil
	p.dictionary = NewPageDictionary()
	p.resources = NewDictionary()

	p.fontResources = nil
	p.fontMap = make(map[Font]string, 15)

	return p
}

func (pf *PageFactory) SetStreamFactory (sf *StreamFactory) {
	pf.StreamFactory = sf
}