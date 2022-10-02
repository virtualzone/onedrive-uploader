package sdk

import (
	"bytes"
	"io"
	"strings"
)

type ProgressReader struct {
	Reader         io.Reader
	OnReadProgress func(r int64)
}

func (r *ProgressReader) Len() int {
	switch v := r.Reader.(type) {
	case *bytes.Buffer:
		return v.Len()
	case *bytes.Reader:
		return v.Len()
	case *strings.Reader:
		return v.Len()
	default:
		return 0
	}
}

/*
	func (r *ProgressReader) Size() int64 {
		return r.Reader.Size()
	}
*/
func (r *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = r.Reader.Read(b)
	if r.OnReadProgress != nil {
		r.OnReadProgress(int64(n))
	}
	return
}

/*
	func (r *ProgressReader) ReadAt(b []byte, off int64) (n int, err error) {
		return r.Reader.ReadAt(b, off)
	}

	func (r *ProgressReader) ReadByte() (byte, error) {
		return r.Reader.ReadByte()
	}

	func (r *ProgressReader) UnreadByte() error {
		return r.Reader.UnreadByte()
	}

	func (r *ProgressReader) ReadRune() (ch rune, size int, err error) {
		return r.Reader.ReadRune()
	}

	func (r *ProgressReader) UnreadRune() error {
		return r.Reader.UnreadRune()
	}

	func (r *ProgressReader) Seek(offset int64, whence int) (int64, error) {
		return r.Reader.Seek(offset, whence)
	}

	func (r *ProgressReader) WriteTo(w io.Writer) (n int64, err error) {
		return r.Reader.WriteTo(w)
	}

	func (r *ProgressReader) Reset(b []byte) {
		r.Reader.Reset(b)
	}
*/
func (r *ProgressReader) Close() error {
	switch v := r.Reader.(type) {
	case io.ReadCloser:
		return v.Close()
	case io.Closer:
		return v.Close()
	default:
		return nil
	}
	//return nil
}
