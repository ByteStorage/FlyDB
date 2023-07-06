package engine

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDB_NewIterator(t *testing.T) {
	opt := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-iterator-1")
	opt.DirPath = dir
	db, err := NewDB(opt)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	iterator := db.NewIterator(config.DefaultIteratorOptions)
	assert.NotNil(t, iterator)
	assert.Equal(t, false, iterator.Valid())
}

func TestDB_Iterator_One_Value(t *testing.T) {
	opt := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-iterator-2")
	opt.DirPath = dir
	db, err := NewDB(opt)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(randkv.GetTestKey(10), randkv.GetTestKey(10))
	assert.Nil(t, err)

	iterator := db.NewIterator(config.DefaultIteratorOptions)
	assert.NotNil(t, iterator)
	assert.True(t, iterator.Valid())
	assert.Equal(t, randkv.GetTestKey(10), iterator.Key())
	value, err := iterator.Value()
	assert.Equal(t, randkv.GetTestKey(10), value)
}

func TestDB_Iterator_Multi_Value(t *testing.T) {
	opt := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-iterator-3")
	opt.DirPath = dir
	db, err := NewDB(opt)
	defer db.Clean()
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put([]byte("abcd"), randkv.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("efjh"), randkv.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("aefg"), randkv.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("cdef"), randkv.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("bcdg"), randkv.RandomValue(10))
	assert.Nil(t, err)

	// Forward iteration
	itertor1 := db.NewIterator(config.DefaultIteratorOptions)
	for itertor1.Rewind(); itertor1.Valid(); itertor1.Next() {
		assert.NotNil(t, itertor1.Key())
	}

	itertor1.Rewind()
	for itertor1.Seek([]byte("c")); itertor1.Valid(); itertor1.Next() {
		assert.NotNil(t, itertor1.Key())
	}

	// Reverse iteration
	iterOpt2 := config.DefaultIteratorOptions
	iterOpt2.Reverse = true
	itertor2 := db.NewIterator(iterOpt2)
	for itertor2.Rewind(); itertor2.Valid(); itertor2.Next() {
		assert.NotNil(t, itertor2.Key())
	}

	itertor2.Rewind()
	for itertor2.Seek([]byte("c")); itertor2.Valid(); itertor2.Next() {
		assert.NotNil(t, itertor2.Key())
	}

	// Specify prefix
	iterOpt3 := config.DefaultIteratorOptions
	iterOpt3.Prefix = []byte("ae")
	itertor3 := db.NewIterator(iterOpt3)
	for itertor3.Rewind(); itertor3.Valid(); itertor3.Next() {
		assert.NotNil(t, itertor3.Key())
	}

}
