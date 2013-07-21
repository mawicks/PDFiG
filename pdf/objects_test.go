package pdf_test

import (
	"fmt"
	"github.com/mawicks/PDFiG/pdf"
	"strconv"
	"testing"
	)

// First define some helper functions
func toString(object pdf.Object, file... pdf.File) string {
	return (&pdf.ObjectStringDecorator{object}).String(file...)
}

// checkObject() requires that the serialization of object matches
// *one* of the elements of expect.
func checkObjectBasic(t *testing.T, descr string, o pdf.Object, file pdf.File, expect... string) {
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
			t.Errorf(`%s produced %s; expected %s`, descr, strconv.Quote(s), strconv.Quote(expect[0]))
		} else {
			t.Errorf(`%s produced %s; expected *one* element of %v`, descr, strconv.Quote(s), expect)
		}
	}
}

func checkObject(t *testing.T, descr string, o pdf.Object, file pdf.File, expect... string) {
	checkObjectBasic(t, descr, o, file, expect...)

	// Make sure that Clone(), Dereference(), Protected(), and Unprotected() can be called without
	// errors and without changing the serialization.
	checkObjectBasic(t, descr + " (Clone())", o.Clone(), file, expect...)
	checkObjectBasic(t, descr + " (Dereference())", o.Dereference(), file, expect...)
	checkObjectBasic(t, descr + " (Protected())", o.Protected(), file, expect...)
	checkObjectBasic(t, descr + " (Unprotected())", o.Unprotected(), file, expect...)
	// Make sure that Unprotected() can be called on objects with
	// protected wrappers without changing the serialization, and
	// vice-versa.
	checkObjectBasic(t, descr + " (Protected() & Unprotected())", o.Protected().Unprotected(), file, expect...)
	checkObjectBasic(t, descr + " (Unprotected()& Protected())", o.Unprotected().Protected(), file, expect...)
}

// Make sure ObjectStringDecorator delegates the Serialize method
func TestObjectDecorator(t *testing.T) {
	n := pdf.NewNull()
	o := pdf.ObjectStringDecorator{n}
	checkObject(t, "ObjectStringDecorator", o, nil, "null")
}

func checkStringFromText(t *testing.T, testValue, expect string, serializer func(pdf.String, pdf.Writer)) {
	s := pdf.NewTextString(testValue)
	s.SetSerializer(serializer)
	hex := toString(s)
	if hex != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, hex)
	}
}

func checkStringFromBytes(t *testing.T, testValue, expect string, serializer func(pdf.String, pdf.Writer)) {
	s := pdf.NewBinaryString([]byte(testValue))
	s.SetSerializer(serializer)
	output := toString(s)
	if output != expect {
		t.Errorf(`NewTextString(%s) produced "%s"`, testValue, output)
	}
}

// Unit tests for each of the object types follow:
func TestNull(t *testing.T) {
	checkObject(t, "NewNull()", pdf.NewNull(), nil, "null")
}

func TestBoolean(t *testing.T) {
	checkObject(t, "NewBoolean(false)", pdf.NewBoolean(false), nil, "false")
	checkObject(t, "NewBoolean(true)", pdf.NewBoolean(true), nil, "true")
}

func TestNumeric(t *testing.T) {
	checkObject(t, "NewNumeric(1)", pdf.NewNumeric(1), nil, "1")
	checkObject(t, "NewNumeric(3.14159)", pdf.NewNumeric(3.14159), nil, "3.14159")
	checkObject(t, "NewNumeric(0.1)", pdf.NewNumeric(0.1), nil, "0.1")
	checkObject(t, "NewNumeric(2147483647)", pdf.NewNumeric(2147483647), nil, "2147483647")
	checkObject(t, "NewNumeric(-2147483648)", pdf.NewNumeric(-2147483648), nil, "-2147483648")
	checkObject(t, "NewNumeric(3.403e+38)", pdf.NewNumeric(3.403e+38), nil, "340282350000000000000000000000000000000")
	checkObject(t, "NewNumeric(-3.403e+38)", pdf.NewNumeric(-3.403e+38), nil, "-340282350000000000000000000000000000000")

	// The PDF spec recommends setting anything below +/-
	// 1.175e-38 to 0 in case a conforming reader uses 32-bit
	// floats.  Here, Adobe is referring to the smallest number
	// that can be represented without losing precision rather
	// than the smallest number that can be represented with a
	// float32. It's odd that Adobe thinks that setting small
	// numbers to zero is better than accepting a representable
	// number with a loss of precision.
	checkObject(t, "NewNumeric(1.176e-38)", pdf.NewNumeric(1.176e-38), nil, "0.00000000000000000000000000000000000001176")
	checkObject(t, "NewNumeric(-1.176e-38)", pdf.NewNumeric(-1.176e-38), nil, "-0.00000000000000000000000000000000000001176")
	checkObject(t, "NewNumeric(1.175e-38)", pdf.NewNumeric(1.175e-38), nil, "0")
	checkObject(t, "NewNumeric(-1.175e-38)", pdf.NewNumeric(-1.175e-38), nil, "0")
}

func TestName(t *testing.T) {
	checkObject(t, `NewName("foo")`, pdf.NewName("foo"), nil, "/foo")
	checkObject(t, `NewName("résumé")`, pdf.NewName("résumé"), nil, "/résumé")
	checkObject(t, `NewName("#foo")`, pdf.NewName("#foo"), nil, "/#23foo")
	checkObject(t, `NewName(" foo")`, pdf.NewName(" foo"), nil, "/#20foo")
	checkObject(t, `NewName("(foo)")`, pdf.NewName("(foo)"), nil, "/#28foo#29")
}

func TestString(t *testing.T) {
	checkObject(t, `NewTextString("foo")`, pdf.NewTextString("foo"), nil, "(foo)")
	checkObject(t, `NewTextString("()\\"`, pdf.NewTextString("()\\"), nil, "(\\(\\)\\\\)")
	checkObject(t, `NewTextString("[]")`, pdf.NewTextString("[]"), nil, "([])")
	checkObject(t, `NewTextString("")`, pdf.NewTextString(""), nil, "()")

	checkStringFromText(t, "[]", "<5B5D>", pdf.HexStringSerializer)
	checkStringFromText(t, "0123", "<30313233>", pdf.HexStringSerializer)
	checkStringFromText(t, "", "<>", pdf.HexStringSerializer)
	checkStringFromText(t, "foo", "(foo)", pdf.AsciiStringSerializer)
	checkStringFromText(t, "\n\r\t", "(\\n\\r\\t)", pdf.AsciiStringSerializer)
	checkStringFromText(t, "\n\r\t\b\f", "(\\376\\377\\000\\n\\000\\r\\000\\t\\000\\b\\000\\f)", pdf.AsciiStringSerializer)
	checkStringFromBytes(t, "\200", "(\\200)", pdf.AsciiStringSerializer)
	checkStringFromBytes(t, "\200", "<80>", pdf.HexStringSerializer)
}

func TestArray(t *testing.T) {
	// Make sure empty array and single element array are serialized properly
	a := pdf.NewArray()
	checkObject(t, "NewArray()", a, nil, "[]")

	a.Add(pdf.NewNumeric(1))
	checkObject(t, "NewArray() with NewNumeric(1)", a, nil, "[1]")

	// Make sure nested arrays are serialized properly
	sa := pdf.NewArray()
	sa.Add(pdf.NewNumeric(2))
	checkObject(t, "NewArray() with NewNumeric(2)", sa, nil, "[2]")
	a.Add(sa)
	checkObject(t, "Add an element to an Array", a, nil, "[1 [2]]")

	// Make sure that clone() produces a copy unaffected by future operations on a.
	c := a.Clone()
	pa := a.Protected().(pdf.ProtectedArray)
	if _,ok := pa.(pdf.Array); ok {
		t.Error ("Protected array can be cast back to Array")
	}

	// Make sure that additions to nested arrays are reflected in enclosing array
	sa.Add(pdf.NewNumeric(3))
	checkObject(t, "Add an element to a subarray", a, nil, "[1 [2 3]]")
	checkObject(t, "Add an element to a subarray", pa, nil, "[1 [2 3]]")

	// Test that At() works and returns references that can be modified.
	saref := a.At(1).(pdf.Array)
	checkObject(t, "Reference returned by pdf.Array.At()", saref, nil, "[2 3]")

	saref.Add(pdf.NewNumeric(4))
	checkObject(t, "Append to subarray reference", saref, nil, "[2 3 4]")
	checkObject(t, "Enclosing array after append to subarray reference", a, nil, "[1 [2 3 4]]")
	checkObject(t, "Enclosing array after append to subarray reference", pa, nil, "[1 [2 3 4]]")

	// Test pdf.Array.Append()
	a.Append(c.(pdf.Array))
	checkObject(t, "Append to Array", a, nil, "[1 [2 3 4] 1 [2]]")
	checkObject(t, "Append to Array", pa, nil, "[1 [2 3 4] 1 [2]]")

	// Check that Unprotected() protected arrays can be modified.
	upa := pa.Unprotected().(pdf.Array)
	upa.Add(pdf.NewNumeric(5))
	checkObject(t, "Add to unprotected Array", upa, nil, "[1 [2 3 4] 1 [2] 5]")

	// Check that subarray references obtained from protected
	// arrays are protected
	sfp := pa.At(1).(pdf.ProtectedArray)
	checkObject(t, "Subarray from protected", sfp, nil, "[2 3 4]")

	if _,ok := sfp.(pdf.Array); ok {
		t.Error ("Subarray obtained from protected array can be cast back to Array")
	}

	// Make sure the original "a" hasn't changed by operations on its protected version.
	checkObject(t, "Array test", a, nil, "[1 [2 3 4] 1 [2]]")

	// Make sure that an earlier clone of "a" hasn't changed
	checkObject(t, "Array test", c, nil, "[1 [2]]")
}

func TestDictionaryGetters (t *testing.T) {
	d := pdf.NewDictionary()

	// Verify that GetBoolean() works.
	d.Add("a", pdf.NewBoolean(false))
	if b,ok := d.GetBoolean("a"); !ok || b {
		t.Error(`GetBoolean failed to retrieve valid value`)
	}
	
	// Verify that GetInt() works.
	d.Add("a", pdf.NewNumeric(1))
	if i,ok := d.GetInt("a"); !ok || i!=1 {
		t.Error(`GetInt() failed to retrieve valid value`)
	}

	// Verify that GetReal() works.
	d.Add("a", pdf.NewNumeric(3.14))
	if x,ok := d.GetReal("a"); !ok || x!=3.14 {
		t.Error(`GetReal() failed to retrieve valid value`);
	}

	// Verify that GetName() works.
	d.Add("a", pdf.NewName("namevalue"))
	if n,ok := d.GetName("a"); !ok || n != "namevalue" {
		t.Error(`GeName() failed to retrieve valid value`)
	}

	// Verify that GetString() works.
	d.Add("a", pdf.NewTextString("string"))
	if n,ok := d.GetString("a"); !ok || string(n) != "string" {
		t.Error(`GetString() failed to retrieve valid value`)
	}

	if d.CheckNameValue("doesntexist", "value") {
		t.Error(`CheckNameValue returned true on empty dictionary`)
	}

	d.Add("a", pdf.NewTextString("string"))
	if d.CheckNameValue("c", "string") {
		t.Error(`CheckNameValue returned true on non-name`)
	}
}

func TestDictionary(t *testing.T) {
	d := pdf.NewDictionary()
	checkObject(t, "NewDictionary", d, nil, "<<>>")

	d.Add("a", pdf.NewNumeric(1))
	checkObject(t, "Dictionary.Add() test", d, nil, "<</a 1>>")

	cd := d.Clone()
	pd := d.Protected().(pdf.ProtectedDictionary)

	if _,ok := pd.(pdf.Dictionary); ok {
		t.Error("Protected dictionary can be cast back to Dictionary")
	}

	// Can't test beyond two entries very easily because the order of entries is not specified
	// and the number of permutations makes it not worth the effort with our simple checkObject() function.
	sa := pdf.NewArray()
	sa.Add(pdf.NewNumeric(2))
	d.Add("b", sa)
	checkObject(t, "Dictionary.Remove() test", d, nil, "<</a 1 /b [2]>>", "<</b [2] /a 1>>")
	// Verify that pd is tracking d
	checkObject(t, "Dictionary.Remove() test", pd, nil, "<</a 1 /b [2]>>", "<</b [2] /a 1>>")

	saref := d.GetArray("b").(pdf.Array)
	saref.Add(pdf.NewNumeric(3))
	checkObject(t, "Dictionary.Remove() test", d, nil, "<</a 1 /b [2 3]>>", "<</b [2 3] /a 1>>")

	// Make sure clone hasn't changed.
	checkObject(t, "Dictionary.Clone() test", cd, nil, "<</a 1>>")

	upd := pd.Unprotected().(pdf.Dictionary)
	upd.Remove("b")
	checkObject(t, "Remove from unprotected dictionary", upd, nil, "<</a 1>>")
	checkObject(t, "Protected dictionary after remove from unprotected counterpart", pd, nil, "<</a 1 /b [2 3]>>", "<</b [2 3] /a 1>>")

	// Begin removing entries to test Remove() method.
	d.Remove("a")
	checkObject(t, "Dictionary.Remove() test", d, nil, "<</b [2 3]>>")
	checkObject(t, "Dictionary.Clone() test", cd, nil, "<</a 1>>")

	d.Remove("b")
	checkObject(t, "Dictionary.Remove() test", d, nil, "<<>>")

	// d is empty now.

}

func TestStream(t *testing.T) {
	s := pdf.NewStream()
	fmt.Fprint(s, "foo")
	checkObject(t, "NewStream", s, nil, "<</Length 3>>\nstream\nfoo\nendstream")
	// Ensure stream can be serialized more than once (i.e., that reading the internal buffer doesn't
	// affect the current read position.
	checkObject(t, "NewStream", s, nil, "<</Length 3>>\nstream\nfoo\nendstream")
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
	checkObjectBasic(t, "Indirect test", i1, f1, "21 37 R")
	checkObjectBasic(t, "Indirect test", i1, f2, "42 23 R")
	checkObjectBasic(t, "Indirect test", i2, f1, "22 38 R")
	checkObjectBasic(t, "Indirect test", i2, f2, "43 24 R")

	// Repeat test to make sure object-file associations have been
	// cached.
	checkObjectBasic(t, "Indirect test", i1, f1, "21 37 R")
	checkObjectBasic(t, "Indirect test", i1, f2, "42 23 R")
	checkObjectBasic(t, "Indirect test", i2, f1, "22 38 R")
	checkObjectBasic(t, "Indirect test", i2, f2, "43 24 R")

	// Now test clones
	checkObjectBasic(t, "Indirect test", c1, f1, "21 37 R")
	checkObjectBasic(t, "Indirect test", c1, f2, "42 23 R")
}

func TestRectangle(t *testing.T) {
	checkObject(t, "Rectangle test", pdf.NewRectangle(1, 2, 3, 4), nil, "[1 2 3 4]")
}
