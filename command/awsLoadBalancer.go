package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/guywithnose/hostBuilder/awsUtil"
	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

// CmdAwsLoadBalancer adds aws load balancer information to the configuration
func CmdAwsLoadBalancer(util awsUtil.AwsInterface) func(*cli.Context) error {
	return func(c *cli.Context) error {
		return CmdAwsLoadBalancerHelper(c, util)
	}
}

// CmdAwsLoadBalancerHelper uses the given awsUtil to add aws load balancer information to the configuration
func CmdAwsLoadBalancerHelper(c *cli.Context, util awsUtil.AwsInterface) error {
	util.SetProfile(c.String("profile"))
	if c.NArg() != 0 {
		return cli.NewExitError("Usage: \"hostBuilder aws loadBalancers\"", 1)
	}

	loadBalancers, err := util.ReadAllLoadBalancers()
	if err != nil {
		return err
	}

	configData, err := loadConfig(c)
	if err != nil {
		return err
	}

	for name, address := range loadBalancers {
		IP, err := resolveAddress(address)
		if err != nil {
			return err
		}

		configData.GlobalIPs[name] = IP
	}

	return config.WriteConfig(c.GlobalString("config"), configData)
}

// CompleteAwsLoadBalancer handles bash autocompletion for the 'aws loadBalancers' command
func CompleteAwsLoadBalancer(util awsUtil.AwsInterface) func(c *cli.Context) {
	return func(c *cli.Context) {
		CompleteAwsLoadBalancerHelper(c, util)
	}
}

// CompleteAwsLoadBalancerHelper handles bash autocompletion for the 'aws instances' command
func CompleteAwsLoadBalancerHelper(c *cli.Context, util awsUtil.AwsInterface) {
	lastParam := os.Args[len(os.Args)-2]
	if lastParam == "--profile" {
		profiles, err := util.ListAllProfiles()
		if err != nil {
			return
		}

		for _, profile := range profiles {
			fmt.Fprintln(c.App.Writer, profile)
		}
	}

	for _, flag := range c.App.Command("loadBalancers").Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
