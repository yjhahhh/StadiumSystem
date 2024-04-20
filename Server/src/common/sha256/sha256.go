package sha256

import (
	"crypto/sha256"
)

func Encode(s string) string {
	b := sha256.Sum256([]byte(s))
	return string(b[:])
}
