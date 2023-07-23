package client

import "github.com/desertbit/grumble"

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
			a.StringList("key", "The keys to get values for", grumble.Default([]string{}))
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
}
