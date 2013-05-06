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
	testDecoder (t, new(pdf.AsciiHexFilter), []byte("3332313>"), []byte("3210"))
	testDecoder (t, new(pdf.AsciiHexFilter), []byte("33323130>"), []byte("3210"))
	testEncoder (t, new(pdf.AsciiHexFilter), []byte("3210"), []byte("33323130>"))
	testEncoder (t, new(pdf.FlateFilter), []byte("foo"), []byte("foo"))

	testRoundTrip (t, new(pdf.AsciiHexFilter), randomBytes(16))
	testRoundTrip (t, new(pdf.AsciiHexFilter), randomBytes(8))
	testRoundTrip (t, new(pdf.AsciiHexFilter), randomBytes(4))
	testRoundTrip (t, new(pdf.AsciiHexFilter), randomBytes(1))
	testRoundTrip (t, new(pdf.AsciiHexFilter), randomBytes(0))

	testRoundTrip (t, new(pdf.FlateFilter), randomBytes(16))

	testRoundTrip (t, new(pdf.LZWFilter), []byte("test test test"))
	testRoundTrip (t, new(pdf.LZWFilter), randomBytes(16))

//	testDecoder (t, pdf.AsciiHexFilter{}, []byte("3332313>"), []byte("3210"))
//	testDecoder (t, pdf.NewAsciiHexReader, []byte("33323130"), []byte("3210"))
//	testRoundTrip (t, pdf.AsciiHexFilter{}, randomBytes(16))
}

