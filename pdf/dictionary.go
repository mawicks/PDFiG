package pdf

// Implements the pdf.Object interface

type Dictionary struct {
	dictionary map[string]Object
}

// Constructor for Dictionary object
func NewDictionary() *Dictionary {
	return &Dictionary{make(map[string]Object, 16)}
}

func (d *Dictionary) Clone() Object {
	newDictionary := NewDictionary()
	for key,value := range d.dictionary {
		newDictionary.dictionary[key] = value.Clone()
	}
	return newDictionary
}

func (d *Dictionary) Add(key string, o Object) {
	d.dictionary[key] = o
}

func (d *Dictionary) Get(key string) Object {
	return d.dictionary[key]
}

func (d *Dictionary) Remove(key string) {
	delete(d.dictionary, key)
}

func (d *Dictionary) Serialize(w Writer, file ...File) {
	w.WriteString("<<")
	haveAny := false
	for key, value := range d.dictionary {
		if haveAny {
			w.WriteByte(' ')
		}
		NewName(key).Serialize(w, file...)
		w.WriteByte(' ')
		value.Serialize(w, file...)
		haveAny = true
	}
	w.WriteString(">>")
}

// Size() returns the number of key-value pairs
func (d *Dictionary) Size() int {
	return len(d.dictionary)
}

func (d *Dictionary) CheckNameValue (key string, expected string) bool {
	if value,ok := d.Get(key).(*Name); ok {
		if ok && value != nil && value.String() == expected {
			return true
		}
	}
	return false
}