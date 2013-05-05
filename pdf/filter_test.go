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

/*
func testRoundTrip (t *testing.T, filter pdf.StreamFilter, data []byte) {
	encoded := filter.Encode(data)
	decoded,ok := filter.Decode(encoded)
	if ok && !bytes.Equal(decoded,data) {
		t.Errorf(`%s encode/decode failed round trip on "%s"`,
			filter.Name(),
			pdf.AsciiFromBytes(data))
	}
}
*/

/*
func testDecoder (t *testing.T, filter pdf.StreamFilter, encoded, expected []byte) {
	decoded,ok := filter.Decode(encoded)
	if ok && !bytes.Equal(decoded, expected) {
		t.Errorf(`AsciHexFilter.Decode("%s") produced "%s"; expected "%s"`,
			pdf.AsciiFromBytes(encoded),
			pdf.AsciiFromBytes(decoded),
			pdf.AsciiFromBytes(expected))
	}
}
*/

func testDecoder(t *testing.T, decoderFactory func(io.Reader) io.Reader, encoded[]byte, expected []byte) {
	decoder := decoderFactory(bytes.NewReader(encoded))
	decoded := new(bytes.Buffer)

	_,err := io.Copy(decoded, decoder)

	if err != nil {
		t.Errorf(`testDecoder: decoder returned error on "%s": "%v"`, pdf.AsciiFromBytes(encoded), err)
	}

	if !bytes.Equal(decoded.Bytes(), expected) {
		t.Errorf(`testDecoder produced "%s"; expected "%s"`,
			pdf.AsciiFromBytes(decoded.Bytes()),
			pdf.AsciiFromBytes(expected))
	}
}

func testEncoder (t *testing.T, encoderFactory func(io.WriteCloser) io.WriteCloser, data[]byte, expected []byte) {
	encoded := pdf.NewBufferCloser()
	encoder := encoderFactory(encoded)
	dataReader := bytes.NewReader(data)

	_,err := io.Copy(encoder, dataReader)
	// Copy doesn't Close()
	encoder.Close()

	if err != nil {
		t.Errorf(`testEncoder: encoder returned error on "%s": "%v"`, pdf.AsciiFromBytes(data), err)
	}

	if !bytes.Equal(encoded.Bytes(), expected) {
		t.Errorf(`testEncoder produced "%s"; expected "%s"`,
			pdf.AsciiFromBytes(encoded.Bytes()),
			pdf.AsciiFromBytes(expected))
	}

}

func testRoundTrip (t *testing.T, encoderFactory func(io.WriteCloser) io.WriteCloser, decoderFactory func(io.Reader) io.Reader, data []byte) {
	encoded := pdf.NewBufferCloser()
	encoder := encoderFactory(encoded)
	_,err := io.Copy(encoder, bytes.NewReader(data))
	if err != nil {
		t.Errorf(`testEncoder: encoder returned error on "%s": "%v"`, pdf.AsciiFromBytes(data), err)
	} else {
		encoder.Close()
	}

	roundTrip := pdf.NewBufferCloser()
	decoder := decoderFactory(bytes.NewReader(encoded.Bytes()))
	_,err = io.Copy(roundTrip, decoder)

	if err != nil {
		t.Errorf(`testEncoder: decoder returned error on "%s": "%v"`, pdf.AsciiFromBytes(encoded.Bytes()), err)
	}

	if !bytes.Equal(data, roundTrip.Bytes()) {
		t.Errorf(`Roundtrip failed on "%s"`,
			pdf.AsciiFromBytes(data))
	}
}


func TestAsciiHexFilter(t *testing.T) {
	testDecoder (t, pdf.NewAsciiHexReader, []byte("3332313>"), []byte("3210"))
	testDecoder (t, pdf.NewAsciiHexReader, []byte("33323130>"), []byte("3210"))

	testEncoder (t, pdf.NewAsciiHexWriter, []byte("3210"), []byte("33323130>"))

	testRoundTrip (t, pdf.NewAsciiHexWriter, pdf.NewAsciiHexReader, randomBytes(16))
	testRoundTrip (t, pdf.NewAsciiHexWriter, pdf.NewAsciiHexReader, randomBytes(8))
	testRoundTrip (t, pdf.NewAsciiHexWriter, pdf.NewAsciiHexReader, randomBytes(4))
	testRoundTrip (t, pdf.NewAsciiHexWriter, pdf.NewAsciiHexReader, randomBytes(1))
	testRoundTrip (t, pdf.NewAsciiHexWriter, pdf.NewAsciiHexReader, randomBytes(0))

	testRoundTrip (t, pdf.NewFlateWriter, pdf.NewFlateReader, randomBytes(1600))

//	testDecoder (t, pdf.AsciiHexFilter{}, []byte("3332313>"), []byte("3210"))
//	testDecoder (t, pdf.NewAsciiHexReader, []byte("33323130"), []byte("3210"))
//	testRoundTrip (t, pdf.AsciiHexFilter{}, randomBytes(16))
}

