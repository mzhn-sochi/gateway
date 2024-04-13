package file

import "io"

type Reader interface {
	io.Reader
	Size() int64
	ContentType() string
}

type FileReader struct {
	size        int64
	contentType string
	io.Reader
}

func (r *FileReader) ContentType() string {
	return r.contentType
}

func NewReader(r io.Reader, size int64, contentType string) *FileReader {
	return &FileReader{
		size:        size,
		Reader:      r,
		contentType: contentType,
	}
}

func (r *FileReader) Size() int64 {
	return r.size
}
