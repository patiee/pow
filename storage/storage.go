package storage

import (
	"fmt"

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
