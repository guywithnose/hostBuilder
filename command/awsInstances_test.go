package command

import (
	"flag"
	"os"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdAwsInstances(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	c := cli.NewContext(nil, set, nil)
	util := new(awsTestUtil)
	instanceHosts := map[string]string{"foo": "127.0.0.1", "bar": "::1"}
	util.instances = instanceHosts
	assert.Nil(t, CmdAwsInstances(util)(c))

	configData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedIPs := instanceHosts
	expectedIPs["baz"] = "10.0.0.4"
	assert.Equal(t, expectedIPs, configData.GlobalIPs)
}

func TestCmdAwsInstancesBadTemplate(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	set.String("template", "{{.badTemplate", "doc")
	c := cli.NewContext(nil, set, nil)
	instanceHosts := map[string]string{}
	util := new(awsTestUtil)
	util.instances = instanceHosts
	assert.EqualError(t, CmdAwsInstances(util)(c), "template: :1: unclosed action")
}

func TestCmdAwsInstancesNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	instanceHosts := map[string]string{}
	util := new(awsTestUtil)
	util.instances = instanceHosts
	assert.EqualError(t, CmdAwsInstances(util)(c), "You must specify a config file")
}

func TestCmdAwsInstancesUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))
	c := cli.NewContext(nil, set, nil)
	instanceHosts := map[string]string{}
	util := new(awsTestUtil)
	util.instances = instanceHosts
	assert.EqualError(t, CmdAwsInstances(util)(c), "Usage: \"hostBuilder aws instances\"")
}

func TestCmdAwsInstancesAwsError(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	util := new(awsTestUtil)
	util.throwError = true
	assert.EqualError(t, CmdAwsInstances(util)(c), "error")
}

func TestCompleteAwsInstances(t *testing.T) {
	os.Args = []string{"aws", "instances", "--profile", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	CompleteAwsInstances(util)(c)

	assert.Equal(t, "profile1\nprofile2\n", writer.String())
}

func TestCompleteAwsInstancesNoProfileParam(t *testing.T) {
	os.Args = []string{"aws", "instances", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
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
	CompleteAwsInstances(util)(c)

	assert.Equal(t, "--profile\n--template\n", writer.String())
}

func TestCompleteAwsInstancesAwsError(t *testing.T) {
	os.Args = []string{"aws", "instances", "--profile", "--bash-completion"}
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	app.Commands = []cli.Command{{Name: "instances"}}
	util := new(awsTestUtil)
	util.profiles = []string{"profile1", "profile2"}
	util.throwError = true
	CompleteAwsInstances(util)(c)

	assert.Equal(t, "", writer.String())
}
