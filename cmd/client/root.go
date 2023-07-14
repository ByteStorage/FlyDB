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
}
