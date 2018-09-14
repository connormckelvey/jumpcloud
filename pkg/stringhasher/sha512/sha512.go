package sha512

import (
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"io"

	sh "github.com/connormckelvey/jumpcloud/pkg/stringhasher"
)

type sha512StringHasher struct {
	hash hash.Hash
}

// New returns a new stringhasher.StringHasher computing the SHA-512 checksum.
func New() sh.StringHasher {
	return &sha512StringHasher{
		hash: sha512.New(),
	}
}

func (h *sha512StringHasher) Write(b []byte) (n int, err error) {
	return h.hash.Write(b)
}

func (h *sha512StringHasher) WriteString(s string) (n int, err error) {
	return io.WriteString(h.hash, s)
}

func (h *sha512StringHasher) Sum(b []byte) []byte {
	return h.hash.Sum(b)
}

func (h *sha512StringHasher) Reset() {
	h.hash.Reset()
}

func (h *sha512StringHasher) Size() int {
	return h.hash.Size()
}

func (h *sha512StringHasher) BlockSize() int {
	return h.hash.BlockSize()
}

func (h *sha512StringHasher) String() string {
	return base64.StdEncoding.EncodeToString(h.hash.Sum(nil))
}
