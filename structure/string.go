package structure

import (
	"encoding/binary"
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"strconv"
	"time"
)

// StringStructure is a structure that stores string data
type StringStructure struct {
	db *engine.DB
}

// NewStringStructure returns a new StringStructure
// It will return a nil StringStructure if the database cannot be opened
// or the database cannot be created
// The database will be created if it does not exist
func NewStringStructure(options config.Options) (*StringStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &StringStructure{db: db}, nil
}

// Set sets the value of a key
// If the key does not exist, it will be created
// If the key exists, it will be overwritten
// If the key is expired, it will be deleted
// If the key is not expired, it will be updated
func (s *StringStructure) Set(key, value []byte, ttl time.Duration) error {
	if value == nil {
		return nil
	}

	// Encode the value
	encValue, err := encodeStringValue(value, ttl)
	if err != nil {
		return err
	}

	// Set the value
	return s.db.Put(key, encValue)
}

// Get gets the value of a key
// If the key does not exist, it will return nil
// If the key exists, it will return the value
// If the key is expired, it will be deleted and return nil
// If the key is not expired, it will be updated and return the value
func (s *StringStructure) Get(key []byte) ([]byte, error) {
	// Get the value
	value, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	//Decode the value
	return decodeStringValue(value)
}

// Del deletes the value of a key
// If the key does not exist, it will return nil
// If the key exists, it will be deleted
// If the key is expired, it will be deleted and return nil
// If the key is not expired, it will be updated and return nil
func (s *StringStructure) Del(key []byte) error {
	// Delete the value
	return s.db.Delete(key)
}

// Type returns the type of a key
// If the key does not exist, it will return ""
// If the key exists, it will return "string"
// If the key is expired, it will be deleted and return ""
// If the key is not expired, it will be updated and return "string"
func (s *StringStructure) Type(key []byte) (string, error) {
	// Get the value
	value, err := s.db.Get(key)
	if err != nil {
		return "", err
	}

	// Decode the value
	_, err = decodeStringValue(value)
	if err != nil {
		return "", err
	}

	// Return the type
	return "string", nil
}

// StrLen returns the length of the value of a key
// If the key does not exist, it will return 0
// If the key exists, it will return the length of the value
// If the key is expired, it will be deleted and return 0
// If the key is not expired, it will be updated and return the length of the value
func (s *StringStructure) StrLen(key []byte) (int, error) {
	// Get the value
	value, err := s.db.Get(key)
	if err != nil {
		return 0, err
	}

	// Decode the value
	value, err = decodeStringValue(value)
	if err != nil {
		return 0, err
	}

	// Return the length of the value
	return len(value), nil
}

// GetSet sets the value of a key and returns its old value
func (s *StringStructure) GetSet(key, value []byte, ttl time.Duration) ([]byte, error) {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return nil, err
	}

	// Set the value
	err = s.Set(key, value, ttl)
	if err != nil {
		return nil, err
	}

	// Return the old value
	return oldValue, nil
}

// Append appends a value to the value of a key
func (s *StringStructure) Append(key, value []byte, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Append the value
	newValue := append(oldValue, value...)

	// Set the value
	return s.Set(key, newValue, ttl)
}

// Incr increments the integer value of a key by 1
func (s *StringStructure) Incr(key []byte, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	oldIntValue, err := strconv.Atoi(string(oldValue))
	if err != nil {
		return err
	}

	// Increment the integer value
	newIntValue := oldIntValue + 1

	// Convert the new integer value to a byte slice
	newValue := []byte(strconv.Itoa(newIntValue))

	// Set the value
	return s.Set(key, newValue, ttl)
}

// IncrBy increments the integer value of a key by the given amount
func (s *StringStructure) IncrBy(key []byte, amount int, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	oldIntValue, err := strconv.Atoi(string(oldValue))
	if err != nil {
		return err
	}

	// Increment the integer value
	newIntValue := oldIntValue + amount

	// Convert the new integer value to a byte slice
	newValue := []byte(strconv.Itoa(newIntValue))

	// Set the value
	return s.Set(key, newValue, ttl)
}

// IncrByFloat increments the float value of a key by the given amount
func (s *StringStructure) IncrByFloat(key []byte, amount float64, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to a float
	oldFloatValue, err := strconv.ParseFloat(string(oldValue), 64)
	if err != nil {
		return err
	}

	// Increment the float value
	newFloatValue := oldFloatValue + amount

	// Convert the new float value to a byte slice
	newValue := []byte(strconv.FormatFloat(newFloatValue, 'f', -1, 64))

	// Set the value
	return s.Set(key, newValue, ttl)
}

// Decr decrements the integer value of a key by 1
func (s *StringStructure) Decr(key []byte, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	oldIntValue, err := strconv.Atoi(string(oldValue))
	if err != nil {
		return err
	}

	// Decrement the integer value
	newIntValue := oldIntValue - 1

	// Convert the new integer value to a byte slice
	newValue := []byte(strconv.Itoa(newIntValue))

	// Set the value
	return s.Set(key, newValue, ttl)
}

// DecrBy decrements the integer value of a key by the given amount
func (s *StringStructure) DecrBy(key []byte, amount int, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	oldIntValue, err := strconv.Atoi(string(oldValue))
	if err != nil {
		return err
	}

	// Decrement the integer value
	newIntValue := oldIntValue - amount

	// Convert the new integer value to a byte slice
	newValue := []byte(strconv.Itoa(newIntValue))

	// Set the value
	return s.Set(key, newValue, ttl)
}

// Exists checks if a key exists
func (s *StringStructure) Exists(key []byte) (bool, error) {
	// Get the value
	_, err := s.Get(key)
	if err != nil {
		if err == _const.ErrKeyNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Expire sets the expiration time of a key
func (s *StringStructure) Expire(key []byte, ttl time.Duration) error {
	// Get the value
	value, err := s.Get(key)
	if err != nil {
		return err
	}

	// Set the value
	return s.Set(key, value, ttl)
}

// encodeStringValue encodes the value
// format: [type][expire][value]
// type: 1 byte
// expire: 8 bytes
// value: n bytes
func encodeStringValue(value []byte, ttl time.Duration) ([]byte, error) {
	// Create a byte slice buf with a length of binary.MaxVarintLen64
	// to hold the encoded value and additional data.
	buf := make([]byte, binary.MaxVarintLen64)

	// Set the first element of buf to represent the data structure type as String.
	buf[0] = String

	// Use the variable bufIndex to keep track of the current index position in buf,
	// starting from 1 to indicate the number of bytes written so far.
	var bufIndex = 1

	// The variable expire is used to store the expiration time, initially set to 0.
	var expire int64 = 0

	// Calculate the expiration time by adding ttl to the current time,
	// convert it to nanoseconds, and store it in the expire variable.
	if ttl != 0 {
		expire = time.Now().Add(ttl).UnixNano()
	}

	// Encode the expiration time expire as a variable-length integer
	// and write it to the sub-slice of byte slice buf starting
	// from the current index position bufIndex.
	bufIndex += binary.PutVarint(buf[bufIndex:], expire)

	// Create a byte slice encValue with a length of bufIndex + len(value)
	// to hold the encoded value and additional data.
	encValue := make([]byte, bufIndex+len(value))

	// Copy the encoded value from the beginning of buf
	// to the corresponding position in encValue.
	copy(encValue[:bufIndex], buf[:bufIndex])

	// Copy the original value value to the remaining positions in encValue,
	// starting from the index bufIndex.
	copy(encValue[bufIndex:], value)

	return encValue, nil
}

var (
	// ErrInvalidValue is returned if the value is invalid.
	ErrInvalidValue = errors.New("Wrong value: invalid value")
	// ErrInvalidType is returned if the type is invalid.
	ErrInvalidType = errors.New("Wrong value: invalid type")
	// ErrKeyExpired is returned if the key is expired.
	ErrKeyExpired = errors.New("Wrong value: key expired")
)

// decodeStringValue decodes the value
// format: [type][expire][value]
// type: 1 byte
// expire: 8 bytes
// value: n bytes
func decodeStringValue(value []byte) ([]byte, error) {
	// Check the length of the value
	if len(value) < 1 {
		return nil, ErrInvalidValue
	}

	// Check the type of the value
	if value[0] != String {
		return nil, ErrInvalidType
	}

	// Use the variable bufIndex to keep track of the current index position in value,
	// starting from 1 to indicate the number of bytes read so far.
	var bufIndex = 1

	// Decode the expiration time expire from the sub-slice of byte slice value
	// starting from the current index position bufIndex.
	expire, n := binary.Varint(value[bufIndex:])

	// Check the number of bytes read
	if n <= 0 {
		return nil, ErrInvalidValue
	}

	// Update the current index position bufIndex by adding the number of bytes read n.
	bufIndex += n

	// Check the expiration time expire
	if expire != 0 && expire < time.Now().UnixNano() {
		return nil, ErrKeyExpired
	}

	// Return the original value value
	return value[bufIndex:], nil
}
