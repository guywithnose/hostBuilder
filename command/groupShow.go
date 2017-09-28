package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

// CmdGroupShow lists the hostnames in a group
func CmdGroupShow(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Usage: \"hostBuilder group show {groupName}\"", 1)
	}

	groupName := c.Args().Get(0)

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	if _, exists := configData.Groups[groupName]; !exists {
		return cli.NewExitError(fmt.Sprintf("Group %s does not exist", groupName), 1)
	}

	hostNames := configData.Groups[groupName]

	sort.Strings(hostNames)

	for _, hostName := range hostNames {
		fmt.Fprintln(c.App.Writer, hostName)
	}

	return nil
}

// CompleteGroupShow handles bash autocompletion for the 'group show' command
func CompleteGroupShow(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	if c.NArg() == 0 {
		fmt.Fprintln(c.App.Writer, strings.Join(sortGroupNames(configData), "\n"))
	}

}
