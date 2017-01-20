package command

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdAwsLoadBalancerHelper(t *testing.T) {
	configFileName := setupBaseConfigFile(t)

	set := flag.NewFlagSet("test", 0)

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	loadBalancers := map[string]string{"foo": "localhost4", "bar": "localhost6"}
	util := new(awsTestUtil)
	util.loadBalancers = loadBalancers
	assert.Nil(t, CmdAwsLoadBalancerHelper(c, util))

	configData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedIPs := map[string]string{"baz": "10.0.0.4", "foo": "127.0.0.1", "bar": "::1"}
	if !reflect.DeepEqual(configData.GlobalIPs, expectedIPs) {
		t.Fatalf("Global IPs was \n%v\n, expected \n%v\n", configData.GlobalIPs, expectedIPs)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdAwsLoadBalancerHelperUnresolvedHostname(t *testing.T) {
	configFileName := setupBaseConfigFile(t)

	set := flag.NewFlagSet("test", 0)

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	loadBalancers := map[string]string{"foo": "notahost"}
	util := new(awsTestUtil)
	util.loadBalancers = loadBalancers
	assert.EqualError(t, CmdAwsLoadBalancerHelper(c, util), "Unable to resolve notahost")

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdAwsLoadBalancerHelperNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	loadBalancers := map[string]string{}
	util := new(awsTestUtil)
	util.loadBalancers = loadBalancers
	assert.EqualError(t, CmdAwsLoadBalancerHelper(c, util), "You must specify a config file")
}

func TestCmdAwsLoadBalancerHelperUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))
	c := cli.NewContext(nil, set, nil)
	loadBalancers := map[string]string{}
	util := new(awsTestUtil)
	util.loadBalancers = loadBalancers
	assert.EqualError(t, CmdAwsLoadBalancerHelper(c, util), "Usage: \"hostBuilder aws loadBalancers\"")
}

func TestCmdAwsLoadBalancerHelperAwsError(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	util := new(awsTestUtil)
	util.throwError = true
	assert.EqualError(t, CmdAwsLoadBalancerHelper(c, util), "error")
}

func TestCompleteAwsLoadBalancersHelper(t *testing.T) {
	os.Args = []string{"aws", "loadBalancers", "--profile", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	writer := new(bytes.Buffer)
	app := cli.NewApp()
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	app.Commands = []cli.Command{{Name: "loadBalancers"}}
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	CompleteAwsLoadBalancerHelper(c, util)

	expectedOutput := "profile1\nprofile2\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}

func TestCompleteAwsLoadBalancersHelperNoProfileParam(t *testing.T) {
	os.Args = []string{"aws", "loadBalancers", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	writer := new(bytes.Buffer)
	app := cli.NewApp()
	app.Writer = writer
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
	CompleteAwsLoadBalancerHelper(c, util)

	expectedOutput := "--profile\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}

func TestCompleteAwsLoadBalancersHelperAwsError(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	writer := new(bytes.Buffer)
	app := cli.NewApp()
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	app.Commands = []cli.Command{{Name: "loadBalancers"}}
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	util.throwError = true
	CompleteAwsLoadBalancerHelper(c, util)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
