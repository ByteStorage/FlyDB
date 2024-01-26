package redis

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/structure"
	"github.com/tidwall/redcon"
)

// Set key value
func HSet(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("hset")
	}

	var ok = 0
	key, field, value := args[0], args[1], args[2]
	res, err := cli.DB[1].(*structure.HashStructure).HSet(string(key), field, value)
	if err != nil {
		return nil, err
	}
	if res {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}

// HGet key field
func HGet(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("hget")
	}

	value, err := cli.DB[1].(*structure.HashStructure).HGet(string(args[0]), args[1])
	if err != nil {
		return nil, err
	}
	return value, nil
}

// HDel key field
func HDel(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("hdel")
	}

	if _, err := cli.DB[1].(*structure.HashStructure).HDel(string(args[0]), args[1]); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// HDelAll key
func HDelAll(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("hdelall")
	}

	if _, err := cli.DB[1].(*structure.HashStructure).HDelAll(string(args[0])); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// HExists key field
func HExists(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("hexists")
	}

	var ok = 0
	res, err := cli.DB[1].(*structure.HashStructure).HExists(string(args[0]), args[1])
	if err != nil {
		return nil, err
	}
	if res {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}

// HExpire key seconds
func HExpire(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("hexpire")
	}

	if _, err := cli.DB[1].(*structure.HashStructure).HExpire(string(args[0]),
		int64(binary.BigEndian.Uint64(args[1]))); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// HLen key
func HLen(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("hlen")
	}

	length, err := cli.DB[1].(*structure.HashStructure).HLen(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(length), nil
}

// HUpdate key field value
func HUpdate(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("hupdate")
	}

	if _, err := cli.DB[1].(*structure.HashStructure).HUpdate(string(args[0]), args[1], args[2]); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// HKeys key
func HKeys(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("hkeys")
	}

	keys, err := cli.DB[1].(*structure.HashStructure).Keys(string(args[0]))
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// HStrlen key field
func HStrlen(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgsError("hstrlen")
	}

	length, err := cli.DB[1].(*structure.HashStructure).HStrLen(string(args[0]), args[1])
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(length), nil
}

// HMove source destination field
func HMove(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("hmove")
	}

	if _, err := cli.DB[1].(*structure.HashStructure).HMove(string(args[0]), string(args[1]),
		int64(binary.BigEndian.Uint64(args[2]))); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

// HSize key
func HSize(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("hsize")
	}

	size, err := cli.DB[1].(*structure.HashStructure).Size(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleString(size), nil
}

// HTTL key
func HTTL(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgsError("httl")
	}

	ttl, err := cli.DB[1].(*structure.HashStructure).TTL(string(args[0]))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(ttl), nil
}

// HIncrBy key field increment
func HIncrBy(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("hincrby")
	}

	value, err := cli.DB[1].(*structure.HashStructure).HIncrBy(string(args[0]), args[1],
		int64(binary.BigEndian.Uint64(args[2])))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(value), nil
}

// HDecrBy key field decrement
func HDecrBy(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, NewWrongNumberOfArgsError("hdecrby")
	}

	value, err := cli.DB[1].(*structure.HashStructure).HDecrBy(string(args[0]), args[1],
		int64(binary.BigEndian.Uint64(args[2])))
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(value), nil
}

// UseHash change to hash db
func UseHash(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 0 {
		return nil, NewWrongNumberOfArgsError("use-hash")
	}
	return redcon.SimpleString("OK"), nil
}
