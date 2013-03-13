package pdf

import "testing"

// Unit tests follow

func TestIsWhite (t *testing.T) {
	for _,b := range []byte("\000\t\n\f\r ") {
		if (!IsWhiteSpace (b)) {
			t.Errorf ("IsWhiteSpace('\\%.3o') failed", b)
		}
	}
}

func TestIsDelimiter (t *testing.T) {
	for _,b := range []byte("%()/<>[]{}") {
		if (!IsDelimiter (b)) {
			t.Errorf ("IsDelimiter('%c') failed", b)
		}
	}
}
