package pdf

import "testing"
import "fmt"

// First define some helper functions
func toString(object Object, file ...File) string {
	return (&ObjectStringDecorator{object}).String(file...)
}

// TestOneObject requires that the serialization of object matches
// *one* of the elements of expect.
func testOneObject(t *testing.T, d string, o Object, file File, expect ...string) {
	matched := false
	s := toString(o, file)
	for _, e := range expect {
		if s == e {
			matched = true
			break
		}
	}
	if !matched {
		if len(expect) == 1 {
			t.Errorf(`%s produced "%s"; expected "%s"`, d, s, expect[0])
		} else {
			t.Errorf(`%s produced "%s"; expected *one* element of %v`, d, s, expect)
		}
	}
}

// Make sure ObjectStringDecorator delegates the Serialize method
func TestObjectDecorator(t *testing.T) {
	n := NewNull()
	o := ObjectStringDecorator{n}
	testOneObject(t, "ObjectStringDecorator", o, nil, "null")
}

func testOneTextStringAsHex(t *testing.T, testValue, expect string) {
	s := NewTextString(testValue)
	s.SetHexOutput()
	hex := toString(s)
	if hex != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, hex)
	}
}

func testOneBinaryStringAsHex(t *testing.T, testValue, expect string) {
	s := NewBinaryString([]byte(testValue))
	s.SetHexOutput()
	hex := toString(s)
	if hex != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, hex)
	}
}

func testOneTextStringAsAscii(t *testing.T, testValue, expect string) {
	s := NewTextString(testValue)
	s.SetAsciiOutput()
	output := toString(s)
	if output != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, output)
	}
}

func testOneBinaryStringAsAscii(t *testing.T, testValue, expect string) {
	s := NewBinaryString([]byte(testValue))
	s.SetAsciiOutput()
	output := toString(s)
	if output != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, output)
	}
}

// Unit tests follow
func TestNull(t *testing.T) {
	testOneObject(t, "NewNull()", NewNull(), nil, "null")
}

func TestBoolean(t *testing.T) {
	testOneObject(t, "NewBoolean(false)", NewBoolean(false), nil, "false")
	testOneObject(t, "NewBoolean(true)", NewBoolean(true), nil, "true")
}

func TestNumeric(t *testing.T) {
	testOneObject(t, "NewNumeric(1)", NewNumeric(1), nil, "1")
	testOneObject(t, "NewNumeric(3.14159)", NewNumeric(3.14159), nil, "3.14159")
	testOneObject(t, "NewNumeric(0.1)", NewNumeric(0.1), nil, "0.1")
	testOneObject(t, "NewNumeric(2147483647)", NewNumeric(2147483647), nil, "2147483647")
	testOneObject(t, "NewNumeric(-2147483648)", NewNumeric(-2147483648), nil, "-2147483648")
	testOneObject(t, "NewNumeric(3.403e+38)", NewNumeric(3.403e+38), nil, "3.4028235e+38")
	testOneObject(t, "NewNumeric(-3.403e+38)", NewNumeric(-3.403e+38), nil, "-3.4028235e+38")

	// The PDF spec recommends setting anything below +/-
	// 1.175e-38 to 0 in case a conforming reader uses 32 bit
	// floats.  Here, Adobe is referring to the smallest number
	// that can be represented without losing precision rather
	// than the smallest number that can be represented with a
	// float32. It's odd that Adobe thinks that setting small
	// numbers to zero is better than accepting a representable
	// number with a loss of precision.

	testOneObject(t, "NewNumeric(1.176e-38)", NewNumeric(1.176e-38), nil, "1.176e-38")
	testOneObject(t, "NewNumeric(-1.176e-38)", NewNumeric(-1.176e-38), nil, "-1.176e-38")
	testOneObject(t, "NewNumeric(1.175e-38)", NewNumeric(1.175e-38), nil, "0")
	testOneObject(t, "NewNumeric(-1.175e-38)", NewNumeric(-1.175e-38), nil, "0")
}

func TestName(t *testing.T) {
	testOneObject(t, `NewName("foo")`, NewName("foo"), nil, "/foo")
	testOneObject(t, `NewName("résumé")`, NewName("résumé"), nil, "/résumé")
	testOneObject(t, `NewName("#foo")`, NewName("#foo"), nil, "/#23foo")
	testOneObject(t, `NewName(" foo")`, NewName(" foo"), nil, "/#20foo")
	testOneObject(t, `NewName("(foo)")`, NewName("(foo)"), nil, "/#28foo#29")
}

func TestString(t *testing.T) {
	testOneObject(t, `NewTextString("foo")`, NewTextString("foo"), nil, "(foo)")
	testOneObject(t, `NewTextString("()\\"`, NewTextString("()\\"), nil, "(\\(\\)\\\\)")
	testOneObject(t, `NewTextString("[]")`, NewTextString("[]"), nil, "([])")
	testOneObject(t, `NewTextString("")`, NewTextString(""), nil, "()")
	testOneTextStringAsHex(t, "[]", "<5B5D>")
	testOneTextStringAsHex(t, "0123", "<30313233>")
	testOneTextStringAsHex(t, "", "<>")
	testOneTextStringAsAscii(t, "foo", "(foo)")
	testOneTextStringAsAscii(t, "\n\r\t", "(\\n\\r\\t)")
	testOneTextStringAsAscii(t, "\n\r\t\b\f", "(\\376\\377\\000\\n\\000\\r\\000\\t\\000\\b\\000\\f)")
	testOneBinaryStringAsAscii(t, "\200", "(\\200)")
	testOneBinaryStringAsHex(t, "\200", "<80>")
}

func TestArray(t *testing.T) {
	a := NewArray()
	testOneObject(t, "NewArray()", a, nil, "[]")

	a.Add(NewNumeric(3.14))
	testOneObject(t, "NewArray() with NewNumeric(3.14)", a, nil, "[3.14]")

	a.Add(NewNumeric(2.718))
	testOneObject(t, "Array test", a, nil, "[3.14 2.718]")

	a.Add(NewName("f o o"))
	testOneObject(t, "Array test", a, nil, "[3.14 2.718 /f#20o#20o]")
}

func TestDictionary(t *testing.T) {
	d := NewDictionary()
	testOneObject(t, "NewDictionary", d, nil, "<<>>")

	d.Add("fee", NewNumeric(3.14))
	testOneObject(t, "Dictionary.Add() test", d, nil, "<</fee 3.14>>")

	// Can't test beyond three entries very easily because the order of entries is not specified
	// and the number of permutations makes it not worth the effort with our simple testOneObject() function.
	d.Add("fi", NewNumeric(2.718))
	testOneObject(t, "Dictionary.Remove() test", d, nil, "<</fee 3.14 /fi 2.718>>", "<</fi 2.718 /fee 3.14>>")

	// Begin removing entries to test Remove() method.
	d.Remove("fee")
	testOneObject(t, "Dictionary.Remove() test", d, nil, "<</fi 2.718>>")

	d.Remove("fi")
	testOneObject(t, "Dictionary.Remove() test", d, nil, "<<>>")
}

func TestStream(t *testing.T) {
	s := NewStream()
	fmt.Fprint(s, "foo")
	testOneObject(t, "NewStream", s, nil, "<</Length 3>>\nstream\nfoo\nendstream")
}

func TestIndirect(t *testing.T) {
	// Two objects
	i1 := NewIndirect()
	i2 := NewIndirect()

	// Two files
	f1 := NewTestFile(21, 37)
	f2 := NewTestFile(42, 23)

	// Test all four combinations
	testOneObject(t, "Indirect test", i1, f1, "21 37 R")
	testOneObject(t, "Indirect test", i1, f2, "42 23 R")
	testOneObject(t, "Indirect test", i2, f1, "22 38 R")
	testOneObject(t, "Indirect test", i2, f2, "43 24 R")

	// Repeat test to make sure object-file associations have been
	// cached.
	testOneObject(t, "Indirect test", i1, f1, "21 37 R")
	testOneObject(t, "Indirect test", i1, f2, "42 23 R")
	testOneObject(t, "Indirect test", i2, f1, "22 38 R")
	testOneObject(t, "Indirect test", i2, f2, "43 24 R")
}

func TestRectangle(t *testing.T) {
	testOneObject(t, "Rectangle test", NewRectangle(1, 2, 3, 4), nil, "[1 2 3 4]")
}
