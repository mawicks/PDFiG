package pdf

// TestFile is a simple file implementing the File interface for use in unit tests.
type mockFile struct {
	nextObjectNumber     uint32
	nextGenerationNumber uint16
}

// Constructor for mockFile object
func NewMockFile(obj uint32, gen uint16) File {
	return &mockFile{obj, gen}
}

// Public methods

// Implements Close() in File interface
func (f *mockFile) Close() {}

// Implements AddObjectAt() in File interface
func (f *mockFile) AddObjectAt(ObjectNumber, Object) {}

// Implements AddObject() in File interface
func (f *mockFile) AddObject(object Object) (reference *Indirect) {
	reference = NewIndirect(f)
	reference.Finalize(object)
	return reference
}

// Implements DeleteObject() in File interface
func (f *mockFile) DeleteObject(*Indirect) {}

// Implements ReserveObjectNumber() in File interface
func (f *mockFile) ReserveObjectNumber(o Object) ObjectNumber {
	result := ObjectNumber{f.nextObjectNumber, f.nextGenerationNumber}
	f.nextObjectNumber += 1
	f.nextGenerationNumber += 1
	return result
}

func (f *mockFile) SetCatalog(i *Indirect) {
}

func (f *mockFile) SetInfo(i *Indirect) {
}

