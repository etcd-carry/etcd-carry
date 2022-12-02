package kv

import (
	"github.com/tecbot/gorocksdb"
	"os"
)

type rocksDBStore struct {
	path string
	db   *gorocksdb.DB
}

func NewRocksDBStore(path string) (KeyValueDB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetCompression(gorocksdb.NoCompression)
	db, err := gorocksdb.OpenDb(opts, path)
	if err != nil {
		return nil, err
	}
	return &rocksDBStore{path: path, db: db}, nil
}

func (s *rocksDBStore) Put(key, value []byte) error {
	writeOpts := gorocksdb.NewDefaultWriteOptions()
	return s.db.Put(writeOpts, key, value)
}

func (s *rocksDBStore) Get(key []byte) ([]byte, error) {
	readOpts := gorocksdb.NewDefaultReadOptions()
	return s.db.GetBytes(readOpts, key)
}

func (s *rocksDBStore) Delete(key []byte) error {
	writeOpts := gorocksdb.NewDefaultWriteOptions()
	return s.db.Delete(writeOpts, key)
}

func (s *rocksDBStore) Iterator() *gorocksdb.Iterator {
	readOpts := gorocksdb.NewDefaultReadOptions()
	return s.db.NewIterator(readOpts)
}

func (s *rocksDBStore) Close() {
	s.db.Close()
}
