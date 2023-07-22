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
