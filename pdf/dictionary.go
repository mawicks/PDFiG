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

func (d *Dictionary) Dereference() Object {
	return d
}

func (d *Dictionary) Add(key string, o Object) {
	d.dictionary[key] = o
}

func (d *Dictionary) Get(key string) Object {
	return d.dictionary[key]
}

// Following routines are convenience functions where the type of the
// dictionary is known to be one of a small number of types.

// GetArray() attempts to retrieve the dictionary entry as an Array,
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *Dictionary) GetArray(key string) (*Array,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if array,ok := value.Dereference().(*Array); ok {
		return array,true
	}
	return nil, false
}

// GetArray() attempts to retrieve the dictionary entry as Boolean
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *Dictionary) GetBoolean(key string) (bool,bool) {
	value := d.Get(key)
	if value == nil {
		return false, false
	}
	if _,ok := value.Dereference().(TrueBoolean); ok {
		return true,true
	}
	if _,ok := value.Dereference().(FalseBoolean); ok {
		return false,true
	}
	return false,false
}

// GetArray() attempts to retrieve the dictionary entry as a
// Dictionary dereferencing as necessary.  The boolean returns value
// indicates whether the entry exists with the expected name and type.
func (d *Dictionary) GetDictionary(key string) (*Dictionary,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if dictionary,ok := value.Dereference().(*Dictionary); ok {
		return dictionary,true
	}
	return nil, false
}

// GetArray() attempts to retrieve the dictionary entry as a an
// Indirect object (no dereferencing is attempted).  The boolean
// returns value indicates whether the entry exists with the expected
// name and type.
func (d *Dictionary) GetIndirect(key string) (*Indirect,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if indirect,ok := value.(*Indirect); ok {
		return indirect,true
	}
	return nil, false
}

// GetArray() attempts to retrieve the dictionary entry as a integer,
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *Dictionary) GetInt(key string) (int,bool) {
	value := d.Get(key)
	if value == nil {
		return 0, false
	}
	if intNumeric,ok := value.Dereference().(*IntNumeric); ok {
		return intNumeric.Value(),true
	}
	return 0, false
}

// GetArray() attempts to retrieve the dictionary entry as a Name,
// dereferencing as necessary.  The string value of the name is
// returned.  The boolean returns value indicates whether the entry
// exists with the expected name and type.
func (d *Dictionary) GetName(key string) (string,bool) {
	value := d.Get(key)
	if value == nil {
		return "", false
	}
	if name,ok := value.Dereference().(*Name); ok {
		return name.String(),true
	}
	return "", false
}

// GetArray() attempts to retrieve the dictionary entry as a real,
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *Dictionary) GetReal(key string) (float32,bool) {
	value := d.Get(key)
	if value == nil {
		return 0.0, false
	}
	if realNumeric,ok := value.Dereference().(*RealNumeric); ok {
		return realNumeric.Value(),true
	}
	return 0, false
}

// GetArray() attempts to retrieve the dictionary entry as a string,
// dereferencing as necessary.  The raw byte sequence is returned.
// The boolean returns value indicates whether the entry exists with
// the expected name and type.
func (d *Dictionary) GetString(key string) ([]byte,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if stringObject,ok := value.Dereference().(*String); ok {
		return stringObject.Bytes(),true
	}
	return nil, false
}

// GetArray() attempts to retrieve the dictionary entry as a Stream,
// dereferencing as necessary.
func (d *Dictionary) GetStream(key string) (*Stream,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if stream,ok := value.Dereference().(*Stream); ok {
		return stream,true
	}
	return nil, false
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

func (d *Dictionary) CheckNameValue (key string, expected string, file... File) bool {
	var rawValue Object

	if rawValue = d.Get(key); rawValue == nil {
		return false
	}

	if value,ok := rawValue.Dereference().(*Name); ok {
		if ok && value != nil && value.String() == expected {
			return true
		}
	}
	return false
}