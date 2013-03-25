/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

// TestFile is a simple file implementing the File interface for use in unit tests.
type TestFile struct {
	nextObjectNumber uint32
	nextGenerationNumber uint16
}

// Constructor for Stream object
func NewTestFile (obj uint32, gen uint16) File {
	return &TestFile{obj,gen}
}

func (f *TestFile) AssignObjectNumber (o Object) ObjectNumber {
	result := ObjectNumber{f.nextObjectNumber,f.nextGenerationNumber}
	f.nextObjectNumber += 1
	f.nextGenerationNumber += 1
	return result
}
