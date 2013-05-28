package pdf

import (
	"io" )

type ExistingPage struct {
	*Dictionary
	reference *Indirect
}

func (ep *ExistingPage) Reader() io.Reader {
	// Try to read the contents on the page.
	if pageStream,ok := ep.GetStream("Contents"); ok {
		return pageStream.Reader()
	}
	if pageStreamArray,ok := ep.GetArray("Contents"); ok {
		n := pageStreamArray.Size()
		readers := make([]io.Reader, n)
		for i:=0; i<n; i++ {
			if streamReference,ok := pageStreamArray.At(i).(*Indirect); ok {
				if stream,ok := streamReference.Dereference().(*Stream); ok {
					readers[i] = stream.Reader()
				}
			}
			if readers[i] == nil {
				return nil
			}
		}
		return io.MultiReader (readers...)
	}
	return nil
}

func (ep *ExistingPage) Rewrite() {
	ep.reference.Write(ep.Dictionary)
}

func (ep *ExistingPage) ensureContentsArray() *Array {
	if pageContentsArray,ok := ep.GetArray("Contents"); ok {
		return pageContentsArray
	}
	if contents, ok := ep.GetIndirect("Contents"); ok {
		contentsArray := NewArray()
		contentsArray.Add(contents)
		ep.Add("Contents", contentsArray)
		return contentsArray
	}
	return nil
}

type ContentsWriter struct {
	*Stream
	streamReference *Indirect
}

func (ep *ExistingPage) PrependWriter(file... File) *ContentsWriter {
	contentsArray := ep.ensureContentsArray()
	newStreamReference := NewIndirect(file...)
	contentsArray.PushFront(newStreamReference)
	return &ContentsWriter{NewStream(),newStreamReference}
}

func (ep *ExistingPage) AppendWriter(file... File) *ContentsWriter {
	contentsArray := ep.ensureContentsArray()
	newStreamReference := NewIndirect(file...)
	contentsArray.Add(newStreamReference)
	return &ContentsWriter{NewStream(),newStreamReference}
}

func (cw *ContentsWriter) Close() {
	cw.streamReference.Write(cw.Stream)
}
