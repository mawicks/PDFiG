package pdf

import (
	"errors" )


type pageTree struct {
	file File
	root *Dictionary
	rootReference *Indirect
	pageCount int
}

func oldPageTree(file File) *pageTree{
	var (
		catalog, pageTreeRoot *Dictionary
		pageTreeRootReference *Indirect
		pageCount *IntNumeric
		ok bool )

	if catalog = file.Catalog(); catalog == nil || !catalog.CheckNameValue("Type","Catalog") {
		panic (errors.New(`Document has no catalog or catalog dictionary type is not "Catalog"`))
	}

	if pageTreeRootReference,ok = catalog.Get("Pages").(*Indirect); !ok {
		panic (errors.New(`/Pages entry missing or is not an indirect reference`))
	}

	if pageTreeRoot,ok = pageTreeRootReference.Dereference(file).(*Dictionary); !ok {
		panic (errors.New(`Page tree root object is not a dictionary`))
	}

	if pageCount,ok = pageTreeRoot.Get("Count").Dereference(file).(*IntNumeric); !ok {
		panic (errors.New(`/Count value is not an integer`))
	}

	pt := new(pageTree)
	pt.file = file
	pt.root = pageTreeRoot
	pt.rootReference = pageTreeRootReference
	pt.pageCount = pageCount.Value()

	return pt
}

