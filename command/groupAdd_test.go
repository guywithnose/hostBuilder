package command

import (
	"flag"
	"io/ioutil"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGroupAdd(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	assert.Nil(t, set.Parse([]string{"foo", "bar"}))
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGroupAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	assert.Equal(t, []string{"baz.com", "goo", "bar"}, modifiedConfigData.Groups["foo"])
}

func TestCmdGroupAddUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder group add {groupName} {hostName}\"")
}

func TestCmdGroupAddNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bar"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdGroupAddBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bar"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdGroupAddFirstGroup(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	configData := &config.HostsConfig{
		Hosts: map[string]config.Host{
			"bar":     {Current: hostIgnore},
			"baz.com": {Current: "baz", Options: map[string]string{"bazz": "10.0.0.7"}},
			"goo":     {Current: "foop", Options: map[string]string{"foop": "10.0.0.8"}},
		},
		GlobalIPs: map[string]string{"baz": "10.0.0.4"},
	}

	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	set := flag.NewFlagSet("test", 0)
	set.String("config", configFile.Name(), "doc")

	defer removeFile(t, configFile.Name())
	assert.Nil(t, set.Parse([]string{"food", "bar"}))
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGroupAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, []string{"bar"}, modifiedConfigData.Groups["food"])
}

func TestCmdGroupAddBadGroupName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"food", "bar"}))
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGroupAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	assert.Equal(t, []string{"bar"}, modifiedConfigData.Groups["food"])
}

func TestCmdGroupAddBadHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo", "bart"}))
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "Hostname bart does not exist")
}

func TestCmdGroupAddAlreadyExists(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo", "goo"}))

	app := cli.NewApp()
	c := cli.NewContext(app, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "Group foo already contains goo")
}

func TestCompleteGroupAddGroupName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	assert.Equal(t, "foo\n", writer.String())
}

func TestCompleteGroupAddHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	assert.Equal(t, "bar\nbaz.com\ngoo\n", writer.String())
}

func TestCompleteGroupAddComplete(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo", "bar"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	assert.Equal(t, "\n", writer.String())
}

func TestCompleteGroupAddNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	assert.Equal(t, "", writer.String())
}
