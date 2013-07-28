package pdf

import (
	"github.com/mawicks/PDFiG/containers"
)

type ProtectedArray interface {
	Object
	Size() int
	At(i int) Object
}

type Array interface {
	ProtectedArray
	// Note that all objects added via Add(), PushFront(), or
	// Append() are Unprotect() before adding them so that owners
	// of the Array also own objects within it.
	Add(o Object)
	PushFront(o Object)
	Append(op ProtectedArray)
}

type array struct {
	array containers.ArrayStack
}

// Constructor for standard implementation of Array
func NewArray() Array {
	return &array{containers.StackArrayDecorator{containers.NewDynamicArray(4)}}
}

// Return value of Clone() can safely be cast to Array.
func (a *array) Clone() Object {
	newArray := NewArray().(*array)
	size := a.Size()
	for i := 0; i<size; i++ {
		newArray.array.PushBack(a.At(i).Clone())
	}
	return newArray
}

// Return value of Dereference can safely be cast to Array.
func (a *array) Dereference() Object {
	return a
}

// Return value of Protect() can safely be cast to ProtectedArray but
// not to Array.
func (a *array) Protect() Object {
	return protectedArray{a}
}

// Return value of Unprotect() can safely be cast to Array.
func (a *array) Unprotect() Object {
	return a
}

func (a *array) Size() int {
	return int(a.array.Size())
}

// The reference returned by At() retains the protection of the object
// being referenced.
func (a *array) At(i int) Object {
	return (*a.array.At(uint(i))).(Object)
}

func (a *array) Add(o Object) {
	a.array.PushBack(o.Unprotect())
}

func (a *array) PushFront(o Object) {
	a.array.PushFront(o.Unprotect())
}

func (a *array) Append(op ProtectedArray) {
	for i:=0; i<op.Size(); i++ {
		a.Add(op.At(i).Unprotect())
	}
}

func (a *array) Serialize(w Writer, file ...File) {
	w.WriteByte('[')
	size := a.Size()
	for i := 0; i < size; i++ {
		if i != 0 {
			w.WriteByte(' ')
		}
		o := a.At(i)
		o.Serialize(w, file...)
	}
	w.WriteByte(']')
}

type protectedArray struct {
	a Array
}

// Return value of Clone() can safely be cast to Array
func (pa protectedArray) Clone() Object {
	return pa.a.Clone()
}

// Return value of Dereferene() can safely be cast to ProtectedArray.
func (pa protectedArray) Dereference() Object {
	return pa
}

// Return value of Protect() can safely be cast to ProtectedArray.
func (pa protectedArray) Protect() Object {
	return pa
}

// Return value of Unprotect() can safely be cast to Array.
func (pa protectedArray) Unprotect() Object {
	newArray := NewArray().(*array)
	size := pa.a.Size()
	for i := 0; i<size; i++ {
		// Note that pa.At(i) is protected, so contained
		// objects remain protected.
		newArray.array.PushBack(pa.At(i))
	}
	return newArray
}

func (pa protectedArray) Size() int {
	return pa.a.Size()
}

// Returned value of At() is always the protected version of the
// object's interface.
func (pa protectedArray) At(i int) Object {
	return pa.a.At(i).Protect()
}

func (pa protectedArray) Serialize(w Writer, file ...File) {
	pa.a.Serialize(w, file...)
}
