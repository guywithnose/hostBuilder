package command

import (
	"fmt"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdGlobalIPRemove adds an ip to the global ip list
func CmdGlobalIPRemove(c *cli.Context) error {
	if c.NArg() != 1 {
		return cli.NewExitError("Usage: \"hostBuilder globalIP remove {Name}\"", 1)
	}

	globalIPName := c.Args().Get(0)

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	if _, exists := configData.GlobalIPs[globalIPName]; !exists {
		return cli.NewExitError(fmt.Sprintf("GlobalIP %s does not exist", globalIPName), 1)
	}

	delete(configData.GlobalIPs, globalIPName)

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteGlobalIPRemove handles bash autocompletion for the 'globalIP remove' command
func CompleteGlobalIPRemove(c *cli.Context) {
	configData, err := loadConfig(c)
	if err != nil {
		return
	}

	IPNames := sortGlobalIPNames(configData)
	for _, IPName := range IPNames {
		fmt.Fprintf(c.App.Writer, "%s:%s\n", IPName, configData.GlobalIPs[IPName])
	}
}
