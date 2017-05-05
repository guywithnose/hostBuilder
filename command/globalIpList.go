package command

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdGlobalIPList adds an ip to the global ip list
func CmdGlobalIPList(c *cli.Context) error {
	if c.NArg() != 0 {
		return cli.NewExitError("Usage: \"hostBuilder globalIP list\"", 1)
	}

	configFile := c.GlobalString("config")
	if configFile == "" {
		return cli.NewExitError("You must specify a config file", 1)
	}

	configData, err := config.LoadConfigFromFile(configFile)
	if err != nil {
		return err
	}

	return printSortedIps(configData, c.App.Writer)
}

func printSortedIps(configData *config.HostsConfig, writer io.Writer) error {
	w := tabwriter.NewWriter(writer, 0, 0, 1, ' ', 0)
	ips := make([]string, 0, len(configData.GlobalIPs))
	for ip := range configData.GlobalIPs {
		ips = append(ips, ip)
	}

	sort.Strings(ips)

	for _, name := range ips {
		fmt.Fprintf(w, "%s\t%s\n", name, configData.GlobalIPs[name])
	}

	return w.Flush()
}
