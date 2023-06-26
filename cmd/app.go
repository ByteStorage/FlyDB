package cmd

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"os"
	"path"
	"strings"
)

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
		return nil
	})
	App.SetPrintASCIILogo(func(a *grumble.App) {
		fmt.Println(strings.Join([]string{`
 __________    __           ___    ___ 
|  |______/   |  |          \  \  /  / 
|  |          |  |           \  \/  /     
|  |________  |  |            \    /
|  |________| |  |             \  /
|  |          |  |             |  |
|  |          |  |_________    |  | 
|__|          |____________|   |__| 
		`}, "\r\n"))
	})
	register(App)
}
