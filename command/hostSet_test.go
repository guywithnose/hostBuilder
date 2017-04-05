package command

import (
	"flag"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdHostSet(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"baz.com", "bazz"}))

	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdHostSet(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	assert.Equal(t, "bazz", modifiedConfigData.Hosts["baz.com"].Current, "baz.com was not set to baz")
}

func TestCmdHostSetUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder host set {hostName} {IPName}\"")
}

func TestCmdHostSetNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com", "bazz"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdHostSetBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com", "bazz"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdHostSetBadHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"food", "baz"}))
	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "HostName food does not exist")
}

func TestCmdHostSetBadIPName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"baz.com", "bar"}))
	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "IPName bar does not exist")
}

func TestCompleteHostSetHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	assert.Equal(t, "bar\nbaz.com\ngoo\n", writer.String())
}

func TestCompleteHostSetIPName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"baz.com"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	assert.Equal(t, "bazz:10.0.0.7\nbaz:10.0.0.4\nignore:\n", writer.String())
}

func TestCompleteHostBadHostname(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"bar.com"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	assert.Equal(t, "", writer.String())
}

func TestCompleteHostSetComplete(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"baz.com", "baz"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	assert.Equal(t, "", writer.String())
}

func TestCompleteHostSetNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	assert.Equal(t, "", writer.String())
}
