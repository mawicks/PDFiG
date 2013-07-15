package pdf

import (
	"bytes"
	"errors"
	"io"
)

// A PageDictionary wraps a Dictionary to simplify access and to limit
// the operations to those that are valid for a page dictionary.
type PageDictionary struct {
	dictionary Dictionary
	hasParent bool
}

func NewPageDictionary() *PageDictionary {
	pd := NewDictionary()
	pd.Add("Type", NewName("Page"))
	return &PageDictionary{pd,false}
}

func (pd *PageDictionary) CloneDictionary() Dictionary {
	dictionary,_ := pd.dictionary.Clone().(Dictionary)
	return dictionary
}

// Return an io.Reader that will read the page's contents.  Stream
// filters are applied and segments from multi-segment streams are
// concatenated.
func (pd *PageDictionary) Reader() io.Reader {
	// If the contents consist of a single stream, return its
	// stream reader.
	if pageStream,ok := pd.dictionary.GetStream("Contents"); ok {
		return pageStream.Reader()
	}
	// Otherwise, the contents should be an array of streams to be
	// concatenated.
	if pageStreamArray,ok := pd.dictionary.GetArray("Contents"); ok {
		n := pageStreamArray.Size()
		if n > 0 {
			readers := make([]io.Reader, 2*n-1)
			nr := 0
			for i:=0; i<n; i++ {
				if i != 0 {
				// Insert whitespace between each stream
					readers[nr] = bytes.NewReader([]byte(" "))
					nr += 1
				}
				if streamReference,ok := pageStreamArray.At(i).(*Indirect); ok {
					if stream,ok := streamReference.Dereference().(Stream); ok {
						readers[nr] = stream.Reader()
					}
				}
				if readers[nr] == nil {
					return nil
				}
				nr += 1
			}
			return io.MultiReader (readers...)
		}
	}
	// Otherwise, the dictionary is partially constructed or the
	// PDF file is invalid.
	return nil
}

// If the dictionary's Contents field is not an array, make it one.
// The dictionary's Contents field should be either an array or an
// indirect object.
func (pd *PageDictionary) ensureContentsIsArray() *Array {

	if pageContentsArray,ok := pd.dictionary.GetArray("Contents"); ok {
		return pageContentsArray
	}

	if contents, ok := pd.dictionary.GetIndirect("Contents"); ok {
		contentsArray := NewArray()
		contentsArray.Add(contents)
		pd.dictionary.Add("Contents", contentsArray)
		return contentsArray
	}
	
	// Dictionary's Contents field is neither an array nor an
	// Indirect object.  Construction of the dictionary may not be
	// finished or the file may be invalid.
	return nil
}

// PrependContents() prepends the passed indirect reference (which
// must reference a stream) in front of the page contents.  The client is
// responsible for ensuring the indirect reference is associated with
// a stream object.
func (pd *PageDictionary) PrependContents(is *Indirect) {
	if contentsArray := pd.ensureContentsIsArray(); contentsArray != nil {
		contentsArray.PushFront(is)
	} else {
		panic (errors.New("Contents dictionary value is neither an array nor an indirect"))
	}
}

// AppendContents() appends the passed indirect reference (which must
// reference a stream) onto the page contents.  The client is
// responsible for ensuring the indirect reference is associated with
// a stream object.
func (pd *PageDictionary) AppendContents(is *Indirect) {
	if contentsArray := pd.ensureContentsIsArray(); contentsArray != nil {
		contentsArray.Add(is)
	} else {
		panic (errors.New("Contents dictionary value is neither an array nor an indirect"))
	}
}

// SetContents() sets the Contents value in the page dictionary to the
// passed indirect reference.  The client is responsible for ensuring
// that the indirect reference is associated with a stream or possibly
// with an array of stream references.
func (pd *PageDictionary) SetContents(is *Indirect) {
	if is == nil {
		panic ("Indirect object is nil")
	}
	pd.dictionary.Add("Contents", is)
}

// SetResources() sets the Resources value in the page dictionary to the
// passed indirect reference.  The client is responsible for ensuring
// that the indirect reference is a valid Resource dictionary.
func (pd *PageDictionary) SetResources(ir *Indirect) {
	if ir == nil {
		panic ("Indirect object is nil")
	}
	pd.dictionary.Add("Resources", ir)
}

// SetParent() sets the Parent value in the page dictionary to the
// passed indirect reference.  The client is responsible for ensuring
// that the indirect reference is a valid page dictionary or pages node
// reference.
func (pd *PageDictionary) SetParent(ip *Indirect) {
	if ip == nil {
		panic ("Indirect object is nil")
	}
	pd.dictionary.Add("Parent", ip)
	pd.hasParent = true
}

func (pd *PageDictionary) setBox (boxname string, llx, lly, urx, ury float64) {
	pd.dictionary.Add(boxname, NewRectangle(llx, lly, urx, ury))
}

func (pd *PageDictionary) SetMediaBox(llx, lly, urx, ury float64) {
	pd.setBox("MediaBox", llx, lly, urx, ury)
}

func (pd *PageDictionary) SetCropBox(llx, lly, urx, ury float64) {
	pd.setBox("CropBox", llx, lly, urx, ury)
}

func (pd *PageDictionary) SetBleedBox(llx, lly, urx, ury float64) {
	pd.setBox("BleedBox", llx, lly, urx, ury)
}

func (pd *PageDictionary) SetTrimBox(llx, lly, urx, ury float64) {
	pd.setBox("TrimBox", llx, lly, urx, ury)
}

func (pd *PageDictionary) SetArtBox(llx, lly, urx, ury float64) {
	pd.setBox("ArtBox", llx, lly, urx, ury)
}

func (pd *PageDictionary) Write(id *Indirect) *Indirect {
	if !pd.hasParent {
		panic("PageDictionary has no Parent")
	}
	id.Write(pd.dictionary)
	return id
}
