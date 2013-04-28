package pdf

type ObjectNumber struct {
	number     uint32
	generation uint16
}

func NewObjectNumber(number uint32, generation uint16) ObjectNumber {
	return ObjectNumber{number,generation}
}

type File interface {
	// WriteObject() adds the passed object to the File.  The
	// returned indirect reference may be used for backward
	// references to the object.  A new object is created
	// either at a new index in the xref or at an old index
	// using a new generation.
	WriteObject(Object) (*Indirect)

	// WriteObjectAt() adds the object to the File at the specified
	// location.  ObjectNumber may have been obtained by an
	// earlier call to ReserveObjectNumber(), or ObjectNumber may
	// be a pre-existing (finalized) object that is being
	// overwritten with a modified copy.
	WriteObjectAt(ObjectNumber, Object)

	// Object() retrieves a finalized object that has already been
	// written to a PDF file.
	Object(ObjectNumber) (Object,error)

	// ReserveObjectNumber() reserves a position (ObjectNumber)
	// for the passed object in the File.
	ReserveObjectNumber(Object) ObjectNumber

	// Info() returns a copy of the Info dictionary.  Caller may
	// modify the copy and use SetInfo() to replace the file's
	// info dictionary
	Info() *Dictionary

	// Catalog() returns a copy of the Info dictionary.  Caller
	// may modify the copy and use SetCatalog() to replace the
	// file's info dictionary
	Catalog() *Dictionary

	// SetCatalog() sets the catalog dictionary
	SetCatalog(*Dictionary)

	// SetInfo() sets the Info dictionary
	SetInfo(DocumentInfo)

	// Trailer() returns a copy of the current contents of the
	// trailer dictionary
	Trailer() *Dictionary

	// DeleteObject() deletes the specified object from the file.
	// It must be an indirect object.
	DeleteObject(*Indirect)

	// Close() writes the xref, trailer, etc., and closes the
	// underlying file.
	Close()
}
