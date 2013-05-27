package pdf

import ("errors")

type pageTree struct {
	root *Dictionary
	rootReference *Indirect
	pageCount uint
}

func existingPageTree(file File) *pageTree {
	var (
		catalog, pageTreeRoot *Dictionary
		pageTreeRootReference *Indirect
		pageCount int
		ok bool )

	if catalog = file.Catalog(); catalog == nil || !catalog.CheckNameValue("Type","Catalog") {
		panic (errors.New(`Document has no catalog or catalog dictionary type is not "Catalog"`))
	}

	if pageTreeRootReference,ok = catalog.GetIndirect("Pages"); !ok {
		panic (errors.New(`/Pages entry missing or is not an indirect reference`))
	}

	if pageTreeRoot,ok = catalog.GetDictionary("Pages"); !ok {
		panic (errors.New(`Page tree root object is not a dictionary`))
	}

	if pageCount,ok = pageTreeRoot.GetInt("Count"); !ok {
		panic (errors.New(`/Count value is not an integer`))
	}

	return &pageTree{pageTreeRoot,pageTreeRootReference,uint(pageCount)}
}

func copyDictionaryEntries(dst, src *Dictionary, list []string) {
	for _,name := range list {
		if dst.Get(name) == nil {
			if value := src.Get(name); value != nil {
				dst.Add(name, value.Clone())
			}
		}
	}
}

// Return the Page dictionary and an indirect object corresponding to
// the nth page in the tree.  The first page is numbered 0.  Any
// inheritable attributes found while descending the tree are copied
// into the dictionary, so the dictionary returned does not exactly
// match the one in the file.
func pageFromTree (node *Dictionary, n uint) *ExistingPage {
	var (
		kids *Array
		ok bool )

	if kids,ok = node.GetArray("Kids"); !ok {
		panic (errors.New(`Page tree node has no "Kids" array`))
	}

	kidCount := kids.Size()
	for i:=0; i<kidCount && n >= 0; i++ {
		var (
			count int
			kid *Dictionary
			kidReference *Indirect
			nodeType string )

		if kidReference,ok = kids.At(i).(*Indirect); !ok {
			panic (errors.New(`Kids array contains an object that isn't an indirect reference.`))
		}
		if kid,ok = kidReference.Dereference().(*Dictionary); !ok {
			panic (errors.New(`Kids array contains an object that isn't a reference to a dictionary.`))
		}
		if nodeType,ok = kid.GetName("Type"); !ok {
			panic (errors.New(`Node in page tree missing /Type entry.`))
		}
		copyDictionaryEntries(kid,node,[]string{"Resources", "MediaBox", "CropBox", "Rotate"})
		switch nodeType {
		case "Pages":
			if count,ok = kid.GetInt("Count"); !ok {
				panic (errors.New(`Page tree node missing /Count`))
			}
			if n < uint(count) {
				return pageFromTree(kid,n)
			}
			n -= uint(count)
		case "Page":
			if n == 0 {
				return &ExistingPage{kid,kidReference}
			}
			n -= 1
		default:
			panic (errors.New(`Unknown page tree node type`))
		}
	}
	return nil
}