package pdf

import (
	"errors"
	"fmt"
)

// IndirectDictionary is used for the common case where you need an
// indirect reference to a dictionary.  IndirectDictionary contains
// both the (direct) dictionary and the indirect reference.  It
// implements Dictionary and Indirect so it may be used interchangably
// as either an Indirect reference or a (direct) Dictionary.  It has a
// Write() method with a different signature than is required for the
// Indirect interface, so it cannot be used as an Indirect.  The
// normal use-case is to use IndirectDictionary as a Dictionary
// interface to construct the desired dictionary.  Once the dictionary
// is completed, call IndirectDictionary.Write() to flush the result
// to the output file and then continue to use IndirectDictionary as a
// Indirect.
type IndirectDictionary struct {
	Dictionary
	indirect Indirect
}

// Since IndirectDictionary delegates to its internal Dictionary, it
// fully implements Dictionary.

func NewIndirectDictionary(file... File) *IndirectDictionary {
	return &IndirectDictionary{NewDictionary(),NewIndirect(file...)}
}

// ObjectNumber() is required for the Indirect interface.
func (id *IndirectDictionary) ObjectNumber(f File) ObjectNumber {
	return id.indirect.ObjectNumber(f) }

// BoundToFile() is required for Indirect interface.
func (id *IndirectDictionary) BoundToFile(f File) bool {
	return id.indirect.BoundToFile(f)
}

// Write() writes the dictionary to the PDF file as an indirect
// object.  The direct object is released so Write() should be called
// only once.  Clients may continue to use the indirect reference
// after calling Write().  This Write() has a different signature from
// the one appearing in Indirect so IndirectDictionary does not
// implement Indirect.
func (id *IndirectDictionary) Write() {
	if id.Dictionary == nil {
		panic(errors.New(fmt.Sprintf(`IndirectDictionary.Write() called more than once`)))
	}
	id.indirect.Write(id.Dictionary)
	id.Dictionary = nil
}

func (id *IndirectDictionary) Serialize(w Writer, file ...File) {
	id.indirect.Serialize(w, file...)
}
























