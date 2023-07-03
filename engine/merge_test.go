package engine

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// merge without any data
func TestDB_Merge(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-merge-1")
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Merge()
	assert.Nil(t, err)
}

// All valid data
func TestDB_Merge2(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-merge-2")
	opts.DataFileSize = 32 * 1024 * 1024
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	for i := 0; i < 50000; i++ {
		err := db.Put(randkv.GetTestKey(i), randkv.RandomValue(1024))
		assert.Nil(t, err)
	}

	err = db.Merge()
	assert.Nil(t, err)

	// Restart check
	err = db.Close()
	assert.Nil(t, err)

	db2, err := NewDB(opts)
	defer func() {
		_ = db2.Close()
	}()
	assert.Nil(t, err)
	keys := db2.GetListKeys()
	assert.Equal(t, 50000, len(keys))

	for i := 0; i < 50000; i++ {
		val, err := db2.Get(randkv.GetTestKey(i))
		assert.Nil(t, err)
		assert.NotNil(t, val)
	}

}

// All invalid data
func TestDB_Merge3(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-merge-3")
	opts.DataFileSize = 32 * 1024 * 1024
	opts.DirPath = dir
	db, err := NewDB(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	for i := 0; i < 50000; i++ {
		err := db.Put(randkv.GetTestKey(i), randkv.RandomValue(1024))
		assert.Nil(t, err)
	}
	for i := 0; i < 50000; i++ {
		err := db.Delete(randkv.GetTestKey(i))
		assert.Nil(t, err)
	}

	err = db.Merge()
	assert.Nil(t, err)

	// Restart check
	err = db.Close()
	assert.Nil(t, err)

	db2, err := NewDB(opts)
	defer func() {
		_ = db2.Close()
	}()
	assert.Nil(t, err)
	keys := db2.GetListKeys()
	assert.Equal(t, 0, len(keys))
}
