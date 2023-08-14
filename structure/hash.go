package structure

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"regexp"
	"strconv"
	"time"
)

type HashMetadata struct {
	dataType        byte  // Represents the data type of the hash object.
	dataSize        int64 // Represents the size of the hash object.
	expire          int64 // Represents the expiration time of the hash object.
	version         int64 // Represents the version number of the hash object.
	counter         int64 // Represents the counter value of the hash object.
	createdTime     int64 // Represents the creation time of the hash object.
	lastUpdatedTime int64 // Represents the last updated time of the hash object.
}

const maxHashMetaSize = 1 + binary.MaxVarintLen64*6

type HashStructure struct {
	db            *engine.DB
	hashValueType string
	HashFieldType string
}

// NewHashStructure Returns a new NewHashStructure
func NewHashStructure(options config.Options) (*HashStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &HashStructure{db: db}, nil
}

// HSet sets the string value of a hash field in the HashStructure.
// It takes the key, field, and value as input and stores the value in the specified hash field.
//
// Parameters:
//
//	k: The key under which the hash is stored.
//	f: The field within the hash where the value will be set.
//	v: The value to be set in the hash field.
//
// Returns:
//
//	bool: A boolean indicating whether the field was newly created (true) or updated (false).
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
// - The function converts the parameters to bytes and ensures they are not empty.
// - It retrieves the existing hash metadata from the database using the given key.
// - If the field does not exist, the hash metadata counter is incremented.
// - The function creates a new HashField containing the field details and encodes it.
// - The function uses a write batch to efficiently commit changes to the database.
// - It returns a boolean indicating whether the field was newly created or updated.
func (hs *HashStructure) HSet(key string, field, value interface{}) (bool, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return false, err
	}

	// Convert the parameters to bytes
	v, err, valueType := interfaceToBytes(value)

	if err != nil {
		return false, err
	}

	// Set the hash value type
	hs.hashValueType = valueType

	// Set the hash field type
	_, _, fieldType := interfaceToBytes(v)
	hs.HashFieldType = fieldType

	// Check the parameters
	if len(k) == 0 || len(f) == 0 || len(v) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// Create a new HashField
	hf := &HashField{
		key:     k,
		field:   f,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	var exist = true

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		exist = false
	}

	// new a write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// If the field is not found, increase the counter
	if !exist {
		hashMeta.counter++
		_ = batch.Put(k, hashMeta.encodeHashMeta())
	}

	// Put the field to the database
	_ = batch.Put(hfBuf, v)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return !exist, nil
}

// HGet gets the string value of a hash field.
// It takes a string key 'k' and a field 'f', and returns the corresponding value and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field name.
//
// Returns:
//
//	interface{}: The value corresponding to the field, or nil if the field doesn't exist.
//	error: An error if occurred during the operation, or nil on success.
//
// Notes:
// - Parameters 'k' and 'f' need to be non-empty.
// - If the counter in the hash table is 0, nil is returned.
// - The function looks up hash metadata based on the provided key 'k'.
// - Creates a new HashField structure for manipulation.
// - Obtains the byte representation of the hash field by encoding the HashField structure.
// - Retrieves the byte data from the database using the HashStructure instance's database object.
// - The retrieved byte data is converted back to the corresponding data type.
// - Returns the value corresponding to the field and any possible error.
func (hs *HashStructure) HGet(key string, field interface{}) (interface{}, error) {
	// check the key ttl
	ttl, _ := hs.TTL(key)
	if ttl == -1 {
		return nil, _const.ErrKeyIsExpired
	}

	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return nil, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return nil, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return nil, err
	}

	// If the counter is 0, return nil
	if hashMeta.counter == 0 {
		return nil, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	value, err := hs.db.Get(hfBuf)
	if err != nil {
		return nil, err
	}

	// Get the value type from the hashValueTypes
	valueType := hs.hashValueType

	// Values of different types need to be converted to corresponding types
	valueToInterface, err := byteToInterface(value, valueType)
	if err != nil {
		return nil, err
	}
	return valueToInterface, nil
}

// HMGet gets the string value of multiple hash fields.
// It takes a string key 'k' and a variadic number of fields 'f'. It returns an array of
// interface{} containing the values corresponding to the provided fields and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: Variable number of fields to retrieve values for.
//
// Returns:
//
//	[]interface{}: An array of interface{} containing values corresponding to the fields.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HMGet(key string, field ...interface{}) ([]interface{}, error) {
	// check the key ttl
	ttl, _ := hs.TTL(key)
	if ttl == -1 {
		return nil, _const.ErrKeyIsExpired
	}

	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)
	var interfaces []interface{}

	for _, fi := range field {
		// Convert the parameters to bytes
		f, err, _ := interfaceToBytes(fi)
		if err != nil {
			return nil, err
		}

		// Check the parameters
		if len(k) == 0 || len(f) == 0 {
			return nil, _const.ErrKeyIsEmpty
		}

		// Find the hash metadata by the given key
		hashMeta, err := hs.findHashMeta(key, Hash)
		if err != nil {
			return nil, err
		}

		// If the counter is 0, return nil
		if hashMeta.counter == 0 {
			return nil, nil
		}

		// Create a new HashField
		hf := &HashField{
			field:   f,
			key:     k,
			version: hashMeta.version,
		}

		// Encode the HashField
		hfBuf := hf.encodeHashField()

		// Get the field from the database
		value, err := hs.db.Get(hfBuf)
		if err != nil {
			return nil, err
		}

		// Get the value type from the hashValueTypes
		valueType := hs.hashValueType

		// Values of different types need to be converted to corresponding types
		valueToInterface, err := byteToInterface(value, valueType)
		if err != nil {
			return nil, err
		}
		interfaces = append(interfaces, valueToInterface)
	}

	return interfaces, nil
}

// HDel deletes one field from a hash.
// It takes a string key 'k' and a field 'f' to be deleted from the hash.
// It returns a boolean indicating the success of the operation and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field to be deleted.
//
// Returns:
//
//	bool: True if the field was deleted successfully, false otherwise.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HDel(key string, field interface{}) (bool, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return false, nil
	}

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Delete the field from the database
	_ = batch.Delete(hfBuf)

	// Decrease the counter
	hashMeta.counter--

	// Put the updated hash metadata to the database
	_ = batch.Put(k, hashMeta.encodeHashMeta())

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HDelAll deletes all fields from a hash.
// It takes a string key 'k' to be deleted from the hash.
// It returns a boolean indicating the success of the operation and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//
// Returns:
//
//	bool: True if the hash was deleted successfully, false otherwise.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HDelAll(key string) (bool, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Check the parameters
	if len(k) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Delete the field from the database
	_ = batch.Delete(k)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HExpire sets a timeout on a hash.
// It takes a string key 'k' and a duration 'ttl' to set the timeout.
// It returns a boolean indicating the success of the operation and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	ttl: The duration of the timeout.
//
// Returns:
//
//	bool: True if the timeout was set successfully, false otherwise.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HExpire(key string, ttl int64) (bool, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Check the parameters
	if len(k) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Put the updated hash metadata to the database
	hashMeta.expire = time.Now().Add(time.Duration(ttl) * time.Second).UnixNano()
	_ = batch.Put(k, hashMeta.encodeHashMeta())

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HExists determines whether a hash field exists or not.
// It takes a string key 'k' and a field 'f' to check for existence.
// It returns a boolean indicating whether the field exists and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field to check for existence.
//
// Returns:
//
//	bool: True if the field exists, false otherwise.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HExists(key string, field interface{}) (bool, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return false, nil
	}

	return true, nil
}

// HLen gets the number of fields contained in a hash.
// It takes a string key 'k' and returns the number of fields in the hash.
//
// Parameters:
//
//	k: The key of the hash table.
//
// Returns:
//
//	int: The number of fields in the hash.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HLen(key string) (int, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Check the parameters
	if len(k) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	return int(hashMeta.counter), nil
}

// HUpdate updates the string value of a hash field.
// It takes a string key 'k', a field 'f', and a value 'v' to update the field's value.
// It returns a boolean indicating the success of the update and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field to be updated.
//	v: The new value to set for the field.
//
// Returns:
//
//	bool: True if the update was successful, false otherwise.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HUpdate(key string, field, value interface{}) (bool, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return false, err
	}

	// Convert the parameters to bytes
	v, err, _ := interfaceToBytes(value)
	if err != nil {
		return false, err
	}
	// Check the parameters
	if len(k) == 0 || len(f) == 0 || len(v) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return false, nil
	}

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Put the field to the database
	_ = batch.Put(hfBuf, v)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HIncrBy increments the integer value of a hash field by the given number.
// It takes a string key 'k', a field 'f', and an increment value 'increment'.
// It returns the updated value after increment and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field whose value needs to be incremented.
//	increment: The value to increment the field by.
//
// Returns:
//
//	int64: The updated value of the field after increment.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HIncrBy(key string, field interface{}, increment int64) (int64, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	value, err := hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return 0, nil
	}

	// Convert the value to int64
	val, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		return 0, err
	}

	// Add the increment to the value
	val += increment

	// Convert the value to string
	value = []byte(strconv.FormatInt(val, 10))

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Put the field to the database
	_ = batch.Put(hfBuf, value)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return 0, err
	}

	return val, nil
}

// HIncrByFloat increments the float value of a hash field by the given number.
// It takes a string key 'k', a field 'f', and an increment value 'increment'.
// It returns the updated value after increment and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field whose value needs to be incremented.
//	increment: The value to increment the field by.
//
// Returns:
//
//	float64: The updated value of the field after increment.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HIncrByFloat(key string, field interface{}, increment float64) (float64, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	value, err := hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return 0, nil
	}

	// Convert the value to float64
	val, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return 0, err
	}

	// Add the increment to the value
	val += increment

	// Convert the value to string
	value = []byte(strconv.FormatFloat(val, 'f', -1, 64))

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Put the field to the database
	_ = batch.Put(hfBuf, value)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return 0, err
	}

	return val, nil
}

// HDecrBy decrements the integer value of a hash field by the given number.
// It takes a string key 'k', a field 'f', and a decrement value 'decrement'.
// It returns the updated value after decrement and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field whose value needs to be decremented.
//	decrement: The value to decrement the field by.
//
// Returns:
//
//	int64: The updated value of the field after decrement.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HDecrBy(key string, field interface{}, decrement int64) (int64, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	value, err := hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return 0, nil
	}

	// Convert the value to int64
	val, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		return 0, err
	}

	// Subtract the decrement from the value
	val -= decrement

	// Convert the value to string
	value = []byte(strconv.FormatInt(val, 10))

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Put the field to the database
	_ = batch.Put(hfBuf, value)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return 0, err
	}

	return val, nil
}

// HStrLen returns the string length of the value associated with a field in the hash.
// It takes a string key 'k' and a field 'f' and returns the length of the field's value.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field whose value length needs to be determined.
//
// Returns:
//
//	int: The length of the field's value.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HStrLen(key string, field interface{}) (int, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	value, err := hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return 0, nil
	}

	return len(value), nil
}

// HMove moves a field from a source hash to a destination hash.
// It takes the source key 'source', the destination key 'destination', and the field 'f' to be moved.
// It returns a boolean indicating the success of the move and any possible error.
//
// Parameters:
//
//	source: The source hash key.
//	destination: The destination hash key.
//	f: The field to be moved.
//
// Returns:
//
//	bool: True if the move was successful, false otherwise.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HMove(source, destination string, field interface{}) (bool, error) {
	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(source) == 0 || len(destination) == 0 || len(f) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given source
	sourceMeta, err := hs.findHashMeta(source, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if sourceMeta.counter == 0 {
		return false, nil
	}

	// Find the hash metadata by the given destination
	destinationMeta, err := hs.findHashMeta(destination, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if destinationMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField for the source
	hf := &HashField{
		field:   f,
		key:     []byte(source),
		version: sourceMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Create a new HashField for the destination
	destinationHf := &HashField{
		field:   f,
		key:     []byte(destination),
		version: destinationMeta.version,
	}

	// Encode the HashField
	destinationHfBuf := destinationHf.encodeHashField()

	// Get the field from the database
	value, err := hs.db.Get(destinationHfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return false, nil
	}

	// Create a new write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Delete the field from the source
	_ = batch.Delete(hfBuf)

	// Put the field to the destination
	_ = batch.Put(hfBuf, value)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HSetNX sets a field in the hash only if the field does not already exist.
// It takes a string key 'k', a field 'f', and a value 'v' to set if the field doesn't exist.
// It returns a boolean indicating whether the field was set and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field to be set if it doesn't exist.
//	v: The value to set for the field.
//
// Returns:
//
//	bool: True if the field was set, false otherwise.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HSetNX(key string, field, value interface{}) (bool, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return false, err
	}

	// Convert the parameters to bytes
	v, err, _ := interfaceToBytes(value)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 || len(v) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		_, err := hs.HSet(key, field, value)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, nil
	}
}

// HTypes returns the type of a field in the hash.
// It takes a string key 'k' and a field 'f'.
// It returns a string indicating the type of the field and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field whose type needs to be determined.
//
// Returns:
//
//	string: The type of the field. Possible values: "hash" (if the field exists), or an empty string.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HTypes(key string, field interface{}) (string, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Convert the parameters to bytes
	f, err, _ := interfaceToBytes(field)
	if err != nil {
		return "", err
	}

	// Check the parameters
	if len(k) == 0 || len(f) == 0 {
		return "", _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return "", err
	}

	// If the counter is 0, return empty string
	if hashMeta.counter == 0 {
		return "", nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   f,
		key:     k,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return "", _const.ErrKeyNotFound
	} else {
		return "hash", nil
	}
}

// Keys returns a list of all field names in the hash stored at the specified key.
// It takes no parameters and returns a slice of strings representing the field names and any possible error.
//
// Returns:
//
// []string: A list of field names in the hash.
// error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) Keys(regx string) ([]string, error) {
	toRegexp := convertToRegexp(regx)
	compile, err := regexp.Compile(toRegexp)
	if err != nil {
		return nil, err
	}
	// Get all the keys from the database
	byteKeys := hs.db.GetListKeys()

	// Create a new slice of strings
	keys := make([]string, 0)

	for _, key := range byteKeys {
		if compile.MatchString(string(key)) {
			// Check if the key has the identifier suffix
			if !keysIdentify(key) {
				fields := hs.GetFields(string(key))
				ok := true
				for _, field := range fields {
					_, err := hs.HGet(string(key), field)
					if err != nil {
						ok = false
						break
					}
				}
				if !ok {
					continue
				}
				keys = append(keys, string(key))
			}
		}

	}

	return keys, nil
}

// GetFields returns a list of all field names in the hash stored at the specified key.
// It takes a string key 'k' and returns a slice of strings representing
// the field names and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//
// Returns:
//
// []string: A list of field names in the hash.
func (hs *HashStructure) GetFields(key string) []string {
	// Get all the keys from the database
	byte_keys := hs.db.GetListKeys()

	// Create a new slice of strings
	fields := make([]string, 0)

	for _, k := range byte_keys {
		// Check if the key has the identifier suffix
		if keysIdentify(k) {
			hf, _ := decodeHashField(k)
			if string(hf.key) == key {
				fields = append(fields, string(hf.field))
			}
		}
	}

	return fields
}

// HGetAllFieldAndValue returns all fields and values of the hash stored at the specified key.
// It takes a string key 'k' and returns a map of strings representing the field names and
// their values and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//
// Returns:
//
// map[string]interface{}: A map of field names and their values.
// error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) HGetAllFieldAndValue(key string) (map[string]interface{}, error) {
	fields := hs.GetFields(key)

	filedAndValue := make(map[string]interface{}, 0)

	for _, field := range fields {
		v, err := hs.HGet(key, field)
		if err != nil {
			return nil, err
		}
		filedAndValue[field] = v
	}

	return filedAndValue, nil
}

// keysIdentify checks if the key has the identifier suffix
//
// Parameters:
//
//	key: The key to check.
//
// Returns:
//
//	bool: True if the key has the identifier suffix, false otherwise.
func keysIdentify(key []byte) bool {
	return bytes.HasSuffix(key, []byte("notk"))
}

// TTL returns the time-to-live (TTL) of a key in the hash.
// It takes a string key 'k' and returns the remaining TTL in seconds and any possible error.
//
// Parameters:
//
//	k: The key for which TTL needs to be determined.
//
// Returns:
//
//	int64: The remaining TTL in seconds. Returns 0 if the key has expired or doesn't exist.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) TTL(key string) (int64, error) {
	// Check the parameters
	if len(key) == 0 {
		return -1, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return -1, err
	}

	// If the counter is 0, return empty string
	if hashMeta.counter == 0 {
		return 0, _const.ErrKeyNotFound
	}

	// Calculate the TTL
	ttl := hashMeta.expire/int64(time.Second) - time.Now().UnixNano()/int64(time.Second)

	// If the TTL is 0, return 0
	if hashMeta.expire == 0 {
		return 0, nil
	}

	// If the TTL is less than 0, return -1
	if ttl <= 0 {
		return -1, _const.ErrKeyIsExpired
	}
	return ttl, nil
}

// Size returns the size of a field in the hash as a formatted string.
// It takes a string key 'k' and one or more fields 'f' (optional).
// It returns a formatted string indicating the size of the field and any possible error.
//
// Parameters:
//
//	k: The key of the hash table.
//	f: The field(s) whose size needs to be determined (optional).
//
// Returns:
//
//	string: A formatted string indicating the size of the field.
//	error: An error if occurred during the operation, or nil on success.
func (hs *HashStructure) Size(key string) (string, error) {
	// Get all fields and values
	keysAndValues, err := hs.HGetAllFieldAndValue(key)
	if err != nil {
		return "", err
	}

	var sizeInBytes int

	// Calculate the size
	for f, v := range keysAndValues {
		sizeInBytes += len(f) + len(v.(string))
	}

	// Convert bytes to corresponding units (KB, MB...)
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	var size string
	switch {
	case sizeInBytes < KB:
		size = fmt.Sprintf("%dB", sizeInBytes)
	case sizeInBytes < MB:
		size = fmt.Sprintf("%.2fKB", float64(sizeInBytes)/KB)
	case sizeInBytes < GB:
		size = fmt.Sprintf("%.2fMB", float64(sizeInBytes)/MB)
	}

	return size, nil
}

// findHashMeta finds the hash metadata by the given key.
func (hs *HashStructure) findHashMeta(key string, dataType DataStructure) (*HashMetadata, error) {
	// Convert the parameters to bytes
	k := stringToBytesWithKey(key)

	// Find the hash metadata by the given key
	meta, err := hs.db.Get(k)
	if err != nil && err != _const.ErrKeyNotFound {
		return nil, err
	}

	var hashMeta *HashMetadata
	var exist = true
	// If the hash metadata is not found, create a new one
	if err == _const.ErrKeyNotFound {
		exist = false
	} else {
		// Decode the hash metadata
		hashMeta = decodeHashMeta(meta)

		// Check the data type
		if hashMeta.dataType != dataType {
			return nil, ErrInvalidType
		}

	}

	// If the hash metadata is not found, create a new one
	if !exist {
		hashMeta = &HashMetadata{
			dataType:        dataType,
			dataSize:        0,
			expire:          0,
			version:         time.Now().UnixNano(),
			counter:         0,
			createdTime:     time.Now().UnixNano(),
			lastUpdatedTime: time.Now().UnixNano(),
		}
	}

	return hashMeta, nil
}

func (hs *HashStructure) Stop() error {
	err := hs.db.Close()
	return err
}

func (hs *HashStructure) Clean() {
	hs.db.Clean()
}

type HashField struct {
	key     []byte
	field   []byte
	version int64
}

// encodeHashField encodes a HashField and returns the byte array and length.
// +-------------+------------+------------+------------+
// |    key      |  field     |  version   | identifier |
// +-------------+------------+------------+------------+
// |  variable   |  variable  |  8 bytes   |   4 byte   |
// +-------------+------------+------------+------------+
func (hf *HashField) encodeHashField() []byte {
	keyLen := int32(len(hf.key))
	fieldLen := int32(len(hf.field))

	// Prepare a buffer to hold the encoded data
	buf := new(bytes.Buffer)

	// Write the key length, key, field length, field, and version to the buffer
	if err := binary.Write(buf, binary.BigEndian, keyLen); err != nil {
		return nil
	}
	if _, err := buf.Write(hf.key); err != nil {
		return nil
	}
	if err := binary.Write(buf, binary.BigEndian, fieldLen); err != nil {
		return nil
	}
	if _, err := buf.Write(hf.field); err != nil {
		return nil
	}
	if err := binary.Write(buf, binary.BigEndian, hf.version); err != nil {
		return nil
	}

	// Write the identifier
	if err := binary.Write(buf, binary.BigEndian, []byte("notk")); err != nil {
		return nil
	}

	return buf.Bytes()

}

// decodeHashField decodes the HashField from a byte buffer.
func decodeHashField(data []byte) (*HashField, error) {
	hf := &HashField{}

	// Create a buffer to hold the data
	buf := bytes.NewBuffer(data)

	var keyLen int32
	if err := binary.Read(buf, binary.BigEndian, &keyLen); err != nil {
		return nil, err
	}

	hf.key = make([]byte, keyLen)
	if _, err := buf.Read(hf.key); err != nil {
		return nil, err
	}

	var fieldLen int32
	if err := binary.Read(buf, binary.BigEndian, &fieldLen); err != nil {
		return nil, err
	}

	hf.field = make([]byte, fieldLen)
	if _, err := buf.Read(hf.field); err != nil {
		return nil, err
	}

	if err := binary.Read(buf, binary.BigEndian, &hf.version); err != nil {
		return nil, err
	}

	return &HashField{
		key:     hf.key,
		field:   hf.field,
		version: hf.version,
	}, nil

}

// EncodeHashMeta encodes a HashMetadata and returns the byte array and length.
// +-------------+------------+------------+------------+---------+----------+----------+---------+---------+
// |  data type  |  data size |   expire   |  version   | counter | created  | updated  |  field  |  value  |
// +-------------+------------+------------+------------+---------+----------+----------+---------+---------+
// |  1 byte     |  variable  |  variable  |  variable  | variable| variable | variable | variable| variable|
// +-------------+------------+------------+------------+---------+----------+----------+---------+---------+
func (meta *HashMetadata) encodeHashMeta() []byte {
	buf := make([]byte, maxHashMetaSize)

	// Store the data type at the first byte
	buf[0] = meta.dataType

	var offset = 1

	// Store the lengths of data size, expire, version, counter, createdTime and lastUpdatedTime
	offset += binary.PutVarint(buf[offset:], meta.dataSize)
	offset += binary.PutVarint(buf[offset:], meta.expire)
	offset += binary.PutVarint(buf[offset:], meta.version)
	offset += binary.PutVarint(buf[offset:], meta.counter)
	offset += binary.PutVarint(buf[offset:], meta.createdTime)
	offset += binary.PutVarint(buf[offset:], meta.lastUpdatedTime)
	return buf[:offset]
}

// DecodeHashMeta decodes the HashMetadata from a byte buffer.
func decodeHashMeta(buf []byte) *HashMetadata {
	var offset = 0
	dataType := buf[offset] // Decode data type
	offset++
	dataSize, n := binary.Varint(buf[offset:]) // Decode data size
	offset += n
	expire, n := binary.Varint(buf[offset:]) // Decode expire
	offset += n
	version, n := binary.Varint(buf[offset:]) // Decode version
	offset += n
	counter, n := binary.Varint(buf[offset:]) // Decode counter
	offset += n
	createdTime, n := binary.Varint(buf[offset:]) // Decode createdTime
	offset += n
	lastUpdatedTime, _ := binary.Varint(buf[offset:]) // Decode lastUpdatedTime
	return &HashMetadata{
		dataType:        dataType,
		dataSize:        dataSize,
		expire:          expire,
		version:         version,
		counter:         counter,
		createdTime:     createdTime,
		lastUpdatedTime: lastUpdatedTime,
	}
}
