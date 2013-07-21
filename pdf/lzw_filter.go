package pdf

import ( //"errors"
	"compress/lzw"
//	"fmt"
	"io")

type LZWFilter struct {
}

const ( lzwDecoderName = "LZWDecode" )

func init () {
	RegisterFilterFactoryFactory(lzwDecoderName,
		func(d ProtectedDictionary) StreamFilterFactory {
			if d != nil {
				if v,ok := d.GetInt("EarlyChange"); ok && v == 0 {
					return new(LZWFilter) }
			}
			return nil
		})
}

func (filter LZWFilter) Name() string {
	return lzwDecoderName
}

func (filter LZWFilter) NewEncoder(writer io.WriteCloser) io.WriteCloser {
	lzwWriter := lzw.NewWriter(writer,lzw.MSB, 8)
	return &LZWWriter{lzwWriter,writer}
}

func (filter LZWFilter) NewDecoder(reader io.Reader) io.Reader {
	lzwReader := lzw.NewReader(reader,lzw.MSB, 8)
	return &LZWReader{lzwReader}
}

func (filter LZWFilter) DecodeParms(file ...File) Object {
	d := NewDictionary()
	// This parameter is necessary due to an incompability between
	// the Go LZW library and the default value in the PDF spec.
	// Unfortunately, this means we cannot decode PDF created with
	// the default value.
	d.Add ("EarlyChange", NewIntNumeric(0))
	return d
}

type LZWWriter struct {
	io.WriteCloser
	output io.WriteCloser
}

func (w *LZWWriter) Write(buffer []byte) (n int, err error) {
	// LZW Write() appears to have an off-by-one error.  It
	// reports one less than the number of bytes actually written.
	n,err = w.WriteCloser.Write(buffer)
	if err == nil {
		n = len(buffer)
	}
	return n,err
}

func (w *LZWWriter) Close() error {
	if err := w.WriteCloser.Close(); err != nil {
		return err
	}
	return w.output.Close()
}

type LZWReader struct {
	io.Reader
}


