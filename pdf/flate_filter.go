package pdf

import ( //"errors"
	"compress/zlib"
//	"fmt"
	"io")

type FlateFilter struct {
	compressionLevel int
}

func (filter *FlateFilter) Name() string {
	return "FlateDecode"
}

func (filter *FlateFilter) SetCompressionLevel(level int) {
	filter.compressionLevel = level
}

func (filter *FlateFilter) NewEncoder(writer io.WriteCloser) io.WriteCloser {
	flateWriter,_ := zlib.NewWriterLevel(writer,filter.compressionLevel)
	return &FlateWriter{flateWriter,writer}
}

func (filter *FlateFilter) NewDecoder(reader io.Reader) io.Reader {
	flateReader,_ := zlib.NewReader(reader)
	return &FlateReader{flateReader}
}

func (filter *FlateFilter) DecodeParms(file ...File) Object {
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


