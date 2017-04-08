package command

import (
	"fmt"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdHostRemove removes a host from the configuration
func CmdHostRemove(c *cli.Context) error {
	if c.NArg() != 2 {
		return cli.NewExitError("Usage: \"hostBuilder host remove {hostName} {IPName}\"", 1)
	}

	hostName := c.Args().Get(0)
	IPName := c.Args().Get(1)

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	if _, exists := configData.Hosts[hostName]; !exists {
		return cli.NewExitError(fmt.Sprintf("Host %s does not exist", hostName), 1)
	}

	if _, exists := configData.Hosts[hostName].Options[IPName]; !exists {
		return cli.NewExitError(fmt.Sprintf("IPName %s does not exist", IPName), 1)
	}

	host := configData.Hosts[hostName]
	delete(host.Options, IPName)
	if len(host.Options) == 0 {
		host.Current = hostIgnore
	}

	configData.Hosts[hostName] = host

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteHostRemove handles bash autocompletion for the 'host remove' command
func CompleteHostRemove(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	if c.NArg() == 0 {
		fmt.Fprintln(c.App.Writer, strings.Join(sortHostNames(configData), "\n"))
	} else if c.NArg() == 1 {
		hostName := c.Args().Get(0)
		if _, exists := configData.Hosts[hostName]; !exists {
			return
		}

		options := sortOptions(configData, hostName)
		for _, option := range options {
			fmt.Fprintf(c.App.Writer, "%s:%s\n", option, configData.Hosts[hostName].Options[option])
		}
	}
}
