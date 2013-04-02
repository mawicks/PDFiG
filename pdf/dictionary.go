package pdf

// Implements the pdf.Object interface

type Dictionary struct {
	dictionary map[string] Object
}

// Constructor for Dictionary object
func NewDictionary () (* Dictionary) {
	return &Dictionary{make(map[string] Object, 16)}
}

func (d *Dictionary) Add (key string, o Object) {
	d.dictionary[key] = o
}

func (d *Dictionary) Remove (key string) {
	delete (d.dictionary, key)
}

func (d *Dictionary) Serialize (w Writer, file... File) {
	w.WriteString("<<")
	haveAny := false;
	for key,value := range d.dictionary {
		if (haveAny) {
			w.WriteByte(' ')
		}
		NewName(key).Serialize(w)
		w.WriteByte(' ')
		value.Serialize(w)
		haveAny = true
	}
	w.WriteString(">>")
}
