package pdf

type DocumentInfo struct {
	*Dictionary
}

func NewDocumentInfo() DocumentInfo {
	return DocumentInfo{NewDictionary()}
}

func (d DocumentInfo) SetTitle(s string) {
	d.Add("Title", NewString(s))
}

func (d DocumentInfo) SetAuthor(s string) {
	d.Add("Author", NewString(s))
}

func (d DocumentInfo) SetSubject(s string) {
	d.Add("Subject", NewString(s))
}

func (d DocumentInfo) SetKeywords(s string) {
	d.Add("Keywords", NewString(s))
}

func (d DocumentInfo) SetCreator(s string) {
	d.Add("Creator", NewString(s))
}

func (d DocumentInfo) SetProducer(s string) {
	d.Add("Producer", NewString(s))
}
