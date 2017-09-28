package command

import (
	"fmt"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdGroupAdd adds a hostname to a group
func CmdGroupAdd(c *cli.Context) error {
	if c.NArg() != 2 {
		return cli.NewExitError("Usage: \"hostBuilder group add {groupName} {hostName}\"", 1)
	}

	groupName := c.Args().Get(0)
	hostName := c.Args().Get(1)

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	if _, exists := configData.Groups[groupName]; !exists {
		if configData.Groups == nil {
			configData.Groups = map[string][]string{}
		}

		configData.Groups[groupName] = []string{hostName}
	} else {
		if _, exists := configData.Hosts[hostName]; !exists {
			return cli.NewExitError(fmt.Sprintf("Hostname %s does not exist", hostName), 1)
		}

		if groupContains(configData, groupName, hostName) {
			return cli.NewExitError(fmt.Sprintf("Group %s already contains %s", groupName, hostName), 1)
		}

		configData.Groups[groupName] = append(configData.Groups[groupName], hostName)
	}

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteGroupAdd handles bash autocompletion for the 'group add' command
func CompleteGroupAdd(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	var options []string
	if c.NArg() == 0 {
		options = sortGroupNames(configData)
	} else if c.NArg() == 1 {
		options = sortHostNames(configData)
	}

	fmt.Fprintln(c.App.Writer, strings.Join(options, "\n"))
}
