package pdf

import ( //"errors"
	"compress/flate"
//	"fmt"
	"io")

type FlateFilter struct {
}

func (filter FlateFilter) Name() string {
	return "FlateDecode"
}

func (filter FlateFilter) NewEncoder(writer io.WriteCloser) io.WriteCloser {
	flateWriter,_ := flate.NewWriter(writer,9)
	return &FlateWriter{flateWriter,writer}
}

func (filter FlateFilter) NewDecoder(reader io.Reader) io.Reader {
	flateReader := flate.NewReader(reader)
	return &FlateReader{flateReader}
}

func (filter FlateFilter) DecodeParms(file ...File) Object {
	return NewNull()
}

type FlateWriter struct {
	io.WriteCloser
	underlyingWriter io.WriteCloser
}

func (fw *FlateWriter) Close() error {
	if err := fw.WriteCloser.Close(); err != nil {
		return err
	}
	return fw.underlyingWriter.Close()
}

type FlateReader struct {
	io.Reader
}


