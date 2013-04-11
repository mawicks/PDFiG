package pdf

import "fmt"

var unicodeToPDFDoc map[rune]byte

func init() {
	var mappings []struct { rune; byte } =  []struct {rune; byte}  {
		{'\u0000', 0x00}, {'\u0001', 0x00}, {'\u0002', 0x00}, {'\u0003', 0x00},
		{'\u0004', 0x00}, {'\u0005', 0x00}, {'\u0006', 0x00}, {'\u0007', 0x00},
		{'\u0008', 0x00}, {'\u000b', 0x00}, {'\u000c', 0x00}, {'\u000e', 0x00},
		{'\u000f', 0x00}, {'\u0010', 0x00}, {'\u0011', 0x00}, {'\u0012', 0x00},
		{'\u0013', 0x00}, {'\u0014', 0x00}, {'\u0015', 0x00}, {'\u0016', 0x00},
		{'\u0017', 0x00}, {'\u02d8', 0x18}, {'\u02c7', 0x19}, {'\u02d9', 0x1b},
		{'\u02dd', 0x1c}, {'\u02db', 0x1d}, {'\u02da', 0x1e}, {'\u02dc', 0x1f},
		{'\u2022', 0x80}, {'\u2020', 0x81}, {'\u2021', 0x82}, {'\u2026', 0x83},
		{'\u2014', 0x84}, {'\u2013', 0x85}, {'\u0192', 0x86}, {'\u2044', 0x87},
		{'\u2039', 0x88}, {'\u203a', 0x89}, {'\u2212', 0x8a}, {'\u2030', 0x8b},
		{'\u201e', 0x8c}, {'\u201c', 0x8d}, {'\u201d', 0x8e}, {'\u2018', 0x8f},
		{'\u2019', 0x90}, {'\u201a', 0x91}, {'\u2122', 0x92}, {'\ufb01', 0x93},
		{'\ufb02', 0x94}, {'\u0141', 0x95}, {'\u0152', 0x96}, {'\u0160', 0x97},
		{'\u0178', 0x98}, {'\u017d', 0x99}, {'\u0131', 0x9a}, {'\u0142', 0x9b},
		{'\u0153', 0x9c}, {'\u0161', 0x9d}, {'\u017e', 0x9e}, {'\u009f', 0x00},
		{'\u20ac', 0xa0}, {'\u00ad', 0x00} }

	unicodeToPDFDoc = make(map[rune]byte,82)
	for _,v := range mappings {
		_,exists := unicodeToPDFDoc[v.rune]
		if (exists) {
			panic (fmt.Sprintf("Duplicate value (%v) in PDFDocEncoding mappings", v.rune))
		}
		unicodeToPDFDoc[v.rune] = v.byte

		if (v.byte != 0x00) {
			_,exists = unicodeToPDFDoc[rune(v.byte)]
			if (exists) {
				panic (fmt.Sprintf("Duplicate value (%x) in PDFDocEncoding mappings", v.byte))
			}
			unicodeToPDFDoc[rune(v.byte)] = 0x00
		}
	}
}

func PDFDocEncoding (s []rune) ([]byte,bool) {
	ok := true;
	result := make([]byte,0, len(s))
	for _,rune := range s {
		subst,exists := unicodeToPDFDoc[rune]
		if (exists) {
			if subst == 0x00 {
				ok = false
			} else {
				result = append(result, subst)
			}
		} else {
			if rune >= 0x100 {
				ok = false
			} else {
				result = append(result, byte(rune))
			}
		}
	}
	return result,ok
}