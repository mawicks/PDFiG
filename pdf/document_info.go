package pdf

import "fmt"

type DocumentInfo struct {
	*Dictionary
	dirty bool
}

func NewDocumentInfo() DocumentInfo {
	fmt.Printf("NewDocumentInfo()\n")
	return DocumentInfo{NewDictionary(), false}
}

func (d DocumentInfo) IsDirty() bool {
	return d.dirty
}

func (d DocumentInfo) SetTitle(s string) {
	d.dirty = true
	d.Add("Title", NewTextString(s))
}

func (d DocumentInfo) SetAuthor(s string) {
	d.dirty = true
	d.Add("Author", NewTextString(s))
}

func (d DocumentInfo) SetSubject(s string) {
	d.dirty = true
	d.Add("Subject", NewTextString(s))
}

func (d DocumentInfo) SetKeywords(s string) {
	d.dirty = true
	d.Add("Keywords", NewTextString(s))
}

func (d DocumentInfo) SetCreator(s string) {
	d.dirty = true
	d.Add("Creator", NewTextString(s))
}

func (d DocumentInfo) SetProducer(s string) {
	d.dirty = true
	d.Add("Producer", NewTextString(s))
}
