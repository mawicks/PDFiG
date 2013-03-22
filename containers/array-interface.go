/*
	Package of containers.
*/
package containers

type Array interface {
	// Return element at specified position
	At (uint) interface {}	
	// Add element at end, growing the array by one
	PushBack (interface{})	
	// Return last element after shrinking the array by one
	PopBack () interface {}
	// Set the size
	SetSize(uint)
	// Return size (number of elements)
	Size() uint		
}
