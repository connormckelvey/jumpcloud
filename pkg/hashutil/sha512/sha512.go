package sha512

import (
	"crypto/sha512"

	"github.com/connormckelvey/jumpcloud/pkg/hashutil"
)

// NewStringWriter returns a new hashutil.HashStringWriter for computing SHA-512
// checksums.
func NewStringWriter() hashutil.HashStringWriter {
	return hashutil.NewStringWriter(sha512.New())
}
