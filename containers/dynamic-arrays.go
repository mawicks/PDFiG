/*
	Package of containers.
*/
package containers

// DynamicArray implements the containters.Array interface. 

// Resizing a DynamicArray does not involve copying.  DynamicArray is
// a sparse array.  Portions of the array are allocated chunks of the
// "clusterSize" parameter passed to the constructor.  The larger the
// cluster size, the faster the access.  Currently, shrinking is not
// very efficient for large cluster sizes.  Currently, PopBack() is
// not very efficient.

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

func (da *DynamicArray) At (i uint) *interface{} {
	if i > da.size {
		panic ("DynamicArray.At(): Value out of range")
	}
	var at func (i uint, tree []interface{}, capacity uint) (result *interface{})

	at = func (i uint, tree []interface{}, capacity uint) (result *interface{}) {
		if capacity <= da.clusterSize {
			result = &tree[i]
		} else {
			subtreeCapacity := capacity / da.clusterSize
			theSubtree := i/subtreeCapacity
			// Create a subtree if necessary
			if tree[theSubtree] == nil {
				tree[theSubtree] =  make([]interface{}, da.clusterSize)
			}
			result = at (i%subtreeCapacity, tree[theSubtree].([]interface{}), subtreeCapacity)
		}
		return result
	}
	
	return at (i, da.tree, da.capacity)
}

func (da *DynamicArray) PushBack (v interface{}) {
	da.SetSize (da.size + 1)
	*da.At (da.size-1) = v
}

func (da *DynamicArray) PopBack () interface{} {
	result := *da.At(da.size-1)
	da.SetSize(da.size-1)
	return result
}

func (da *DynamicArray) SetSize (newSize uint) {
//	fmt.Printf ("\nshrinkOrGrow: newSize=%v, da.capacity=%v, da.clusterSize=%v, da.tree=%v\n",
//		newSize, da.capacity,da.clusterSize,da.tree)
	// Shrink if necessary
	for newCap := da.capacity/da.clusterSize; newCap >= newSize && newCap > 1; newCap /= da.clusterSize {
		da.tree = da.tree[0].([]interface{})
		da.capacity = newCap
	}
	// Grow if necessary
	for ; da.capacity < newSize; da.capacity *= da.clusterSize  {
		newArray := make([]interface{}, da.clusterSize)
		newArray[0] = da.tree
		da.tree = newArray
	}
	// Release unused elements after shrinking
	if newSize < da.size {
		var release func (size uint, tree[]interface{}, capacity uint)

		release = func (lastKeepItem uint, tree[]interface{}, capacity uint) {
			subtreeCapacity := capacity/da.clusterSize
			lastOccupiedSubtree := lastKeepItem/subtreeCapacity
			for i:=lastOccupiedSubtree+1; i<da.clusterSize; i++ {
				tree[i] = nil
			}
			if subtreeCapacity > 1 && tree[lastOccupiedSubtree] != nil {
				release (lastKeepItem%subtreeCapacity, tree[lastOccupiedSubtree].([]interface{}), subtreeCapacity)
			}
		}

		release (newSize-1, da.tree, da.capacity)
	}
	da.size = newSize
//	fmt.Printf ("Shrink or grow returning: size=%v, tree=%v, capacity=%v\n", da.size, da.tree, da.capacity)
}

func (da *DynamicArray) Size() uint {
	return da.size
}

