package structure

import (
	"bytes"
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/bits-and-blooms/bitset"
)

type BitmapStructure struct {
	db *engine.DB
}
type BitOperation string

const (
	BitAndOperation BitOperation = "AND"
	BitOrOperation  BitOperation = "OR"
	BitXorOperation BitOperation = "XOR"
	BitNotOperation BitOperation = "NOT"
)

func NewBitmap(options config.Options) (*BitmapStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &BitmapStructure{db: db}, nil
}

func (b *BitmapStructure) SetBit(key []byte, off uint) error {
	bit := bitset.New(1)
	value, err := b.GetBits(key)
	if err == _const.ErrKeyNotFound {
		buf, _ := bitsetToByteArray(bit)
		err := b.db.Put(key, buf)
		if err != nil {
			return err
		}
		value, err = b.GetBits(key)
	} else if err != nil {
		return err
	}
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	bit.Set(off)
	buf, err := bitsetToByteArray(bit)
	return b.db.Put(key, buf)
}
func (b *BitmapStructure) SetBits(key []byte, off ...uint) error {
	bit := bitset.New(1)
	value, err := b.GetBits(key)
	if err == _const.ErrKeyNotFound {
		buf, _ := bitsetToByteArray(bit)
		err := b.db.Put(key, buf)
		if err != nil {
			return err
		}
		value, err = b.GetBits(key)
	} else if err != nil {
		return err
	}
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	for _, u := range off {
		bit.Set(u)
	}
	buf, err := bitsetToByteArray(bit)
	return b.db.Put(key, buf)
}

func (b *BitmapStructure) GetBits(key []byte) ([]byte, error) {
	value, err := b.db.Get(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BitmapStructure) GetBit(key []byte, off uint) (bool, error) {
	value, err := b.db.Get(key)
	if err != nil {
		return false, err
	}
	bit := bitset.New(1)
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return false, err
	}

	return bit.Test(off), nil
}

func (b *BitmapStructure) DelBit(key []byte, off uint) error {
	value, err := b.db.Get(key)
	if err != nil {
		return err
	}
	bit := bitset.New(1)
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	bit.Clear(off)
	buff, err := bitsetToByteArray(bit)
	if err != nil {
		return err
	}
	err = b.db.Put(key, buff)
	if err != nil {
		return err
	}

	return nil
}
func (b *BitmapStructure) DelBits(key []byte, off ...uint) error {
	value, err := b.db.Get(key)
	if err != nil {
		return err
	}
	bit := bitset.New(1)
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	// delete all the provided ranges
	for _, u := range off {
		bit.Clear(u)
	}

	buff, err := bitsetToByteArray(bit)
	if err != nil {
		return err
	}
	err = b.db.Put(key, buff)
	if err != nil {
		return err
	}

	return nil
}
func (b *BitmapStructure) BitCount(key []byte, start, end uint) (uint, error) {
	value, err := b.db.Get(key)
	if err != nil {
		return 0, err
	}
	bit := bitset.New(1)
	err = bit.UnmarshalBinary(value)
	if err != nil {
		return 0, err
	}
	idx := start
	var count uint

	for i, e := bit.NextSet(idx); e; i, e = bit.NextSet(i + 1) {
		if i > end {
			break
		}
		count++
	}

	return count, nil
}
func (b *BitmapStructure) BitOp(op BitOperation, destKey []byte, keys ...[]byte) error {
	if len(keys) == 0 {
		return errors.New("no keys specified")
	}

	value, err := b.db.Get(keys[0])
	if err != nil {
		return err
	}
	baseBit := bitset.New(1)
	err = baseBit.UnmarshalBinary(value)
	if err != nil {
		return err
	}
	for i := 1; i < len(keys); i++ {
		value, err := b.db.Get(keys[i])
		if err != nil {
			return err
		}
		bit := bitset.New(1)
		err = bit.UnmarshalBinary(value)
		if err != nil {
			return err
		}
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

	buf, err := bitsetToByteArray(baseBit)
	if err != nil {
		return err
	}
	err = b.db.Put(destKey, buf)
	if err != nil {
		return err
	}

	return nil
}
func (b *BitmapStructure) getBits(key []byte) ([]byte, error) {
	value, err := b.db.Get(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func bitsetToByteArray(set *bitset.BitSet) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := set.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
