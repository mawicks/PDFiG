/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bufio"

// Implements the pdf.Object interface

type Dictionary struct {
	dictionary map[string] Object
}


// Constructor for Name object
func NewDictionary () (* Dictionary) {
	return &Dictionary{make(map[string] Object, 16)}
}

func (d *Dictionary) Add (key string, o Object) {
	d.dictionary[key] = o
}

func (d *Dictionary) Remove (key string) {
	delete (d.dictionary, key)
}

func (d *Dictionary) Serialize (f *bufio.Writer) {
	f.WriteString("<<")
	haveAny := false;
	for key,value := range d.dictionary {
		if (haveAny) {
			f.WriteByte(' ')
		}
		NewName(key).Serialize(f)
		f.WriteByte(' ')
		value.Serialize(f)
		haveAny = true
	}
	f.WriteString(">>")
}
