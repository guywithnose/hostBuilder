package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/guywithnose/hostBuilder/awsUtil"
	"github.com/guywithnose/hostBuilder/command"
	"github.com/urfave/cli"
)

var profileFlag = cli.StringFlag{
	Name:   "profile, p",
	Usage:  "The AWS credentials profile to use",
	EnvVar: "HOST_BUILDER_AWS_PROFILE",
	Value:  "default",
}

var forceFlag = cli.BoolFlag{
	Name:  "force",
	Usage: "Overwrite existing",
}

// GlobalFlags defines flags that apply to all commands
var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "config, c",
		Usage:  "The path to your config file",
		EnvVar: "HOST_BUILDER_CONFIG_FILE",
	},
}

// Commands defines the commands that can be called on hostBuilder
var Commands = []cli.Command{
	{
		Name:         "createConfig",
		Aliases:      []string{"c"},
		Usage:        "Create a config file from an existing hosts file",
		Action:       command.CmdCreateConfig,
		BashComplete: command.CompleteCreateConfig,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostsFile, hosts",
				Usage: "The path to your hosts file",
				Value: "/etc/hosts",
			},
		},
	},
	{
		Name:         "build",
		Aliases:      []string{"b"},
		Usage:        "Builds your host file",
		Action:       command.CmdBuild,
		BashComplete: command.CompleteBuild,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "output, o",
				Usage:  "The path to write your hosts file",
				EnvVar: "HOST_BUILDER_OUTPUT_FILE",
			},
			cli.BoolFlag{
				Name:  "oneLinePerIP",
				Usage: "Put all hosts for an IP on the same line",
			},
		},
	},
	{
		Name:         "globalIP",
		Aliases:      []string{"gl"},
		Usage:        "Add things to the configuration",
		Category:     "Config",
		BashComplete: RootCompletion,
		Subcommands: []cli.Command{
			{
				Name:         "add",
				Aliases:      []string{"a"},
				Usage:        "Add a global IP to the configuration",
				Action:       command.CmdGlobalIPAdd,
				BashComplete: command.CompleteGlobalIPAdd,
				Flags:        []cli.Flag{forceFlag},
			},
			{
				Name:         "remove",
				Aliases:      []string{"r"},
				Usage:        "Remove a global IP from the configuration",
				Action:       command.CmdGlobalIPRemove,
				BashComplete: command.CompleteGlobalIPRemove,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List global IPs in the configuration",
				Action:  command.CmdGlobalIPList,
			},
		},
	},
	{
		Name:         "host",
		Aliases:      []string{"ho"},
		Usage:        "Modify hosts",
		Category:     "Config",
		BashComplete: RootCompletion,
		Subcommands: []cli.Command{
			{
				Name:         "add",
				Aliases:      []string{"a"},
				Usage:        "Add an IP to a hostname",
				Action:       command.CmdHostAdd,
				BashComplete: command.CompleteHostAdd,
				Flags:        []cli.Flag{forceFlag},
			},
			{
				Name:         "remove",
				Aliases:      []string{"r"},
				Usage:        "Remove an IP from a hostname",
				Action:       command.CmdHostRemove,
				BashComplete: command.CompleteHostRemove,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List the available hostnames",
				Action:  command.CmdHostList,
			},
			{
				Name:         "show",
				Aliases:      []string{"sh"},
				Usage:        "Describe the IPs on a hostname",
				Action:       command.CmdHostShow,
				BashComplete: command.CompleteHostShow,
			},
			{
				Name:         "set",
				Aliases:      []string{"se"},
				Usage:        "Set a hostname to a specific ip",
				Action:       command.CmdHostSet,
				BashComplete: command.CompleteHostSet,
			},
		},
	},
	{
		Name:         "group",
		Aliases:      []string{"gr"},
		Usage:        "Modify groups",
		Category:     "Config",
		BashComplete: RootCompletion,
		Subcommands: []cli.Command{
			{
				Name:         "add",
				Aliases:      []string{"a"},
				Usage:        "Add a hostname to a group",
				Action:       command.CmdGroupAdd,
				BashComplete: command.CompleteGroupAdd,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List available groups",
				Action:  command.CmdGroupList,
			},
			{
				Name:         "show",
				Aliases:      []string{"sh"},
				Usage:        "List the hostnames in a group",
				Action:       command.CmdGroupShow,
				BashComplete: command.CompleteHostShow,
			},
			{
				Name:         "set",
				Aliases:      []string{"se"},
				Usage:        "Set the hostnames in a group to a global ip",
				Action:       command.CmdGroupSet,
				BashComplete: command.CompleteGroupSet,
			},
		},
	},
	{
		Name:         "aws",
		Aliases:      []string{"a"},
		Usage:        "Add information from AWS to the configuration",
		Category:     "Config",
		BashComplete: RootCompletion,
		Subcommands: []cli.Command{
			{
				Name:         "loadBalancers",
				Aliases:      []string{"l", "lb"},
				Usage:        "Add load balancer information to the configuration",
				Action:       command.CmdAwsLoadBalancer(new(awsUtil.AwsUtil)),
				BashComplete: command.CompleteAwsLoadBalancer(new(awsUtil.AwsUtil)),
				Flags:        []cli.Flag{profileFlag},
			},
			{
				Name:         "instances",
				Aliases:      []string{"i"},
				Usage:        "Add instance information to the configuration",
				Action:       command.CmdAwsInstances(new(awsUtil.AwsUtil)),
				BashComplete: command.CompleteAwsInstances(new(awsUtil.AwsUtil)),
				Flags: []cli.Flag{
					profileFlag,
					cli.StringFlag{
						Name:   "template, t",
						Usage:  "The template to use for naming instance ips",
						EnvVar: "HOST_BUILDER_INSTANCE_TEMPLATE",
						Value:  "{{.InstanceId}}",
					},
				},
			},
		},
	},
}

// CommandNotFound runs when hostBuilder is invoked with an invalid command
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(c.App.Writer, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

// RootCompletion prints the list of root commands as the root completion method
// This is similar to the deafult method, but it excludes aliases
func RootCompletion(c *cli.Context) {
	lastParam := os.Args[len(os.Args)-2]
	if lastParam == "--config" {
		fmt.Println("fileCompletion")
		return
	}

	for _, command := range c.App.Commands {
		if command.Hidden {
			continue
		}

		fmt.Fprintf(c.App.Writer, "%s:%s\n", command.Name, command.Usage)
	}

	for _, flag := range c.App.Flags {
		name := strings.Split(flag.GetName(), ",")[0]
		if !c.IsSet(name) {
			fmt.Fprintf(c.App.Writer, "--%s\n", name)
		}
	}
}
