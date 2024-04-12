package file

import "io"

type Reader interface {
	io.Reader
	Size() int64
}

type FileReader struct {
	size int64
	io.Reader
}

func NewReader(r io.Reader, size int64) *FileReader {
	return &FileReader{
		size:   size,
		Reader: r,
	}
}

func (r *FileReader) Size() int64 {
	return r.size
}
