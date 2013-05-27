package containers

import "testing"

func testRetrieveSlot (t *testing.T, s string, da Array, i uint, v interface{}) {
	if *(da.At(i)) != v {
		t.Errorf ("%s: *(da.At(%d)) == %v; expected %v", s, i, *(da.At(i)), v)
	}
}

func fillArray (da Array, size int) {
	for i:=0; i<size; i++ {
		*da.At(uint(i)) = i+1
	}
}

func checkArrayForFill (t *testing.T, s string, da Array, max int) {
	for i:=0; i<max; i++ {
		testRetrieveSlot (t, s, da, uint(i), i+1)
	}
}

func checkArrayForNull (t *testing.T, s string, da Array, start, size int) {
	for i:=start; i<size; i++ {
		testRetrieveSlot (t, s, da, uint(i), nil)
	}
}

func TestDynamicArray (t *testing.T) {
	var d Array
	
	for csize:= 2; csize<6; csize++ {
		d = NewDynamicArray(uint(csize))

		for size:=64; size>=0; size-- {
			d.SetSize(64)
			fillArray (d, 64)

			d.SetSize(uint(size))
			checkArrayForFill (t, "Shrink doesn't preserve values", d, size)

			d.SetSize(64)
			checkArrayForFill (t, "Expand doesn't preserve values", d, size)
			checkArrayForNull (t, "Shrink didn't clear unused values", d, size, 64)
		}
	}

	s := &StackArrayDecorator{NewDynamicArray (3)}
	for i:=0; i<10; i++ {
		s.PushBack(i+1);
	}

	for i:=10; i>0; i-- {
		v := s.PopBack()
		if v != i {
			t.Errorf ("PushBack/PopBack: PopBack() == %v; expected %v", v, i)
		}
	}

	s = &StackArrayDecorator{NewDynamicArray (3)}
	for i:=0; i<10; i++ {
		s.PushFront(i+1)
	}
	
	for i:=10; i>0; i-- {
		v := s.PopFront()
		if v != i {
			t.Errorf ("PushFront/PopFront: PopBack() = %v; expected %v", v, i)
		}
	}
}

