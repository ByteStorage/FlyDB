package client

import (
	"github.com/desertbit/grumble"
)

func register(app *grumble.App) {

	app.AddCommand(&grumble.Command{
		Name: "put",
		Help: "put data",
		Run:  stringPutData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "get",
		Help: "get data",
		Run:  stringGetData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "delete",
		Help: "delete key",
		Run:  stringDeleteKey,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "getset",
		Help: "get curvalue set newvalue",
		Run:  stringGetSet,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "exist",
		Help: "exist key",
		Run:  stringExists,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "MGet",
		Help: "Gets the value of multiple keys simultaneously. MGet key1 key2 keyN",
		Run:  stringMGet,
		Args: func(a *grumble.Args) {
			a.StringList("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "MSet",
		Help: "Set multiple key pairs at the same time. MSet key1 value1 key2 value2 .. keyN valueN",
		Run:  stringMSet,
		Args: func(a *grumble.Args) {
			a.StringList("keyvalue", "key value", grumble.Default(""))
		},
	})
}
