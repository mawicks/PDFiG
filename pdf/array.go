/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "io"
import "maw/containers"

type Array struct {
	array containers.Array
}

// Constructor for Name object
func NewArray () (* Array) {
	return &Array{containers.NewDynamicArray(4)}
}

func (a *Array) Add (o Object) {
	a.array.PushBack(o)
}

func (a *Array) Serialize (f io.Writer) {
	f.Write([]byte{'['})
	for i:=0; i<int(a.array.Size()); i++ {
		if i != 0 {
			f.Write([]byte{' '})
		}
		o := (*a.array.At(uint(i))).(Object)
		o.Serialize(f)
	}
	f.Write([]byte{']'})
}
