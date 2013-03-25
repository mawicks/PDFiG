/*
	Package for creating, reading, and editing PDF files.
*/
package pdf

import "bufio"
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

func (a *Array) Serialize (f *bufio.Writer, file... File) {
	f.WriteByte('[')
	for i:=0; i<int(a.array.Size()); i++ {
		if i != 0 {
			f.WriteByte(' ')
		}
		o := (*a.array.At(uint(i))).(Object)
		o.Serialize(f)
	}
	f.WriteByte(']')
}
