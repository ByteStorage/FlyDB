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
func interfaceToBytes(value interface{}) ([]byte, error, string) {
	switch value := value.(type) {

	case string:
		return []byte(value), nil, "string"
	case int:
		return []byte(strconv.Itoa(value)), nil, "int"
	case int64:
		return []byte(strconv.FormatInt(value, 10)), nil, "int64"
	case float64:
		return []byte(strconv.FormatFloat(value, 'f', -1, 64)), nil, "float64"
	case bool:
		return []byte(strconv.FormatBool(value)), nil, "bool"
	case []byte:
		return value, nil, "[]byte"
	default:
		return nil, errors.New("unsupported type"), ""
	}
}

func byteToInterface(value []byte, dataType string) (interface{}, error) {
	switch dataType {
	case "string":
		return string(value), nil
	case "int":
		return strconv.Atoi(string(value))
	case "int64":
		return strconv.ParseInt(string(value), 10, 64)
	case "float64":
		return strconv.ParseFloat(string(value), 64)
	case "bool":
		return strconv.ParseBool(string(value))
	case "[]byte":
		return value, nil
	default:
		return nil, errors.New("unsupported type")
	}
}
