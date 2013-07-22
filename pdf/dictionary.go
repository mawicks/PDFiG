package pdf

// Implements the pdf.Object interface

type ProtectedDictionary interface {
	Object	// A Dictionary must implement the Object interface.

	// Get() returns the object stored under the specified key.
	// or nil if the key doesn't exist.
	Get(key string) Object

	// Following methods are for convenience when the client
	// expects the dictionary entry to belong to one of a small
	// number of types.  They also provide dereferencing of
	// indirect objects and return the indicated type rather than
	// the indirect reference.

	// GetArray() attempts to retrieve the dictionary entry as an
	// ProtectedArray, dereferencing as necessary.  It returns nil
	// if the entry doesn't exist or isn't a ProtectedArray.
	// ProtectedArray is a subset of the Array interface.  The
	// returned ProtectedArray may or may not implement the full
	// Array interface.
	GetArray(key string) ProtectedArray

	// GetBoolean() attempts to retrieve the dictionary entry as a
	// Boolean, dereferencing as necessary.  The second boolean
	// return value indicates whether the entry exists and is a
	// boolean.  If so, the first return value is the value of the
	// Boolean.
	GetBoolean(key string) (bool,bool)

	// GetDictionary() attempts to retrieve the dictionary entry
	// as a ProtectedDictionary, dereferencing as necessary.  It
	// returns nil if the entry doesn't exist or isn't a
	// ProtectedDictionary.  ProtectedDictionary is a subset of
	// the Dictionary interface. The returned ProtectedDictionary
	// may or may not implement the full Dictionary interface.
	GetDictionary(key string) ProtectedDictionary

	// GetIndirect() attempts to retrieve the dictionary entry as
	// a ProtectedIndirect object (no dereferencing is attempted).
	// The return value is nil if the entry doesn't exist or isn't
	// a ProtectedIndirect.  ProtectedIndirect is a subset of the
	// Indirect interface.  The returned ProtectedIndirect may or
	// may not implement the full Indirect interface.
	GetIndirect(key string) ProtectedIndirect

	// GetInt() attempts to retrieve the dictionary entry as an
	// integer, dereferencing as necessary.  The boolean return
	// value indicates whether the entry exists and is an Int.  If
	// so, the returned int represents the value of the integer.
	GetInt(key string) (int,bool)

	// GetName() attempts to retrieve the dictionary entry as a
	// Name, dereferencing as necessary.  The string value of the
	// name is returned.  The boolean return value indicates
	// whether the entry exists and is a name.  If so, the
	// returned string is the name, expressed as a string.
	GetName(key string) (string,bool)

	// GetReal() attempts to retrieve the dictionary entry as a
	// real, dereferencing as necessary.  The boolean return value
	// indicates whether the entry exists and is a Real.  If so,
	// the returned float32 is the value of the Real.
	GetReal(key string) (float32,bool)

	// GetStream() attempts to retrieve the dictionary entry as a
	// ProtectedStream, dereferencing as necessary.  It returns
	// nil if the entry doesn't exist or isn't a ProtectedStream.
	// ProtectedStream is a subset of the Stream interface.  The
	// return value may or may not implement the full Stream
	// interface.
	GetStream(key string) ProtectedStream

	// GetString() attempts to retrieve the dictionary entry as a
	// string, dereferencing as necessary.  The boolean return
	// value indicates whether the entry exists and is a String.
	// If so, the returned slice represents the byte sequence
	// represented by the string.
	GetString(key string) ([]byte,bool)

	// CheckNameValue() returns true if the value associated with
	// the key is a name represented by the expected string (after
	// applying name decoding).
	CheckNameValue (key string, expected string, file... File) bool

	// Size() returns the number of key-value pairs
	Size() int

	// Keys() returns a slice of strings representing the 
	// names in the dictionary
	Keys() []string
}

type Dictionary interface {
	ProtectedDictionary
	// Add() stores an object under the specified key.
	// The protection of the added Object is preserved.
	Add(key string, object Object)
	
	// Remove() removes the key and value stored under the
	// specified key.
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

func (d *dictionary) Protect() Object {
	return protectedDictionary{d}
}

func (d *dictionary) Unprotect() Object {
	return d
}

func (d *dictionary) Add(key string, o Object) {
	d.dictionary[key] = o
}

func (d *dictionary) Get(key string) Object {
	return d.dictionary[key]
}

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
	if b,ok := value.Dereference().(Boolean); ok {
		return b.Value(),true
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
	if stringObject,ok := value.Dereference().(ProtectString); ok {
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

func (d *dictionary) Keys() []string {
	keys := make([]string, len(d.dictionary))
	var i int = 0
	for k,_ := range d.dictionary {
		keys[i] = k
		i += 1
	}
	return keys;
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

// The value returned by Clone() can safely be cast to Dictionary.
func (pd protectedDictionary) Clone() Object {
	return pd.d.Clone()
}

// The value returned by Dereference() can safely be cast to
// ProtectedDictionary.
func (pd protectedDictionary) Dereference() Object {
	return pd
}

// The value returned by Protect() can safely be cast to
// ProtectedDictionary.
func (pd protectedDictionary) Protect() Object {
	return pd
}

// The value returned by Unprotect() can safely be cast to Dictionary.
func (pd protectedDictionary) Unprotect() Object {
	newDictionary := NewDictionary().(*dictionary)
	for _,key := range pd.d.Keys() {
		newDictionary.dictionary[key] = pd.d.Get(key).Protect()
	}
	return newDictionary
}

func (pd protectedDictionary) Serialize(w Writer, file... File) {
	pd.d.Serialize(w, file...)
}

func (pd protectedDictionary) Get(key string) Object {
	return pd.d.Get(key).Protect()
}

func (pd protectedDictionary) GetArray(key string) ProtectedArray {
	if a := pd.d.GetArray(key); a != nil {
		return a.Protect().(ProtectedArray)
	}
	return nil
}

func (pd protectedDictionary) GetBoolean(key string) (bool,bool) {
	return pd.d.GetBoolean(key)
}

func (pd protectedDictionary) GetDictionary(key string) ProtectedDictionary {
	if d := pd.d.GetDictionary(key); d != nil {
		return d.Protect().(ProtectedDictionary)
	}
	return nil
}

func (pd protectedDictionary) GetIndirect(key string) ProtectedIndirect {
	if i := pd.d.GetIndirect(key); i != nil {
		return i.Protect().(ProtectedIndirect)
	}
	return nil
}

func (pd protectedDictionary) GetInt(key string) (int,bool) {
	return pd.d.GetInt(key)
}

func (pd protectedDictionary) GetName(key string) (string,bool) {
	return pd.d.GetName(key)
}

func (pd protectedDictionary) GetReal(key string) (float32,bool) {
	return pd.d.GetReal(key)
}

func (pd protectedDictionary) GetStream(key string) ProtectedStream {
	if s := pd.d.GetStream(key); s != nil {
		return s.Protect().(ProtectedStream)
	}
	return nil
}

func (pd protectedDictionary) GetString(key string) ([]byte,bool) {
	return pd.d.GetString(key)
}

func (pd protectedDictionary) Size() int {
	return pd.d.Size()
}

func (pd protectedDictionary) Keys() []string {
	return pd.d.Keys()
}

func (pd protectedDictionary) CheckNameValue(key string, expected string, file... File) bool {
	return pd.d.CheckNameValue(key,expected,file...)
}
