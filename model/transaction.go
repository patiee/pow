package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// Hash returns transaction hash
func (t *Transaction) Hash() Hash {
	return sha256.Sum256([]byte(t.String()))
}

// Serialize the transaction into a byte slice
func (t *Transaction) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Idx
	if err := binary.Write(&buf, binary.LittleEndian, t.Idx); err != nil {
		return nil, fmt.Errorf("failed to write index: %v", err)
	}

	// Sender
	sender := make([]byte, 32)
	copy(sender, t.Sender)
	if err := binary.Write(&buf, binary.LittleEndian, sender); err != nil {
		return nil, fmt.Errorf("failed to write sender: %v", err)
	}

	// Receiver
	receiver := make([]byte, 32)
	copy(receiver, t.Receiver)
	if err := binary.Write(&buf, binary.LittleEndian, receiver); err != nil {
		return nil, fmt.Errorf("failed to write receiver: %v", err)
	}

	// Amount
	amount := make([]byte, 32)
	copy(amount, []byte(t.Amount))
	if err := binary.Write(&buf, binary.LittleEndian, amount); err != nil {
		return nil, fmt.Errorf("failed to write amount: %v", err)
	}

	// Signature
	signature := make([]byte, 32)
	copy(signature, t.Signature)
	if err := binary.Write(&buf, binary.LittleEndian, signature); err != nil {
		return nil, fmt.Errorf("failed to write signature: %v", err)
	}

	// Nonce
	nonce := int64(t.Nonce)
	if err := binary.Write(&buf, binary.LittleEndian, nonce); err != nil {
		return nil, fmt.Errorf("failed to write nonce: %v", err)
	}

	// Return serialized transaction as bytes
	return buf.Bytes(), nil
}

// DeserializeTransaction from a byte slice
func DeserializeTransaction(data []byte) (*Transaction, error) {
	buf := bytes.NewReader(data)
	transaction := &Transaction{}

	// Read Idx
	var idx int64
	if err := binary.Read(buf, binary.LittleEndian, &idx); err != nil {
		return nil, fmt.Errorf("failed to read idx: %v", err)
	}

	// Read Sender
	sender := make([]byte, 32)
	if err := binary.Read(buf, binary.LittleEndian, sender); err != nil {
		return nil, fmt.Errorf("failed to read sender: %v", err)
	}
	transaction.Sender = sender

	// Read Receiver
	receiver := make([]byte, 32)
	if err := binary.Read(buf, binary.LittleEndian, &receiver); err != nil {
		return nil, fmt.Errorf("failed to read receiver: %v", err)
	}
	transaction.Receiver = receiver

	// Read Amount
	amount := make([]byte, 32)
	if err := binary.Read(buf, binary.LittleEndian, &amount); err != nil {
		return nil, fmt.Errorf("failed to read amount: %v", err)
	}
	transaction.Amount = string(amount)

	// Read Signature
	signature := make([]byte, 32)
	if err := binary.Read(buf, binary.LittleEndian, &signature); err != nil {
		return nil, fmt.Errorf("failed to read signature: %v", err)
	}
	transaction.Signature = signature

	// Read Nonce
	var nonce int64
	if err := binary.Read(buf, binary.LittleEndian, &nonce); err != nil {
		return nil, fmt.Errorf("failed to read nonce: %v", err)
	}
	transaction.Nonce = nonce

	return transaction, nil
}
