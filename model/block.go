package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/patiee/pow/config"
)

// Hash returns block hash as big number
func (b *Block) Hash() *big.Int {
	hash := sha256.Sum256([]byte(b.String()))
	return new(big.Int).SetBytes(hash[:])
}

// HashHex returns block hash as hex
func (b *Block) HashHex() string {
	hash := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(hash[:])
}

// MerkleRoot computes the Merkle root from a list of transactions
func MerkleRoot(txs []*Transaction) string {
	if len(txs) == 0 {
		return ""
	}

	hashes := make([]string, len(txs))
	for i, tx := range txs {
		hashes[i] = tx.HashHex()
	}

	for len(hashes) > 1 {
		tempHashes := []string{}
		for i := 0; i < len(hashes); i += 2 {
			if i+1 >= len(hashes) {
				hashes = append(hashes, hashes[i])
			}
			combined := hashes[i] + hashes[i+1]
			hash := sha256.Sum256([]byte(combined))
			tempHashes = append(tempHashes, hex.EncodeToString(hash[:]))
		}
		hashes = tempHashes
	}

	return hashes[0]
}

func calculateTarget(difficulty string) (*big.Int, error) {
	difficultyBig, ok := new(big.Int).SetString(difficulty, 10)
	if !ok {
		return nil, fmt.Errorf("invalid difficulty: %s", difficulty)
	}
	target := new(big.Int).Div(config.MaxDifficultyTarget, difficultyBig) // Target = maxTarget / difficulty
	return target, nil
}

func (b *Block) MineBlock() (string, error) {
	target, err := calculateTarget(b.Difficulty)
	if err != nil {
		return "", err
	}

	for {
		hash := b.Hash()
		if hash.Cmp(target) == -1 {
			fmt.Printf("Block mined! Nonce: %d, Hash: %s\n", b.Nonce, hash)
			return fmt.Sprintf("%x", hash), nil
		}
		b.Nonce++
	}
}
