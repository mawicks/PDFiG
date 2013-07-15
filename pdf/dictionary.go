package pdf

// Implements the pdf.Object interface

type ReadOnlyDictionary interface {
	Object	// A Dictionary must implement the Object interface.

	// Get() returns the object stored under the specified key.
	Get(key string) Object

	// GetArray() attempts to retrieve the dictionary entry as an Array,
	// dereferencing as necessary.  The boolean returns value indicates
	// whether the entry exists with the expected name and type.
	GetArray(key string) (*Array, bool)

	// GetBoolean() attempts to retrieve the dictionary entry as a Boolean
	// dereferencing as necessary.  The boolean second return value indicates
	// whether a boolean entry exists.  If so, the first return value is the
	// value of the Boolean.
	GetBoolean(key string) (bool,bool)

	// GetDictionary() attempts to retrieve the dictionary entry as a
	// Dictionary dereferencing as necessary.  The boolean returns value
	// indicates whether the entry exists with the expected name and type.
	GetDictionary(key string) (ReadOnlyDictionary,bool)

	// GetIndirect() attempts to retrieve the dictionary entry as an
	// Indirect object (no dereferencing is attempted).  The boolean
	// return value indicates whether the entry exists with the expected
	// name and type.
	GetIndirect(key string) (*Indirect,bool)

	// GetInt() attempts to retrieve the dictionary entry as an
	// integer, dereferencing as necessary.  The boolean returns
	// value indicates whether the entry exists with the expected
	// name and type.
	GetInt(key string) (int,bool)

	// GetName() attempts to retrieve the dictionary entry as a Name,
	// dereferencing as necessary.  The string value of the name is
	// returned.  The boolean returns value indicates whether the entry
	// exists with the expected name and type.
	GetName(key string) (string,bool)

	// GetReal() attempts to retrieve the dictionary entry as a real,
	// dereferencing as necessary.  The boolean returns value indicates
	// whether the entry exists with the expected name and type.
	GetReal(key string) (float32,bool)

	// GetStream() attempts to retrieve the dictionary entry as a Stream,
	// dereferencing as necessary.
	GetStream(key string) (Stream,bool)

	// GetString() attempts to retrieve the dictionary entry as a string,
	// dereferencing as necessary.  The raw byte sequence is returned.
	// The boolean return value indicates whether the entry exists with
	// the expected name and type.
	GetString(key string) ([]byte,bool)

	// CheckNameValue() determines whether the value of associated with the
	// is a name corresponding to the expected string (after applying
	// name decoding).
	CheckNameValue (key string, expected string, file... File) bool

	// Size() returns the number of key-value pairs
	Size() int
}

type Dictionary interface {
	ReadOnlyDictionary
	// Add() stores an object under the specified key.
	Add(key string, object Object)
	
	// Remove() removes the key and value stored under the specified key.
	Remove(key string)
}

type dictionary struct {
	dictionary map[string]Object
}

// Constructor for Dictionary object
func NewDictionary() Dictionary {
	return &dictionary{make(map[string]Object, 16)}
	
}

func (d *dictionary) Clone() Object {
	newDictionary := NewDictionary().(*dictionary)
	for key,value := range d.dictionary {
		newDictionary.dictionary[key] = value.Clone()
	}
	return newDictionary
}

func (d *dictionary) Dereference() Object {
	return d
}

func (d *dictionary) Add(key string, o Object) {
	d.dictionary[key] = o
}

func (d *dictionary) Get(key string) Object {
	return d.dictionary[key]
}

// Following routines are convenience functions where the type of the
// dictionary is known to be one of a small number of types.

// GetArray() attempts to retrieve the dictionary entry as an Array,
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *dictionary) GetArray(key string) (*Array,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if array,ok := value.Dereference().(*Array); ok {
		return array,true
	}
	return nil, false
}

// GetBoolean() attempts to retrieve the dictionary entry as Boolean
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *dictionary) GetBoolean(key string) (bool,bool) {
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

// GetDictionary() attempts to retrieve the dictionary entry as a
// Dictionary dereferencing as necessary.  The boolean returns value
// indicates whether the entry exists with the expected name and type.
func (d *dictionary) GetDictionary(key string) (ReadOnlyDictionary,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if dictionary,ok := value.Dereference().(*dictionary); ok {
		return dictionary,true
	}
	return nil, false
}

// GetIndirect() attempts to retrieve the dictionary entry as a an
// Indirect object (no dereferencing is attempted).  The boolean
// returns value indicates whether the entry exists with the expected
// name and type.
func (d *dictionary) GetIndirect(key string) (*Indirect,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if indirect,ok := value.(*Indirect); ok {
		return indirect,true
	}
	return nil, false
}

// GetInt() attempts to retrieve the dictionary entry as a integer,
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *dictionary) GetInt(key string) (int,bool) {
	value := d.Get(key)
	if value == nil {
		return 0, false
	}
	if intNumeric,ok := value.Dereference().(*IntNumeric); ok {
		return intNumeric.Value(),true
	}
	return 0, false
}

// GetName() attempts to retrieve the dictionary entry as a Name,
// dereferencing as necessary.  The string value of the name is
// returned.  The boolean returns value indicates whether the entry
// exists with the expected name and type.
func (d *dictionary) GetName(key string) (string,bool) {
	value := d.Get(key)
	if value == nil {
		return "", false
	}
	if name,ok := value.Dereference().(*Name); ok {
		return name.String(),true
	}
	return "", false
}

// GetReal() attempts to retrieve the dictionary entry as a real,
// dereferencing as necessary.  The boolean returns value indicates
// whether the entry exists with the expected name and type.
func (d *dictionary) GetReal(key string) (float32,bool) {
	value := d.Get(key)
	if value == nil {
		return 0.0, false
	}
	if realNumeric,ok := value.Dereference().(*RealNumeric); ok {
		return realNumeric.Value(),true
	}
	return 0, false
}

// GetString() attempts to retrieve the dictionary entry as a string,
// dereferencing as necessary.  The raw byte sequence is returned.
// The boolean returns value indicates whether the entry exists with
// the expected name and type.
func (d *dictionary) GetString(key string) ([]byte,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if stringObject,ok := value.Dereference().(*String); ok {
		return stringObject.Bytes(),true
	}
	return nil, false
}

// GetStream() attempts to retrieve the dictionary entry as a Stream,
// dereferencing as necessary.
func (d *dictionary) GetStream(key string) (Stream,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if stream,ok := value.Dereference().(Stream); ok {
		return stream,true
	}
	return nil, false
}

func (d *dictionary) Remove(key string) {
	delete(d.dictionary, key)
}

func (d *dictionary) Serialize(w Writer, file ...File) {
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
func (d *dictionary) Size() int {
	return len(d.dictionary)
}

func (d *dictionary) CheckNameValue (key string, expected string, file... File) bool {
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