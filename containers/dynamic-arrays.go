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
	// a positive integral power of clusterSize
	capacity uint	

	// Number of elements stored (<= capacity)
	size uint

	// If this is a leaf node, "array" contains user-stored elements.
	// If this is not a leaf node, "array" is an array of DynamicArray
	// each of size capacity/clusterSize.
	// Nodes are leaf nodes if and only if capacity == clusterSize
	tree []interface{}
}

func NewDynamicArray (clusterSize uint) *DynamicArray {
	if clusterSize <= 1 {
		panic ("NewDynamicArray(): clusterSize too small")
	}
	return &DynamicArray{clusterSize, clusterSize, 0, make([]interface{}, clusterSize)}
}

func recursiveAt (i uint, tree []interface{}, capacity uint, clusterSize uint) (result *interface{}) {
	if capacity <= clusterSize {
		result = &tree[i]
	} else {
		subtreeCapacity := capacity / clusterSize
		whichSubtree := i/subtreeCapacity
		// Create a subtree if necessary
		if tree[whichSubtree] == nil {
			tree[whichSubtree] =  make([]interface{}, clusterSize)
		}
		result = recursiveAt (i%subtreeCapacity, tree[whichSubtree].([]interface{}), subtreeCapacity, clusterSize)
	}
	return result
}

func recursiveRelease (size uint, tree[]interface{}, capacity uint, clusterSize uint) {
	subtreeCapacity := capacity / clusterSize
	lastPartialSubtree := size/subtreeCapacity
	for i:=lastPartialSubtree+1; i<clusterSize; i++ {
		tree[i] = nil
	}
	if subtreeCapacity > 1 && lastPartialSubtree < clusterSize && tree[lastPartialSubtree] != nil {
		recursiveRelease (size%subtreeCapacity, tree[lastPartialSubtree].([]interface{}), subtreeCapacity, clusterSize)
	}
}

func (da *DynamicArray) shrinkOrGrow (newSize uint) {
	shrunk := false
	// Shrink...
	for newCap := da.capacity/da.clusterSize; newCap >= newSize && newCap > 1; newCap /= da.clusterSize {
		da.tree = da.tree[0].([]interface{})
		da.capacity = newCap
		shrunk = true
	}
	// or grow
	for ; da.capacity < newSize; da.capacity *= da.clusterSize  {
		newArray := make([]interface{}, da.clusterSize)
		newArray[0] = da.tree
		da.tree = newArray
	}
	if shrunk {
		recursiveRelease (newSize, da.tree, da.capacity, da.clusterSize)
	}
}

func (da *DynamicArray) At (i uint) *interface{} {
	if i > da.size {
		panic ("DynamicArray.At(): Value out of range")
	}
	return recursiveAt (i, da.tree, da.capacity, da.clusterSize)
}

func (da *DynamicArray) Size() uint {
	return da.size
}

func (da *DynamicArray) SetSize (newSize uint) {
	da.shrinkOrGrow (newSize)
	da.size = newSize
}