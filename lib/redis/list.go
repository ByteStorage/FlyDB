package redis

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/structure"
	"github.com/tidwall/redcon"
)

// LPush key value
func LPush(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) < 2 {
		return nil, NewWrongNumberOfArgsError("lpush")
	}

	key := string(args[0])
	values := make([]string, len(args)-1)
	for i, v := range args[1:] {
		values[i] = string(v)
	}

	if err := cli.DB[2].(*structure.ListStructure).LPush(key, values, 0); err != nil {
		return nil, err
	}
	return redcon.SimpleInt(int64(len(values))), nil
}

// RPush key value
func RPush(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) < 2 {
		return nil, NewWrongNumberOfArgsError("rpush")
	}

	key := string(args[0])
	values := make([]string, len(args)-1)
	for i, v := range args[1:] {
		values[i] = string(v)
	}

	if err := cli.DB[2].(*structure.ListStructure).RPush(key, values, 0); err != nil {
		return nil, err
	}
	return redcon.SimpleInt(int64(len(values))), nil
}

// LPop key
func LPop(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("lpop")
	}

	value, err := cli.DB[2].(*structure.ListStructure).LPop(string(args[0]))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// RPop key
func RPop(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("rpop")
	}

	value, err := cli.DB[2].(*structure.ListStructure).RPop(string(args[0]))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// LIndex key index
func LIndex(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("lindex")
	}

	value, err := cli.DB[2].(*structure.ListStructure).LIndex(string(args[0]), int(binary.BigEndian.Uint64(args[1])))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// LRange key start stop
func LRange(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("lrange")
	}

	values, err := cli.DB[2].(*structure.ListStructure).LRange(string(args[0]), int(binary.BigEndian.Uint64(args[1])), int(binary.BigEndian.Uint64(args[2])))
	if err != nil {
		return nil, err
	}
	return values, nil
}

// LRem key count value
func LRem(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("lrem")
	}

	count := int(binary.BigEndian.Uint64(args[1]))
	if err := cli.DB[2].(*structure.ListStructure).LRem(string(args[0]),
		count, string(args[2])); err != nil {
		return nil, err
	}

	return redcon.SimpleString("OK"), nil
}

// LLen key
func LLen(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("llen")
	}

	length, err := cli.DB[2].(*structure.ListStructure).LLen(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(int64(length)), nil
}

// LSet key index value
func LSet(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("lset")
	}

	index := int(binary.BigEndian.Uint64(args[1]))
	if err := cli.DB[2].(*structure.ListStructure).LSet(string(args[0]), index, string(args[2]), 0); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// LTrim key start stop
func LTrim(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("ltrim")
	}

	if err := cli.DB[2].(*structure.ListStructure).LTrim(string(args[0]), int(binary.BigEndian.Uint64(args[1])), int(binary.BigEndian.Uint64(args[2]))); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// LKeys key
func LKeys(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("lkeys")
	}

	keys, err := cli.DB[2].(*structure.ListStructure).Keys(string(args[0]))
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// LSize key
func LSize(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("lsize")
	}

	size, err := cli.DB[2].(*structure.ListStructure).Size(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleString(size), nil
}

// UseList change to list db
func UseList(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 0 {
		return nil, NewWrongNumberOfArgsError("use-list")
	}
	return redcon.SimpleString("OK"), nil
}
