package storage

import (
	"fmt"
	"math/big"

	"github.com/patiee/pow/model"
)

type Storage struct {
	block       *KVStorage
	transaction *KVStorage
	wallet      *KVStorage
}

func NewStorage(path string) *Storage {
	return &Storage{
		block:       New("block", path),
		transaction: New("tx", path),
		wallet:      New("wallet", path),
	}
}

func (s *Storage) AddBlock(block *model.Block) error {
	blockBytes, err := block.Serialize()
	if err != nil {
		return fmt.Errorf("could not serialize block: %v", err)
	}
	blockHash := block.Hash()
	return s.block.Set(blockHash[:], blockBytes)
}

func (s *Storage) GetBlock(hash []byte) (*model.Block, error) {
	blockBytes, err := s.block.Get(hash)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize block: %v", err)
	}
	return model.DeserializeBlock(blockBytes)
}

func (s *Storage) AddTransaction(transaction *model.Transaction) error {
	blockBytes, err := transaction.Serialize()
	if err != nil {
		return fmt.Errorf("could not serialize block: %v", err)
	}
	blockHash := transaction.Hash()
	return s.transaction.Set(blockHash[:], blockBytes)
}

func (s *Storage) GetTransaction(hash []byte) (*model.Transaction, error) {
	transactionBytes, err := s.transaction.Get(hash)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize transaction: %v", err)
	}
	return model.DeserializeTransaction(transactionBytes)
}

func (s *Storage) SetWallet(wallet []byte, value *big.Int) error {
	valueBytes := value.Bytes()
	err := s.wallet.Update(wallet, valueBytes)
	if err != nil {
		return s.wallet.Set(wallet, valueBytes)
	}

	return nil
}

func (s *Storage) GetWallet(wallet []byte) (*big.Int, error) {
	valueBytes, err := s.wallet.Get(wallet)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(valueBytes), nil
}
