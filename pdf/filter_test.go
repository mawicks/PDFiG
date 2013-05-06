package pdf_test

import (
	"github.com/mawicks/PDFiG/pdf"
//	"fmt"
	"io"
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

func testDecoder(t *testing.T, filter pdf.StreamFilter, encoded[]byte, expected []byte) {
	decoder := filter.NewDecoder(bytes.NewReader(encoded))
	decoded := new(bytes.Buffer)

	_,err := io.Copy(decoded, decoder)

	if err != nil {
		t.Errorf(`testDecoder: %s decoder returned error on "%s": "%v"`, filter.Name(), pdf.AsciiFromBytes(encoded), err)
	}

	if !bytes.Equal(decoded.Bytes(), expected) {
		t.Errorf(`testDecoder: %s decoder produced "%s"; expected "%s"`,
			filter.Name(),
			pdf.AsciiFromBytes(decoded.Bytes()),
			pdf.AsciiFromBytes(expected))
	}
}

func testEncoder (t *testing.T, filter pdf.StreamFilter, data[]byte, expected []byte) {
	encoded := pdf.NewBufferCloser()
	encoder := filter.NewEncoder(encoded)
	dataReader := bytes.NewReader(data)

	_,err := io.Copy(encoder, dataReader)
	// Copy doesn't Close()
	encoder.Close()

	if err != nil {
		t.Errorf(`testEncoder: %s encoder returned error on "%s": "%v"`, filter.Name(), pdf.AsciiFromBytes(data), err)
	}

	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Errorf(`testEncoder: %s encoder produced "%s"; expected "%s"`,
			filter.Name(),
			pdf.AsciiFromBytes(encoded.Bytes()),
			pdf.AsciiFromBytes(expected))
	}

}

func testRoundTrip (t *testing.T, filter pdf.StreamFilter, data []byte) {
	encoded := pdf.NewBufferCloser()
	encoder := filter.NewEncoder(encoded)
	_,err := io.Copy(encoder, bytes.NewReader(data))
	if err != nil {
		t.Errorf(`testEncoder: %s encoder returned error on "%s": "%v"`, filter.Name(), pdf.AsciiFromBytes(data), err)
	} else {
		encoder.Close()
	}

	roundTrip := pdf.NewBufferCloser()
	decoder := filter.NewDecoder(bytes.NewReader(encoded.Bytes()))
	_,err = io.Copy(roundTrip, decoder)

	if err != nil {
		t.Errorf(`testEncoder: %s decoder returned error on "%s": "%v"`, filter.Name(), pdf.AsciiFromBytes(encoded.Bytes()), err)
	}

	if !bytes.Equal(data, roundTrip.Bytes()) {
		t.Errorf(`Roundtrip for %s failed on "%s"`,
			filter.Name(),
			pdf.AsciiFromBytes(data))
	}
}


func TestFilters(t *testing.T) {
	// Test some specific cases that are easy enough to type
	testDecoder (t, new(pdf.AsciiHexFilter), []byte("3332313>"), []byte("3210"))
	testDecoder (t, new(pdf.AsciiHexFilter), []byte("33323130>"), []byte("3210"))
	testEncoder (t, new(pdf.AsciiHexFilter), []byte("3210"), []byte("33323130>"))

	// Then make sure random sequences can make the round trip.
	flateFilter := new(pdf.FlateFilter)
	flateFilter.SetCompressionLevel(9)
	for i:=1; i<65536; i*=8 {
		r := randomBytes (i-1)
		testRoundTrip (t, new(pdf.AsciiHexFilter), r)
		testRoundTrip (t, flateFilter, r)
		testRoundTrip (t, new(pdf.LZWFilter), r)
	}
}

