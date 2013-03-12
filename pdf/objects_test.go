package pdf

import "bytes"
import "testing"

func toString (object Object) string {
	var buffer bytes.Buffer

	object.Serialize (&buffer)
	return buffer.String()
}

func TestNull(t *testing.T) {
	if s := toString(&Null{}); s != "null" {
		t.Errorf ("null.Serialize() produced \"%s\"", s)
	}
}

func TestBoolean(t *testing.T) {
	if s := toString(NewBoolean(false)); s != "false" {
		t.Errorf ("NewBoolean(false).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewBoolean(true)); s != "true" {
		t.Errorf ("NewBoolean(true).Serialize() produced \"%s\"", s)
	}		
}

func TestNumeric(t *testing.T) {
	if s := toString(NewNumeric(1)); s != "1" {
		t.Errorf ("NewNumeric(1).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewNumeric(3.14159)); s != "3.14159" {
		t.Errorf ("NewNumeric(3.14159).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewNumeric(0.1)); s != "0.1" {
		t.Errorf ("NewNumeric(3.14159).Serialize() produced \"%s\"", s)
	}

	if s := toString(NewNumeric(2147483647)); s != "2147483647" {
		t.Errorf ("NewNumeric(2147483647).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewNumeric(-2147483648)); s != "-2147483648" {
		t.Errorf ("NewNumeric(-2147483648).Serialize() produced \"%s\"", s)
	}		


	if s := toString(NewNumeric(3.403e+38)); s != "3.4028235e+38" {
		t.Errorf ("NewNumeric(3.403e+38).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewNumeric(-3.403e+38)); s != "-3.4028235e+38" {
		t.Errorf ("NewNumeric(-3.403e+38).Serialize() produced \"%s\"", s)
	}		

	// According to PDF spec, anything below +/- 1.175e-38 should
	// be set to 0 in case the reader uses 32 bit floats.  This is
	// not the smallest 32-bit representable 32-bit number, but it
	// is the smallest number without losing precision.
	if s := toString(NewNumeric(1.176e-38)); s != "1.176e-38" {
		t.Errorf ("NewNumeric(1.175e-38).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewNumeric(-1.176e-38)); s != "-1.176e-38" {
		t.Errorf ("NewNumeric(-1.175e-38).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewNumeric(1.175e-38)); s != "0" {
		t.Errorf ("NewNumeric(1.0e-40).Serialize() produced \"%s\"", s)
	}		

	if s := toString(NewNumeric(-1.175e-40)); s != "0" {
		t.Errorf ("NewNumeric(-1.0e-40).Serialize() produced \"%s\"", s)
	}		


}



