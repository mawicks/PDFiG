package pdf

type Rectangle struct {
	ProtectedArray
}

func NewRectangle(llx, lly, urx, ury float64) *Rectangle {
	result := NewArray()
	result.Add(NewNumeric(llx))
	result.Add(NewNumeric(lly))
	result.Add(NewNumeric(urx))
	result.Add(NewNumeric(ury))
	return &Rectangle{result}
}
