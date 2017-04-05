package command

import (
	"flag"
	"os"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdAwsLoadBalancer(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	c := cli.NewContext(nil, set, nil)
	util := new(awsTestUtil)
	util.loadBalancers = map[string]string{"foo": "localhost4", "bar": "localhost6"}
	assert.Nil(t, CmdAwsLoadBalancer(util)(c))

	configData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedIPs := map[string]string{"baz": "10.0.0.4", "foo": "127.0.0.1", "bar": "::1"}
	assert.Equal(t, expectedIPs, configData.GlobalIPs)
}

func TestCmdAwsLoadBalancerUnresolvedHostname(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	c := cli.NewContext(nil, set, nil)
	loadBalancers := map[string]string{"foo": "notahost"}
	util := new(awsTestUtil)
	util.loadBalancers = loadBalancers
	assert.EqualError(t, CmdAwsLoadBalancer(util)(c), "Unable to resolve notahost")
}

func TestCmdAwsLoadBalancerNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	loadBalancers := map[string]string{}
	util := new(awsTestUtil)
	util.loadBalancers = loadBalancers
	assert.EqualError(t, CmdAwsLoadBalancer(util)(c), "You must specify a config file")
}

func TestCmdAwsLoadBalancerUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))
	c := cli.NewContext(nil, set, nil)
	loadBalancers := map[string]string{}
	util := new(awsTestUtil)
	util.loadBalancers = loadBalancers
	assert.EqualError(t, CmdAwsLoadBalancer(util)(c), "Usage: \"hostBuilder aws loadBalancers\"")
}

func TestCmdAwsLoadBalancerAwsError(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	util := new(awsTestUtil)
	util.throwError = true
	assert.EqualError(t, CmdAwsLoadBalancer(util)(c), "error")
}

func TestCompleteAwsLoadBalancers(t *testing.T) {
	os.Args = []string{"aws", "loadBalancers", "--profile", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	app.Commands = []cli.Command{{Name: "loadBalancers"}}
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	CompleteAwsLoadBalancer(util)(c)

	assert.Equal(t, "profile1\nprofile2\n", writer.String())
}

func TestCompleteAwsLoadBalancersNoProfileParam(t *testing.T) {
	os.Args = []string{"aws", "loadBalancers", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	app.Commands = []cli.Command{
		{
			Name: "loadBalancers",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "profile, p",
				},
			},
		},
	}
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	CompleteAwsLoadBalancer(util)(c)

	assert.Equal(t, "--profile\n", writer.String())
}

func TestCompleteAwsLoadBalancersAwsError(t *testing.T) {
	os.Args = []string{"aws", "loadBalancers", "--profile", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	app.Commands = []cli.Command{{Name: "loadBalancers"}}
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	util.throwError = true
	CompleteAwsLoadBalancer(util)(c)

	assert.Equal(t, "", writer.String())
}
