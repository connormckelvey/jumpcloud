package hashutil

import (
	"hash"
	"io"
)

// HashStringWriter is an interface that thinly the hash.Hash interface.
// It provides and additional convenience methods for working with the
// embedded hash.Hash
type HashStringWriter interface {
	hash.Hash
	WriteString(string) (int, error)
}

type hashStringWriter struct {
	hash.Hash
}

func NewStringWriter(h hash.Hash) HashStringWriter {
	return &hashStringWriter{
		Hash: h,
	}
}

func (h *hashStringWriter) WriteString(s string) (n int, err error) {
	return io.WriteString(h.Hash, s)
}
