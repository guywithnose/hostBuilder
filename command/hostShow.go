package command

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// CmdHostShow describes the current state of a host
func CmdHostShow(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Usage: \"hostBuilder host show {hostName}\"", 1)
	}

	hostName := c.Args().Get(0)

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	if _, exists := configData.Hosts[hostName]; !exists {
		return cli.NewExitError(fmt.Sprintf("Hostname %s does not exist", hostName), 1)
	}

	found := false
	numOptions := len(configData.Hosts[hostName].Options)
	pluralSuffix := "s"
	if numOptions == 1 {
		pluralSuffix = ""
	}

	fmt.Fprintf(c.App.Writer, "%d Option%s:\n", numOptions, pluralSuffix)
	for _, option := range sortOptions(configData, hostName) {
		IP := configData.Hosts[hostName].Options[option]
		if option == configData.Hosts[hostName].Current {
			fmt.Fprintf(c.App.Writer, "*%s => %s*\n", option, IP)
			found = true
		} else {
			fmt.Fprintf(c.App.Writer, "%s => %s\n", option, IP)
		}
	}

	if !found {
		if IP, exists := configData.GlobalIPs[configData.Hosts[hostName].Current]; exists {
			fmt.Fprintf(c.App.Writer, "Current: Global IP %s => %s\n", configData.Hosts[hostName].Current, IP)
		} else {
			fmt.Fprintf(c.App.Writer, "Current: %s", configData.Hosts[hostName].Current)
			if configData.Hosts[hostName].Current != "ignore" {
				fmt.Fprint(c.App.Writer, " (Warning: no associated IP please validate your config)")
			}

			fmt.Fprintln(c.App.Writer, "")
		}
	}

	return nil
}

// CompleteHostShow handles bash autocompletion for the 'host show' command
func CompleteHostShow(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	if c.NArg() == 0 {
		fmt.Fprintln(c.App.Writer, strings.Join(sortHostNames(configData), "\n"))
	}
}
