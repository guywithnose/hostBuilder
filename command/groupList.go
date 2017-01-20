package command

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// CmdGroupList lists the groups in the configuration
func CmdGroupList(c *cli.Context) error {
	if c.NArg() != 0 {
		return cli.NewExitError("Usage: \"hostBuilder group list\"", 1)
	}

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	groupNames := sortGroupNames(configData)

	fmt.Fprintln(c.App.Writer, strings.Join(groupNames, "\n"))

	return nil
}
