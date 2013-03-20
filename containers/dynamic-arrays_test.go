package containers

import "fmt"
import "testing"

func TestNothing(t *testing.T) {
	var i interface{} = nil
	var j interface{} = 0
	var k *interface{} = nil
	var l interface{}
	fmt.Printf ("i=%v; j=%v; k=%v; l=%v\n", i, j, k, l)
}

func TestDynamicArray (t *testing.T) {
	d := NewDynamicArray(10)
	d.SetSize(100)

	*(d.At(0)) = 1
	fmt.Printf ("d.At(0) = %v\n", *(d.At(0)))

	*(d.At(1)) = 2
	fmt.Printf ("d.At(1) = %v\n", *(d.At(1)))

	*(d.At(9)) = 10
	fmt.Printf ("d.At(9) = %v\n", *(d.At(9)))

	*(d.At(10)) = 11
	fmt.Printf ("d.At(10) = %v\n", *(d.At(10)))

	*(d.At(99)) = 100
	fmt.Printf ("d.At(99) = %v\n", *(d.At(99)))
}