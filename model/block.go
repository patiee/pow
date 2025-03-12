package model

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/patiee/pow/config"
)

// CalculateHash computes the block's hash based on its contents
func (b *Block) CalculateHash() *big.Int {
	hash := sha256.Sum256([]byte(b.String()))
	return new(big.Int).SetBytes(hash[:])
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
		hash := b.CalculateHash()
		if hash.Cmp(target) == -1 {
			fmt.Printf("Block mined! Nonce: %d, Hash: %s\n", b.Nonce, hash)
			return fmt.Sprintf("%x", hash), nil
		}
		b.Nonce++
	}
}
