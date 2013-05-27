package containers

type Array interface {
	// Return element at specified position
	At (uint) *interface {}
	// Set the size
	SetSize(uint)
	// Return size (number of elements)
	Size() uint
}

type Stack interface {
	PushBack (v interface{})
	PushFront (v interface{})
	PopBack() interface{}
	PopFront() interface{}
}

type ArrayStack interface {
	Array
	Stack
}

type StackArrayDecorator struct {
	Array
}

func (sad StackArrayDecorator) PushBack (v interface{}) {
	sad.SetSize (sad.Size() + 1)
	*sad.At (sad.Size()-1) = v
}

func (sad StackArrayDecorator) PushFront (v interface{}) {
	sad.SetSize (sad.Size() + 1)
	for i:= sad.Size()-1; i!=0; i-- {
		*sad.At(i) = *sad.At(i-1)
	}
	*sad.At(0) = v
}

func (sad StackArrayDecorator) PopBack () interface{} {
	result := *sad.At(sad.Size()-1)
	sad.SetSize(sad.Size()-1)
	return result
}

func (sad StackArrayDecorator) PopFront () interface{} {
	result := *sad.At(0)
	s := sad.Size()
	for i:=uint(1); i<s; i++ {
		*sad.At(i-1) = *sad.At(i)
	}
	sad.SetSize(s-1)
	return result
}
