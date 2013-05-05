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

type FlateWriter struct {
	io.WriteCloser
	output io.WriteCloser
}

func NewFlateWriter(writer io.WriteCloser) io.WriteCloser {
	flateWriter,_ := flate.NewWriter(writer,9)
	return &FlateWriter{flateWriter,writer}
}

func (fw *FlateWriter) Close() error {
	if err := fw.WriteCloser.Close(); err != nil {
		return err
	}
	return fw.output.Close()
}

type FlateReader struct {
	io.Reader
}

func NewFlateReader(reader io.Reader) io.Reader {
	flateReader := flate.NewReader(reader)
	return &FlateReader{flateReader}
}

