package structure

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
)

// StringStructureInternal is the internal version of StringStructure using shared DB instance
type StringStructureInternal struct {
	db        *engine.DB
	valueType string
}

// newStringStructureInternal creates a new StringStructureInternal
func newStringStructureInternal(db *engine.DB) *StringStructureInternal {
	return &StringStructureInternal{
		db: db,
	}
}

// Set sets a key-value pair
func (s *StringStructureInternal) Set(k string, v interface{}, ttl int64) error {
	key := stringToBytesWithKey(k)
	value, err, valueType := interfaceToBytes(v)

	if err != nil {
		return err
	}

	// Set value type
	s.valueType = valueType

	if value == nil {
		return nil
	}

	// If key exists, check current TTL
	var currentTTL int64 = 0
	existingValue, err := s.db.Get(key)
	if err == nil {
		// Key exists, get current TTL
		_, expire, _ := decodeStringValue(existingValue)
		if expire > 0 {
			now := time.Now().UnixNano()
			if expire > now {
				currentTTL = (expire - now) / int64(time.Second)
			}
		}
	}

	// If TTL is -1 and key exists, keep existing TTL
	if ttl == -1 && existingValue != nil {
		ttl = currentTTL
	}

	// Encode value
	encValue, err := encodeStringValue(value, integerToDuration(ttl))
	if err != nil {
		return err
	}

	// Set value
	return s.db.Put(key, encValue)
}

// Get gets the value associated with the key
func (s *StringStructureInternal) Get(k string) (interface{}, error) {
	key := stringToBytesWithKey(k)

	// Get value
	value, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	interValue, _, err := decodeStringValue(value)
	if err != nil {
		return nil, err
	}

	valueType := s.valueType
	if valueType == "" {
		// When type is not specified, return value as string
		return string(interValue), nil
	}

	valueToInterface, err := byteToInterface(interValue, valueType)
	if err != nil {
		return nil, err
	}

	return valueToInterface, nil
}

// Del deletes a key
func (s *StringStructureInternal) Del(k string) error {
	key := stringToBytesWithKey(k)
	return s.db.Delete(key)
}

// Exists checks if the key exists
func (s *StringStructureInternal) Exists(k string) (bool, error) {
	key := stringToBytesWithKey(k)
	value, err := s.db.Get(key)
	if err != nil {
		if err == _const.ErrKeyNotFound {
			return false, nil
		}
		return false, err
	}

	// Decode value to check expiration time
	_, expire, err := decodeStringValue(value)
	if err != nil {
		return false, err
	}

	// Check if expired
	if expire > 0 && expire <= time.Now().UnixNano() {
		// Asynchronously delete expired key
		go func() {
			_ = s.db.Delete(key)
		}()
		return false, nil
	}

	return true, nil
}

// Type returns the type of the key
func (s *StringStructureInternal) Type(k string) (string, error) {
	key := stringToBytesWithKey(k)
	// Get value
	value, err := s.db.Get(key)
	if err != nil {
		return "", err
	}

	// Decode value
	_, _, err = decodeStringValue(value)
	if err != nil {
		return "", err
	}

	// Return type
	return "string", nil
}

// Expire sets the expiration time for a key
func (s *StringStructureInternal) Expire(k string, ttl int64) error {
	key := stringToBytesWithKey(k)

	// Get current value
	value, err := s.db.Get(key)
	if err != nil {
		return err
	}

	// Decode value
	rawValue, _, err := decodeStringValue(value)
	if err != nil {
		return err
	}

	// Re-encode with new TTL
	encValue, err := encodeStringValue(rawValue, integerToDuration(ttl))
	if err != nil {
		return err
	}

	// Update value
	return s.db.Put(key, encValue)
}

// TTL gets the remaining time to live for a key
func (s *StringStructureInternal) TTL(k string) (int64, error) {
	key := stringToBytesWithKey(k)

	// Get value
	value, err := s.db.Get(key)
	if err != nil {
		if err == _const.ErrKeyNotFound {
			return -2, nil // Key not found
		}
		return -1, err
	}

	// Decode value
	_, expire, err := decodeStringValue(value)
	if err != nil {
		return -1, err
	}

	// If no expiration time set
	if expire == 0 {
		return -1, nil
	}

	// Calculate remaining time
	now := time.Now().UnixNano()
	remaining := (expire - now) / int64(time.Second)

	if remaining <= 0 {
		// Key expired, delete it
		_ = s.db.Delete(key)
		return -2, nil
	}

	return remaining, nil
}

// Keys gets all keys matching the given pattern
func (s *StringStructureInternal) Keys(pattern string) ([]string, error) {
	// Compile regex pattern
	regex, err := regexp.Compile(convertToRegexp(pattern))
	if err != nil {
		return nil, err
	}

	// Create result slice
	result := make([]string, 0)

	// Create iterator to scan all keys
	iter := s.db.NewIterator(config.DefaultIteratorOptions)
	defer iter.Close()

	// Iterate through all keys
	for iter.Rewind(); iter.Valid(); iter.Next() {
		key := iter.Key()
		keyStr := string(key)

		// Check if key matches pattern
		if regex.MatchString(keyStr) {
			// Check if key expired
			value, err := iter.Value()
			if err != nil {
				continue
			}

			_, expire, err := decodeStringValue(value)
			if err != nil {
				continue
			}

			// If key not expired, add to result
			if expire == 0 || expire > time.Now().UnixNano() {
				result = append(result, keyStr)
			}
		}
	}

	return result, nil
}

// GetSet sets the key's value and returns the old value
func (s *StringStructureInternal) GetSet(key string, value interface{}, ttl int64) (interface{}, error) {
	// Get old value
	oldValue, err := s.Get(key)
	if err != nil && err != _const.ErrKeyNotFound {
		return nil, err
	}

	// Set new value
	err = s.Set(key, value, ttl)
	if err != nil {
		return nil, err
	}

	// If key doesn't exist, return nil
	if err == _const.ErrKeyNotFound {
		return nil, nil
	}

	// Return old value
	return oldValue, nil
}

// Incr increments the key's integer value by 1
func (s *StringStructureInternal) Incr(key string, ttl int64) error {
	return s.IncrBy(key, 1, ttl)
}

// IncrBy increments the key's integer value by the specified amount
func (s *StringStructureInternal) IncrBy(key string, amount int, ttl int64) error {
	// Get current value
	currentValue, err := s.Get(key)
	if err != nil && err != _const.ErrKeyNotFound {
		return err
	}

	var newVal int

	// If key doesn't exist or is empty
	if err == _const.ErrKeyNotFound || currentValue == nil {
		newVal = amount
	} else {
		// Convert to integer
		currentInt, err := convertToInt(currentValue)
		if err != nil {
			return err
		}
		newVal = currentInt + amount
	}

	// Set new value
	return s.Set(key, newVal, ttl)
}

// Decr decrements the key's integer value by 1
func (s *StringStructureInternal) Decr(key string, ttl int64) error {
	return s.DecrBy(key, 1, ttl)
}

// DecrBy decrements the key's integer value by the specified amount
func (s *StringStructureInternal) DecrBy(key string, amount int, ttl int64) error {
	return s.IncrBy(key, -amount, ttl)
}

// StrLen returns the length of the key's value
func (s *StringStructureInternal) StrLen(key string) (int, error) {
	// Get value
	value, err := s.Get(key)
	if err != nil {
		if err == _const.ErrKeyNotFound {
			return 0, nil
		}
		return 0, err
	}

	if value == nil {
		return 0, nil
	}

	// Convert to string
	strValue, err := interfaceToString(value)
	if err != nil {
		return 0, err
	}

	return len(strValue), nil
}

// Append appends a string to the key's value
func (s *StringStructureInternal) Append(key string, v interface{}, ttl int64) error {
	// Get old value
	oldValue, err := s.Get(key)
	if err != nil && err != _const.ErrKeyNotFound {
		return err
	}

	// Handle value conversion
	value, err, valueType := interfaceToBytes(v)
	if err != nil {
		return err
	}

	// Set value type
	s.valueType = valueType

	var newValue []byte

	// If key doesn't exist
	if err == _const.ErrKeyNotFound || oldValue == nil {
		newValue = value
	} else {
		// Convert old value to byte array
		oldBytes, err, _ := interfaceToBytes(oldValue)
		if err != nil {
			return err
		}

		// Concatenate
		newValue = append(oldBytes, value...)
	}

	// Set new value
	return s.Set(key, string(newValue), ttl)
}

// SetNX sets the key's value only if the key does not exist
func (s *StringStructureInternal) SetNX(key string, value interface{}, ttl int64) (bool, error) {
	// Check if key exists
	exists, err := s.Exists(key)
	if err != nil {
		return false, err
	}

	// If key exists, return false directly
	if exists {
		return false, nil
	}

	// Set value
	err = s.Set(key, value, ttl)
	if err != nil {
		return false, err
	}

	return true, nil
}

// MSetNX sets multiple key-value pairs only if none of the keys exist
func (s *StringStructureInternal) MSetNX(pairs ...interface{}) (bool, error) {
	if len(pairs)%2 != 0 {
		return false, fmt.Errorf("wrong number of arguments for MSetNX")
	}

	// Check all keys don't exist
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return false, fmt.Errorf("key must be string")
		}

		exists, err := s.Exists(key)
		if err != nil {
			return false, err
		}

		if exists {
			return false, nil
		}
	}

	// Set all key-value pairs
	batch := s.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i].(string)
		value := pairs[i+1]

		k := stringToBytesWithKey(key)
		v, err, valueType := interfaceToBytes(value)
		if err != nil {
			return false, err
		}

		// Set value type
		s.valueType = valueType

		if v == nil {
			continue
		}

		// Encode value
		encValue, err := encodeStringValue(v, 0) // No expiration time
		if err != nil {
			return false, err
		}

		// Add key-value pair to batch
		err = batch.Put(k, encValue)
		if err != nil {
			return false, err
		}
	}

	// Commit batch
	err := batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// Scan scans keys matching the given pattern starting from cursor, returns up to specified count of keys
func (s *StringStructureInternal) Scan(cursor int, pattern string, count int) (int, []string, error) {
	if count <= 0 {
		count = 10 // 默认值
	}

	// Compile regex pattern
	var regex *regexp.Regexp
	if pattern != "" {
		var err error
		regex, err = regexp.Compile(convertToRegexp(pattern))
		if err != nil {
			return 0, nil, err
		}
	}

	// Create result slice
	result := make([]string, 0, count)

	// Create iterator to scan all keys
	iter := s.db.NewIterator(config.DefaultIteratorOptions)
	defer iter.Close()

	// Position at cursor
	if cursor > 0 {
		// Simplified processing: skip first cursor keys
		iter.Rewind()
		for i := 0; i < cursor && iter.Valid(); i++ {
			iter.Next()
		}
	} else {
		iter.Rewind()
	}

	// Collect matching keys
	newCursor := cursor
	for iter.Valid() && len(result) < count {
		newCursor++
		key := iter.Key()
		keyStr := string(key)

		// Check if key matches pattern
		if regex == nil || regex.MatchString(keyStr) {
			// Check if key expired
			value, err := iter.Value()
			if err != nil {
				iter.Next()
				continue
			}

			_, expire, err := decodeStringValue(value)
			if err != nil {
				iter.Next()
				continue
			}

			// If key not expired, add to result
			if expire == 0 || expire > time.Now().UnixNano() {
				result = append(result, keyStr)
			}
		}

		iter.Next()
	}

	// If reached end, return 0 as new cursor
	if !iter.Valid() {
		newCursor = 0
	}

	return newCursor, result, nil
}

// GetRange gets a substring of the string
func (s *StringStructureInternal) GetRange(key string, start, end int) (string, error) {
	// Get value
	value, err := s.Get(key)
	if err != nil {
		return "", err
	}

	if value == nil {
		return "", nil
	}

	// Convert to string
	strValue, err := interfaceToString(value)
	if err != nil {
		return "", err
	}

	strLen := len(strValue)

	// Handle negative indices
	if start < 0 {
		start = strLen + start
		if start < 0 {
			start = 0
		}
	}

	if end < 0 {
		end = strLen + end
		if end < 0 {
			end = 0
		}
	}

	// Ensure indices are within valid range
	if start >= strLen || start > end {
		return "", nil
	}

	if end >= strLen {
		end = strLen - 1
	}

	// Return substring
	return strValue[start : end+1], nil
}

// SetRange overwrites part of the string
func (s *StringStructureInternal) SetRange(key string, offset int, value string) (int, error) {
	// Get current value
	currentVal, err := s.Get(key)
	if err != nil && err != _const.ErrKeyNotFound {
		return 0, err
	}

	var currentStr string
	if err == _const.ErrKeyNotFound || currentVal == nil {
		currentStr = ""
	} else {
		currentStr, err = interfaceToString(currentVal)
		if err != nil {
			return 0, err
		}
	}

	// If offset exceeds string length, pad with zero bytes
	if offset > len(currentStr) {
		padding := strings.Repeat("\u0000", offset-len(currentStr))
		currentStr = currentStr + padding
	}

	// Build new string
	var newStr string
	if offset == 0 {
		newStr = value + currentStr[len(value):]
	} else if offset+len(value) > len(currentStr) {
		newStr = currentStr[:offset] + value
	} else {
		newStr = currentStr[:offset] + value + currentStr[offset+len(value):]
	}

	// Set new value
	err = s.Set(key, newStr, 0) // Don't change expiration time
	if err != nil {
		return 0, err
	}

	return len(newStr), nil
}

// BitCount counts the number of bits set to 1 in the string
func (s *StringStructureInternal) BitCount(key string, start, end int) (int, error) {
	// Get value
	value, err := s.Get(key)
	if err != nil {
		return 0, err
	}

	if value == nil {
		return 0, nil
	}

	// Convert to byte array
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		strVal, err := interfaceToString(value)
		if err != nil {
			return 0, err
		}
		bytes = []byte(strVal)
	}

	byteLen := len(bytes)

	// Handle negative indices
	if start < 0 {
		start = byteLen + start
		if start < 0 {
			start = 0
		}
	}

	if end < 0 {
		end = byteLen + end
		if end < 0 {
			end = 0
		}
	}

	// Ensure indices are within valid range
	if start >= byteLen || start > end {
		return 0, nil
	}

	if end >= byteLen {
		end = byteLen - 1
	}

	// Calculate the number of bits set to 1
	count := 0
	for i := start; i <= end; i++ {
		b := bytes[i]
		// Calculate the number of bits set to 1 in the byte
		for j := 0; j < 8; j++ {
			if (b & (1 << j)) != 0 {
				count++
			}
		}
	}

	return count, nil
}

// MGet gets values for multiple keys
func (s *StringStructureInternal) MGet(keys ...string) ([]interface{}, error) {
	result := make([]interface{}, len(keys))

	for i, key := range keys {
		value, err := s.Get(key)
		if err != nil {
			if err == _const.ErrKeyNotFound {
				result[i] = nil
				continue
			}

			// Type conversion error should not interrupt the entire operation
			if strings.Contains(err.Error(), "cannot convert to") {
				result[i] = nil
				continue
			}

			return nil, err
		}
		result[i] = value
	}

	return result, nil
}

// GetBit gets the bit value at specified offset in the string
func (s *StringStructureInternal) GetBit(key string, offset int) (int, error) {
	// Get value
	value, err := s.Get(key)
	if err != nil {
		if err == _const.ErrKeyNotFound {
			return 0, nil
		}
		return 0, err
	}

	if value == nil {
		return 0, nil
	}

	// Convert to byte array
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		strVal, err := interfaceToString(value)
		if err != nil {
			return 0, err
		}
		bytes = []byte(strVal)
	}

	// Calculate byte index and bit index
	byteIndex := offset / 8
	bitIndex := offset % 8

	// If offset exceeds string length, return 0
	if byteIndex >= len(bytes) {
		return 0, nil
	}

	// Get bit value
	bit := (bytes[byteIndex] >> bitIndex) & 1

	return int(bit), nil
}

// SetBit sets the bit value at specified offset in the string
func (s *StringStructureInternal) SetBit(key string, offset int, value int) (int, error) {
	if value != 0 && value != 1 {
		return 0, fmt.Errorf("bit value must be 0 or 1")
	}

	// Get current value
	currentVal, err := s.Get(key)
	if err != nil && err != _const.ErrKeyNotFound {
		return 0, err
	}

	var bytes []byte
	if err == _const.ErrKeyNotFound || currentVal == nil {
		bytes = []byte{}
	} else {
		switch v := currentVal.(type) {
		case string:
			bytes = []byte(v)
		case []byte:
			bytes = v
		default:
			strVal, err := interfaceToString(currentVal)
			if err != nil {
				return 0, err
			}
			bytes = []byte(strVal)
		}
	}

	// Calculate byte index and bit index
	byteIndex := offset / 8
	bitIndex := offset % 8

	// If offset exceeds string length, expand byte array
	if byteIndex >= len(bytes) {
		newBytes := make([]byte, byteIndex+1)
		copy(newBytes, bytes)
		bytes = newBytes
	}

	// Get old bit value
	oldBit := (bytes[byteIndex] >> bitIndex) & 1

	// Set new bit value
	if value == 1 {
		bytes[byteIndex] |= (1 << bitIndex)
	} else {
		bytes[byteIndex] &= ^(1 << bitIndex)
	}

	// Save modified value
	err = s.Set(key, bytes, 0) // Don't change expiration time
	if err != nil {
		return 0, err
	}

	return int(oldBit), nil
}
