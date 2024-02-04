package redis

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/structure"
	"github.com/tidwall/redcon"
)

// Set key value
func Set(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("set")
	}

	key, value := args[0], args[1]
	if err := cli.DB[0].(*structure.StringStructure).Set(string(key), string(value), 0); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// Get key
func Get(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("get")
	}

	value, err := cli.DB[0].(*structure.StringStructure).Get(string(args[0]))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Del key
func Del(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("del")
	}

	if err := cli.DB[0].(*structure.StringStructure).Del(string(args[0])); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// GetSet key value
func GetSet(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("getset")
	}

	value, err := cli.DB[0].(*structure.StringStructure).GetSet(string(args[0]), string(args[1]), 0)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Append key value
func Append(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("append")
	}

	if err := cli.DB[0].(*structure.StringStructure).Append(string(args[0]),
		string(args[1]), 0); err != nil {
		return nil, err
	}

	return redcon.SimpleString("OK"), nil
}

// Strlen key
func Strlen(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("strlen")
	}

	length, err := cli.DB[0].(*structure.StringStructure).StrLen(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(length), nil
}

// Incr key
func Incr(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("incr")
	}

	if err := cli.DB[0].(*structure.StringStructure).Incr(string(args[0]),
		int64(binary.BigEndian.Uint64(args[1]))); err != nil {
		return nil, err
	}

	return redcon.SimpleString("OK"), nil
}

// IncrBy key increment
func IncrBy(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("incrby")
	}

	if err := cli.DB[0].(*structure.StringStructure).Incr(string(args[0]),
		int64(binary.BigEndian.Uint64(args[1]))); err != nil {
		return nil, err
	}

	return redcon.SimpleString("OK"), nil
}

// Decr key
func Decr(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("decr")
	}

	if err := cli.DB[0].(*structure.StringStructure).Decr(string(args[0]),
		int64(binary.BigEndian.Uint64(args[1]))); err != nil {
		return nil, err
	}

	return redcon.SimpleString("OK"), nil
}

// DecrBy key decrement
func DecrBy(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("decrby")
	}

	if err := cli.DB[0].(*structure.StringStructure).Decr(string(args[0]),
		int64(binary.BigEndian.Uint64(args[1]))); err != nil {
		return nil, err
	}

	return redcon.SimpleString("OK"), nil
}

// Keys pattern
func Keys(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("keys")
	}

	keys, err := cli.DB[0].(*structure.StringStructure).Keys(string(args[0]))
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// Exists key
func Exists(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("exists")
	}

	ok, err := cli.DB[0].(*structure.StringStructure).Exists(string(args[0]))
	if err != nil {
		return nil, err
	}

	if ok {
		return redcon.SimpleInt(1), nil
	} else {
		return redcon.SimpleInt(0), nil
	}
}

// Expire key seconds
func Expire(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("expire")
	}

	if err := cli.DB[0].(*structure.StringStructure).Expire(string(args[0]),
		int64(binary.BigEndian.Uint64(args[1]))); err != nil {
		return nil, err
	}
	return redcon.SimpleInt(1), nil
}

// Persist key
func Persist(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("persist")
	}

	if err := cli.DB[0].(*structure.StringStructure).Persist(string(args[0])); err != nil {
		return nil, err
	}
	return redcon.SimpleInt(1), nil
}

// TTL key
func TTL(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("ttl")
	}

	ttl, err := cli.DB[0].(*structure.StringStructure).TTL(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(ttl), nil
}

// Size key
func Size(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("size")
	}

	size, err := cli.DB[0].(*structure.StringStructure).Size(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleString(size), nil
}

// UseString change to string db
func UseString(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 0 {
		return nil, NewWrongNumberOfArgsError("use-str")
	}
	return redcon.SimpleString("OK"), nil
}
