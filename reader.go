package rkcel

import (
	"io"
)

type CacheReader struct {
	src  io.Reader
	buf  []byte
	pos  int
	done bool
}

func NewCacheReader(src io.Reader) *CacheReader {
	return &CacheReader{
		src:  src,
		buf:  make([]byte, 0),
		pos:  0,
		done: false,
	}
}

func (r *CacheReader) Read(p []byte) (int, error) {
	if r.pos < len(r.buf) {
		n := copy(p, r.buf[r.pos:])
		r.pos += n
		return n, nil
	}

	if r.done {
		return 0, io.EOF
	}

	n, err := r.src.Read(p)
	if n > 0 {
		r.buf = append(r.buf, p[:n]...)
		r.pos += n
	}
	if err == io.EOF {
		r.done = true
	}
	return n, err
}

func (r *CacheReader) Rewind() {
	r.pos = 0
}
