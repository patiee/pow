package storage

import (
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type KVStorage struct {
	name string
	path string

	db   *leveldb.DB
	lock sync.RWMutex
}

func New(name, path string) (*KVStorage, error) {
	path = fmt.Sprintf("%s/%s", path, name)
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s db: %v", name, err)
	}
	return &KVStorage{
		name: name,
		path: path,
		db:   db,
	}, nil
}

// Close connection to leveldb
func (kv *KVStorage) Close() {
	kv.db.Close()
}

// Set a new value by key in storage
func (kv *KVStorage) Set(key []byte, value []byte) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	if err := kv.db.Put(key, value, nil); err != nil {
		return fmt.Errorf("could not set new value to %s kv storage: %v", kv.name, err)
	}

	return nil
}

// Get a value by key from storage
func (kv *KVStorage) Get(key []byte) ([]byte, error) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	value, err := kv.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return value, fmt.Errorf("key not found")
}
