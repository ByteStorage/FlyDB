package structure

import (
	"encoding/binary"
	"errors"
	"sort"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
)

type BitmapStructure struct {
	db *engine.DB
}

// NewBitmapStructure returns a new BitmapStructure
// It will return a nil BitmapStructure if the database cannot be opened
// or the database cannot be created
// The database will be created if it does not exist
func NewBitmapStructure(options config.Options) (*BitmapStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &BitmapStructure{db: db}, nil
}

// SetBit Set the bit at the specified offset in the bitmap to the specified value
// If the key does not exist, it will be created
// If the bitmap length is not sufficient, it will be extended
func (b *BitmapStructure) SetBit(key []byte, offset uint, value bool) error {
	// Get bitmap
	bitmap, length, err := b.getBitmapFromDB(key, true)

	if err != nil {
		return err
	}

	// If the bitmap length is not sufficient and value is 1, extend it
	if offset >= uint(length) {
		if value {
			length = offset + 1
			newSize := len2Size(length)
			newBitmap := make([]byte, newSize)
			copy(newBitmap, bitmap)
			bitmap = newBitmap
		} else { // If the value is 0, no operation is needed
			return nil
		}
	}
	b.setBit(bitmap, offset, value)

	return b.setBitmapToDB(key, bitmap, length)
}

// SetBits Set the bit at the specified offset in the bitmap to the specified value
// If the key does not exist, it will be created
// If the bitmap length is not sufficient, it will be extended
func (b *BitmapStructure) SetBits(key []byte, args ...interface{}) error {
	// Validate arguments
	if len(args) == 0 || len(args)%2 != 0 {
		return ErrInvalidArgs
	}

	// Get bitmap
	bitmap, length, err := b.getBitmapFromDB(key, true)
	if err != nil {
		return err
	}

	offsets := make([]uint, 0, len(args)/2)
	values := make([]bool, 0, len(args)/2)
	maxNotFalseOffset := uint(0)
	for i := 0; i < len(args); i += 2 {
		offsets = append(offsets, uint(args[i].(int)))
		values = append(values, args[i+1].(bool))
		if args[i+1].(bool) && uint(args[i].(int)) > maxNotFalseOffset {
			maxNotFalseOffset = uint(args[i].(int))
		}
	}

	// If the bitmap length is not sufficient, extend it
	if maxNotFalseOffset >= length {
		length = maxNotFalseOffset + 1
		newSize := len2Size(length)
		newBitmap := make([]byte, newSize)
		copy(newBitmap, bitmap)
		bitmap = newBitmap
	}

	for i := 0; i < len(args)/2; i++ {
		offset := offsets[i]
		if offset < length {
			b.setBit(bitmap, offset, values[i])
		}
	}

	return b.setBitmapToDB(key, bitmap, length)
}

// GetBit Get the value of the bit at the specified offset in the bitmap
// If the key does not exist, it returns an error
func (b *BitmapStructure) GetBit(key []byte, offset uint) (bool, error) {
	bitmap, length, err := b.getBitmapFromDB(key, false)
	if err != nil {
		return false, err
	}

	// If the offset is out of range, default value is false
	if offset >= length {
		return false, nil
	}

	return b.getBit(bitmap, offset), nil
}

// GetBits Get the values of a group of bits at the specified offsets in the bitmap
// If the key does not exist, it returns an error
func (b *BitmapStructure) GetBits(key []byte, offsets ...uint) ([]bool, error) {
	if len(offsets) == 0 {
		return nil, ErrInvalidArgs
	}

	bitmap, length, err := b.getBitmapFromDB(key, false)
	if err != nil {
		return nil, err
	}

	result := make([]bool, len(offsets))

	for idx, offset := range offsets {
		// If the offset is out of range, default value is false
		if offset >= length {
			result[idx] = false
		} else {
			result[idx] = b.getBit(bitmap, offset)
		}
	}

	return result, nil
}

// BitCount Count the number of bits set to 1 in the specified range of the bitmap
// If the key does not exist, it returns an error
func (b *BitmapStructure) BitCount(key []byte, start uint, end uint) (int, error) {
	bitmap, length, err := b.getBitmapFromDB(key, false)
	if err != nil {
		return 0, err
	}

	total1 := 0

	// Iterate over the range
	for offset := start; offset <= end; offset++ {
		if offset < length {
			if b.getBit(bitmap, offset) {
				total1++
			}
		}
	}

	return total1, nil
}

// BitOp count the number of bits set to 1 in the specified range of the bitmap
// If the key does not exist, it returns an error
func (b *BitmapStructure) BitOp(operation []byte, destkey []byte, keys ...[]byte) error {
	// Check the validity of the parameters
	if string(operation) != "AND" && string(operation) != "OR" && string(operation) != "XOR" && string(operation) != "NOT" {
		return ErrInvalidArgs
	}

	if string(operation) == "NOT" { // Handle single and double parameters separately
		if len(keys) != 1 {
			return ErrInvalidValue
		}
		bitmap, length, err := b.getBitmapFromDB(keys[0], false)
		if err != nil {
			return err
		}
		result := make([]byte, len2Size(length))

		for i := uint(0); i < length; i++ {
			b.setBit(result, i, !b.getBit(bitmap, i))
		}
		return b.setBitmapToDB(destkey, result, length)
	} else {
		if len(keys) != 2 {
			return ErrInvalidValue
		}
		bitmap1, length1, err := b.getBitmapFromDB(keys[0], false)
		if err != nil {
			return err
		}
		bitmap2, length2, err := b.getBitmapFromDB(keys[1], false)
		if err != nil {
			return err
		}

		if string(operation) == "AND" {
			length := min(length1, length2)
			result := make([]byte, len2Size(length))

			for i := uint(0); i < length; i++ {
				b.setBit(result, i, b.getBit(bitmap1, i) && b.getBit(bitmap2, i))
			}

			return b.setBitmapToDB(destkey, result, length)
		} else if string(operation) == "OR" {
			length := max(length1, length2)
			minLength := min(length1, length2)
			result := make([]byte, len2Size(length))

			for i := uint(0); i < minLength; i++ {
				b.setBit(result, i, b.getBit(bitmap1, i) || b.getBit(bitmap2, i))
			}
			// Handle the remaining bits separately
			if length1 > length2 {
				for i := minLength; i < length; i++ {
					b.setBit(result, i, b.getBit(bitmap1, i))
				}
			} else {
				for i := minLength; i < length; i++ {
					b.setBit(result, i, b.getBit(bitmap2, i))
				}
			}
			return b.setBitmapToDB(destkey, result, length)
		} else if string(operation) == "XOR" {
			length := max(length1, length2)
			minLength := min(length1, length2)
			result := make([]byte, len2Size(length))

			for i := uint(0); i < minLength; i++ {
				b.setBit(result, uint(i), b.getBit(bitmap1, i) != b.getBit(bitmap2, i))
			}
			// Handle the remaining bits separately
			if length1 > length2 {
				for i := minLength; i < length; i++ {
					b.setBit(result, uint(i), b.getBit(bitmap1, i))
				}
			} else {
				for i := minLength; i < length; i++ {
					b.setBit(result, uint(i), b.getBit(bitmap2, i))
				}
			}
			return b.setBitmapToDB(destkey, result, length)
		}
	}
	return nil
}

// BitDel delete the bit at the specified offset in the bitmap
// If the key does not exist, it returns an error
// There is room for optimization
func (b *BitmapStructure) BitDel(key []byte, offset uint) error {

	bitmap, length, err := b.getBitmapFromDB(key, false)
	if err != nil {
		return err
	}
	// If offset is greater than or equal to length, no operation is needed
	if offset >= length {
		return nil
	}

	length--
	// Shift the bits one by one
	for i := offset; i < length; i++ {
		b.setBit(bitmap, i, b.getBit(bitmap, i+1))
	}

	return b.setBitmapToDB(key, bitmap, length)
}

// BitDels delete a group of bits at the specified offsets in the bitmap
// If the key does not exist, it returns an error
// There is room for optimization
func (b *BitmapStructure) BitDels(key []byte, offsets ...uint) error {
	// Check the parameters
	if len(offsets) == 0 {
		return ErrInvalidValue
	}

	bitmap, length, err := b.getBitmapFromDB(key, false)
	if err != nil {
		return err
	}

	sort.Slice(offsets, func(i, j int) bool {
		return offsets[i] < offsets[j]
	})

	// Number of deletions so far
	cntDel := uint(0)
	// Current offset to delete
	offsetIndex := 0

	length -= uint(len(offsets))
	for i := uint(0); i < length; i++ {
		if offsetIndex < len(offsets) && i+cntDel == offsets[offsetIndex] {
			cntDel++
			offsetIndex++
			i--
			continue
		}
		b.setBit(bitmap, i, b.getBit(bitmap, i+cntDel))
	}

	return b.setBitmapToDB(key, bitmap, length)
}

var (
	// ErrListEmpty is returned if the list is empty.
	ErrInvalidArgs = errors.New("Error Args: The args are wrong")
)

// Set the value of a specific bit, assuming the value is valid
func (b *BitmapStructure) setBit(bitmap []byte, offset uint, value bool) {
	index := offset / 8
	bit := uint(offset % 8)

	if value {
		// Set the bit at the offset to 1
		bitmap[index] |= 1 << bit
	} else {
		// Set the bit at the offset to 0
		bitmap[index] &= ^(1 << bit)
	}
}

// Get the value of a specific bit, assuming the value is valid
func (b *BitmapStructure) getBit(bitmap []byte, offset uint) bool {
	index := offset / 8
	bit := uint(offset % 8)

	value := bitmap[index] & (1 << bit)

	return value > 0
}

// getBitmapFromDB retrieves data from the database. When isKeyCanNotExist is true, it returns an empty slice if the key doesn't exist instead of an error.
func (b *BitmapStructure) getBitmapFromDB(key []byte, isKeyCanNotExist bool) ([]byte, uint, error) {
	if isKeyCanNotExist {
		// Get data corresponding to the key from the database
		dbData, err := b.db.Get(key)
		// Since the key might not exist, we need to handle ErrKeyNotFound separately as it is a valid case
		if err != nil && err != _const.ErrKeyNotFound {
			return nil, 0, err
		}
		// Deserialize the data into a bitmap
		bitmapArr, length, err := b.decodeBitmap(dbData)
		if err != nil {
			if len(dbData) != 0 {
				return nil, 0, err
			} else {
				bitmapArr = make([]byte, 0)
				length = 0
			}
		}
		return bitmapArr, length, nil
	} else {
		// Get data corresponding to the key from the database
		dbData, err := b.db.Get(key)
		if err != nil {
			return nil, 0, err
		}
		// Deserialize the data into a bitmap
		bitmapArr, length, err := b.decodeBitmap(dbData)
		if err != nil {
			return nil, 0, err
		}
		return bitmapArr, length, nil
	}
}

// setBitmapToDB stores the data into the database.
func (b *BitmapStructure) setBitmapToDB(key []byte, bm []byte, length uint) error {
	// Serialize into a binary array
	encValue, err := b.encodeBitmap(bm, length)
	if err != nil {
		return err
	}
	// Store in the database
	return b.db.Put(key, encValue)
}

// encodeBitmap encodes the value
// format: [type][length][value]
// length: the number of bits to save
// value: bitmap data
func (b *BitmapStructure) encodeBitmap(data []byte, length uint) ([]byte, error) {
	// Calculate the actual length to save
	dataSize := len2Size(length)

	// Either data or length is invalid
	if uint(len(data)) < dataSize {
		return nil, ErrInvalidValue
	}

	buf := make([]byte, 1+binary.MaxVarintLen64+int(dataSize))

	// Set the first element of buf to represent the data structure type as Bitmap.
	buf[0] = Bitmap

	bufIndex := 1
	bufIndex += binary.PutVarint(buf[bufIndex:], int64(length))

	// Append the data to the end
	bufIndex += copy(buf[bufIndex:], data[:dataSize])

	return buf[:bufIndex], nil
}

// decodeBitmap decodes the value
// format: [type][length][value]
// length: the number of bits to save
// value: bitmap data
func (b *BitmapStructure) decodeBitmap(value []byte) ([]byte, uint, error) {
	// Check the length of the value
	if len(value) < 2 {
		return nil, 0, ErrInvalidValue
	}

	// Check the type of the value
	if value[0] != Bitmap {
		return nil, 0, ErrInvalidType
	}

	valueLen := len(value)

	nowIndex := 1

	length, lenOfLen := binary.Varint(value[nowIndex:])

	// Check the number of bytes read
	if lenOfLen <= 0 {
		return nil, 0, ErrInvalidValue
	}

	nowIndex += lenOfLen

	// Actual size
	dataSize := len2Size(uint(length))

	// Either data or length is invalid
	if int(dataSize) > valueLen-nowIndex {
		return nil, 0, ErrInvalidValue
	}

	// Create an array to store the bitmap result
	result := make([]byte, dataSize)

	copy(result, value[nowIndex:])

	return result, uint(length), nil
}

func len2Size(len uint) uint {

	return uint((int(len)-1)/8 + 1)
}

func max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}
