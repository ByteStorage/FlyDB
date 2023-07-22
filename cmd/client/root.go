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
		Name: "type",
		Help: "get the type of a key",
		Run:  stringGetType,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "strlen",
		Help: "str value length in string-structure",
		Run:  stringStrLen,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "getset",
		Help: "get curvalue set newvalue in string-structure",
		Run:  stringGetSet,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "append",
		Help: "append value in string-structure",
		Run:  stringAppend,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "incr",
		Help: "increment the value of a key in string-structure",
		Run:  stringIncr,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "incrby",
		Help: "increment the value of a key by a specific amount in string-structure",
		Run:  stringIncrBy,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int64("amount", "amount", grumble.Default(0))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "incrbyfloat",
		Help: "increment the value of a key by a floating-point amount in string-structure",
		Run:  stringIncrByFloat,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Float64("amount", "amount", grumble.Default(0.0))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "decr",
		Help: "decrement the value of a key",
		Run:  stringDecr,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "decrby",
		Help: "decrement the value of a key by a specific amount in string-structure",
		Run:  stringDecrBy,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int64("amount", "amount", grumble.Default(0))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "exist",
		Help: "Check if the given key exists in string-structure",
		Run:  stringExists,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "expire",
		Help: "Set the expiration time for the key, which will no longer be available after expiration. Unit in seconds in string-structure",
		Run:  stringExpire,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.Int64("ttl", "time of this key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "persist",
		Help: "Remove the expiration time of the given key so that it never expires in string-structure",
		Run:  stringPersist,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "MGet",
		Help: "Gets the value of multiple keys simultaneously in string-structure",
		Run:  stringMGet,
		Args: func(a *grumble.Args) {
			a.StringList("key", "key", grumble.Default(""))
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
			a.StringList("values", "values", grumble.Default([]string{}))
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

}
