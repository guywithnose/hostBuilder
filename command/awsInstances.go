package command

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/guywithnose/hostBuilder/awsUtil"
	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdAwsInstances adds aws instance information to the configuration
func CmdAwsInstances(util awsUtil.AwsInterface) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return CmdAwsInstancesHelper(c, util)
	}
}

// CmdAwsInstancesHelper uses the given awsUtil to add aws instance information to the configuration
func CmdAwsInstancesHelper(c *cli.Context, util awsUtil.AwsInterface) error {
	util.SetProfile(c.String("profile"))
	if c.NArg() != 0 {
		return cli.NewExitError("Usage: \"hostBuilder aws instances\"", 1)
	}

	templ, err := template.New("").Parse(c.String("template"))
	if err != nil {
		return err
	}

	instances, err := util.ReadAllInstances(templ)
	if err != nil {
		return err
	}

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	for name, IP := range instances {
		configData.GlobalIPs[name] = IP
	}

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteAwsInstances handles bash autocompletion for the 'aws instances' command
func CompleteAwsInstances(util awsUtil.AwsInterface) func(c *cli.Context) {
	return func(c *cli.Context) {
		CompleteAwsInstancesHelper(c, util)
	}
}

// CompleteAwsInstancesHelper handles bash autocompletion for the 'aws instances' command
func CompleteAwsInstancesHelper(c *cli.Context, util awsUtil.AwsInterface) {
	lastParam := os.Args[len(os.Args)-2]
	if lastParam == "--profile" {
		profiles, err := util.ListAllProfiles()
		if err != nil {
			return
		}

		for _, profile := range profiles {
			fmt.Fprintln(c.App.Writer, profile)
		}

		return
	}

	for _, flag := range c.App.Command("instances").Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
