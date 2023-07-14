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

}
