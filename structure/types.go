package structure

import (
	"errors"
	"strconv"
)

type DataStructure = byte

const (
	// String is a string data structure
	String DataStructure = iota + 1
	// Hash is a hash data structure
	Hash
	// List is a list data structure
	List
	// Set is a set data structure
	Set
	// ZSet is a zset data structure
	ZSet
	// bitmap is a bitmap data structure
	Bitmap
	// Stream is a stream data structure
	Stream
	// Expire is a expire data structure
	Expire
)

// interfaceToBytes converts an interface to a byte slice
func interfaceToBytes(value interface{}) ([]byte, error) {
	switch value := value.(type) {
	case string:
		return []byte(value), nil
	case int:
		return []byte(strconv.Itoa(value)), nil
	case int64:
		return []byte(strconv.FormatInt(value, 10)), nil
	case float64:
		return []byte(strconv.FormatFloat(value, 'f', -1, 64)), nil
	case bool:
		return []byte(strconv.FormatBool(value)), nil
	case []byte:
		return value, nil
	default:
		return nil, errors.New("unsupported type")
	}
}
