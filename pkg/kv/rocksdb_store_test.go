package kv

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var dbPath = "test"

func TestPut(t *testing.T) {
	db, err := NewRocksDBStore(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 1; i < 5; i++ {
		err = db.Put([]byte(fmt.Sprintf("data-%d", i)), []byte(fmt.Sprintf("value-%d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 1; i < 5; i++ {
		value, err := db.Get([]byte(fmt.Sprintf("data-%d", i)))
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(value, []byte(fmt.Sprintf("value-%d", i))) {
			t.Fatal("should not happen")
		}
	}

	if err = db.Delete([]byte("data-1")); err != nil {
		t.Fatal(err)
	}

	it := db.Iterator()
	for it.SeekToFirst(); it.Valid(); it.Next() {
		if bytes.Equal(it.Key().Data(), []byte("data-1")) {
			t.Fatal("should not happen")
		}
	}

	if _, err = os.Stat(dbPath); err != nil {
		t.Fatal("stat rocksdb test dir error")
	} else {
		os.RemoveAll(dbPath)
	}
}
