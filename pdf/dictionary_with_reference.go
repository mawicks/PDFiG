package pdf

import (
	"errors"
	"fmt"
)

// DictionaryWithReference is used for the common case where you need an
// indirect reference to a dictionary.  DictionaryWithReference contains
// both the (direct) dictionary and the indirect reference.  It
// implements Dictionary and Indirect so it may be used interchangably
// as either an Indirect reference or a (direct) Dictionary.  It has a
// Write() method with a different signature than is required for the
// Indirect interface, so it cannot be used as an Indirect.  The
// normal use-case is to use DictionaryWithReference as a Dictionary
// interface to construct the desired dictionary.  Once the dictionary
// is completed, call DictionaryWithReference.Write() to flush the result
// to the output file and then continue to use DictionaryWithReference as a
// Indirect.
type DictionaryWithReference struct {
	Dictionary
	indirect Indirect
}

// Likewise, ProtectedDictionaryWithReference implements ProtectedDicationary
type ProtectedDictionaryWithReference struct {
	ProtectedDictionary
	indirect Indirect
}

// Since DictionaryWithReference delegates to its internal Dictionary, it
// fully implements Object and Dictionary.

func NewDictionaryWithReference(file... File) *DictionaryWithReference {
	return &DictionaryWithReference{NewDictionary(),NewIndirect(file...)}
}

func (id *DictionaryWithReference) Clone() Object {
	var clonedDictionary Dictionary = nil
	if id.Dictionary != nil {
		clonedDictionary = id.Dictionary.Clone().(Dictionary)
	}
	return &DictionaryWithReference{clonedDictionary,id.indirect.Clone().(Indirect)}
}

func (id *DictionaryWithReference) Protect() Object {
	var  protectedDictionary ProtectedDictionary = nil
	if id.Dictionary != nil {
		protectedDictionary = id.Dictionary.Protect().(ProtectedDictionary)
	}
	return &ProtectedDictionaryWithReference{protectedDictionary,id.indirect.Protect().(Indirect)}
}

func (id *DictionaryWithReference) Unprotect() Object {
	return id
}

// ObjectNumber() is required for the Indirect interface.
func (id *DictionaryWithReference) ObjectNumber(f File) ObjectNumber {
	return id.indirect.ObjectNumber(f) }

// BoundToFile() is required for Indirect interface.
func (id *DictionaryWithReference) BoundToFile(f File) bool {
	return id.indirect.BoundToFile(f)
}

// Write() writes the dictionary to the PDF file as an indirect
// object.  The direct object is released so Write() should be called
// only once.  Clients may continue to use the indirect reference
// after calling Write().  This Write() has a different signature from
// the one appearing in Indirect so DictionaryWithReference does not
// implement Indirect.
func (id *DictionaryWithReference) Write() {
	if id.Dictionary == nil {
		panic(errors.New(fmt.Sprintf(`DictionaryWithReference.Write() called more than once`)))
	}
	id.indirect.Write(id.Dictionary)
	id.Dictionary = nil
}

func (id *DictionaryWithReference) Serialize(w Writer, file... File) {
	id.indirect.Serialize(w, file...)
}

func (id *ProtectedDictionaryWithReference) Clone() Object {
	var protectedDictionary ProtectedDictionary
	if id.ProtectedDictionary != nil {
		protectedDictionary = id.ProtectedDictionary.Clone().(ProtectedDictionary)
	}
	return &ProtectedDictionaryWithReference{protectedDictionary,id.indirect.Clone().(Indirect)}
}

func (id *ProtectedDictionaryWithReference) Protect() Object {
	return id
}

func (id *ProtectedDictionaryWithReference) Unprotect() Object {
	var dictionary Dictionary
	if id.ProtectedDictionary != nil {
		dictionary = id.ProtectedDictionary.Unprotect().(Dictionary)
	}
	return &DictionaryWithReference{dictionary,id.indirect.Unprotect().(Indirect)}
}

func (id *ProtectedDictionaryWithReference) Serialize(w Writer, file... File) {
	id.indirect.Serialize(w, file...)
}














