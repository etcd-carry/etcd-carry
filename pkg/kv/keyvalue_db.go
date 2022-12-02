package kv

import "github.com/tecbot/gorocksdb"

type KeyValueDB interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Iterator() *gorocksdb.Iterator
	Close()
}
