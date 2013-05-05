package pdf

import ("errors"
//	"fmt"
	"io")


type AsciiHexFilter struct {
}

func (fileter AsciiHexFilter) Name() string {
	return "ASCIIHexDecode"
}

type AsciiHexWriter struct {
	writer io.WriteCloser
	count int
}

func NewAsciiHexWriter(writer io.WriteCloser) io.WriteCloser {
	return &AsciiHexWriter{writer,0}
}

func (ahw *AsciiHexWriter) Write(buffer []byte) (n int, err error) {
	var m int
	for n=0; n<len(buffer) && err == nil; n++ {
		m,err = ahw.writer.Write([]byte{HexDigit(buffer[n]/16),HexDigit(buffer[n]%16)})
		ahw.count += m
		if ahw.count != 0 && ahw.count%40 == 0 && err == nil {
			m,err = ahw.Write([]byte{'\n'})
			ahw.count += m
		}
	}
	return n,err
}

func (ahw *AsciiHexWriter) Close() error {
	if _,err := ahw.writer.Write([]byte{'>'}); err != nil {
		return err
	}
	return ahw.writer.Close()
}

type AsciiHexReader struct {
	reader io.Reader
	err error
}

func NewAsciiHexReader(reader io.Reader) io.Reader {
	return &AsciiHexReader{reader,nil}
}

func (ahr *AsciiHexReader) Read(buffer []byte) (n int, err error) {
	var (
		m,count int
		nextByte byte)

	next := make([]byte, 1)
	for n=0; n<len(buffer) && ahr.err == nil; {
		advance := func () {
			buffer[n] = nextByte
			n += 1
			nextByte = 0
		}
		m,err = ahr.reader.Read (next)
		switch {
		case m == 1:
			switch {
			case IsHexDigit(next[0]):
				nextByte = nextByte*16 + ParseHexDigit(next[0])
				count += 1
				if count % 2 == 0 {
					advance()
				}
			case next[0] == '>':
				nextByte = nextByte*16
				if count % 2 == 1 {
					advance()
				}
				ahr.err = io.EOF
			case IsWhiteSpace(next[0]):
				// Do nothing
			default:
				ahr.err = errors.New("AsciiHexReader:  Invalid character")
			}
		default:
			if err == io.EOF {
				ahr.err = errors.New(`Unexpected end of stream (no trailing ">")`)
			}
		}
	}
	return n,ahr.err
}
