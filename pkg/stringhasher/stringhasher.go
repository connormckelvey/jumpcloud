package stringhasher

import (
	"hash"
)

// StringHasher is an interface that thinly wraps the embedded hash.Hash
// interface. It provides additional methods for hashing strings and returning
// hash checksum data as base64 encoded strings
type StringHasher interface {
	hash.Hash
	WriteString(string) (int, error)
	String() string
}
