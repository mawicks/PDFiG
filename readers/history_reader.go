package readers

import "io"

// ByteScannerReader is a combination of io.ByteScanner and io.Reader.
// The UnreadByte() method is useful in lexical analyzers and so is
// io.Reader.  None of the interfaces supplied with Go appear to have
// both of these methods.
type ByteScannerReader interface {
	io.ByteScanner
	io.Reader }

// HistoryReader is a decorator for a ByteScannerReader interface that
// saves the history of the previous "n" reads in a circular buffer.
// It is useful, for example, with lexical scanners so that when an
// error occurs, the HistoryReader can provide the last "n" bytes
// leading up to the error.
type HistoryReader struct {
	reader ByteScannerReader
	buffer []byte
	end, size uint
	capacity uint
}

// NewHistoryReader() creates a new HistoryReader from a
// ByteScannerReader with a circular buffer of the requested capacity.
func NewHistoryReader (reader ByteScannerReader,capacity uint) *HistoryReader {
	return &HistoryReader{
		reader: reader,
		buffer: make([]byte, capacity),
		end: 0,
		size: 0,
		capacity: capacity}
}

// GetHistory() returns the contents of the circular history buffer.
// It is the only method added to HistoryReader that distinguishes it
// from ByteScannerReader.
func (d *HistoryReader) GetHistory() []byte {
	history := make([]byte,d.size)
	beginning := d.end + d.capacity - d.size
	for i:=uint(0); i<d.size; i++ {
		history[i] = d.buffer[(beginning+i)%d.capacity]
	}
	return history
}

func (d *HistoryReader) Read(b []byte) (n int, err error) {
	n,err = d.reader.Read(b)
	for i:=0; i<n; i++ {
		d.buffer[d.end] = b[i]
		d.end = (d.end+1) % d.capacity
	}
	d.size += uint(n)
	if (d.size > d.capacity) {
		d.size = d.capacity
	}
	return
}

func (d *HistoryReader) ReadByte() (b byte, err error) {
	b,err = d.reader.ReadByte()
	if err == nil {
		d.buffer[d.end] = b
		d.end = (d.end+1) % d.capacity
		d.size += 1
		if (d.size > d.capacity) {
			d.size = d.capacity
		}
	}
	return
}

func (d *HistoryReader) UnreadByte() (err error) {
	err = d.reader.UnreadByte()
	if (err == nil) {
		d.end = (d.end+d.capacity-1) % d.capacity
		if (d.size > 0) {
			d.size = d.size - 1
		}
	}
	return
}
