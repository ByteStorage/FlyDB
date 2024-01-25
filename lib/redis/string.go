package redis

import (
	"github.com/ByteStorage/FlyDB/structure"
	"github.com/tidwall/redcon"
)

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

func UseString(cli *FlyDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 0 {
		return nil, NewWrongNumberOfArgsError("use-str")
	}
	return redcon.SimpleString("OK"), nil
}
