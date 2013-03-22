package containers

import "testing"

func testRetrieveSlot (t *testing.T, s string, da *DynamicArray, i uint, v interface{}) {
	if *(da.At(i)) != v {
		t.Errorf ("%s: *(da.At(%d)) == %v; expected %v", s, i, *(da.At(i)), v)
	}
}

func fillArray (da *DynamicArray, size int) {
	for i:=0; i<size; i++ {
		*da.At(uint(i)) = i+1
	}
}

func checkArrayForFill (t *testing.T, s string, da *DynamicArray, max int) {
	for i:=0; i<max; i++ {
		testRetrieveSlot (t, s, da, uint(i), i+1)
	}
}

func checkArrayForNull (t *testing.T, s string, da *DynamicArray, start, size int) {
	for i:=start; i<size; i++ {
		testRetrieveSlot (t, s, da, uint(i), nil)
	}
}

func TestDynamicArray (t *testing.T) {
	var d *DynamicArray

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

	d = NewDynamicArray (3)
	for i:=0; i<10; i++ {
		d.PushBack(i+1);
	}

	for i:=10; i>0; i-- {
		v := d.PopBack()
		if v != i {
			t.Errorf ("PushBack/PopBack: PopBack() == %v; expected %v", v, i)
		}
	}
}

