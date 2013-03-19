/*
	Package of containers.
*/
package containers

// All PDF objects implement the pdf.Object inteface
type Array interface {
	// Return element at specified position
	At (uint) interface {}	
	// Add element at end, growing the array by one
	PushBack (interface{})	
	// Return last element after shrinking the array by one
	PopBack () interface {}
	// Return size (number of elements)
	Size() uint		
	// Set the size
	SetSize(uint)
}

type DynamicArray struct {
	// Cluster size
	clusterSize uint

	// Total capacity (without a reallocation) Capacity is always
	// a positive integral power of clusterSize or it's 0.
	capacity uint	

	// Number of elements stored (<= capacity)
	size uint

	// If this is a leaf node, "array" contains user-stored elements.
	// If this is not a leaf node, "array" is an array of DynamicArray
	// each of size capacity/clusterSize.
	// Nodes are leaf nodes if and only if capacity == clusterSize
	array []interface{}
}

func NewDynamicArray (clusterSize uint) *DynamicArray {
	return &DynamicArray{clusterSize, clusterSize, 0, make([]interface{}, clusterSize)}
}

func recursiveAt (i uint, clusterSize uint, capacity uint, array []interface{}) (result interface{}) {
	if capacity <= clusterSize {
		result = array[i]
	} else {
		nextCapacity := capacity / clusterSize
		result = recursiveAt (i % nextCapacity, clusterSize, nextCapacity, array[i/nextCapacity].([]interface{}))
	}
	return result
}

func (da *DynamicArray) shrinkOrGrow (newSize uint) {
	// Shrink...
	for newCap := da.capacity/da.clusterSize; newCap >= newSize && newCap > 1; newCap /= da.clusterSize {
		da.array = da.array[0].([]interface{})
		da.capacity = newCap
	}
	// or grow
	for ; da.capacity < newSize; da.capacity *= da.clusterSize  {
		newArray := make([]interface{}, da.clusterSize)
		newArray[0] = da.array
		da.array = newArray
	}
}

func (da *DynamicArray) At (i uint) interface{} {
	return recursiveAt (i, da.clusterSize, da.capacity, da.array)
}

func (da *DynamicArray) Size() uint {
	return da.size
}

func (da *DynamicArray) SetSize (newSize uint) {
	da.shrinkOrGrow (newSize)
}