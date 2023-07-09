package structure

// Importing necessary packages
import (
	"bytes"
	"errors"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/bits-and-blooms/bitset"
)

// BitmapStructure is a structure that holds a reference to the database engine
type BitmapStructure struct {
	db *engine.DB
}

// BitSet is a structure that holds a reference to a BitSet
type BitSet struct {
	b *bitset.BitSet
}

// BitOperation is a type used to define the operations that can be performed on bits
type BitOperation string

// Constants representing the different bit operations
const (
	// BitAndOperation performs the AND operation.
	BitAndOperation BitOperation = "AND"
	// BitOrOperation performs the OR operation.
	BitOrOperation BitOperation = "OR"
	// BitXorOperation performs the XOR operation.
	BitXorOperation BitOperation = "XOR"
	// BitNotOperation performs the NOT operation.
	BitNotOperation BitOperation = "NOT"
)

// NewBitmap function initializes a new BitmapStructure with the provided options
func NewBitmap(options config.Options) (*BitmapStructure, error) {
	// Create a new database with the provided options
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	// Return a new BitmapStructure with the created database
	return &BitmapStructure{db: db}, nil
}

// SetBit function sets a bit at the specified offset in the bitmap for the provided key
func (b *BitmapStructure) SetBit(k string, off uint) error {
	// Convert the key to bytes
	key := stringToBytesWithKey(k)
	// Create a new bitset
	bit := bitset.New(1)
	// Get the current bits for the key
	value, err := b.getBits(key)
	// If the key is not found, initialize it with the new bitset
	if err == _const.ErrKeyNotFound {
		buf, _ := bitsetToByteArray(bit)
		err = b.db.Put(key, buf)
		if err != nil {
			return err
		}
		value, err = b.getBits(key)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	// Unmarshal the current bits into the bitset
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	// Set the bit at the specified offset
	bit.Set(off)
	// Convert the bitset to a byte array
	buf, err := bitsetToByteArray(bit)
	if err != nil {
		return err
	}
	// Store the updated bitset in the database
	return b.db.Put(key, buf)
}

// SetBits function sets multiple bits at the specified offsets in the bitmap for the provided key
func (b *BitmapStructure) SetBits(k string, off ...uint) error {
	// Convert the key to bytes
	key := stringToBytesWithKey(k)
	// Create a new bitset
	bit := bitset.New(1)
	// Get the current bits for the key
	value, err := b.getBits(key)
	if err != nil && !errors.Is(err, _const.ErrKeyNotFound) {
		return err
	}
	// If the key is not found, initialize it with the new bitset
	if err == _const.ErrKeyNotFound {
		buf, _ := bitsetToByteArray(bit)
		err = b.db.Put(key, buf)
		if err != nil {
			return err
		}
		value, err = b.getBits(key)
		if err != nil {
			return err
		}
	}

	// Unmarshal the current bits into the bitset
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	// Set the bits at the specified offsets
	for _, u := range off {
		bit.Set(u)
	}
	// Convert the bitset to a byte array
	buf, err := bitsetToByteArray(bit)
	if err != nil {
		return err
	}
	// Store the updated bitset in the database
	return b.db.Put(key, buf)
}

// GetBits function retrieves the bits for the provided key
func (b *BitmapStructure) GetBits(k string) (*BitSet, error) {
	// Convert the key to bytes
	key := stringToBytesWithKey(k)
	// Get the bits from the database
	value, err := b.db.Get(key)
	if err != nil {
		return nil, err
	}
	// Create a new bitset from the retrieved bits
	bit, err := newBitSet(value)
	if err != nil {
		return nil, err
	}
	// Return the bitset
	return bit, nil
}

// GetBit function retrieves a bit at the specified offset in the bitmap for the provided key
func (b *BitmapStructure) GetBit(k string, off uint) (bool, error) {
	// Convert the key to bytes
	key := stringToBytesWithKey(k)
	// Get the bits from the database
	value, err := b.db.Get(key)
	if err != nil {
		return false, err
	}
	// Create a new bitset
	bit := bitset.New(1)
	// Unmarshal the retrieved bits into the bitset
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return false, err
	}
	// Return the bit at the specified offset
	return bit.Test(off), nil
}

// DelBit function deletes a bit at the specified offset in the bitmap for the provided key
func (b *BitmapStructure) DelBit(k string, off uint) error {
	// Convert the key to bytes
	key := stringToBytesWithKey(k)
	// Get the bits from the database
	value, err := b.db.Get(key)
	if err != nil {
		return err
	}
	// Create a new bitset
	bit := bitset.New(1)
	// Unmarshal the retrieved bits into the bitset
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	// Clear the bit at the specified offset
	bit.Clear(off)
	// Convert the bitset to a byte array
	buff, err := bitsetToByteArray(bit)
	if err != nil {
		return err
	}
	// Store the updated bitset in the database
	err = b.db.Put(key, buff)
	if err != nil {
		return err
	}
	// Return nil if no errors occurred
	return nil
}

// DelBits function deletes multiple bits at the specified offsets in the bitmap for the provided key
func (b *BitmapStructure) DelBits(k string, off ...uint) error {
	// Convert the key to bytes
	key := stringToBytesWithKey(k)
	// Get the bits from the database
	value, err := b.db.Get(key)
	if err != nil {
		return err
	}
	// Create a new bitset
	bit := bitset.New(1)
	// Unmarshal the retrieved bits into the bitset
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	// Clear the bits at the specified offsets
	for _, u := range off {
		bit.Clear(u)
	}
	// Convert the bitset to a byte array
	buff, err := bitsetToByteArray(bit)
	if err != nil {
		return err
	}
	// Store the updated bitset in the database
	err = b.db.Put(key, buff)
	if err != nil {
		return err
	}
	// Return nil if no errors occurred
	return nil
}

// BitCount function counts the number of set bits in the specified range in the bitmap for the provided key
func (b *BitmapStructure) BitCount(k string, start, end uint) (uint, error) {
	// Convert the key to bytes
	key := stringToBytesWithKey(k)
	// Get the bits from the database
	value, err := b.db.Get(key)
	if err != nil {
		return 0, err
	}
	// Create a new bitset
	bit := bitset.New(1)
	// Unmarshal the retrieved bits into the bitset
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return 0, err
	}
	// Initialize the count to 0
	var count uint
	// Iterate over the set bits in the specified range and increment the count
	for i, e := bit.NextSet(start); e; i, e = bit.NextSet(i + 1) {
		if i > end {
			break
		}
		count++
	}
	// Return the count
	return count, nil
}

// BitOp function performs the specified bit operation on the provided keys and stores
// the result in the destination key
func (b *BitmapStructure) BitOp(op BitOperation, destKey string, keys ...string) error {
	// Check if any keys were provided
	if len(keys) == 0 {
		return errors.New("no keys specified")
	}
	// Get the bits for the first key
	value, err := b.db.Get(stringToBytesWithKey(keys[0]))
	if err != nil {
		return err
	}
	// Create a new bitset and unmarshal the retrieved bits into it
	baseBit := bitset.New(1)
	err = baseBit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	// Iterate over the remaining keys
	for i := 1; i < len(keys); i++ {
		// Get the bits for the current key
		value, err := b.db.Get(stringToBytesWithKey(keys[i]))
		if err != nil {
			return err
		}
		// Create a new bitset and unmarshal the retrieved bits into it
		bit := bitset.New(1)
		err = bit.UnmarshalBinary(value)
		if err != nil {
			return err
		}
		// Perform the specified bit operation
		switch op {
		case BitOrOperation:
			baseBit = baseBit.Union(bit)
		case BitAndOperation:
			baseBit = baseBit.Intersection(bit)
		case BitXorOperation:
			baseBit = baseBit.SymmetricDifference(bit)
		case BitNotOperation:
			baseBit = baseBit.Difference(bit)
		}
	}
	// Convert the resulting bitset to a byte array
	buf, err := bitsetToByteArray(baseBit)
	if err != nil {
		return err
	}
	// Store the result in the destination key
	err = b.db.Put(stringToBytesWithKey(destKey), buf)
	if err != nil {
		return err
	}
	// Return nil if no errors occurred
	return nil
}

// getBits function retrieves the bits for the provided key
func (b *BitmapStructure) getBits(key []byte) ([]byte, error) {
	// Get the bits from the database
	value, err := b.db.Get(key)
	if err != nil {
		return nil, err
	}
	// Return the retrieved bits
	return value, nil
}

// bitsetToByteArray function converts a bitset to a byte array
func bitsetToByteArray(set *bitset.BitSet) ([]byte, error) {
	// Create a new buffer
	buf := new(bytes.Buffer)
	// Write the bitset to the buffer
	_, err := set.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	// Return the bytes from the buffer
	return buf.Bytes(), nil
}

// newBitSet function creates a new bitset from the provided bytes
func newBitSet(d []byte) (*BitSet, error) {
	// Create a new bitset
	b := bitset.New(1)
	// Unmarshal the provided bytes into the bitset
	err := b.UnmarshalBinary(d)
	if err != nil {
		return nil, err
	}
	// Return the bitset
	return &BitSet{b: b}, nil
}

// At function checks if the bit at the specified position is set
func (b *BitSet) At(pos uint) bool {
	return b.b.Test(pos)
}

// Next function finds the next set bit in the bitset
func (b *BitSet) Next(pos uint) (uint, bool) {
	return b.b.NextSet(pos + 1)
}
