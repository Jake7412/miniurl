// Package miniurl provides URL shortening utilities.
package miniurl

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash returns a 32-character hex string derived from the MD5 of input.
func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
