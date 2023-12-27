package offsetlimitreader

import (
	"io"
)

// TODO: Test coverage
// TODO: coding style

type OffsetLimitReader struct {
	src    io.Reader
	pos    int64
	start  int64
	length int64
}

func New(r io.Reader, start int64, length int64) *OffsetLimitReader {
	return &OffsetLimitReader{
		src:    r,
		start:  start,
		length: length,
	}
}

func (r *OffsetLimitReader) seekToStart() error {
	chunkSize := int64(1024)
	for {
		l := r.start - r.pos
		if l <= 0 {
			break
		}
		if l > chunkSize {
			l = chunkSize
		}

		// TODO: re-use buffer where possible
		buf := make([]byte, l)
		n, err := r.src.Read(buf)
		r.pos += int64(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OffsetLimitReader) Read(p []byte) (int, error) {
	if r.pos < r.start {
		err := r.seekToStart()
		if err != nil {
			return 0, err
		}
	}

	var n int
	var err error
	left := r.length - (r.pos - r.start)

	if left == 0 {
		return 0, io.EOF
	}

	if len(p) > int(left) {
		// Read request would exceed limit for this reader,
		// so use a temporary buffer to read less data from the source.
		buf := make([]byte, left)
		n, err = r.src.Read(buf)
		if n > 0 {
			copy(p, buf)
		}
	} else {
		n, err = r.src.Read(p)
	}
	r.pos += int64(n)
	return n, err
}
