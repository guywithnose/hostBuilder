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
func CmdAwsLoadBalancer(c *cli.Context) error {
	return CmdAwsInstancesHelper(c, awsUtil.NewAwsUtil(c.String("profile")))
}

// CmdAwsLoadBalancerHelper uses the given awsUtil to add aws load balancer information to the configuration
func CmdAwsLoadBalancerHelper(c *cli.Context, util awsUtil.AwsInterface) error {
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
func CompleteAwsLoadBalancer(c *cli.Context) {
	CompleteAwsLoadBalancerHelper(c, awsUtil.NewAwsUtil(""))
}

// CompleteAwsLoadBalancerHelper handles bash autocompletion for the 'aws instances' command
func CompleteAwsLoadBalancerHelper(c *cli.Context, util awsUtil.AwsInterface) {
	lastParam := os.Args[len(os.Args)-2]

	profiles, err := util.ListAllProfiles()
	if err != nil {
		return
	}

	if lastParam == "--profile" {
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
