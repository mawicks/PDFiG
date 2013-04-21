package pdf_test

import (
	"github.com/mawicks/PDFiG/pdf"
	"strconv"
	"testing"
	"fmt" )

// First define some helper functions
func toString(object pdf.Object, file ...pdf.File) string {
	return (&pdf.ObjectStringDecorator{object}).String(file...)
}

// TestOneObject requires that the serialization of object matches
// *one* of the elements of expect.
func testOneObject(t *testing.T, d string, o pdf.Object, file pdf.File, expect ...string) {
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
			t.Errorf(`%s produced %s; expected %s`, d, strconv.Quote(s), strconv.Quote(expect[0]))
		} else {
			t.Errorf(`%s produced %s; expected *one* element of %v`, d, strconv.Quote(s), expect)
		}
	}
}

// Make sure ObjectStringDecorator delegates the Serialize method
func TestObjectDecorator(t *testing.T) {
	n := pdf.NewNull()
	o := pdf.ObjectStringDecorator{n}
	testOneObject(t, "ObjectStringDecorator", o, nil, "null")
}

func testOneTextStringAsHex(t *testing.T, testValue, expect string) {
	s := pdf.NewTextString(testValue)
	s.SetHexOutput()
	hex := toString(s)
	if hex != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, hex)
	}
}

func testOneBinaryStringAsHex(t *testing.T, testValue, expect string) {
	s := pdf.NewBinaryString([]byte(testValue))
	s.SetHexOutput()
	hex := toString(s)
	if hex != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, hex)
	}
}

func testOneTextStringAsAscii(t *testing.T, testValue, expect string) {
	s := pdf.NewTextString(testValue)
	s.SetAsciiOutput()
	output := toString(s)
	if output != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, output)
	}
}

func testOneBinaryStringAsAscii(t *testing.T, testValue, expect string) {
	s := pdf.NewBinaryString([]byte(testValue))
	s.SetAsciiOutput()
	output := toString(s)
	if output != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, output)
	}
}

// Unit tests follow
func TestNull(t *testing.T) {
	testOneObject(t, "NewNull()", pdf.NewNull(), nil, "null")
	testOneObject(t, "NewNull().Clone()", pdf.NewNull().Clone(), nil, "null")
}

func TestBoolean(t *testing.T) {
	testOneObject(t, "NewBoolean(false)", pdf.NewBoolean(false), nil, "false")
	testOneObject(t, "NewBoolean(true)", pdf.NewBoolean(true), nil, "true")
	testOneObject(t, "NewBoolean(false).Clone()", pdf.NewBoolean(false).Clone(), nil, "false")
	testOneObject(t, "NewBoolean(true).Clone()", pdf.NewBoolean(true).Clone(), nil, "true")
}

func TestNumeric(t *testing.T) {
	testOneObject(t, "NewNumeric(1)", pdf.NewNumeric(1), nil, "1")
	testOneObject(t, "NewNumeric(3.14159)", pdf.NewNumeric(3.14159), nil, "3.14159")
	testOneObject(t, "NewNumeric(0.1)", pdf.NewNumeric(0.1), nil, "0.1")
	testOneObject(t, "NewNumeric(2147483647)", pdf.NewNumeric(2147483647), nil, "2147483647")
	testOneObject(t, "NewNumeric(-2147483648)", pdf.NewNumeric(-2147483648), nil, "-2147483648")
	testOneObject(t, "NewNumeric(3.403e+38)", pdf.NewNumeric(3.403e+38), nil, "340282350000000000000000000000000000000")
	testOneObject(t, "NewNumeric(-3.403e+38)", pdf.NewNumeric(-3.403e+38), nil, "-340282350000000000000000000000000000000")

	testOneObject(t, "NewNumeric(1).Clone()", pdf.NewNumeric(1).Clone(), nil, "1")
	testOneObject(t, "NewNumeric(3.14159).Clone()", pdf.NewNumeric(3.14159).Clone(), nil, "3.14159")
	testOneObject(t, "NewNumeric(0.1).Clone()", pdf.NewNumeric(0.1).Clone(), nil, "0.1")
	testOneObject(t, "NewNumeric(2147483647).Clone()", pdf.NewNumeric(2147483647).Clone(), nil, "2147483647")
	testOneObject(t, "NewNumeric(-2147483648).Clone()", pdf.NewNumeric(-2147483648).Clone(), nil, "-2147483648")
	testOneObject(t, "NewNumeric(3.403e+38).Clone()", pdf.NewNumeric(3.403e+38), nil, "340282350000000000000000000000000000000")
	testOneObject(t, "NewNumeric(-3.403e+38).Clone()", pdf.NewNumeric(-3.403e+38), nil, "-340282350000000000000000000000000000000")

	// The PDF spec recommends setting anything below +/-
	// 1.175e-38 to 0 in case a conforming reader uses 32 bit
	// floats.  Here, Adobe is referring to the smallest number
	// that can be represented without losing precision rather
	// than the smallest number that can be represented with a
	// float32. It's odd that Adobe thinks that setting small
	// numbers to zero is better than accepting a representable
	// number with a loss of precision.

	testOneObject(t, "NewNumeric(1.176e-38)", pdf.NewNumeric(1.176e-38), nil, "0.00000000000000000000000000000000000001176")
	testOneObject(t, "NewNumeric(-1.176e-38)", pdf.NewNumeric(-1.176e-38), nil, "-0.00000000000000000000000000000000000001176")
	testOneObject(t, "NewNumeric(1.175e-38)", pdf.NewNumeric(1.175e-38), nil, "0")
	testOneObject(t, "NewNumeric(-1.175e-38)", pdf.NewNumeric(-1.175e-38), nil, "0")
}

func TestName(t *testing.T) {
	testOneObject(t, `NewName("foo")`, pdf.NewName("foo"), nil, "/foo")
	testOneObject(t, `NewName("résumé")`, pdf.NewName("résumé"), nil, "/résumé")
	testOneObject(t, `NewName("#foo")`, pdf.NewName("#foo"), nil, "/#23foo")
	testOneObject(t, `NewName(" foo")`, pdf.NewName(" foo"), nil, "/#20foo")
	testOneObject(t, `NewName("(foo)")`, pdf.NewName("(foo)"), nil, "/#28foo#29")

	testOneObject(t, `NewName("foo").Clone()`, pdf.NewName("foo").Clone(), nil, "/foo")
	testOneObject(t, `NewName("résumé").Clone()`, pdf.NewName("résumé").Clone(), nil, "/résumé")
	testOneObject(t, `NewName("#foo").Clone()`, pdf.NewName("#foo").Clone(), nil, "/#23foo")
	testOneObject(t, `NewName(" foo").Clone()`, pdf.NewName(" foo").Clone(), nil, "/#20foo")
	testOneObject(t, `NewName("(foo)".Clone())`, pdf.NewName("(foo)").Clone(), nil, "/#28foo#29")
}

func TestString(t *testing.T) {
	testOneObject(t, `NewTextString("foo")`, pdf.NewTextString("foo"), nil, "(foo)")
	testOneObject(t, `NewTextString("()\\"`, pdf.NewTextString("()\\"), nil, "(\\(\\)\\\\)")
	testOneObject(t, `NewTextString("[]")`, pdf.NewTextString("[]"), nil, "([])")
	testOneObject(t, `NewTextString("")`, pdf.NewTextString(""), nil, "()")
	testOneTextStringAsHex(t, "[]", "<5B5D>")
	testOneTextStringAsHex(t, "0123", "<30313233>")
	testOneTextStringAsHex(t, "", "<>")
	testOneTextStringAsAscii(t, "foo", "(foo)")
	testOneTextStringAsAscii(t, "\n\r\t", "(\\n\\r\\t)")
	testOneTextStringAsAscii(t, "\n\r\t\b\f", "(\\376\\377\\000\\n\\000\\r\\000\\t\\000\\b\\000\\f)")
	testOneBinaryStringAsAscii(t, "\200", "(\\200)")
	testOneBinaryStringAsHex(t, "\200", "<80>")

	testOneObject(t, `NewTextString("foo").Clone()`, pdf.NewTextString("foo").Clone(), nil, "(foo)")
	testOneObject(t, `NewTextString("()\\".Clone()`, pdf.NewTextString("()\\").Clone(), nil, "(\\(\\)\\\\)")
	testOneObject(t, `NewTextString("[]").Clone()`, pdf.NewTextString("[]").Clone(), nil, "([])")
	testOneObject(t, `NewTextString("").Clone()`, pdf.NewTextString("").Clone(), nil, "()")
}

func TestArray(t *testing.T) {
	a := pdf.NewArray()
	testOneObject(t, "NewArray()", a, nil, "[]")

	a.Add(pdf.NewNumeric(3.14))
	testOneObject(t, "NewArray() with NewNumeric(3.14)", a, nil, "[3.14]")

	a.Add(pdf.NewNumeric(2.718))
	testOneObject(t, "Array test", a, nil, "[3.14 2.718]")
	b := a.Clone()
	testOneObject(t, "Array test", b, nil, "[3.14 2.718]")

	a.Add(pdf.NewName("f o o"))
	testOneObject(t, "Array test", a, nil, "[3.14 2.718 /f#20o#20o]")

	// Make sure clone hasn't changed
	testOneObject(t, "Array test", b, nil, "[3.14 2.718]")
}

func TestDictionary(t *testing.T) {
	d := pdf.NewDictionary()
	testOneObject(t, "NewDictionary", d, nil, "<<>>")

	d.Add("fee", pdf.NewNumeric(3.14))
	testOneObject(t, "Dictionary.Add() test", d, nil, "<</fee 3.14>>")

	b := d.Clone()
	testOneObject(t, "Dictionary.Clone() test", b, nil, "<</fee 3.14>>")

	// Can't test beyond three entries very easily because the order of entries is not specified
	// and the number of permutations makes it not worth the effort with our simple testOneObject() function.
	d.Add("fi", pdf.NewNumeric(2.718))
	testOneObject(t, "Dictionary.Remove() test", d, nil, "<</fee 3.14 /fi 2.718>>", "<</fi 2.718 /fee 3.14>>")

	// Make sure clone hasn't changed.
	testOneObject(t, "Dictionary.Clone() test", b, nil, "<</fee 3.14>>")

	// Begin removing entries to test Remove() method.
	d.Remove("fee")
	testOneObject(t, "Dictionary.Remove() test", d, nil, "<</fi 2.718>>")
	testOneObject(t, "Dictionary.Clone() test", b, nil, "<</fee 3.14>>")

	d.Remove("fi")
	testOneObject(t, "Dictionary.Remove() test", d, nil, "<<>>")
}

func TestStream(t *testing.T) {
	s := pdf.NewStream()
	fmt.Fprint(s, "foo")
	testOneObject(t, "NewStream", s, nil, "<</Length 3>>\nstream\nfoo\nendstream")
	// Ensure stream can be serialized more than once (i.e., that reading the internal buffer doesn't
	// affect the current read position.
	testOneObject(t, "NewStream", s, nil, "<</Length 3>>\nstream\nfoo\nendstream")
	b := s.Clone()
	testOneObject(t, "Stream.Clone()", b, nil, "<</Length 3>>\nstream\nfoo\nendstream")
}

func TestIndirect(t *testing.T) {
	// Two objects
	i1 := pdf.NewIndirect()
	i2 := pdf.NewIndirect()

	// Two files
	f1 := pdf.NewMockFile(21, 37)
	f2 := pdf.NewMockFile(42, 23)

	c1 := i1.Clone()

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

	// Now test clones
	testOneObject(t, "Indirect test", c1, f1, "21 37 R")
	testOneObject(t, "Indirect test", c1, f2, "42 23 R")
}

func TestRectangle(t *testing.T) {
	testOneObject(t, "Rectangle test", pdf.NewRectangle(1, 2, 3, 4), nil, "[1 2 3 4]")
}
