package redis

import (
	"github.com/ByteStorage/FlyDB/structure"
	"github.com/tidwall/redcon"
)

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

func UseHash(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 0 {
		return nil, NewWrongNumberOfArgsError("use-hash")
	}

	cli.DB[0].(*structure.StringStructure).Close()
	return redcon.SimpleString("OK"), nil
}
