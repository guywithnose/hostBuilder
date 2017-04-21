package command

import (
	"fmt"
	"io"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdHostSet sets the current IP on a hostname
func CmdHostSet(c *cli.Context) error {
	if c.NArg() != 2 {
		return cli.NewExitError("Usage: \"hostBuilder host set {hostName} {IPName}\"", 1)
	}

	hostName := c.Args().Get(0)
	IPName := c.Args().Get(1)

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	err = validateParameters(configData, hostName, IPName)
	if err != nil {
		return err
	}

	host := configData.Hosts[hostName]
	host.Current = IPName
	configData.Hosts[hostName] = host

	return config.WriteConfig(c.GlobalString("config"), configData)
}

func validateParameters(configData *config.HostsConfig, hostName, IPName string) error {
	if _, exists := configData.Hosts[hostName]; !exists {
		return cli.NewExitError(fmt.Sprintf("HostName %s does not exist", hostName), 1)
	}

	if _, exists := configData.Hosts[hostName].Options[IPName]; !exists && IPName != hostIgnore {
		if _, exists := configData.GlobalIPs[IPName]; !exists {
			return cli.NewExitError(fmt.Sprintf("IPName %s does not exist", IPName), 1)
		}
	}

	return nil
}

// CompleteHostSet handles bash autocompletion for the 'host set' command
func CompleteHostSet(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	if c.NArg() == 0 {
		fmt.Fprintln(c.App.Writer, strings.Join(sortHostNames(configData), "\n"))
	} else if c.NArg() == 1 {
		hostName := c.Args().Get(0)
		printIPs(configData, hostName, c.App.Writer)
	}
}

func printIPs(configData *config.HostsConfig, hostName string, writer io.Writer) {
	if _, exists := configData.Hosts[hostName]; !exists {
		return
	}

	for _, option := range sortOptions(configData, hostName) {
		fmt.Fprintf(writer, "%s:%s\n", option, configData.Hosts[hostName].Options[option])
	}
	for _, globalIPName := range sortGlobalIPNames(configData) {
		fmt.Fprintf(writer, "%s:%s\n", globalIPName, configData.GlobalIPs[globalIPName])
	}

	fmt.Fprintln(writer, "ignore:")
}
