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

func TestCmdAwsInstancesHelper(t *testing.T) {
	configFileName := setupBaseConfigFile(t)

	set := flag.NewFlagSet("test", 0)

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	instanceHosts := map[string]string{"foo": "127.0.0.1", "bar": "::1"}
	util := new(awsTestUtil)
	util.instances = instanceHosts
	assert.Nil(t, CmdAwsInstancesHelper(c, util))

	configData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedIPs := map[string]string{"baz": "10.0.0.4", "foo": "127.0.0.1", "bar": "::1"}
	if !reflect.DeepEqual(configData.GlobalIPs, expectedIPs) {
		t.Fatalf("Global IPs was \n%v\n, expected \n%v\n", configData.GlobalIPs, expectedIPs)
	}
	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdAwsInstancesHelperBadTemplate(t *testing.T) {
	configFileName := setupBaseConfigFile(t)

	set := flag.NewFlagSet("test", 0)

	set.String("config", configFileName, "doc")
	set.String("template", "{{.badTemplate", "doc")
	c := cli.NewContext(nil, set, nil)
	instanceHosts := map[string]string{}
	util := new(awsTestUtil)
	util.instances = instanceHosts
	assert.EqualError(t, CmdAwsInstancesHelper(c, util), "template: :1: unclosed action")

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdAwsInstancesHelperNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	instanceHosts := map[string]string{}
	util := new(awsTestUtil)
	util.instances = instanceHosts
	assert.EqualError(t, CmdAwsInstancesHelper(c, util), "You must specify a config file")
}

func TestCmdAwsInstancesHelperUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))
	c := cli.NewContext(nil, set, nil)
	instanceHosts := map[string]string{}
	util := new(awsTestUtil)
	util.instances = instanceHosts
	assert.EqualError(t, CmdAwsInstancesHelper(c, util), "Usage: \"hostBuilder aws instances\"")
}

func TestCmdAwsInstancesHelperAwsError(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	util := new(awsTestUtil)
	util.throwError = true
	assert.EqualError(t, CmdAwsInstancesHelper(c, util), "error")
}

func TestCompleteAwsInstancesHelper(t *testing.T) {
	os.Args = []string{"aws", "instances", "--profile", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	writer := new(bytes.Buffer)
	app := cli.NewApp()
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	CompleteAwsInstancesHelper(c, util)

	expectedOutput := "profile1\nprofile2\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}

func TestCompleteAwsInstancesHelperNoProfileParam(t *testing.T) {
	os.Args = []string{"aws", "instances", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	writer := new(bytes.Buffer)
	app := cli.NewApp()
	app.Writer = writer
	app.Commands = []cli.Command{
		{
			Name: "instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "profile, p",
				},
				cli.StringFlag{
					Name: "template, t",
				},
			},
		},
	}
	c := cli.NewContext(app, set, nil)
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	CompleteAwsInstancesHelper(c, util)

	expectedOutput := "--profile\n--template\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}

func TestCompleteAwsInstancesHelperAwsError(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	writer := new(bytes.Buffer)
	app := cli.NewApp()
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	app.Commands = []cli.Command{{Name: "instances"}}
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	util.throwError = true
	CompleteAwsInstancesHelper(c, util)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
