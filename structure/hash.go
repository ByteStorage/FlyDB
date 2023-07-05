package structure

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
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

