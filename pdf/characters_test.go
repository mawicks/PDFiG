package pdf

import "testing"
import "strings"

// Unit tests follow

func removeIf (f func (b byte) bool, s string) string {
	return strings.Map (
		func (r rune) rune {
			if f(byte(r)) {
				return r
			}
			return -1
		},
		s)
}

func TestCharacterSets (t *testing.T) {
	arrayOfAllLatin1Runes := make([]rune, 256, 256)

	for i,_ := range arrayOfAllLatin1Runes {
		arrayOfAllLatin1Runes[i] = rune(i)
	}

	stringOfAllLatin1Runes := string(arrayOfAllLatin1Runes)

	justDelimiters := removeIf (IsDelimiter, stringOfAllLatin1Runes)
	justWhite := removeIf (IsWhiteSpace, stringOfAllLatin1Runes)

	if justDelimiters != "%()/<>[]{}" {
		t.Errorf ("Incorrect delimiter character list: \"%s\"", justDelimiters)
	}

	if justWhite != "\000\t\n\f\r " {
		t.Errorf ("Incorrect delimiter character list: \"%s\"", justWhite)
	}
}
