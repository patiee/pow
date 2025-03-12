package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
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

// SerializeBlock serializes the Block into a byte slice
func (block *Block) SerializeBlock() ([]byte, error) {
	var buf bytes.Buffer

	// Version
	version := config.BlockVersion
	if err := binary.Write(&buf, binary.LittleEndian, version); err != nil {
		return nil, fmt.Errorf("failed to write version: %v", err)
	}

	// Previous Hash
	prevHash := make([]byte, 32)
	copy(prevHash, []byte(block.PreviousHash))
	if err := binary.Write(&buf, binary.LittleEndian, prevHash); err != nil {
		return nil, fmt.Errorf("failed to write previous hash: %v", err)
	}

	// Merkle Root
	merkleRoot := make([]byte, 32)
	copy(merkleRoot, []byte(block.MerkleRoot))
	if err := binary.Write(&buf, binary.LittleEndian, merkleRoot); err != nil {
		return nil, fmt.Errorf("failed to write merkle root: %v", err)
	}

	// Timestamp
	timestamp := int64(block.Timestamp)
	if err := binary.Write(&buf, binary.LittleEndian, timestamp); err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %v", err)
	}

	// Difficulty
	difficulty := new(big.Int)
	difficulty, ok := difficulty.SetString(block.Difficulty, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse difficulty")
	}
	difficultyBytes := difficulty.Bytes()

	difficultyLength := int32(len(difficultyBytes))
	if err := binary.Write(&buf, binary.LittleEndian, difficultyLength); err != nil {
		return nil, fmt.Errorf("failed to write difficulty length: %v", err)
	}

	if err := binary.Write(&buf, binary.LittleEndian, difficultyBytes); err != nil {
		return nil, fmt.Errorf("failed to write difficulty value: %v", err)
	}

	// Nonce
	nonce := int64(block.Nonce)
	if err := binary.Write(&buf, binary.LittleEndian, nonce); err != nil {
		return nil, fmt.Errorf("failed to write nonce: %v", err)
	}

	// Transactions
	for _, tx := range block.Transactions {
		txHash := make([]byte, 32)
		copy(txHash, tx)
		if err := binary.Write(&buf, binary.LittleEndian, txHash); err != nil {
			return nil, fmt.Errorf("failed to write transaction hash: %v", err)
		}
	}

	// Return serialized block as bytes
	return buf.Bytes(), nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	buf := bytes.NewReader(data)
	block := &Block{}

	// Read Version
	var version int32
	if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %v", err)
	}

	// Read Previous Hash
	prevHash := make([]byte, 32)
	if err := binary.Read(buf, binary.LittleEndian, &prevHash); err != nil {
		return nil, fmt.Errorf("failed to read previous hash: %v", err)
	}
	block.PreviousHash = string(bytes.Trim(prevHash, "\x00"))

	// Read Merkle Root
	merkleRoot := make([]byte, 32)
	if err := binary.Read(buf, binary.LittleEndian, &merkleRoot); err != nil {
		return nil, fmt.Errorf("failed to read merkle root: %v", err)
	}
	block.MerkleRoot = string(bytes.Trim(merkleRoot, "\x00"))

	// Read Timestamp
	var timestamp int64
	if err := binary.Read(buf, binary.LittleEndian, &timestamp); err != nil {
		return nil, fmt.Errorf("failed to read timestamp: %v", err)
	}
	block.Timestamp = timestamp

	// Read Difficulty
	var difficultyLength int32
	if err := binary.Read(buf, binary.LittleEndian, &difficultyLength); err != nil {
		return nil, fmt.Errorf("failed to read difficulty length: %v", err)
	}

	difficultyBytes := make([]byte, difficultyLength)
	if err := binary.Read(buf, binary.LittleEndian, &difficultyBytes); err != nil {
		return nil, fmt.Errorf("failed to read difficulty value: %v", err)
	}

	difficulty := new(big.Int).SetBytes(difficultyBytes)
	block.Difficulty = difficulty.String()

	// Read Nonce
	var nonce int64
	if err := binary.Read(buf, binary.LittleEndian, &nonce); err != nil {
		return nil, fmt.Errorf("failed to read nonce: %v", err)
	}
	block.Nonce = nonce

	// Read Transactions
	block.Transactions = [][]byte{}
	for buf.Len() > 0 {
		txHash := make([]byte, 32)
		if err := binary.Read(buf, binary.LittleEndian, &txHash); err != nil {
			return nil, fmt.Errorf("failed to read transaction hash: %v", err)
		}
		block.Transactions = append(block.Transactions, bytes.Trim(txHash, "\x00"))
	}

	return block, nil
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
