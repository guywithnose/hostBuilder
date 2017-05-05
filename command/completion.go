package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

// RootCompletion prints the list of root commands as the root completion method
// This is similar to the default method, but it excludes aliases
func RootCompletion(c *cli.Context) {
	lastParam := os.Args[len(os.Args)-2]
	if lastParam == "--config" {
		fmt.Fprintln(c.App.Writer, "fileCompletion")
		return
	}

	for _, command := range c.App.Commands {
		if command.Hidden {
			continue
		}

		fmt.Fprintf(c.App.Writer, "%s:%s\n", command.Name, command.Usage)
	}

	for _, flag := range c.App.Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
