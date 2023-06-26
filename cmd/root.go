package cmd

import (
	"github.com/desertbit/grumble"
)

func register(app *grumble.App) {
	app.AddCommand(&grumble.Command{
		Name: "start",
		Help: "start server",
		Run:  startServer,
	})

	app.AddCommand(&grumble.Command{
		Name: "stop",
		Help: "stop server",
		Run:  stopServer,
	})

	app.AddCommand(&grumble.Command{
		Name: "clean",
		Help: "clean server",
		Run:  cleanServer,
	})

	app.AddCommand(&grumble.Command{
		Name: "put",
		Help: "put data",
		Run:  putData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
			a.String("value", "value", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "get",
		Help: "get data",
		Run:  getData,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "delete",
		Help: "delete data",
		Run:  deleteDate,
		Args: func(a *grumble.Args) {
			a.String("key", "key", grumble.Default(""))
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "keys",
		Help: "list keys",
		Run:  getKeys,
	})

}
