package model

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// Hash returns transaction hash as big number
func (t *Transaction) Hash() *big.Int {
	hash := sha256.Sum256([]byte(t.String()))
	return new(big.Int).SetBytes(hash[:])
}

// HashHex returns transaction hash as hex
func (t *Transaction) HashHex() string {
	hash := sha256.Sum256([]byte(t.String()))
	return hex.EncodeToString(hash[:])
}
