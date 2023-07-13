package structure

import (
	"os"
	"testing"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/randkv"
	"github.com/stretchr/testify/assert"
)

var bitmapErr error

func initBitmap() *BitmapStructure {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestBitmapStructure")
	opts.DirPath = dir
	bitmap, _ := NewBitmapStructure(opts)
	return bitmap
}

func TestBitmapStructure_SetBit(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test SetBit function
	bitmapErr = bitmap.SetBit(string(randkv.GetTestKey(1)), 0, true)
	assert.Nil(t, bitmapErr)
}

func TestBitmapStructure_SetBits(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test SetBits function
	bitmapErr = bitmap.SetBits(string(randkv.GetTestKey(1)), 0, true, 1, false)
	assert.Nil(t, bitmapErr)
}

func TestBitmapStructure_GetBit(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test GetBit function
	bitmapErr = bitmap.SetBit(string(randkv.GetTestKey(1)), 0, true)
	assert.Nil(t, bitmapErr)
	value, err := bitmap.GetBit(string(randkv.GetTestKey(1)), 0)
	assert.Nil(t, err)
	assert.Equal(t, true, value)
}

func TestBitmapStructure_GetBits(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test GetBits function
	bitmapErr = bitmap.SetBits(string(randkv.GetTestKey(1)), 0, true, 1, false)
	assert.Nil(t, bitmapErr)
	values, err := bitmap.GetBits(string(randkv.GetTestKey(1)), 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, []bool{true, false}, values)
}

func TestBitmapStructure_BitCount(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test BitCount function
	bitmapErr = bitmap.SetBits(string(randkv.GetTestKey(1)), 0, true, 1, false, 2, true)
	assert.Nil(t, bitmapErr)
	count, err := bitmap.BitCount(string(randkv.GetTestKey(1)), 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}

func TestBitmapStructure_BitOp(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test BitOp function
	bitmapErr = bitmap.SetBits(string(randkv.GetTestKey(1)), 0, true, 1, false)
	assert.Nil(t, bitmapErr)
	bitmapErr = bitmap.SetBits(string(randkv.GetTestKey(2)), 0, false, 1, true)
	assert.Nil(t, bitmapErr)
	bitmapErr = bitmap.BitOp([]byte("AND"), string(randkv.GetTestKey(3)), string(randkv.GetTestKey(1)), string(randkv.GetTestKey(2)))
	assert.Nil(t, bitmapErr)
	values, err := bitmap.GetBits(string(randkv.GetTestKey(3)), 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, []bool{false, false}, values)
}

func TestBitmapStructure_BitDel(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test BitDel function
	bitmapErr = bitmap.SetBits(string(randkv.GetTestKey(1)), 0, true, 1, false, 2, true)
	assert.Nil(t, bitmapErr)
	bitmapErr = bitmap.BitDel(string(randkv.GetTestKey(1)), 1)
	assert.Nil(t, bitmapErr)
	values, err := bitmap.GetBits(string(randkv.GetTestKey(1)), 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, []bool{true, true}, values)
}

func TestBitmapStructure_BitDels(t *testing.T) {
	bitmap := initBitmap()
	defer bitmap.db.Clean()

	// Test BitDels function
	bitmapErr = bitmap.SetBits(string(randkv.GetTestKey(1)), 0, true, 1, false, 2, true)
	assert.Nil(t, bitmapErr)
	bitmapErr = bitmap.BitDels(string(randkv.GetTestKey(1)), 0, 1)
	assert.Nil(t, bitmapErr)
	values, err := bitmap.GetBits(string(randkv.GetTestKey(1)), 0)
	assert.Nil(t, err)
	assert.Equal(t, []bool{true}, values)
}

func TestBitmapStructure_Integration(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestBitmapStructure")
	opts.DirPath = dir
	bitmap, _ := NewBitmapStructure(opts)
	defer bitmap.db.Clean()

	key := []byte("testKey")

	// Test with large amount of data
	for i := uint(0); i < 10000; i++ {
		err := bitmap.SetBit(string(key), i, i%2 == 0)
		assert.Nil(t, err)
	}

	for i := uint(0); i < 10000; i++ {
		value, err := bitmap.GetBit(string(key), i)
		assert.Nil(t, err)
		assert.Equal(t, i%2 == 0, value)
	}

	// Test boundary conditions
	err := bitmap.SetBit(string(key), 0, true)
	assert.Nil(t, err)
	value, err := bitmap.GetBit(string(key), 0)
	assert.Nil(t, err)
	assert.True(t, value)

	err = bitmap.SetBit(string(key), 1000000, true)
	assert.Nil(t, err)
	value, err = bitmap.GetBit(string(key), 1000000)
	assert.Nil(t, err)
	assert.True(t, value)

	// Test error handling
	err = bitmap.SetBit(string(key), 1000001, true)
	assert.Nil(t, err)
	value, err = bitmap.GetBit(string(key), 1000002)
	assert.Nil(t, err)
	assert.False(t, value)

	// Test BitOp
	key1 := []byte("testKey1")
	key2 := []byte("testKey2")
	destKey := []byte("destKey")

	err = bitmap.SetBits(string(key1), 0, true, 1, false, 2, true)
	assert.Nil(t, err)
	err = bitmap.SetBits(string(key2), 0, false, 1, true, 2, false)
	assert.Nil(t, err)

	err = bitmap.BitOp([]byte("AND"), string(destKey), string(key1), string(key2))
	assert.Nil(t, err)
	values, err := bitmap.GetBits(string(destKey), 0, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, []bool{false, false, false}, values)

	err = bitmap.BitOp([]byte("OR"), string(destKey), string(key1), string(key2))
	assert.Nil(t, err)
	values, err = bitmap.GetBits(string(destKey), 0, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, []bool{true, true, true}, values)

	err = bitmap.BitOp([]byte("XOR"), string(destKey), string(key1), string(key2))
	assert.Nil(t, err)
	values, err = bitmap.GetBits(string(destKey), 0, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, []bool{true, true, true}, values)

	err = bitmap.BitOp([]byte("NOT"), string(destKey), string(key1))
	assert.Nil(t, err)
	values, err = bitmap.GetBits(string(destKey), 0, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, []bool{false, true, false}, values)
}
