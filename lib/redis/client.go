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

// FlyDBSupportCommands is the map of all supported redis commands
var FlyDBSupportCommands = map[string]CmdHandler{
	// string
	"use-str": UseString,
	"set":     Set,
	"get":     Get,
	"del":     Del,
	"getset":  GetSet,
	"append":  Append,
	"strlen":  Strlen,
	"incr":    Incr,
	"decr":    Decr,
	"incrby":  IncrBy,
	"decrby":  DecrBy,
	"keys":    Keys,
	"exists":  Exists,
	"expire":  Expire,
	"persist": Persist,
	"ttl":     TTL,
	"size":    Size,

	// hash
	"use-hash": UseHash,
	"hset":     HSet,
	"hget":     HGet,
	"hdel":     HDel,
	"hdelall":  HDelAll,
	"hexists":  HExists,
	"hexpire":  HExpire,
	"hlen":     HLen,
	"hupdate":  HUpdate,
	"hkeys":    HKeys,
	"hstrlen":  HStrlen,
	"hmove":    HMove,
	"hsize":    HSize,
	"httl":     HTTL,
	"hincrby":  HIncrBy,
	"hdecrby":  HDecrBy,
}

// ClientCommands is the handler for all redis commands
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
