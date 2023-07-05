package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initHashDB() *HashStructure {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestHashStructure")
	opts.DirPath = dir
	hash, _ := NewHashStructure(opts)
	return hash
}

func TestHashStructure_HGet(t *testing.T) {
	hash := initHashDB()

	ok1, err := hash.HSet(randkv.GetTestKey(1), []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	v1 := randkv.RandomValue(10)
	ok2, err := hash.HSet(randkv.GetTestKey(1), []byte("field1"), v1)
	assert.Nil(t, err)
	assert.False(t, ok2)
	value1, err := hash.HGet(randkv.GetTestKey(1), []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, value1, v1)

	v2 := randkv.RandomValue(10)
	ok3, err := hash.HSet(randkv.GetTestKey(1), []byte("field2"), v2)
	assert.Nil(t, err)
	assert.True(t, ok3)
	value2, err := hash.HGet(randkv.GetTestKey(1), []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, value2, v2)

	_, err = hash.HGet(randkv.GetTestKey(1), []byte("field3"))
	assert.Equal(t, err, _const.ErrKeyNotFound)

}

func TestHashStructure_HDel(t *testing.T) {
	hash := initHashDB()

	ok, err := hash.HDel(randkv.GetTestKey(1), []byte("field1"))
	assert.Nil(t, err)
	assert.False(t, ok)

	ok1, err := hash.HSet(randkv.GetTestKey(1), []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HDel(randkv.GetTestKey(1), []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet(randkv.GetTestKey(1), []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSet(randkv.GetTestKey(1), []byte("field2"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok4)

	ok5, err := hash.HDel(randkv.GetTestKey(1), []byte("field1"), []byte("field2"))
	assert.Nil(t, err)
	assert.True(t, ok5)

}

func TestHashStructure_HExists(t *testing.T) {
	hash := initHashDB()

	ok1, err := hash.HSet(randkv.GetTestKey(1), []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HExists(randkv.GetTestKey(1), []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HExists(randkv.GetTestKey(1), []byte("field2"))
	assert.Nil(t, err)
	assert.False(t, ok3)

}
