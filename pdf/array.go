package pdf

import "github.com/mawicks/PDFiG/containers"

type ReadOnlyArray interface {
	Object
	Size() int
	At(i int) Object
}

type Array interface {
	ReadOnlyArray
	Add(o Object)
	PushFront(o Object)
	Append(op Array)
}

type array struct {
	array containers.ArrayStack
}

// Constructor for standard implementation of Array
func NewArray() Array {
	return &array{containers.StackArrayDecorator{containers.NewDynamicArray(4)}}
}

func (a *array) Clone() Object {
	newArray := NewArray().(*array)
	size := a.Size()
	for i := 0; i<size; i++ {
		o := a.At(i)
		newArray.array.PushBack(o.Clone())
	}
	return newArray
}

func (a *array) Dereference() Object {
	return a
}

func (a *array) Add(o Object) {
	a.array.PushBack(o)
}

func (a *array) PushFront(o Object) {
	a.array.PushFront(o)
}

func (a *array) Append(op Array) {
	for i:=0; i<op.Size(); i++ {
		a.Add(op.At(i))
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
