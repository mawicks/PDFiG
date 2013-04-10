package pdf

type DocumentInfo struct {
	*Dictionary
}

func NewDocumentInfo() DocumentInfo {
	return DocumentInfo{NewDictionary()}
}

func (d DocumentInfo) SetTitle(s string) {
	d.Add("Title", NewTextString(s))
}

func (d DocumentInfo) SetAuthor(s string) {
	d.Add("Author", NewTextString(s))
}

func (d DocumentInfo) SetSubject(s string) {
	d.Add("Subject", NewTextString(s))
}

func (d DocumentInfo) SetKeywords(s string) {
	d.Add("Keywords", NewTextString(s))
}

func (d DocumentInfo) SetCreator(s string) {
	d.Add("Creator", NewTextString(s))
}

func (d DocumentInfo) SetProducer(s string) {
	d.Add("Producer", NewTextString(s))
}
