package client

import (
	"github.com/desertbit/grumble"
	"math"
)

func register(app *grumble.App) {
	app.AddCommand(&grumble.Command{
		Name: "put",
		Help: "put data in string-structure",
		Run:  stringPutData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "get",
		Help: "get data from string-structure",
		Run:  stringGetData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "delete",
		Help: "delete key in string-structure",
		Run:  stringDeleteKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "strlen",
		Help: "get the length of the value stored in a key in string-structure",
		Run:  stringStrLen,
		Args: func(a *grumble.Args) {
			a.String("key", "The key whose value length to retrieve", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "type",
		Help: "get the type of the value stored in a key in string-structure",
		Run:  stringGetType,
		Args: func(a *grumble.Args) {
			a.String("key", "The key whose value type to retrieve", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "getset",
		Help: "set the value of a key and return its old value in string-structure",
		Run:  stringGetSet,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to set", grumble.Default(""))
			a.String("value", "The new value to set", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "append",
		Help: "append a value to a key in string-structure",
		Run:  stringAppend,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to append to", grumble.Default(""))
			a.String("value", "The value to append", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "incr",
		Help: "increment the integer value of a key in string-structure",
		Run:  stringIncr,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to increment", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "incrby",
		Help: "increment the integer value of a key by a specific amount in string-structure",
		Run:  stringIncrBy,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to increment", grumble.Default(""))
			a.Int64("amount", "The amount to increment by", grumble.Default(1))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "incrbyfloat",
		Help: "increment the float value of a key by a specific amount in string-structure",
		Run:  stringIncrByFloat,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to increment", grumble.Default(""))
			a.Float64("amount", "The amount to increment by", grumble.Default(1.0))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "decr",
		Help: "decrement the integer value of a key in string-structure",
		Run:  stringDecr,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to decrement", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "decrby",
		Help: "decrement the integer value of a key by a specific amount in string-structure",
		Run:  stringDecrBy,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to decrement", grumble.Default(""))
			a.Int64("amount", "The amount to decrement by", grumble.Default(1))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "exists",
		Help: "check if a key exists in string-structure",
		Run:  stringExists,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to check for existence", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "expire",
		Help: "set a timeout on a key in string-structure",
		Run:  stringExpire,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to set a timeout on", grumble.Default(""))
			a.Int64("ttl", "The time-to-live (TTL) in seconds", grumble.Default(0))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "persist",
		Help: "remove the timeout on a key, making it persist in string-structure",
		Run:  stringPersist,
		Args: func(a *grumble.Args) {
			a.String("key", "The key to make persistent", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "mget",
		Help: "get the values of multiple keys in string-structure",
		Run:  stringMGet,
		Args: func(a *grumble.Args) {
			a.StringList("key", "The keys to get values for", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "mset",
		Help: "Set multiple key-value pairs in string-structure",
		Run:  stringMSet,
		Args: func(a *grumble.Args) {
			a.StringList("key-value", "key-value pairs (e.g., key1 value1 key2 value2)", grumble.Default(""))
		},
	})

	// Command for stringMSetNX
	app.AddCommand(&grumble.Command{
		Name: "msetnx",
		Help: "Set multiple key-value pairs if the keys do not exist in string-structure",
		Run:  stringMSetNX,
		Args: func(a *grumble.Args) {
			a.StringList("key-value", "key-value pairs (e.g., key1 value1 key2 value2)", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HSet",
		Help: "put data in hash-structure",
		Run:  hashHSetData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HGet",
		Help: "get data from hash-structure",
		Run:  hashHGetData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HDel",
		Help: "delete key in hash-structure",
		Run:  hashHDelKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HExists",
		Help: "Check if a field exists in a hash in hash-structure",
		Run:  hashHExistsKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HLen",
		Help: "Get the length of a hash in hash-structure",
		Run:  hashHLenKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HUpdate",
		Help: "Update a field in a hash in hash-structure",
		Run:  hashHUpdateKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HIncrby",
		Help: "Increment the value of a field in a hash by an integer in hash-structure",
		Run:  hashHIncrByKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
			a.Int64("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HIncrbyfloat",
		Help: "Increment the value of a field in a hash by a float in hash-structure",
		Run:  hashHIncrByFloatKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
			a.Float64("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HDecrby",
		Help: "Decrement the value of a field in a hash by an integer in hash-structure",
		Run:  hashHDecrByKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
			a.Int64("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HStrlen",
		Help: "Get the length of the string value of a field in a hash in hash-structure",
		Run:  hashHStrLenKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HMove",
		Help: "Move a field from one hash to another in hash-structure",
		Run:  hashHMoveKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
			a.String("dest", "dest", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HSetnx",
		Help: "Set the value of a field in a hash if it does not exist in hash-structure",
		Run:  hashHSetNXKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "HType",
		Help: "Get the type of a field in a hash in hash-structure",
		Run:  hashHType,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("field", "field", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LPush",
		Help: "Inserts a value at the head of a list in list-structure",
		Run:  stringLPushData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LPushs",
		Help: "Inserts multiple values at the head of a list in list-structure",
		Run:  stringLPushsData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.StringList("values", "values", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "RPush",
		Help: "Inserts a value at the tail of a list in list-structure",
		Run:  stringRPushData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "RPushs",
		Help: "Push elements to the end of a list in list-structure",
		Run:  stringRPushsData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.StringList("values", "values", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LPop",
		Help: "Removes and returns the first element of a list in list-structure",
		Run:  stringLPopData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "RPop",
		Help: "Removes and returns the last element of a list in list-structure",
		Run:  stringRPopData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LRange",
		Help: "Returns a range of elements from a list in list-structure",
		Run:  stringLRangeData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("start", "start", grumble.Default(0))
			a.Int("stop", "stop", grumble.Default(-1))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LLen",
		Help: "Returns the length of a list in list-structure",
		Run:  stringLLenData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LRem",
		Help: "Remove elements from a list in list-structure",
		Run:  stringLRemData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("count", "count", grumble.Default(0))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LIndex",
		Help: "Get the element at a specific index in a list in list-structure",
		Run:  stringLIndexData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("index", "index", grumble.Default(0))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LSet",
		Help: "Set the value of an element at a specific index in a list in list-structure",
		Run:  stringLSetData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("index", "index", grumble.Default(0))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "LTrim",
		Help: "Trim a list to a specified range of elements in list-structure",
		Run:  stringLTrimData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("start", "start", grumble.Default(0))
			a.Int("stop", "stop", grumble.Default(0))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zadd",
		Help: "Add the value from a zset in zset-structure",
		Run:  ZSetAdd,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("score", "score", grumble.Default(0))
			a.String("member", "member", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zadds",
		Help: "Add multiple values from a zset in zset-structure",
		Run:  ZSetAdds,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.StringList("members", "members", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zrem",
		Help: "Remove the value from a zset in zset-structure",
		Run:  ZSetRem,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("member", "member", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zrems",
		Help: "Remove multiple values from a zset in zset-structure",
		Run:  ZSetRems,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.StringList("members", "members", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zscore",
		Help: "Get the score of a member in a zset in zset-structure",
		Run:  ZSetScore,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("member", "member", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zrank",
		Help: "Get the rank of a member in a zset in zset-structure",
		Run:  ZSetRank,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("member", "member", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zrevrank",
		Help: "Get the reverse rank of a member in a zset in zset-structure",
		Run:  ZSetRevRank,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("member", "member", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zrange",
		Help: "Get a range of members from a zset in zset-structure",
		Run:  ZSetRange,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("start", "start", grumble.Default(0))
			a.Int("stop", "stop", grumble.Default(-1))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zcount",
		Help: "Count the number of members in a zset within a score range in zset-structure",
		Run:  ZSetCount,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("min", "min", grumble.Default(math.MinInt32))
			a.Int("max", "max", grumble.Default(math.MaxInt32))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zrevrange",
		Help: "Get a reverse range of members from a zset in zset-structure",
		Run:  ZSetRevRange,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int("start", "start", grumble.Default(0))
			a.Int("stop", "stop", grumble.Default(-1))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zcard",
		Help: "Get the number of members in a zset in zset-structure",
		Run:  ZSetCard,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "Zincrby",
		Help: "Increment the score of a member in a zset in zset-structure",
		Run:  ZSetIncrBy,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("member", "member", grumble.Default(""))
			a.Int("increment", "increment", grumble.Default(0))
		},
	})
}
