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
	Bitmap DataStructure = iota + 1
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

//func byteToInterface(value []byte, dataType string) (interface{}, error) {
//	switch dataType {
//	case "string":
//		return string(value), nil
//	case "int":
//		return strconv.Atoi(string(value))
//	case "int64":
//		return strconv.ParseInt(string(value), 10, 64)
//	case "float64":
//		return strconv.ParseFloat(string(value), 64)
//	case "bool":
//		return strconv.ParseBool(string(value))
//	case "[]byte":
//		return value, nil
//	default:
//		return nil, errors.New("unsupported type")
//	}
//}

func byteToInterface(value []byte, dataType string) (interface{}, error) {
	switch dataType {
	case "string":
		return string(value), nil
	case "int":
		intValue, err := strconv.Atoi(string(value))
		if err != nil {
			return nil, errors.New("cannot convert to int: " + err.Error())
		}
		return intValue, nil
	case "int64":
		int64Value, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return nil, errors.New("cannot convert to int64: " + err.Error())
		}
		return int64Value, nil
	case "float64":
		float64Value, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return nil, errors.New("cannot convert to float64: " + err.Error())
		}
		return float64Value, nil
	case "bool":
		boolValue, err := strconv.ParseBool(string(value))
		if err != nil {
			return nil, errors.New("cannot convert to bool: " + err.Error())
		}
		return boolValue, nil
	case "[]byte":
		return value, nil
	default:
		return nil, errors.New("unsupported type")
	}
}

// convertToInt converts an interface to an int
func convertToInt(Value interface{}) (int, error) {
	switch value := Value.(type) {
	case int:
		return value, nil
	case int32:
		return int(value), nil
	case int64:
		return int(value), nil
	case float32:
		return int(value), nil
	case float64:
		return int(value), nil
	case uint:
		return int(value), nil
	case uint8:
		return int(value), nil
	case uint16:
		return int(value), nil
	case uint32:
		return int(value), nil
	case uint64:
		return int(value), nil
	case string:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return intValue, nil
	case []byte:
		strValue := string(value)
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			return 0, err
		}
		return intValue, nil
	default:
		return 0, errors.New("unsupported type")
	}
}

// convertToFloat converts an interface to a float64
func convertToFloat(Value interface{}) (float64, error) {
	switch value := Value.(type) {
	case int:
		return float64(value), nil
	case uint:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case string:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err
		}
		return floatValue, nil
	case []byte:
		strValue := string(value)
		floatValue, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return 0, err
		}
		return floatValue, nil
	default:
		return 0, errors.New("unsupported type")
	}
}
