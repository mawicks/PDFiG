package pdf

import (
	"errors" )


type pageTree struct {
	file File
	root *Dictionary
}

func oldPageTree(file File) *pageTree{
	var (
		catalog, pageTreeRoot *Dictionary
		pageTreeRootIndirect *Indirect
		object Object
		ok bool
		err error )

	if catalog = file.Catalog(); catalog == nil || !catalog.CheckNameValue("Type","Catalog") {
		panic (errors.New(`Document has no catalog or catalog dictionary type is not "Catalog"`))
	}

	if pageTreeRootIndirect,ok = catalog.Get("Pages").(*Indirect) ; !ok {
		panic (errors.New(`/Pages entry is not an indirect reference`))
	}

	if object,err = file.Object(pageTreeRootIndirect.ObjectNumber(file)); err != nil {
		panic (errors.New(`Unable to read page tree dictioanry`))
	}

	if pageTreeRoot,ok = object.(*Dictionary); !ok {
		panic (errors.New(`Page tree root object is not a dictionary`))
	}

	pt := new(pageTree)
	pt.file = file
	pt.root = pageTreeRoot
	return pt
}

