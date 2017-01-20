package command

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// CmdHostList lists the hostNames in the configuration
func CmdHostList(c *cli.Context) error {
	if c.NArg() != 0 {
		return cli.NewExitError("Usage: \"hostBuilder host list\"", 1)
	}

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	hostNames := sortHostNames(configData)

	fmt.Fprintln(c.App.Writer, strings.Join(hostNames, "\n"))

	return nil
}
