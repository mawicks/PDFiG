package pdf

// Implements the pdf.Object interface

type ProtectedDictionary interface {
	Object	// A Dictionary must implement the Object interface.

	// Objects retrieved from a protected or an unprotected array
	// implemented the protected interface for the object type.
	// Objects retrieved from an unprotected array are guaranteed
	// to be unprotected and may be cast to the unprotected
	// versions of their respective interfaces.

	// Get() returns the object stored under the specified key.
	// or nil if the key doesn't exist.
	Get(key string) Object

	// GetArray() attempts to retrieve the dictionary entry as an
	// Array, dereferencing as necessary.  It returns nil if the
	// entry doesn't exist or isn't a ProtectedArray
	// ProtectedArray.  ProtectedArray is a subset of the Array
	// interface.  The returned ProtectedArray may or may not
	// implement the full Array interface.
	GetArray(key string) ProtectedArray

	// GetBoolean() attempts to retrieve the dictionary entry as a
	// Boolean, dereferencing as necessary.  The second boolean
	// return value indicates whether a boolean entry exists.  If
	// so, the first return value is the value of the Boolean.
	GetBoolean(key string) (bool,bool)

	// GetDictionary() attempts to retrieve the dictionary entry
	// as a Dictionary, dereferencing as necessary.  It returns
	// nil if the entry doesn't exist or isn't a Protected
	// Dictionary.  ProtectedDictionary is a subset of the
	// Dictionary interface. The returned ProtectedDictionary may
	// or may not implement the full Dictionary interface.
	GetDictionary(key string) ProtectedDictionary

	// GetIndirect() attempts to retrieve the dictionary entry as
	// an Indirect object (no dereferencing is attempted).  The
	// return value is nil if the entry doesn't exist or isn't
	// a ProtectedIndirect.  ProtectedIndirect is a
	// subset of the Indirect interface.  The returned
	// ProtectedIndirect may or may not implement the full Indirect
	// interface.
	GetIndirect(key string) ProtectedIndirect

	// GetInt() attempts to retrieve the dictionary entry as an
	// integer, dereferencing as necessary.  The boolean returns
	// value indicates whether the entry exists with the expected
	// name and type.
	GetInt(key string) (int,bool)

	// GetName() attempts to retrieve the dictionary entry as a
	// Name, dereferencing as necessary.  The string value of the
	// name is returned.  The boolean return value indicates
	// whether the entry exists with the expected name and type.
	GetName(key string) (string,bool)

	// GetReal() attempts to retrieve the dictionary entry as a real,
	// dereferencing as necessary.  The boolean returns value indicates
	// whether the entry exists with the expected name and type.
	GetReal(key string) (float32,bool)

	// GetStream() attempts to retrieve the dictionary entry as a
	// ProtectedStream, dereferencing as necessary.  It returns
	// nil if the entry doesn't exist or isn't a ProtectedStream.
	// ProtectedStream is a subset of the Stream interface.  The
	// return value may or may not implement the full Stream
	// interface.
	GetStream(key string) ProtectedStream

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
	ProtectedDictionary
	// Add() stores an object under the specified key.
	// Note that all objects Add()ed to a Dictionary are
	// Unprotected() before adding them so that owners of the
	// Dictionary also own objects within it.
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

func (d *dictionary) Protected() Object {
	return protectedDictionary{d}
}

func (d *dictionary) Unprotected() Object {
	return d
}

func (d *dictionary) Add(key string, o Object) {
	d.dictionary[key] = o.Unprotected()
}

func (d *dictionary) Get(key string) Object {
	return d.dictionary[key]
}

// Following routines are convenience functions where the type of the
// dictionary is known to be one of a small number of types.

// GetArray() attempts to retrieve the dictionary entry as an Array,
// dereferencing as necessary.  The return value is nil if the key
// doesn't exist or if the value is not compatible with
// ProtectedArray.  If the return value is non nil, then the returned
// ProtectedArray can safely be casted up to Array.
func (d *dictionary) GetArray(key string) ProtectedArray {
	value := d.Get(key)
	if value == nil {
		return nil
	}
	if array,ok := value.Dereference().(ProtectedArray); ok {
		return array
	}
	return nil
}

// GetBoolean() attempts to retrieve the dictionary entry as Boolean,
// dereferencing as necessary.  The second boolean return value
// indicates whether the entry exists with the expected name and type.
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
// Dictionary, dereferencing as necessary.  The boolean return value
// indicates whether the entry exists with the expected name and type.
// The returned ProtectedDictionary may or may not implement the
// full Dictionary interface.
func (d *dictionary) GetDictionary(key string) ProtectedDictionary {
	value := d.Get(key)
	if value == nil {
		return nil
	}
	if dictionary,ok := value.Dereference().(ProtectedDictionary); ok {
		return dictionary
	}
	return nil
}

// GetIndirect() attempts to retrieve the dictionary entry as a an
// Indirect object (no dereferencing is attempted).  The retun value is nil
// if the key doesn't exist or the value isn't compatible with Indirect.
// The returned ProtectedIndirect may safely be cast to Indirect.
func (d *dictionary) GetIndirect(key string) ProtectedIndirect {
	value := d.Get(key)
	if value == nil {
		return nil
	}
	if indirect,ok := value.(ProtectedIndirect); ok {
		return indirect
	}
	return nil
}

// GetInt() attempts to retrieve the dictionary entry as a integer,
// dereferencing as necessary.  The boolean return value indicates
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
// returned.  The boolean return value indicates whether the entry
// exists with the expected name and type.
func (d *dictionary) GetName(key string) (string,bool) {
	value := d.Get(key)
	if value == nil {
		return "", false
	}
	if name,ok := value.Dereference().(Name); ok {
		return name.String(),true
	}
	return "", false
}

// GetReal() attempts to retrieve the dictionary entry as a real,
// dereferencing as necessary.  The boolean return value indicates
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
// The boolean return value indicates whether the entry exists with
// the expected name and type.
func (d *dictionary) GetString(key string) ([]byte,bool) {
	value := d.Get(key)
	if value == nil {
		return nil, false
	}
	if stringObject,ok := value.Dereference().(ProtectedString); ok {
		return stringObject.Bytes(),true
	}
	return nil, false
}

// GetStream() attempts to retrieve the dictionary entry as a Stream,
// dereferencing as necessary.  The return value is nil if the key
// doesn't exist or if the value isn't compatible with Stream.  The
// returned ProtectedStream may safely be cast to Stream.
func (d *dictionary) GetStream(key string) ProtectedStream {
	value := d.Get(key)
	if value == nil {
		return nil
	}
	if stream,ok := value.Dereference().(ProtectedStream); ok {
		return stream
	}
	return nil
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

	if value,ok := rawValue.Dereference().(Name); ok {
		if ok && value != nil && value.String() == expected {
			return true
		}
	}
	return false
}

type protectedDictionary struct {
	d Dictionary
}

// The value returned by Dereference() can safely be cast to
// Dictionary.
func (rod protectedDictionary) Clone() Object {
	return rod.d.Clone()
}

// The value returned by Dereference() can safely be cast to
// ProtectedDictionary.
func (rod protectedDictionary) Dereference() Object {
	return rod
}

// The value returned by Protected() can safely be cast to
// ProtectedDictionary.
func (rod protectedDictionary) Protected() Object {
	return rod
}

// The value returned by Unprotected() can safely be cast to
// ProtectedDictionary or Dictionary.
func (rod protectedDictionary) Unprotected() Object {
	return rod.d.Clone()
}

func (rod protectedDictionary) Serialize(w Writer, file... File) {
	rod.d.Serialize(w, file...)
}

func (rod protectedDictionary) Get(key string) Object {
	return rod.d.Get(key).Protected()
}

func (rod protectedDictionary) GetArray(key string) ProtectedArray {
	a := rod.d.GetArray(key)
	if  a != nil {
		return a.Protected().(ProtectedArray)
	} else {
		return nil
	}
}

func (rod protectedDictionary) GetBoolean(key string) (bool,bool) {
	return rod.d.GetBoolean(key)
}

func (rod protectedDictionary) GetDictionary(key string) ProtectedDictionary {
	d := rod.d.GetDictionary(key)
	if d != nil {
		return d.Protected().(ProtectedDictionary)
	} else {
		return nil
	}
}

func (rod protectedDictionary) GetIndirect(key string) ProtectedIndirect {
	i := rod.d.GetIndirect(key)
	if i != nil {
		return i.Protected().(ProtectedIndirect)
	} else {
		return nil
	}
}

func (rod protectedDictionary) GetInt(key string) (int,bool) {
	return rod.d.GetInt(key)
}

func (rod protectedDictionary) GetName(key string) (string,bool) {
	return rod.d.GetName(key)
}

func (rod protectedDictionary) GetReal(key string) (float32,bool) {
	return rod.d.GetReal(key)
}

func (rod protectedDictionary) GetStream(key string) ProtectedStream {
	s := rod.d.GetStream(key)
	if s != nil {
		return s.Protected().(ProtectedStream)
	} else {
		return nil
	}
}

func (rod protectedDictionary) GetString(key string) ([]byte,bool) {
	return rod.d.GetString(key)
}

func (rod protectedDictionary) Size() int {
	return rod.d.Size()
}

func (rod protectedDictionary) CheckNameValue(key string, expected string, file... File) bool {
	return rod.d.CheckNameValue(key,expected,file...)
}
