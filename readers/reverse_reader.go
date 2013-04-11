package readers

import "io"

type ReverseReader struct {
	readSeeker io.ReadSeeker
}

func NewReverseReader (readSeeker io.ReadSeeker) io.Reader {
	readSeeker.Seek(0,2)
	return &ReverseReader{readSeeker}
}

func (d *ReverseReader) Read (b []byte) (n int, err error) {
	position,_ := d.readSeeker.Seek(0,1)

	if position == 0 {
		return 0, io.EOF
	}

	var min int
	if position < int64(len(b)) {
		min = int(position)
	} else {
		min = len(b)
	}

	_,err = d.readSeeker.Seek(int64(-min),1)

	if err == nil {
		readBuffer := make([]byte, min, min)
		n,err = d.readSeeker.Read (readBuffer)

		if err == nil {
			_,err = d.readSeeker.Seek(int64(-n), 1)
			for i,v := range readBuffer {
				b[min-i-1] = v
			}
		}
	}
	return n,err
}