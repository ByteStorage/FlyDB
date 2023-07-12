package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var err error

func initdb() (*StringStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestStringStructure_Get")
	opts.DirPath = dir
	str, _ := NewStringStructure(opts)
	return str, &opts
}

func TestStringStructure_Get(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	err = str.Set("1", randkv.RandomValue(100), 2*time.Second)
	assert.Nil(t, err)

	value1, err := str.Get("1")
	assert.Nil(t, err)
	assert.NotNil(t, value1)

	time.Sleep(3 * time.Second)

	value2, err := str.Get("1")
	assert.Equal(t, err, ErrKeyExpired)
	assert.Nil(t, value2)

	_, err = str.Get("3")
	assert.Equal(t, err, _const.ErrKeyNotFound)
}

func TestStringStructure_Del(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", randkv.RandomValue(100), 0)
	assert.Nil(t, err)

	err = str.Del("1")
	assert.Nil(t, err)

	_, err = str.Get("1")
	assert.Equal(t, err, _const.ErrKeyNotFound)
}

func TestStringStructure_Type(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", randkv.RandomValue(100), 0)
	assert.Nil(t, err)

	keyType, err := str.Type("1")

	TypeString := "string"
	assert.Equal(t, keyType, TypeString)
	assert.Nil(t, err)
}

func TestStringStructure_StrLen(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", randkv.RandomValue(100), 0)
	assert.Nil(t, err)

	strLen, err := str.StrLen("1")
	assert.Nil(t, err)
	assert.Equal(t, strLen, 112)
}

func TestStringStructure_GetSet(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	value1, _ := str.Get("1")

	value2, err := str.GetSet("1", randkv.RandomValue(100), 2*time.Second)
	assert.Nil(t, err)
	assert.NotNil(t, value2)
	assert.Equal(t, value1, value2)
}

func TestStringStructure_Append(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", randkv.RandomValue(10), 0)
	assert.Nil(t, err)

	err = str.Append("1", randkv.RandomValue(5), 0)
	assert.Nil(t, err)
}

func TestStringStructure_Incr(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", "1", 0)
	assert.Nil(t, err)

	err := str.Incr("1", 0)
	assert.Nil(t, err)
	v1, _ := str.Get("1")
	assert.Equal(t, v1, "2")

	err = str.Incr("1", 0)
	v2, _ := str.Get("1")
	assert.Equal(t, v2, "3")

	err = str.Set("1", 1, 0)
	assert.Nil(t, err)

	err = str.Incr("1", 0)
	assert.Nil(t, err)
	v3, _ := str.Get("1")
	assert.Equal(t, v3, "2")

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	err = str.Incr("1", 0)
	assert.Nil(t, err)
	v4, _ := str.Get("1")
	assert.Equal(t, v4, "2")
}

func TestStringStructure_IncrBy(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", 1, 0)
	assert.Nil(t, err)

	err := str.IncrBy("1", 10, 0)
	assert.Nil(t, err)

	v1, _ := str.Get("1")
	assert.Equal(t, v1, "11")

	err = str.IncrBy("1", 10, 0)
	assert.Nil(t, err)

	v2, _ := str.Get("1")
	assert.Equal(t, v2, "21")

	err = str.Set("1", "1", 0)
	assert.Nil(t, err)

	err = str.IncrBy("1", 10, 0)
	assert.Nil(t, err)

	v2, _ = str.Get("1")
	assert.Equal(t, v2, "11")

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	err = str.IncrBy("1", 10, 0)
	assert.Nil(t, err)
	v4, _ := str.Get("1")
	assert.Equal(t, v4, "11")
}

func TestStringStructure_IncrByFloat(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", "1", 0)
	assert.Nil(t, err)

	err := str.IncrByFloat("1", 1.1, 0)
	assert.Nil(t, err)
	v1, _ := str.Get("1")
	assert.Equal(t, v1, "2.1")

	err = str.IncrByFloat("1", 1.1, 0)
	assert.Nil(t, err)
	v2, _ := str.Get("1")
	assert.Equal(t, v2, "3.2")

	err = str.Set("1", 1, 0)
	assert.Nil(t, err)

	err = str.IncrByFloat("1", 1.1, 0)
	assert.Nil(t, err)
	v3, _ := str.Get("1")
	assert.Equal(t, v3, "2.1")

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	err = str.IncrByFloat("1", 1.1, 0)
	assert.Nil(t, err)
	v4, _ := str.Get("1")
	assert.Equal(t, v4, "2.1")
}

func TestStringStructure_Decr(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", "1", 0)
	assert.Nil(t, err)

	err := str.Decr("1", 0)
	assert.Nil(t, err)
	v1, _ := str.Get("1")
	assert.Equal(t, v1, "0")

	err = str.Decr("1", 0)
	v2, _ := str.Get("1")
	assert.Equal(t, v2, "-1")

	err = str.Set("1", 1, 0)
	assert.Nil(t, err)

	err = str.Decr("1", 0)
	assert.Nil(t, err)
	v3, _ := str.Get("1")
	assert.Equal(t, v3, "0")

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	err = str.Decr("1", 0)
	assert.Nil(t, err)
	v4, _ := str.Get("1")
	assert.Equal(t, v4, "0")
}

func TestStringStructure_DecrBy(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", 1, 0)
	assert.Nil(t, err)

	err := str.DecrBy("1", 10, 0)
	assert.Nil(t, err)
	v1, _ := str.Get("1")
	assert.Equal(t, v1, "-9")

	err = str.DecrBy("1", 10, 0)
	assert.Nil(t, err)
	v2, _ := str.Get("1")
	assert.Equal(t, v2, "-19")

	err = str.Set("1", "1", 0)
	assert.Nil(t, err)

	err = str.DecrBy("1", 10, 0)
	assert.Nil(t, err)
	v3, _ := str.Get("1")
	assert.Equal(t, v3, "-9")

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	err = str.DecrBy("1", 10, 0)
	assert.Nil(t, err)
	v4, _ := str.Get("1")
	assert.Equal(t, v4, "-9")
}

func TestStringStructure_Exists(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	ok1, err := str.Exists("1")
	assert.Nil(t, err)
	assert.Equal(t, ok1, true)

	ok2, err := str.Exists("1")
	assert.Nil(t, err)
	assert.Equal(t, ok2, true)
}

func TestStringStructure_Expire(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	err = str.Expire("1", 1*time.Second)
	assert.Nil(t, err)
	v1, err := str.Get("1")
	assert.Nil(t, err)
	assert.Equal(t, v1, "1")

	time.Sleep(2 * time.Second)
	v2, err := str.Get("1")
	assert.Equal(t, err, ErrKeyExpired)
	assert.Equal(t, v2, nil)
}

func TestStringStructure_Persist(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err = str.Set("1", []byte("1"), 0)
	assert.Nil(t, err)

	err = str.Expire("1", 1*time.Second)
	assert.Nil(t, err)
	v1, err := str.Get("1")
	assert.Nil(t, err)
	assert.Equal(t, v1, "1")

	err = str.Persist("1")
	assert.Nil(t, err)
	v2, err := str.Get("1")
	assert.Nil(t, err)
	assert.Equal(t, v2, "1")
}

func TestStringStructure_MGet(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err := str.Set("key1", "value1", 0)
	assert.Nil(t, err)
	err = str.Set("key2", "value2", 0)
	assert.Nil(t, err)

	keys := []string{"key1", "key2"} // Simulating keys to be retrieved
	values, err := str.MGet(keys...)
	assert.Nil(t, err)
	assert.Equal(t, len(values), len(keys))

	expectedValues := []interface{}{"value1", "value2", nil} // Expected values based on the keys
	for i, value := range values {
		assert.Equal(t, value, expectedValues[i])
	}
}

func TestStringStructure_MSet(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	err := str.MSet("key1", "value1", "key2", "value2", "key3", "value3")
	assert.Nil(t, err)

	value1, err := str.Get("key1")
	assert.Nil(t, err)
	assert.Equal(t, value1, "value1")

	value2, err := str.Get("key2")
	assert.Nil(t, err)
	assert.Equal(t, value2, "value2")

	value3, err := str.Get("key3")
	assert.Nil(t, err)
	assert.Equal(t, value3, "value3")
}

func TestStringStructure_MSetNX(t *testing.T) {
	str, _ := initdb()
	defer str.db.Clean()

	// Test case: All keys and values are new, should return true
	success, err := str.MSetNX("key1", "value1", "key2", "value2", "key3", "value3")
	assert.Nil(t, err)
	assert.True(t, success)

	// Test case: At least one key already exists, should return false
	err = str.Set("key1", "existingValue", 0)
	assert.Nil(t, err)

	success, err = str.MSetNX("key1", "value1", "key2", "value2", "key3", "value3")
	assert.Nil(t, err)
	assert.False(t, success)

	value1, err := str.Get("key1")
	assert.Nil(t, err)
	assert.Equal(t, value1, "existingValue")

	value2, err := str.Get("key2")
	assert.Nil(t, err)
	assert.Equal(t, value2, "value2")

	value3, err := str.Get("key3")
	assert.Nil(t, err)
	assert.Equal(t, value3, "value3")
}
