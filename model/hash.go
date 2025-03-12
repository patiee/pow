package model

type Hash [32]byte

func ZeroHash() [32]byte {
	return [32]byte{}
}
