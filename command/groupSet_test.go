package command

import (
	"flag"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGroupSet(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	assert.Nil(t, set.Parse([]string{"foo", "baz"}))
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGroupSet(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	assert.Equal(t, "baz", modifiedConfigData.Hosts["goo"].Current, "goo was not set to baz")
}

func TestCmdGroupSetUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder group set {groupName} {globalIPName}\"")
}

func TestCmdGroupSetNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "baz"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdGroupSetBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "baz"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdGroupSetBadGroupName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"food", "baz"}))
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "Group food does not exist")
}

func TestCmdGroupSetBadHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo", "barz"}))
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "Global IP barz does not exist")
}

func TestCompleteGroupSetGroupName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	assert.Equal(t, "foo\n", writer.String())
}

func TestCompleteGroupSetGlobalIPs(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	assert.Equal(t, "baz:10.0.0.4\n", writer.String())
}

func TestCompleteGroupSetComplete(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo", "baz"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	assert.Equal(t, "", writer.String())
}

func TestCompleteGroupSetNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	assert.Equal(t, "", writer.String())
}
