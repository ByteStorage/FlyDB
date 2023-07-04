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

func initdb() *StringStructure {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestStringStructure_Get")
	opts.DirPath = dir
	str, _ := NewStringStructure(opts)
	return str
}

func TestStringStructure_Get(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	err = str.Set(randkv.GetTestKey(2), randkv.RandomValue(100), 2*time.Second)
	assert.Nil(t, err)

	value1, err := str.Get(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, value1)

	time.Sleep(3 * time.Second)

	value2, err := str.Get(randkv.GetTestKey(2))
	assert.Equal(t, err, ErrKeyExpired)
	assert.Nil(t, value2)

	_, err = str.Get(randkv.GetTestKey(3))
	assert.Equal(t, err, _const.ErrKeyNotFound)
}

func TestStringStructure_Del(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), randkv.RandomValue(100), 0)
	assert.Nil(t, err)

	err = str.Del(randkv.GetTestKey(1))
	assert.Nil(t, err)

	_, err = str.Get(randkv.GetTestKey(1))
	assert.Equal(t, err, _const.ErrKeyNotFound)
}

func TestStringStructure_Type(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), randkv.RandomValue(100), 0)
	assert.Nil(t, err)

	keyType, err := str.Type(randkv.GetTestKey(1))

	TypeString := "string"
	assert.Equal(t, keyType, TypeString)
	assert.Nil(t, err)
}

func TestStringStructure_StrLen(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), randkv.RandomValue(100), 0)
	assert.Nil(t, err)

	strLen, err := str.StrLen(randkv.GetTestKey(1))
	assert.Nil(t, err)
	assert.Equal(t, strLen, 112)
}

func TestStringStructure_GetSet(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), randkv.RandomValue(100), 0)
	assert.Nil(t, err)
	value1, _ := str.Get(randkv.GetTestKey(1))

	value2, err := str.GetSet(randkv.GetTestKey(1), randkv.RandomValue(100), 2*time.Second)
	assert.Nil(t, err)
	assert.NotNil(t, value2)
	assert.Equal(t, value1, value2)
}

func TestStringStructure_Append(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), randkv.RandomValue(10), 0)
	//assert.Nil(t, err)
	s1, _ := str.Get(randkv.GetTestKey(1))
	t.Log(string(s1))

	err = str.Append(randkv.GetTestKey(1), randkv.RandomValue(5), 0)
	//assert.Nil(t, err)
	s2, _ := str.Get(randkv.GetTestKey(1))
	t.Log(string(s2))
}

func TestStringStructure_Incr(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), []byte("1"), 0)
	assert.Nil(t, err)

	err := str.Incr(randkv.GetTestKey(1), 0)
	assert.Nil(t, err)
	v1, _ := str.Get(randkv.GetTestKey(1))
	assert.Equal(t, string(v1), "2")

	err = str.Incr(randkv.GetTestKey(1), 0)
	assert.Nil(t, err)
	v2, _ := str.Get(randkv.GetTestKey(1))
	assert.Equal(t, string(v2), "3")
}

func TestStringStructure_IncrBy(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), []byte("1"), 0)
	assert.Nil(t, err)

	err := str.IncrBy(randkv.GetTestKey(1), 10, 0)
	assert.Nil(t, err)
	v1, _ := str.Get(randkv.GetTestKey(1))
	assert.Equal(t, string(v1), "11")

	err = str.IncrBy(randkv.GetTestKey(1), 10, 0)
	assert.Nil(t, err)
	v2, _ := str.Get(randkv.GetTestKey(1))
	assert.Equal(t, string(v2), "21")
}

func TestStringStructure_IncrByFloat(t *testing.T) {
	str := initdb()

	err = str.Set(randkv.GetTestKey(1), []byte("1"), 0)
	assert.Nil(t, err)

	err := str.IncrByFloat(randkv.GetTestKey(1), 1.1, 0)
	assert.Nil(t, err)
	v1, _ := str.Get(randkv.GetTestKey(1))
	assert.Equal(t, string(v1), "2.1")

	err = str.IncrByFloat(randkv.GetTestKey(1), 1.1, 0)
	assert.Nil(t, err)
	v2, _ := str.Get(randkv.GetTestKey(1))
	assert.Equal(t, string(v2), "3.2")
}
