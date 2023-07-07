package structure

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"strconv"
	"time"
)

// NewHashStructure Returns a new NewHashStructure
// If the database cannot be opened, it returns a nil NewHashStructure
// Otherwise, the database cannot be created
// If the database does not exist, it will be created

// HSet Sets the string value of a hash field.
// This operation is atomic
// If key does not exist, a new key holding a hash is created
// If field already exists in the hash, it is overwritten
// If key does not hold a hash, an error is returned

// HGet Returns the value associated with field in the hash stored at key
// If the hash table does not exist, a nil slice is returned
// If the key does not exist, a nil slice is returned
// If the field does not exist, a nil slice is returned

// HDel Removes the specified fields from the hash stored at key
// This operation is atomic
// Specified fields that do not exist within this hash are ignored
// If key does not exist, it is treated as an empty hash and this command returns error

// HExists Returns if field is an existing field in the hash stored at key
// Returns true if the hash stored at key contains field
// Returns false if the hash does not contain field, or key does not exist

// HIncrBy Increments the number stored at field in the hash stored at key by increment
// The range of values supported by HINCRBY is limited to 64 bit signed integers

// HIncrByFloat Increments the specified field of a hash stored at key, and representing a floating point number
// An error is returned if one of the following conditions occur:
// The field contains a value of the wrong type (not a string)
// The current field content or the specified increment are not parsable as a float number

// HLens Returns the number of fields contained in the hash stored at key
// If the key does not exist, 0 is returned
// If the key is not a hash, an error is returned

// HStrLen Returns the string length of the value associated with field in the hash stored at key
// If the key does not exist, 0 is returned
// If the field is not present or the key does not exist, 0 is returned

// HMove Moves field from the hash stored at source to the hash stored at destination
// This operation is atomic
// In every given moment the key will either exist or not exist
// Even if the field already exists in the destination hash, it is overwritten

// HUpdate Updates the hash stored at key by only replacing fields with new values from the given hash
// Other fields are left untouched
// If key does not exist, a new key holding a hash is created
// If key does not hold a hash, an error is returned

// HSetNX Sets field in the hash stored at key to value, only if field does not yet exist
// If key does not exist, a new key holding a hash is created
// If field already exists, this operation has no effect
// If key does not hold a hash, an error is returned

// HTypes Returns if field is an existing field in the hash stored at key
// Returns "hash" if the hash stored at key contains field and the type is hash
// Returns "" if the hash does not contain field, or key does not exist

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
}

// NewHashStructure Returns a new NewHashStructure
func NewHashStructure(options config.Options) (*HashStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &HashStructure{db: db}, nil
}

// HSet sets the string value of a hash field.
func (hs *HashStructure) HSet(k string, f, v interface{}) (bool, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return false, err
	}

	// Convert the parameters to bytes
	value, err, valueType := interfaceToBytes(v)

	if err != nil {
		return false, err
	}

	// Set the hash value type
	hs.hashValueType = valueType

	// Check the parameters
	if len(key) == 0 || len(field) == 0 || len(value) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return false, err
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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
		_ = batch.Put(key, hashMeta.encodeHashMeta())
	}

	// Put the field to the database
	_ = batch.Put(hfBuf, value)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return !exist, nil
}

// HGet gets the string value of a hash field.
func (hs *HashStructure) HGet(k string, f interface{}) (interface{}, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return nil, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return nil, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return nil, err
	}

	// If the counter is 0, return nil
	if hashMeta.counter == 0 {
		return nil, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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

// HDel deletes one field from a hash.
func (hs *HashStructure) HDel(k string, f interface{}) (bool, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return false, nil
	}

	// new a write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Delete the field from the database
	_ = batch.Delete(hfBuf)

	// Decrease the counter
	hashMeta.counter--

	// Put the hash metadata to the database
	_ = batch.Put(key, hashMeta.encodeHashMeta())

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HExists determines whether a hash field exists or not.
func (hs *HashStructure) HExists(k string, f interface{}) (bool, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return false
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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
func (hs *HashStructure) HLen(k string) (int, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Check the parameters
	if len(key) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
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
func (hs *HashStructure) HUpdate(k string, f, v interface{}) (bool, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return false, err
	}

	// Convert the parameters to bytes
	value, err, _ := interfaceToBytes(v)
	if err != nil {
		return false, err
	}
	// Check the parameters
	if len(key) == 0 || len(field) == 0 || len(value) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		return false, nil
	}

	// new a write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Put the field to the database
	_ = batch.Put(hfBuf, value)

	// Commit the write batch
	err = batch.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HIncrBy increments the integer value of a hash field by the given number.
func (hs *HashStructure) HIncrBy(k string, f interface{}, increment int64) (int64, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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

	// new a write batch
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
func (hs *HashStructure) HIncrByFloat(k string, f interface{}, increment float64) (float64, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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

	// new a write batch
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
func (hs *HashStructure) HDecrBy(k string, f interface{}, decrement int64) (int64, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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

	// new a write batch
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

// HStrLen returns the string length of the value associated with field in the hash stored at key.
func (hs *HashStructure) HStrLen(k string, f interface{}) (int, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return 0, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return 0, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return 0, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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

// HMove moves field from the hash stored at source to the hash stored at destination.
func (hs *HashStructure) HMove(source, destination string, f interface{}) (bool, error) {
	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(source) == 0 || len(destination) == 0 || len(field) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given source
	sourceMeta, err := hs.findHashMeta(source, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return 0
	if sourceMeta.counter == 0 {
		return false, nil
	}

	// Find the hash metadata by the given destination
	destinationMeta, err := hs.findHashMeta(destination, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return 0
	if destinationMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     []byte(source),
		version: sourceMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Create a new HashField
	destinationHf := &HashField{
		field:   field,
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

	// new a write batch
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

// HSetNX sets field in the hash stored at key to value, only if field does not yet exist.
func (hs *HashStructure) HSetNX(k string, f, v interface{}) (bool, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return false, err
	}

	// Convert the parameters to bytes
	value, err, _ := interfaceToBytes(v)
	if err != nil {
		return false, err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 || len(value) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return false, err
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
		version: hashMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Get the field from the database
	_, err = hs.db.Get(hfBuf)
	if err != nil && err == _const.ErrKeyNotFound {
		_, err := hs.HSet(k, field, value)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, nil
	}

}

// HTypes returns if field is an existing hash key in the hash stored at key.
func (hs *HashStructure) HTypes(k string, f interface{}) (string, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Convert the parameters to bytes
	field, err, _ := interfaceToBytes(f)
	if err != nil {
		return "", err
	}

	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
		return "", _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(k, Hash)
	if err != nil {
		return "", err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return "", nil
	}

	// Create a new HashField
	hf := &HashField{
		field:   field,
		key:     key,
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

// findHashMeta finds the hash metadata by the given key.
func (hs *HashStructure) findHashMeta(k string, dataType DataStructure) (*HashMetadata, error) {
	// Convert the parameters to bytes
	key := stringToBytesWithKey(k)

	// Find the hash metadata by the given key
	meta, err := hs.db.Get(key)
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

		// Check the expiration time
		if hashMeta.expire > 0 && hashMeta.expire < time.Now().UnixNano() {
			exist = false
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

type HashField struct {
	field   []byte
	key     []byte
	version int64
}

// encodeHashField encodes a HashField and returns the byte array and length.
// +-------------+------------+------------+
// |  field      |  key       |  version   |
// +-------------+------------+------------+
// |  variable   |  variable  |  8 bytes   |
// +-------------+------------+------------+
func (hf *HashField) encodeHashField() []byte {
	buf := make([]byte, len(hf.field)+len(hf.key)+8)

	// offset is the offset of the buf
	var offset = 0

	// copy the field to buf
	offset += copy(buf[offset:], hf.field)

	// copy the key to buf
	offset += copy(buf[offset:], hf.key)

	// copy the version to buf
	binary.BigEndian.PutUint64(buf[offset:], uint64(hf.version))

	return buf[:offset+8]
}

// decodeHashField decodes the HashField from a byte buffer.
func decodeHashField(buf []byte) *HashField {
	var offset = 0

	// get the field from buf
	field := buf[offset:]

	// get the key from buf
	offset += len(field)
	key := buf[offset:]

	// get the version from buf
	offset += len(key)
	version := int64(binary.BigEndian.Uint64(buf[offset:]))

	return &HashField{
		field:   field,
		key:     key,
		version: version,
	}
}

// EncodeHashMeta encodes a HashMetadata and returns the byte array and length.
// +-------------+------------+------------+--------------+-------+---------+---------+---------+---------+
// |  data type   |  data size  |    expire  |    version   | counter | created | updated |  field  |  value  |
// +-------------+------------+------------+--------------+-------+---------+---------+---------+---------+
// |  1 byte      |  variable   |  variable  |   variable   | variable| variable| variable| variable| variable|
// +-------------+------------+------------+--------------+-------+---------+---------+---------+---------+
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
