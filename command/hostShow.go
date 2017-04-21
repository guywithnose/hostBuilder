package command

import (
	"fmt"
	"io"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
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

	found := printOptions(configData, hostName, c.App.Writer)

	if !found {
		pringtGlobalIPInfo(configData, hostName, c.App.Writer)
	}

	return nil
}

func pringtGlobalIPInfo(configData *config.HostsConfig, hostName string, writer io.Writer) {
	if IP, exists := configData.GlobalIPs[configData.Hosts[hostName].Current]; exists {
		fmt.Fprintf(writer, "Current: Global IP %s => %s\n", configData.Hosts[hostName].Current, IP)
	} else {
		fmt.Fprintf(writer, "Current: %s", configData.Hosts[hostName].Current)
		if configData.Hosts[hostName].Current != hostIgnore {
			fmt.Fprint(writer, " (Warning: no associated IP please validate your config)")
		}

		fmt.Fprintln(writer, "")
	}
}

func printOptions(configData *config.HostsConfig, hostName string, writer io.Writer) bool {
	found := false
	numOptions := len(configData.Hosts[hostName].Options)
	pluralSuffix := "s"
	if numOptions == 1 {
		pluralSuffix = ""
	}

	fmt.Fprintf(writer, "%d Option%s:\n", numOptions, pluralSuffix)
	for _, option := range sortOptions(configData, hostName) {
		IP := configData.Hosts[hostName].Options[option]
		if option == configData.Hosts[hostName].Current {
			fmt.Fprintf(writer, "*%s => %s*\n", option, IP)
			found = true
		} else {
			fmt.Fprintf(writer, "%s => %s\n", option, IP)
		}
	}

	return found
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
