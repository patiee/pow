package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
)

type KVStorage struct {
	name string
	path string

	lock sync.RWMutex
}

func New(name, path string) *KVStorage {
	return &KVStorage{
		name: name,
		path: fmt.Sprintf("%s/%s", path, name),
	}
}

// Set a new value by key in storage
func (kv *KVStorage) Set(key []byte, value []byte) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	// Open the file for appending
	file, err := os.OpenFile(kv.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := write(file, key); err != nil {
		return err
	}
	if err := write(file, value); err != nil {
		return err
	}

	return nil
}

// Get a value by key from storage
func (kv *KVStorage) Get(key []byte) ([]byte, error) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	file, err := os.OpenFile(kv.path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for {
		storedKey, err := read(file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if bytes.Equal(storedKey, key) {
			value, err := read(file)
			if err != nil {
				return nil, err
			}
			return value, nil
		} else {
			// Skip the value if key doesn’t match
			_, err = read(file)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("key not found")
}

// Update a value by key from storage
func (kv *KVStorage) Update(key, newValue []byte) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	file, err := os.OpenFile(kv.path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		storedKey, err := read(file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if bytes.Equal(storedKey, key) {
			valuePosition, err := file.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
			_, err = read(file)
			if err != nil {
				return err
			}

			remaining, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			// Write new value
			file.Seek(valuePosition, io.SeekStart)
			if err := write(file, newValue); err != nil {
				return err
			}

			// Write remaining data
			if _, err := file.Write(remaining); err != nil {
				return err
			}

			// Truncate if new value is shorter
			newEnd := valuePosition + int64(4+len(newValue)) + int64(len(remaining))
			if err := file.Truncate(newEnd); err != nil {
				return err
			}

			return nil
		} else {
			// Skip the value if key doesn’t match
			_, err = read(file)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("key not found")
}

// write a byte slice with length prefix in binary format
func write(w io.Writer, key []byte) error {
	if err := binary.Write(w, binary.LittleEndian, int32(len(key))); err != nil {
		return err
	}
	if _, err := w.Write(key); err != nil {
		return err
	}
	return nil
}

// read a byte slice with length prefix from binary format
func read(r io.Reader) ([]byte, error) {
	var length int32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	data := make([]byte, length)
	if _, err := r.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
