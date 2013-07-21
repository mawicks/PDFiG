package pdf

import "github.com/mawicks/PDFiG/containers"

type ProtectedArray interface {
	Object
	Size() int
	// The object returned by At() is protected if and only if 
	// array is protected.  In other words, an object
	// retrieved from an unprotected array is guaranteed not
	// to be protected.
	At(i int) Object
}

type Array interface {
	ProtectedArray
	// Note that all objects added via Add(), PushFront(), or
	// Append() are Unprotected() before adding them so that owners
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
		o := a.At(i)
		newArray.array.PushBack(o.Clone())
	}
	return newArray
}

// Return value of Clone() can safely be cast to Array.
func (a *array) Dereference() Object {
	return a
}

// Return value of Protected() can safely be cast to ProtectedArray but not to
// Array.
func (a *array) Protected() Object {
	return protectedArray{a}
}

// Return value of Unprotected() can safely be cast to ProtectedArray
// or Array.
func (a *array) Unprotected() Object {
	return a
}

func (a *array) Add(o Object) {
	a.array.PushBack(o.Unprotected())
}

func (a *array) PushFront(o Object) {
	a.array.PushFront(o.Unprotected())
}

func (a *array) Append(op ProtectedArray) {
	for i:=0; i<op.Size(); i++ {
		a.Add(op.At(i).Unprotected())
	}
}

func (a *array) Size() int {
	return int(a.array.Size())
}

func (a *array) At(i int) Object {
	return (*a.array.At(uint(i))).(Object)
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

// Return value of Clone() can safely be cast to ProtectedArray.
func (roa protectedArray) Clone() Object {
	return roa
}

// Return value of Dereferene() can safely be cast to ProtectedArray.
func (roa protectedArray) Dereference() Object {
	return roa
}

// Return value of Protected() can safely be cast to ProtectedArray.
func (roa protectedArray) Protected() Object {
	return roa
}

// Return value of Protected() can safely be cast to ProtectedArray.
func (roa protectedArray) Unprotected() Object {
	return roa.a.Clone()
}

func (roa protectedArray) Size() int {
	return roa.a.Size()
}

func (roa protectedArray) At(i int) Object {
	return roa.a.At(i).Protected()
}

func (roa protectedArray) Serialize(w Writer, file ...File) {
	roa.a.Serialize(w, file...)
}




















