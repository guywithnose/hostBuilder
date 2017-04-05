package command

import (
	"flag"
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
