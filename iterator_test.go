package flydb

import (
	"github.com/qishenonly/flydb/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDB_NewIterator(t *testing.T) {
	opt := DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-iterator-1")
	opt.DirPath = dir
	db, err := NewFlyDB(opt)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	iterator := db.NewIterator(DefaultIteratorOptions)
	assert.NotNil(t, iterator)
	assert.Equal(t, false, iterator.Valid())
}

func TestDB_Iterator_One_Value(t *testing.T) {
	opt := DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-iterator-2")
	opt.DirPath = dir
	db, err := NewFlyDB(opt)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(utils.GetTestKey(10), utils.GetTestKey(10))
	assert.Nil(t, err)

	iterator := db.NewIterator(DefaultIteratorOptions)
	assert.NotNil(t, iterator)
	assert.True(t, iterator.Valid())
	assert.Equal(t, utils.GetTestKey(10), iterator.Key())
	value, err := iterator.Value()
	assert.Equal(t, utils.GetTestKey(10), value)
}

func TestDB_Iterator_Multi_Value(t *testing.T) {
	opt := DefaultOptions
	dir, _ := os.MkdirTemp("", "flydb-iterator-3")
	opt.DirPath = dir
	db, err := NewFlyDB(opt)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put([]byte("abcd"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("efjh"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("aefg"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("cdef"), utils.RandomValue(10))
	assert.Nil(t, err)
	err = db.Put([]byte("bcdg"), utils.RandomValue(10))
	assert.Nil(t, err)

	// 正向迭代
	itertor1 := db.NewIterator(DefaultIteratorOptions)
	for itertor1.Rewind(); itertor1.Valid(); itertor1.Next() {
		//t.Log("key => ", string(itertor1.Key()))
		assert.NotNil(t, itertor1.Key())
	}

	itertor1.Rewind()
	for itertor1.Seek([]byte("c")); itertor1.Valid(); itertor1.Next() {
		t.Log("key => ", string(itertor1.Key()))
		assert.NotNil(t, itertor1.Key())
	}

	// 反向迭代
	iterOpt2 := DefaultIteratorOptions
	iterOpt2.Reverse = true
	itertor2 := db.NewIterator(iterOpt2)
	for itertor2.Rewind(); itertor2.Valid(); itertor2.Next() {
		//t.Log("key => ", string(itertor2.Key()))
		assert.NotNil(t, itertor2.Key())
	}

	itertor2.Rewind()
	for itertor2.Seek([]byte("c")); itertor2.Valid(); itertor2.Next() {
		//t.Log("key => ", string(itertor2.Key()))
		assert.NotNil(t, itertor2.Key())
	}

	// 指定 prefix
	iterOpt3 := DefaultIteratorOptions
	iterOpt3.Prefix = []byte("ae")
	itertor3 := db.NewIterator(iterOpt3)
	for itertor3.Rewind(); itertor3.Valid(); itertor3.Next() {
		//t.Log("key => ", string(itertor3.Key()))
		assert.NotNil(t, itertor3.Key())
	}

}
