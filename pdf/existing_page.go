package pdf

type ExistingPage struct {
	*PageDictionary
}

func (ep *ExistingPage) Rewrite() {
	ep.PageDictionary.Write()
}











