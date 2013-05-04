package pdf_test

import (
	"github.com/mawicks/PDFiG/pdf"
	"bytes"
	"math/rand"
	"testing" )

func randomBytes(n int) []byte {
	result := make([]byte, n, n)
	for i:=0; i<n; i++ {
		result[i] = byte(rand.Int() & 0xff)
	}
	return result
}

func testRoundTrip (t *testing.T, filter pdf.StreamFilter, data []byte) {
	encoded := filter.Encode(data)
	decoded,ok := filter.Decode(encoded)
	if ok && !bytes.Equal(decoded,data) {
		t.Errorf(`%s encode/decode failed round trip on "%s"`,
			filter.Name(),
			pdf.AsciiFromBytes(data))
	}
}

func testDecoder (t *testing.T, filter pdf.StreamFilter, encoded, expected []byte) {
	decoded,ok := filter.Decode(encoded)
	if ok && !bytes.Equal(decoded, expected) {
		t.Errorf(`AsciHexFilter.Decode("%s") produced "%s"; expected "%s"`,
			pdf.AsciiFromBytes(encoded),
			pdf.AsciiFromBytes(decoded),
			pdf.AsciiFromBytes(expected))
	}
}

func TestAsciiHexFilter(t *testing.T) {
	testDecoder (t, pdf.AsciiHexFilter{}, []byte("3332313>"), []byte("3210"))
	testRoundTrip (t, pdf.AsciiHexFilter{}, randomBytes(16))
}

