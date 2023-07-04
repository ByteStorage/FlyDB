package client

import (
	"errors"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"os"
	"path"
	"strings"
)

var addr string

// App FlyDB command app
var App = grumble.New(&grumble.Config{
	Name:                  "FlyDB Cli",
	Description:           "A command of FlyDB",
	HistoryFile:           path.Join(os.TempDir(), ".FlyDB_Cli.history"),
	HistoryLimit:          10000,
	ErrorColor:            color.New(color.FgRed, color.Bold, color.Faint),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: false,
	HelpSubCommands:       true,
	Prompt:                "flydb $> ",
	PromptColor:           color.New(color.FgBlue, color.Bold),
	Flags:                 func(f *grumble.Flags) {},
})

func init() {
	App.OnInit(func(a *grumble.App, fm grumble.FlagMap) error {
		if len(os.Args) != 1 {
			fmt.Println("usage: flydb-cli [addr]")
			return errors.New("usage: flydb-cli [addr]")
		}
		addr = os.Args[1]
		return nil
	})
	App.SetPrintASCIILogo(func(a *grumble.App) {
		fmt.Println(strings.Join([]string{`
              ______    __             ____     ____ 
             / ____/   / /   __  __   / __ \   / __ )
            / /       / /   / / / /  / / / /  / /_/ /
           / /_      / /   / / / /  / / / /  / __ |
          / __/     / /   / /_/ /  / / / /  / / / / 
         / /       / /    \__, /  / /_/ /  / /_/ /  
        /_/       /_/    ,__/ /  /_____/  /_____/  
                        /____/                                              
`}, "\r\n"))
	})
	register(App)
}
