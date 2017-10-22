package command

import (
	"fmt"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdGroupSet sets the hostnames in a group to a global ip
func CmdGroupSet(c *cli.Context) error {
	if c.NArg() != 2 {
		return cli.NewExitError("Usage: \"hostBuilder group set {groupName} {globalIPName}\"", 1)
	}

	groupName := c.Args().Get(0)
	globalIPName := c.Args().Get(1)

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	if _, exists := configData.Groups[groupName]; !exists {
		return cli.NewExitError(fmt.Sprintf("Group %s does not exist", groupName), 1)
	}

	if _, exists := configData.GlobalIPs[globalIPName]; !exists && globalIPName != "ignore" {
		return cli.NewExitError(fmt.Sprintf("Global IP %s does not exist", globalIPName), 1)
	}

	for _, hostName := range configData.Groups[groupName] {
		host := configData.Hosts[hostName]
		host.Current = globalIPName
		configData.Hosts[hostName] = host
	}

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteGroupSet handles bash autocompletion for the 'group set' command
func CompleteGroupSet(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	if c.NArg() == 0 {
		fmt.Fprintln(c.App.Writer, strings.Join(sortGroupNames(configData), "\n"))
	} else if c.NArg() == 1 {
		globalIPs := sortGlobalIPNames(configData)
		for _, IPName := range globalIPs {
			fmt.Fprintf(c.App.Writer, "%s:%s\n", IPName, configData.GlobalIPs[IPName])
		}
	}
}
