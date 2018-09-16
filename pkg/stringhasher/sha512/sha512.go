package sha512

import (
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"io"

	sh "github.com/connormckelvey/jumpcloud/pkg/stringhasher"
)

type sha512StringHasher struct {
	hash.Hash
}

func HashAndBase64Encode(str string) string {
	hasher := New()
	hasher.WriteString(str)
	return hasher.String()
}

// New returns a new stringhasher.StringHasher computing the SHA-512 checksum.
func New() sh.StringHasher {
	return &sha512StringHasher{
		sha512.New(),
	}
}

func (h *sha512StringHasher) WriteString(s string) (n int, err error) {
	return io.WriteString(h.Hash, s)
}

func (h *sha512StringHasher) String() string {
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
