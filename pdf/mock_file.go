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
	return NewIndirect(f).Finalize(object)
}

// Implements Object() in File interface
func (f *mockFile) Object(o ObjectNumber) (Object,error) {
	return nil,nil
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

func (f *mockFile) SetCatalog(i *Dictionary) {
}

func (f *mockFile) SetInfo(i DocumentInfo) {
}

func (f *mockFile) Info() *Dictionary {
	return nil
}

func (f *mockFile) Catalog() *Dictionary {
	return nil
}

func (f *mockFile) Trailer() *Dictionary {
	return nil
}

