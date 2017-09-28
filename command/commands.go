package command

import (
	"github.com/guywithnose/hostBuilder/awsUtil"
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
		Action:       CmdCreateConfig,
		BashComplete: CompleteCreateConfig,
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
		Action:       CmdBuild,
		BashComplete: CompleteBuild,
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
				Action:       CmdGlobalIPAdd,
				BashComplete: CompleteGlobalIPAdd,
				Flags:        []cli.Flag{forceFlag},
			},
			{
				Name:         "remove",
				Aliases:      []string{"r"},
				Usage:        "Remove a global IP from the configuration",
				Action:       CmdGlobalIPRemove,
				BashComplete: CompleteGlobalIPRemove,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List global IPs in the configuration",
				Action:  CmdGlobalIPList,
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
				Action:       CmdHostAdd,
				BashComplete: CompleteHostAdd,
				Flags:        []cli.Flag{forceFlag},
			},
			{
				Name:         "remove",
				Aliases:      []string{"r"},
				Usage:        "Remove an IP from a hostname",
				Action:       CmdHostRemove,
				BashComplete: CompleteHostRemove,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List the available hostnames",
				Action:  CmdHostList,
			},
			{
				Name:         "show",
				Aliases:      []string{"sh"},
				Usage:        "Describe the IPs on a hostname",
				Action:       CmdHostShow,
				BashComplete: CompleteHostShow,
			},
			{
				Name:         "set",
				Aliases:      []string{"se"},
				Usage:        "Set a hostname to a specific ip",
				Action:       CmdHostSet,
				BashComplete: CompleteHostSet,
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
				Action:       CmdGroupAdd,
				BashComplete: CompleteGroupAdd,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List available groups",
				Action:  CmdGroupList,
			},
			{
				Name:         "show",
				Aliases:      []string{"sh"},
				Usage:        "List the hostnames in a group",
				Action:       CmdGroupShow,
				BashComplete: CompleteGroupShow,
			},
			{
				Name:         "set",
				Aliases:      []string{"se"},
				Usage:        "Set the hostnames in a group to a global ip",
				Action:       CmdGroupSet,
				BashComplete: CompleteGroupSet,
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
				Action:       CmdAwsLoadBalancer(new(awsUtil.AwsUtil)),
				BashComplete: CompleteAwsLoadBalancer(new(awsUtil.AwsUtil)),
				Flags:        []cli.Flag{profileFlag},
			},
			{
				Name:         "instances",
				Aliases:      []string{"i"},
				Usage:        "Add instance information to the configuration",
				Action:       CmdAwsInstances(new(awsUtil.AwsUtil)),
				BashComplete: CompleteAwsInstances(new(awsUtil.AwsUtil)),
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
