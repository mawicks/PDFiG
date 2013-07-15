package pdf

type ExistingPage struct {
	*PageDictionary
	reference Indirect
}

func (ep *ExistingPage) Rewrite() {
	ep.PageDictionary.Write(ep.reference)
}

