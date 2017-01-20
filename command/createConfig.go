package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/guywithnose/hostBuilder/hosts"
	"github.com/urfave/cli"
)

// CmdCreateConfig creates a config file from an existing hosts file
func CmdCreateConfig(c *cli.Context) error {
	configFile := c.GlobalString("config")
	hostsFile := c.String("hostsFile")
	if configFile == "" {
		return cli.NewExitError("You must specify a config file", 1)
	}

	hosts, err := hosts.ReadHostsFile(hostsFile)
	if err != nil {
		return err
	}

	return config.WriteConfig(configFile, config.BuildConfigFromHosts(hosts))
}

// CompleteCreateConfig handles bash autocompletion for the 'createConfig' command
func CompleteCreateConfig(c *cli.Context) {
	lastParam := os.Args[len(os.Args)-2]
	if lastParam == "--hostsFile" {
		fmt.Println("fileCompletion")
		return
	}

	for _, flag := range c.App.Command("createConfig").Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
