package command

import (
	"fmt"
	"io"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdHostAdd adds an IP to a hostname
func CmdHostAdd(c *cli.Context) error {
	configData, configErr := loadConfig(c)
	if configErr != nil {
		return configErr
	}

	if c.NArg() == 3 {
		return addHost(configData, c, c.Bool("force"), c.App.ErrWriter)
	}

	if c.NArg() == 2 {
		return addGlobalIPHost(configData, c)
	}

	return cli.NewExitError("Usage: \"hostBuilder host add {hostName} ({address} {IPName}|{globalIpName})\"", 1)
}

func addGlobalIPHost(configData *config.HostsConfig, c *cli.Context) error {
	hostName := c.Args().Get(0)
	globalIPName := c.Args().Get(1)

	if _, exists := configData.Hosts[hostName]; exists {
		return cli.NewExitError(fmt.Sprintf("IP %s already exists", hostName), 1)
	}

	configData.Hosts[hostName] = config.Host{Current: globalIPName}

	return config.WriteConfig(c.GlobalString("config"), configData)
}

func addHost(configData *config.HostsConfig, c *cli.Context, force bool, errWriter io.Writer) error {
	hostName := c.Args().Get(0)
	address, err := resolveAddress(c.Args().Get(1))
	if err != nil {
		return err
	}

	IPName := c.Args().Get(2)

	if _, exists := configData.Hosts[hostName]; !exists {
		configData.Hosts[hostName] = config.Host{Current: IPName, Options: map[string]string{IPName: address}}
	} else {
		if current, exists := configData.Hosts[hostName].Options[IPName]; exists {
			if force {
				fmt.Fprintf(errWriter, "Warning: Overwriting %s (%s => %s)", IPName, current, address)
			} else {
				return cli.NewExitError(fmt.Sprintf("IP %s already exists", hostName), 1)
			}
		}

		configData.Hosts[hostName].Options[IPName] = address
	}

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteHostAdd handles bash autocompletion for the 'host add' command
func CompleteHostAdd(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	if c.NArg() == 0 {
		fmt.Fprintln(c.App.Writer, strings.Join(sortHostNames(configData), "\n"))
	} else if c.NArg() == 1 {
		IPs, IPMap := sortAllIPs(configData)
		for _, IP := range IPs {
			fmt.Fprintf(c.App.Writer, "%s:%s\n", IP, IPMap[IP])
		}
	} else if c.NArg() == 2 {
		fmt.Fprintln(c.App.Writer, strings.Join(sortAllOptions(configData), "\n"))
	}

	flagCompletion(c)
}

func flagCompletion(c *cli.Context) {
	for _, flag := range c.App.Command("add").Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
