package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/guywithnose/hostBuilder/hosts"
	"github.com/urfave/cli"
)

// CmdBuild builds the hostfile from a configuration file
func CmdBuild(c *cli.Context) error {
	outputFile := c.String("output")
	if outputFile == "" {
		return cli.NewExitError("You must specify an output file", 1)
	}

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	return hosts.OutputHostLines(outputFile, configData, c.Bool("oneLinePerIP"))
}

// CompleteBuild handles bash autocompletion for the 'build' command
func CompleteBuild(c *cli.Context) {
	lastParam := os.Args[len(os.Args)-2]
	if lastParam == "--output" {
		fmt.Println("fileCompletion")
		return
	}

	for _, flag := range c.App.Command("build").Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
