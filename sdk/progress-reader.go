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

func (r *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = r.Reader.Read(b)
	if r.OnReadProgress != nil {
		r.OnReadProgress(int64(n))
	}
	return
}

func (r *ProgressReader) Close() error {
	switch v := r.Reader.(type) {
	case io.ReadCloser:
		return v.Close()
	case io.Closer:
		return v.Close()
	default:
		return nil
	}
}
