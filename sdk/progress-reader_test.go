package sdk

import (
	"bytes"
	"io"
	"testing"
)

func TestProgressReaderString(t *testing.T) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	b := []byte(s)
	total := int64(0)
	reader := &ProgressReader{
		Reader: bytes.NewReader(b),
		OnReadProgress: func(r int64) {
			total += r
		},
	}
	var res bytes.Buffer
	buf := make([]byte, 4)
	for {
		n, err := reader.Read(buf)
		res.Write(buf[:n])
		if err == io.EOF {
			break
		}
	}
	checkTestString(t, s, res.String())
	checkTestInt(t, len(s), int(total))
}

func TestProgressReaderNil(t *testing.T) {
	total := int64(0)
	reader := &ProgressReader{
		Reader: bytes.NewReader(nil),
		OnReadProgress: func(r int64) {
			total += r
		},
	}
	var res bytes.Buffer
	buf := make([]byte, 4)
	for {
		n, err := reader.Read(buf)
		res.Write(buf[:n])
		if err == io.EOF {
			break
		}
	}
	checkTestString(t, "", res.String())
	checkTestInt(t, 0, int(total))
}
