package command

import (
	"fmt"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdGlobalIPAdd adds an ip to the global ip list
func CmdGlobalIPAdd(c *cli.Context) error {
	if c.NArg() != 2 {
		return cli.NewExitError("Usage: \"hostBuilder globalIP add {Name} {address}\"", 1)
	}

	name := c.Args().Get(0)
	address, err := resolveAddress(c.Args().Get(1))
	if err != nil {
		return err
	}

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	if current, exists := configData.GlobalIPs[name]; exists {
		if c.Bool("force") {
			fmt.Fprintf(c.App.ErrWriter, "Warning: Overwriting %s (%s => %s)", name, current, address)
		} else {
			return cli.NewExitError(fmt.Sprintf("Global IP %s already exists", name), 1)
		}
	}

	configData.GlobalIPs[name] = address

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteGlobalIPAdd handles bash autocompletion for the 'globalIP add' command
func CompleteGlobalIPAdd(c *cli.Context) {
	for _, flag := range c.App.Command("add").Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
