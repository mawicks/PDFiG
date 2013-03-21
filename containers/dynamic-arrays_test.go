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

func testStoreAndRetrieveSlot (t *testing.T, da *DynamicArray, i uint, v interface{}) {
	*(da.At(i)) = v
	if *(da.At(i)) != v {
		t.Errorf ("*(da.At(%d))==%v after *(da.At(%d)) = %v", i, *(da.At(i)), i, v)
	}
}

func testRetrieveSlot (t *testing.T, s string, da *DynamicArray, i uint, v interface{}) {
	if *(da.At(i)) != v {
		t.Errorf ("%s: *(da.At(%d)) == %v; expected %v", s, i, *(da.At(i)), v)
	}
}

func TestDynamicArray (t *testing.T) {
	d := NewDynamicArray(2)
	d.SetSize(9)
	// Structure is as follows: 
	// Layer 1  0=========  1=========  2=========  3=========
	// Layer 2  0  1  2  3  4  5  6  7  8  x  x  x  x  x  x  x

	var i uint
	for i=0; i<9; i++ {
		testStoreAndRetrieveSlot (t, d, i, i+1)
	}

	for i=0; i<9; i++ {
		testRetrieveSlot (t, "test1", d, i, i+1)
	}

	d.SetSize(5)
	d.SetSize(4)
	d.SetSize(3)
	d.SetSize(20)
	for i=0; i<3; i++ {
		testRetrieveSlot (t, "test2", d, i, i+1)
	}
	for i=3; i<20; i++ {
		testRetrieveSlot (t, "test3", d, i, nil)
	}
	d.SetSize(0)
	d.SetSize(32)
	for i=0; i<32; i++ {
		testRetrieveSlot (t, "test4", d, i, nil)
	}
	d.SetSize(0)
}

