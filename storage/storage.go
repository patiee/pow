package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type KVStorage struct {
	name string
}

func New(name string) *KVStorage {
	return &KVStorage{
		name: name,
	}
}

// Set a new value by key in storage
func (kv *KVStorage) Set(key string, value []byte) error {
	// Open the file for appending
	file, err := os.OpenFile(kv.name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := setKey(file, key); err != nil {
		return err
	}
	if err := setValue(file, value); err != nil {
		return err
	}

	return nil
}

// Get a value by key from storage
func (kv *KVStorage) Get(key string) ([]byte, error) {
	file, err := os.OpenFile(kv.name, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for {
		storedKey, err := readKey(file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if storedKey == key {
			value, err := readValue(file)
			if err != nil {
				return nil, err
			}
			return value, nil
		}
	}

	return nil, fmt.Errorf("key not found")
}

// setKey writes a string with length prefix in binary format
func setKey(w io.Writer, s string) error {
	if err := binary.Write(w, binary.LittleEndian, int32(len(s))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(s)); err != nil {
		return err
	}
	return nil
}

// setValue writes a byte slice in binary format
func setValue(w io.Writer, data []byte) error {
	if err := binary.Write(w, binary.LittleEndian, int32(len(data))); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

// readKey reads a string with length prefix from binary format
func readKey(r io.Reader) (string, error) {
	var length int32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	data := make([]byte, length)
	if _, err := r.Read(data); err != nil {
		return "", err
	}
	return string(data), nil
}

// readValue reads a byte slice from binary format
func readValue(r io.Reader) ([]byte, error) {
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
