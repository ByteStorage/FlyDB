package structure

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
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

type HashStructure struct {
	db *engine.DB
}

// findHashMeta finds the hash metadata by the given key.
func (hs *HashStructure) findHashMeta(key []byte, dataType DataStructure) (*HashMetadata, error) {
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

func NewHashStructure(options config.Options) (*HashStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &HashStructure{db: db}, nil
}

// HSet sets the string value of a hash field.
func (hs *HashStructure) HSet(key, field, value []byte) (bool, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 || len(value) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
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
func (hs *HashStructure) HGet(key, field []byte) ([]byte, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
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

	return value, nil
}

// HDel deletes one or more hash fields.
func (hs *HashStructure) HDel(key []byte, fields ...[]byte) (bool, error) {
	// Check the parameters
	if len(key) == 0 || len(fields) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
	if err != nil {
		return false, err
	}

	// If the counter is 0, return 0
	if hashMeta.counter == 0 {
		return false, nil
	}

	// Create a new HashField
	hf := &HashField{
		key:     key,
		version: hashMeta.version,
	}

	var count int64

	// new a write batch
	batch := hs.db.NewWriteBatch(config.DefaultWriteBatchOptions)

	// Delete the fields one by one
	for _, field := range fields {
		// If the field is not found, continue
		hf.field = field
		hfBuf := hf.encodeHashField()
		_, err = hs.db.Get(hfBuf)
		if err != nil && err == _const.ErrKeyNotFound {
			continue
		}

		// Delete the field
		_ = batch.Delete(hfBuf)

		// Decrease the counter
		hashMeta.counter--
		count++
	}

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
func (hs *HashStructure) HExists(key, field []byte) (bool, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
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
func (hs *HashStructure) HLen(key []byte) (int, error) {
	// Check the parameters
	if len(key) == 0 {
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
func (hs *HashStructure) HUpdate(key, field, value []byte) (bool, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 || len(value) == 0 {
		return false, _const.ErrKeyIsEmpty
	}

	// Find the hash metadata by the given key
	hashMeta, err := hs.findHashMeta(key, Hash)
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
func (hs *HashStructure) HIncrBy(key, field []byte, increment int64) (int64, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
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
func (hs *HashStructure) HIncrByFloat(key, field []byte, increment float64) (float64, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
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
func (hs *HashStructure) HDecrBy(key, field []byte, decrement int64) (int64, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
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
func (hs *HashStructure) HStrLen(key, field []byte) (int, error) {
	// Check the parameters
	if len(key) == 0 || len(field) == 0 {
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
func (hs *HashStructure) HMove(source, destination, field []byte) (bool, error) {
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
		key:     source,
		version: sourceMeta.version,
	}

	// Encode the HashField
	hfBuf := hf.encodeHashField()

	// Create a new HashField
	destinationHf := &HashField{
		field:   field,
		key:     destination,
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
