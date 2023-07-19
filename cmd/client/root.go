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
}
