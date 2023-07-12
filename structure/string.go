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
	db        *engine.DB
	valueType string
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

// stringToBytesWithKey converts a string to a byte slice
func stringToBytesWithKey(key string) []byte {
	return []byte(key)
}

// Set sets the value of a key
// If the key does not exist, it will be created
// If the key exists, it will be overwritten
// If the key is expired, it will be deleted
// If the key is not expired, it will be updated
// func (s *StringStructure) Set(key, value []byte, ttl time.Duration) error {
func (s *StringStructure) Set(k string, v interface{}, ttl time.Duration) error {
	key := stringToBytesWithKey(k)
	value, err, valueType := interfaceToBytes(v)

	if err != nil {
		return err
	}

	// Set the value type
	s.valueType = valueType

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
func (s *StringStructure) Get(k string) (interface{}, error) {
	key := stringToBytesWithKey(k)

	// Get the value
	value, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	interValue, err := decodeStringValue(value)
	if err != nil {
		return nil, err
	}

	valueType := s.valueType

	valueToInterface, err := byteToInterface(interValue, valueType)
	if err != nil {
		return nil, err
	}

	return valueToInterface, nil
}

// Del deletes the value of a key
// If the key does not exist, it will return nil
// If the key exists, it will be deleted
// If the key is expired, it will be deleted and return nil
// If the key is not expired, it will be updated and return nil
func (s *StringStructure) Del(k string) error {
	key := stringToBytesWithKey(k)
	// Delete the value
	return s.db.Delete(key)
}

// Type returns the type of a key
// If the key does not exist, it will return ""
// If the key exists, it will return "string"
// If the key is expired, it will be deleted and return ""
// If the key is not expired, it will be updated and return "string"
func (s *StringStructure) Type(k string) (string, error) {
	key := stringToBytesWithKey(k)
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
func (s *StringStructure) StrLen(k string) (int, error) {
	key := stringToBytesWithKey(k)
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
func (s *StringStructure) GetSet(key string, value interface{}, ttl time.Duration) (interface{}, error) {
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
func (s *StringStructure) Append(key string, v interface{}, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	value, err, _ := interfaceToBytes(v)
	if err != nil {
		return err
	}

	// Convert the old value to a byte slice
	oldValueType := oldValue.([]byte)

	// Append the value
	newValue := append(oldValueType, value...)

	// Set the value
	return s.Set(key, newValue, ttl)
}

// Incr increments the integer value of a key by 1
func (s *StringStructure) Incr(key string, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	intValue, err := convertToInt(oldValue)
	if err != nil {
		return err
	}

	// Increment the integer value
	newIntValue := intValue + 1

	// Convert the new integer value to a byte slice
	newValue := strconv.Itoa(newIntValue)

	// Set the value
	return s.Set(key, newValue, ttl)
}

// IncrBy increments the integer value of a key by the given amount
func (s *StringStructure) IncrBy(key string, amount int, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	intValue, err := convertToInt(oldValue)
	if err != nil {
		return err
	}

	newIntValue := intValue + amount

	newValue := strconv.Itoa(newIntValue)

	return s.Set(key, newValue, ttl)
}

// IncrByFloat increments the float value of a key by the given amount
func (s *StringStructure) IncrByFloat(key string, amount float64, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to a byte slice
	floatValue, err := convertToFloat(oldValue)
	if err != nil {
		return err
	}

	// Increment the float value
	newFloatValue := floatValue + amount

	// Convert the new float value to a byte slice
	newValue := strconv.FormatFloat(newFloatValue, 'f', -1, 64)

	// Set the value
	return s.Set(key, newValue, ttl)
}

// Decr decrements the integer value of a key by 1
func (s *StringStructure) Decr(key string, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	intValue, err := convertToInt(oldValue)
	if err != nil {
		return err
	}

	// Decrement the integer value
	newIntValue := intValue - 1

	newValue := strconv.Itoa(newIntValue)

	// Set the value
	return s.Set(key, newValue, ttl)
}

// DecrBy decrements the integer value of a key by the given amount
func (s *StringStructure) DecrBy(key string, amount int, ttl time.Duration) error {
	// Get the old value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	intValue, err := convertToInt(oldValue)
	if err != nil {
		return err
	}

	// Decrement the integer value
	newIntValue := intValue - amount

	newValue := strconv.Itoa(newIntValue)

	// Set the value
	return s.Set(key, newValue, ttl)
}

// Exists checks if a key exists
func (s *StringStructure) Exists(key string) (bool, error) {
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
func (s *StringStructure) Expire(key string, ttl time.Duration) error {
	// Get the value
	oldValue, err := s.Get(key)
	if err != nil {
		return err
	}

	// Convert the old value to an integer
	intValue, err := convertToInt(oldValue)
	if err != nil {
		return err
	}

	newValue := strconv.Itoa(intValue)

	// Set the value
	return s.Set(key, newValue, ttl)
}

// Persist removes the expiration time of a key
func (s *StringStructure) Persist(key string) error {
	// Get the value
	value, err := s.Get(key)
	if err != nil {
		return err
	}

	// Set the value
	return s.Set(key, value, 0)
}

func (s *StringStructure) MGet(keys ...string) ([]interface{}, error) {
	// Create a slice to store the values
	values := make([]interface{}, len(keys))

	// Get the value for each key
	for i, key := range keys {
		value, err := s.Get(key)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return values, nil
}

func (s *StringStructure) MSet(pairs ...interface{}) error {
	if len(pairs)%2 != 0 {
		return errors.New("Wrong number of arguments")
	}

	// Create a map to store the key-value pairs
	data := make(map[string]interface{})

	// Extract key-value pairs from the input arguments and store them in the map
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return errors.New("Invalid key")
		}

		value := pairs[i+1]
		data[key] = value
	}

	// Set each key-value pair in the map
	for key, value := range data {
		if err := s.Set(key, value, 0); err != nil {
			return err
		}
	}

	return nil
}

// MSetNX sets multiple key-value pairs only if none of the specified keys exist
func (s *StringStructure) MSetNX(pairs ...interface{}) (bool, error) {
	if len(pairs)%2 != 0 {
		return false, errors.New("Wrong number of arguments")
	}

	// Check if any of the specified keys already exist
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return false, errors.New("Invalid key")
		}

		exists, err := s.Exists(key)
		if err != nil {
			return false, err
		}

		if exists {
			return false, nil
		}
	}

	// Create a map to store the key-value pairs
	data := make(map[string]interface{})

	// Extract key-value pairs from the input arguments and store them in the map
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return false, errors.New("Invalid key")
		}

		value := pairs[i+1]
		data[key] = value
	}

	// Set each key-value pair in the map
	for key, value := range data {
		if err := s.Set(key, value, 0); err != nil {
			return false, err
		}
	}

	return true, nil
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

func (s *StringStructure) Stop() error {
	err := s.db.Close()
	return err
}

func (s *StringStructure) Clean() {
	s.db.Clean()
}
