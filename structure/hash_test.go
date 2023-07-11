package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initHashDB() (*HashStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestHashStructure")
	opts.DirPath = dir
	hash, _ := NewHashStructure(opts)
	return hash, &opts
}

func TestHashStructure_HGet(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	v1 := randkv.RandomValue(10)
	ok2, err := hash.HSet("1", []byte("field1"), v1)
	assert.Nil(t, err)
	assert.False(t, ok2)
	value1, err := hash.HGet("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, value1, v1)

	v2 := randkv.RandomValue(10)
	ok3, err := hash.HSet("1", []byte("field2"), v2)
	assert.Nil(t, err)
	assert.True(t, ok3)
	value2, err := hash.HGet("1", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, value2, v2)

	_, err = hash.HGet("1", []byte("field3"))
	assert.Equal(t, err, _const.ErrKeyNotFound)

}

func TestHashStructure_HMGet(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	v1 := randkv.RandomValue(10)
	ok2, err := hash.HSet("1", []byte("field1"), v1)
	assert.Nil(t, err)
	assert.False(t, ok2)

	v2 := randkv.RandomValue(10)
	ok3, err := hash.HSet("1", []byte("field2"), v2)
	assert.Nil(t, err)
	assert.True(t, ok3)

	mulVal, err := hash.HMGet("1", []byte("field1"), []byte("field2"))
	assert.Equal(t, v1, mulVal[0])
	assert.Equal(t, v2, mulVal[1])

	_, err = hash.HGet("1", []byte("field3"))
	assert.Equal(t, err, _const.ErrKeyNotFound)

}

func TestHashStructure_HDel(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok, err := hash.HDel("1", []byte("field1"))
	assert.Nil(t, err)
	assert.False(t, ok)

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HDel("1", []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	v, err := hash.HGet("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Nil(t, v)

	ok3, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSet("1", []byte("field2"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok4)

}

func TestHashStructure_HExists(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HExists("1", []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HExists("1", []byte("field2"))
	assert.Nil(t, err)
	assert.False(t, ok3)

}

func TestHashStructure_HLen(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok3)

	l, err := hash.HLen("1")
	assert.Nil(t, err)
	assert.Equal(t, l, 3)
}

func TestHashStructure_HUpdate(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HUpdate("1", []byte("field1"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok4)

	ok5, err := hash.HUpdate("1", []byte("field2"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.True(t, ok5)

	ok6, err := hash.HUpdate("1", []byte("field4"), randkv.RandomValue(10))
	assert.Nil(t, err)
	assert.False(t, ok6)

	l, err := hash.HLen("1")
	assert.Nil(t, err)
	assert.Equal(t, l, 3)
}

func TestHashStructure_HIncrBy(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	v1, err := hash.HIncrBy("1", []byte("field1"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v1, int64(11))

	v2, err := hash.HIncrBy("1", []byte("field2"), -1)
	assert.Nil(t, err)
	assert.Equal(t, v2, int64(9))

	v3, err := hash.HIncrBy("1", []byte("field3"), 0)
	assert.Nil(t, err)
	assert.Equal(t, v3, int64(10))

	v4, err := hash.HIncrBy("1", []byte("field4"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v4, int64(0))

}

func TestHashStructure_HIncrByFloat(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	v1, err := hash.HIncrByFloat("1", []byte("field1"), 1.1)
	assert.Nil(t, err)
	assert.Equal(t, v1, float64(11.1))

	v2, err := hash.HIncrByFloat("1", []byte("field2"), -1.1)
	assert.Nil(t, err)
	assert.Equal(t, v2, float64(8.9))

	v3, err := hash.HIncrByFloat("1", []byte("field3"), 0)
	assert.Nil(t, err)
	assert.Equal(t, v3, float64(10))

	v4, err := hash.HIncrByFloat("1", []byte("field4"), 1.1)
	assert.Nil(t, err)
	assert.Equal(t, v4, float64(0))

}

func TestHashStructure_HDecrBy(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	v1, err := hash.HDecrBy("1", []byte("field1"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v1, int64(9))

	v2, err := hash.HDecrBy("1", []byte("field2"), 10)
	assert.Nil(t, err)
	assert.Equal(t, v2, int64(0))

	v3, err := hash.HDecrBy("1", []byte("field3"), 0)
	assert.Nil(t, err)
	assert.Equal(t, v3, int64(10))

	v4, err := hash.HDecrBy("1", []byte("field4"), 1)
	assert.Nil(t, err)
	assert.Equal(t, v4, int64(0))

}

func TestHashStructure_HStrLen(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	l1, err := hash.HStrLen("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, l1, 4)

	l2, err := hash.HStrLen("1", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, l2, 3)

	l3, err := hash.HStrLen("1", []byte("field3"))
	assert.Nil(t, err)
	assert.Equal(t, l3, 2)

	l4, err := hash.HStrLen("1", []byte("field4"))
	assert.Nil(t, err)
	assert.Equal(t, l4, 0)

}

func TestHashStructure_HMove(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("111-1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("111-100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("111-10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSet("2", []byte("field1"), []byte("222-1000"))
	assert.Nil(t, err)
	assert.True(t, ok4)

	ok5, err := hash.HSet("2", []byte("field2"), []byte("222-100"))
	assert.Nil(t, err)
	assert.True(t, ok5)

	ok6, err := hash.HSet("2", []byte("field3"), []byte("222-10"))
	assert.Nil(t, err)
	assert.True(t, ok6)

	ok7, err := hash.HMove("2", "1", []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, ok7)

	ok8, err := hash.HMove("2", "1", []byte("field2"))
	assert.Nil(t, err)
	assert.True(t, ok8)

	ok9, err := hash.HMove("2", "1", []byte("field3"))
	assert.Nil(t, err)
	assert.True(t, ok9)

	ok10, err := hash.HMove("2", "1", []byte("field4"))
	assert.Nil(t, err)
	assert.False(t, ok10)

	v1, err := hash.HGet("2", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, v1, []byte("111-1000"))

	v2, err := hash.HGet("2", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, v2, []byte("111-100"))

	v3, err := hash.HGet("2", []byte("field3"))
	assert.Nil(t, err)
	assert.Equal(t, v3, []byte("111-10"))

}

func TestHashStructure_HSetNX(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSetNX("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSetNX("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSetNX("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	ok4, err := hash.HSetNX("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.False(t, ok4)

	ok5, err := hash.HSetNX("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.False(t, ok5)

	ok6, err := hash.HSetNX("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.False(t, ok6)

}

func TestHashStructure_HTypes(t *testing.T) {
	hash, _ := initHashDB()
	defer hash.db.Clean()

	ok1, err := hash.HSet("1", []byte("field1"), []byte("1000"))
	assert.Nil(t, err)
	assert.True(t, ok1)

	ok2, err := hash.HSet("1", []byte("field2"), []byte("100"))
	assert.Nil(t, err)
	assert.True(t, ok2)

	ok3, err := hash.HSet("1", []byte("field3"), []byte("10"))
	assert.Nil(t, err)
	assert.True(t, ok3)

	type1, err := hash.HTypes("1", []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, type1, "hash")

	type2, err := hash.HTypes("1", []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, type2, "hash")

	type3, err := hash.HTypes("1", []byte("field3"))
	assert.Nil(t, err)
	assert.Equal(t, type3, "hash")

	type4, err := hash.HTypes("1", []byte("field4"))
	assert.Equal(t, "", type4)
	assert.Equal(t, err, _const.ErrKeyNotFound)
}
