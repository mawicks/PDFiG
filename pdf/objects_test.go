package pdf

import "bytes"
import "testing"

// First define some helper functions

func toString (object Object) string {
	var buffer bytes.Buffer

	object.Serialize (&buffer)
	return buffer.String()
}

func testOneObject (t *testing.T, d string, o Object, expect string) {
	if s := toString(o); s != expect {
		t.Errorf (`%s produced "%s"; expected "%s"`, d, s, expect)
	}
}

func testOneStringAsHex (t *testing.T, testValue, expect string) {
	s := NewString(testValue)
	s.SetHexOutput()
	hex := toString(s)
	if hex != expect {
		t.Errorf (`NewString(%s) produced "%s"`, testValue, hex)
	}
}

func testOneStringAsAscii (t *testing.T, testValue, expect string) {
	s := NewString(testValue)
	s.SetAsciiOutput()
	output := toString(s)
	if output != expect {
		t.Errorf (`NewString(%s) produced "%s"`, testValue, output)
	}
}

// Unit tests follow

func TestNull(t *testing.T) {
	expect := "null"
	if s := toString(&Null{}); s != expect {
		t.Errorf (`null.Serialize() produced "%s"; expected "%s"`, s, expect)
	}
}


func TestBoolean(t *testing.T) {
	testOneObject (t, "NewBoolean(false)", NewBoolean(false), "false")
	testOneObject (t, "NewBoolean(true)", NewBoolean(true), "true")
}


func TestNumeric(t *testing.T) {
	testOneObject(t, "NewNumeric(1)", NewNumeric(1), "1")
	testOneObject(t, "NewNumeric(3.14159)", NewNumeric(3.14159), "3.14159")
	testOneObject(t, "NewNumeric(0.1)", NewNumeric(0.1), "0.1")
	testOneObject(t, "NewNumeric(2147483647)", NewNumeric(2147483647), "2147483647")
	testOneObject(t, "NewNumeric(-2147483648)", NewNumeric(-2147483648), "-2147483648")
	testOneObject(t, "NewNumeric(3.403e+38)", NewNumeric(3.403e+38), "3.4028235e+38")
	testOneObject(t, "NewNumeric(-3.403e+38)", NewNumeric(-3.403e+38), "-3.4028235e+38")

	// The PDF spec recommends setting anything below +/-
	// 1.175e-38 to 0 in case a conforming reader uses 32 bit
	// floats.  Here, Adobe is referring to the smallest number
	// that can be represented without losing precision rather
	// than the smallest number that can be represented with a
	// float32. It's odd that Adobe thinks that setting small
	// numbers to zero is better than accepting a representable
	// number with a loss of precision.

	testOneObject(t, "NewNumeric(1.176e-38)", NewNumeric(1.176e-38), "1.176e-38")
	testOneObject(t, "NewNumeric(-1.176e-38)", NewNumeric(-1.176e-38), "-1.176e-38")
	testOneObject(t, "NewNumeric(1.175e-38)", NewNumeric(1.175e-38), "0")
	testOneObject(t, "NewNumeric(-1.175e-38)", NewNumeric(-1.175e-38), "0")
}

func TestName (t *testing.T) {
	testOneObject (t, `NewName("foo")`, NewName("foo"), "/foo")
	testOneObject (t, `NewName("résumé")`, NewName("résumé"), "/résumé")
	testOneObject (t, `NewName("#foo")`, NewName("#foo"), "/#23foo")
	testOneObject (t, `NewName(" foo")`, NewName(" foo"), "/#20foo")
	testOneObject (t, `NewName("(foo)")`, NewName("(foo)"), "/#28foo#29")
}

func TestString (t *testing.T) {
	testOneObject (t, `NewString("foo")`, NewString("foo"), "(foo)")
	testOneObject (t, `NewString("()\\"`, NewString("()\\"), "(\\(\\)\\\\)")
	testOneObject (t, `NewString("[]")`, NewString("[]"), "([])")
	testOneObject (t, `NewString("")`, NewString(""), "()")
	testOneStringAsHex (t, "[]", "<5B5D>")
	testOneStringAsHex (t, "0123", "<30313233>")
	testOneStringAsHex (t, "", "<>")
	testOneStringAsAscii (t, "foo", "(foo)")
	testOneStringAsAscii (t, "\200", "(\\200)")
	testOneStringAsAscii (t, "\n\r\t\b\f", "(\\n\\r\\t\\b\\f)")
}

func TestArray (t *testing.T) {
	a := NewArray()
	testOneObject (t, "NewArray()", a, "[]");

	a.Add (NewNumeric(3.14))
	testOneObject (t, "NewArray() with NewNumeric(3.14)", a, "[3.14]");

	a.Add (NewNumeric(2.718))
	testOneObject (t, "Array test", a, "[3.14 2.718]");

	a.Add (NewName("f o o"))
	testOneObject (t, "Array test", a, "[3.14 2.718 /f#20o#20o]");
}

func TestDictionary (t *testing.T) {
	d := NewDictionary()
	testOneObject (t, "NewDictionary", d, "<<>>");

	d.Add ("fee", NewNumeric(3.14))
	testOneObject (t, "Dictionary.Add() test", d, "<</fee 3.14>>");

	// Can't test beyond one entry reliably because order of entries is not defined
	// We can add a new entry and remove one of the earlier entries.
	d.Add ("fi", NewNumeric(2.718))
	d.Remove ("fee")
	testOneObject (t, "Dictionary.Remove() test", d, "<</fi 2.718>>")

	d.Remove ("fi")
	testOneObject (t, "Dictionary.Remove() test", d, "<<>>")
}

