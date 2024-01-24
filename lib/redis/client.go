package redis

import (
	"fmt"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"strings"

	"github.com/tidwall/redcon"
)

type FlyDBClient struct {
	DB     map[int]interface{}
	Server *FlyDBServer
}

func NewWrongNumberOfArgsError(cmd string) error {
	return fmt.Errorf("ERR wrong number of arguments for '%s' command", cmd)
}

type CmdHandler func(cli *FlyDBClient, args [][]byte) (interface{}, error)

var FlyDBSupportCommands = map[string]CmdHandler{
	// string
	"use-string": UseString,
	"set":        Set,
	"get":        Get,

	// hash
	"use-hash": UseHash,
	"hset":     HSet,
	"hget":     HGet,
}

func ClientCommands(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))
	cmdFunc, ok := FlyDBSupportCommands[command]
	if !ok {
		conn.WriteError("Err unsupported command: '" + command + "'")
		return
	}

	cli, _ := conn.Context().(*FlyDBClient)
	switch command {
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		_ = conn.Close()
	default:
		res, err := cmdFunc(cli, cmd.Args[1:])
		if err != nil {
			if err == _const.ErrKeyNotFound {
				conn.WriteNull()
			} else {
				conn.WriteError(err.Error())
			}
			return
		}
		conn.WriteAny(res)
	}
}
