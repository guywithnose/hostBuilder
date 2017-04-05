package command

import (
	"flag"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdHostRemove(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	assert.Nil(t, set.Parse([]string{"goo", "foop"}))

	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdHostRemove(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedHost := config.Host{Current: "ignore", Options: map[string]string{}}
	assert.Equal(t, expectedHost, modifiedConfigData.Hosts["goo"])
}

func TestCmdHostRemoveUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)

	c := cli.NewContext(nil, set, nil)
	assert.EqualError(t, CmdHostRemove(c), "Usage: \"hostBuilder host remove {hostName} {IPName}\"")
}

func TestCmdHostRemoveNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo", "foop"}))
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdHostRemove(c), "You must specify a config file")
}

func TestCmdHostRemoveBadHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	assert.Nil(t, set.Parse([]string{"goop", "foop"}))

	c := cli.NewContext(nil, set, nil)
	assert.EqualError(t, CmdHostRemove(c), "Host goop does not exist")
}

func TestCmdHostRemoveBadIPName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)

	assert.Nil(t, set.Parse([]string{"goo", "foo"}))
	c := cli.NewContext(nil, set, nil)
	assert.EqualError(t, CmdHostRemove(c), "IPName foo does not exist")
}

func TestCompleteHostRemoveHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	assert.Equal(t, "bar\nbaz.com\ngoo\n", writer.String())
}

func TestCompleteHostRemoveBadHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"goop"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	assert.Equal(t, "", writer.String())
}

func TestCompleteHostRemoveIPName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"goo"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	assert.Equal(t, "foop:10.0.0.8\n", writer.String())
}

func TestCompleteHostRemoveComplete(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"goo", "foop"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	assert.Equal(t, "", writer.String())
}

func TestCompleteHostRemoveNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	assert.Equal(t, "", writer.String())
}
