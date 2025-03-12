package model

import (
	"crypto/sha256"
)

// Hash returns transaction hash
func (t *Transaction) Hash() Hash {
	return sha256.Sum256([]byte(t.String()))
}
