package model

import (
	"crypto/ed25519"

	"github.com/tyler-smith/go-bip39"
)

type PrivateKey ed25519.PrivateKey
type PublicKey ed25519.PublicKey

// GenerateMnemonic creates random mnemonic
func GenerateMnemonic() (string, error) {
	// 128 bits = 12 words in the mnemonic
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// PrivateKeyFromMnemonic creates ed25519 from mnemonic
func PrivateKeyFromMnemonic(mnemonic string) (*ed25519.PrivateKey, error) {
	seed := bip39.NewSeed(mnemonic, "")
	privateKey := ed25519.NewKeyFromSeed(seed[:32])
	return &privateKey, nil
}
