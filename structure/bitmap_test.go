package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/bits-and-blooms/bitset"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initBitmap() *BitmapStructure {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestBitmapStructure")
	opts.DirPath = dir
	bitmap, _ := NewBitmap(opts)
	return bitmap
}

func TestNewBitmap(t *testing.T) {

}

func TestBitmapStructure_SetBit(t *testing.T) {
	// Setup the bitmap structure
	bm := initBitmap()
	bs := bitset.New(1)
	err := bm.SetBit("1", 7)
	assert.NoError(t, err)
	val, err := bm.GetBits("1")
	err = bs.UnmarshalBinary(val)

	assert.NoError(t, err)
	assert.True(t, bs.Test(7))

}
func TestBitmapStructure_SetBits(t *testing.T) {
	// Setup the bitmap structure
	bm := initBitmap()
	bs := bitset.New(1)
	err := bm.SetBits("1", 7, 4, 2, 9)
	assert.NoError(t, err)
	val, err := bm.GetBits("1")
	err = bs.UnmarshalBinary(val)

	assert.NoError(t, err)
	assert.True(t, bs.Test(7))
	assert.True(t, bs.Test(9))
	assert.True(t, bs.Test(4))
	assert.True(t, bs.Test(2))
	assert.False(t, bs.Test(1))
	assert.False(t, bs.Test(10))

}

func TestBitmapStructure_GetBits(t *testing.T) {
	// Setup the bitmap structure
	bm := initBitmap()

	err := bm.SetBit("1", 2)
	assert.NoError(t, err)
	val, err := bm.GetBits("1")
	assert.NoError(t, err)
	bs := bitset.New(1)
	err = bs.UnmarshalBinary(val)
	assert.NoError(t, err)
	assert.True(t, bs.Test(2))

	// add one other bit
	err = bm.SetBit("1", 7)
	val, err = bm.GetBits("1")
	assert.NoError(t, err)
	err = bs.UnmarshalBinary(val)
	assert.NoError(t, err)
	assert.True(t, bs.Test(7))
	assert.False(t, bs.Test(8))
}

func TestBitmapStructure_GetBit(t *testing.T) {
	// Setup the bitmap structure
	bm := initBitmap()

	err := bm.SetBit("1", 2)
	assert.NoError(t, err)
	val, err := bm.GetBit("1", 2)
	assert.NoError(t, err)
	assert.True(t, val)
	val, _ = bm.GetBit("1", 1)
	assert.False(t, val)
}

func TestBitmapStructure_DelBit(t *testing.T) {
	bm := initBitmap()

	err := bm.SetBits("1", 2, 3, 4, 5, 6)
	assert.NoError(t, err)
	val, _ := bm.GetBit("1", 2)
	assert.True(t, val)
	err = bm.DelBit("1", 2)
	assert.NoError(t, err)
	val, err = bm.GetBit("1", 2)
	assert.NoError(t, err)
	assert.False(t, val)
	val, _ = bm.GetBit("1", 4)
	assert.True(t, val)

}
func TestBitmapStructure_DelBits(t *testing.T) {
	bm := initBitmap()

	err := bm.SetBits("1", 2, 3, 4, 5, 6)
	assert.NoError(t, err)
	err = bm.DelBits("1", 4, 5)
	assert.NoError(t, err)
	val, err := bm.GetBit("1", 2)
	assert.NoError(t, err)
	assert.True(t, val)
	val, _ = bm.GetBit("1", 4)
	assert.False(t, val)
	val, _ = bm.GetBit("1", 5)
	assert.False(t, val)
	val, _ = bm.GetBit("1", 6)
	assert.True(t, val)

}

func TestBitmapStructure_BitCount(t *testing.T) {
	bm := initBitmap()

	err := bm.SetBits("1", 2, 3, 4, 5, 6, 10, 11, 12, 13, 14, 15)
	assert.NoError(t, err)
	count, err := bm.BitCount("1", 3, 10)
	assert.NoError(t, err)
	assert.Equal(t, uint(5), count)
	// check over the end
	count, err = bm.BitCount("1", 10, 100)
	assert.NoError(t, err)
	assert.Equal(t, uint(6), count)
}

func TestBitmapStructure_BitOp(t *testing.T) {
	type bitmapStruct struct {
		key string
		arr []uint
	}

	tests := []struct {
		name     string
		op       BitOperation
		input    []bitmapStruct
		expected bitmapStruct
	}{
		{
			name: "bitset OR operation",
			op:   BitOrOperation,
			input: []bitmapStruct{
				{"1", []uint{2, 3}},
				{"2", []uint{0, 6, 12, 17}},
				{"3", []uint{1, 4, 10, 11, 19}},
				{"4", []uint{10, 12, 14}},
			},
			expected: bitmapStruct{
				"5",
				[]uint{0, 1, 2, 3, 4, 6, 10, 11, 12, 14, 17, 19},
			},
		},
		{
			name: "bitset AND operation",
			op:   BitAndOperation,
			input: []bitmapStruct{
				{"1", []uint{2, 3, 6}},
				{"2", []uint{0, 6, 12, 17}},
				{"3", []uint{1, 4, 10, 11, 6, 19}},
				{"4", []uint{10, 12, 14, 6}},
			},
			expected: bitmapStruct{
				"5",
				[]uint{6},
			},
		},
		{
			name: "bitset NOT operation",
			op:   BitNotOperation,
			input: []bitmapStruct{
				{"1", []uint{0, 1, 2, 3, 12}},
				{"2", []uint{0, 6, 12, 17}},
				{"3", []uint{1, 4, 10, 11, 19}},
				{"4", []uint{10, 12, 14}},
			},
			expected: bitmapStruct{
				"5",
				[]uint{2, 3},
			},
		},
		{
			name: "bitset XOR operation",
			op:   BitXorOperation,
			input: []bitmapStruct{
				{"1", []uint{2, 3}},
				{"2", []uint{0, 2, 6, 12, 17}},
				{"3", []uint{1, 4, 10, 11, 19}},
				{"4", []uint{3, 10, 12, 14}},
			},
			expected: bitmapStruct{
				"5",
				[]uint{0, 1, 4, 6, 11, 14, 17, 19},
			},
		},
	}

	for _, test := range tests {
		bm := initBitmap()
		var keys []string
		for _, b := range test.input {
			err = bm.SetBits(b.key, b.arr...)
			assert.NoError(t, err)
			keys = append(keys, b.key)
		}
		err := bm.BitOp(test.op, test.expected.key, keys...)
		assert.NoError(t, err)
		bits, err := bm.getBits([]byte(test.expected.key))

		actualBit := bitset.New(1)
		err = actualBit.UnmarshalBinary(bits)
		expectedBit := bitset.New(1)
		for _, u := range test.expected.arr {
			expectedBit.Set(u)
		}
		assert.Equal(t, expectedBit.String(), actualBit.String())
	}

}
