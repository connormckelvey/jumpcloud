package sha512

import (
	"crypto/sha512"

	"github.com/connormckelvey/jumpcloud/pkg/hashutil"
)

// New returns a new hashutil.HashStringWriter computing the SHA-512 checksum.
func NewStringWriter() hashutil.HashStringWriter {
	return hashutil.NewStringWriter(sha512.New())
}
