package pdf

import "maw/containers"

type Array struct {
	array containers.Array
}

// Constructor for Name object
func NewArray() *Array {
	return &Array{containers.NewDynamicArray(4)}
}

func (a *Array) Add(o Object) {
	a.array.PushBack(o)
}

func (a *Array) Serialize(w Writer, file ...File) {
	w.WriteByte('[')
	for i := 0; i < int(a.array.Size()); i++ {
		if i != 0 {
			w.WriteByte(' ')
		}
		o := (*a.array.At(uint(i))).(Object)
		o.Serialize(w)
	}
	w.WriteByte(']')
}
